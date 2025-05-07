package task

type Result struct {
	Content string `json:"content"`
	Success bool   `json:"success"`
	TaskID  uint   `json:"task_id"`
}
