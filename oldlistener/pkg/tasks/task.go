package tasks

type Task struct {
	id       uint32
	TaskType string            `json:"type"`
	Options  map[string]string `json:"options"`
}
