package task

type Ping struct {
	ID uint `json:"id"`
}

func NewPing(ID uint) *Ping {
	return &Ping{
		ID: ID,
	}
}

func (p *Ping) Do() Result {
	return Result{Content: "pong", Success: true}
}

func (p *Ping) GetID() uint {
	return p.ID
}
