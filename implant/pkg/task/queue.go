package task

type Queue struct {
	tasks []Tasker
}

func NewQueue() *Queue {
	return &Queue{tasks: make([]Tasker, 0)}
}

func (q *Queue) Enqueue(t Tasker) {
	q.tasks = append(q.tasks, t)
}

// Dequeue returns the front Tasker and removes it from the queue
func (q *Queue) Dequeue() Tasker {
	if len(q.tasks) == 0 {
		return nil
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task
}
func (q *Queue) IsEmpty() bool {
	return len(q.tasks) == 0
}
