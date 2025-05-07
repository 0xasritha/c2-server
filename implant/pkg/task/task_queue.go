package task

type TaskQueue struct {
	tasks []Tasker
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{tasks: make([]Tasker, 0)}
}

func (q *TaskQueue) Enqueue(t Tasker) {
	q.tasks = append(q.tasks, t)
}

// Dequeue returns the front Tasker and removes it from the queue
func (q *TaskQueue) Dequeue() Tasker {
	if len(q.tasks) == 0 {
		return nil
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task
}
func (q *TaskQueue) IsEmpty() bool {
	return len(q.tasks) == 0
}
