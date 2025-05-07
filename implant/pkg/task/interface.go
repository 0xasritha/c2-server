package task

type Tasker interface {
	Execute() Result
	GetID() uint
}

type TaskMetadata struct {
	ID uint `json:"id"`
}
