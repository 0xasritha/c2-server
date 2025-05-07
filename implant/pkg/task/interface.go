package task

type Tasker interface {
	Do() Result
	GetID() uint
}

type Metadata struct {
}
