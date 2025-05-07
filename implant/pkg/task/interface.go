package task

type Tasker interface {
	Do() Result // TODO: should Do() return *Result?
	GetID() uint
}

type Metadata struct {
}
