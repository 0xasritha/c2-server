package results

type Result struct {
	id        uint32
	Content   []string `json:"contents"`
	IsSuccess bool     `json:"success"`
	TaskId    uint32   `json:"task_id"`
}
