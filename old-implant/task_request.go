package main

type TaskRequest struct {
	ID        uint           // what is the JSON for this?
	Completed bool           `json:"completed"`
	TaskType  string         `json:"type"`
	Options   map[string]any `json:"options"`
}
