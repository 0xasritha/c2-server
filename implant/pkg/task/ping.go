package task

type Ping struct {
	TaskMetadata
}

func NewPing(ID uint) *Ping {
	return &Ping{
		TaskMetadata{
			ID: ID,
		},
	}
}

func (p *Ping) Execute() Result {
	return Result{Content: "pong", Success: true}
}

func (p *Ping) GetID() uint {
	return p.ID
}
