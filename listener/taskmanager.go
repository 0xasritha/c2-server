package main

type Task struct {
	uuid     uint32
	TaskType string            `json:"type"`
	Options  map[string]string `json:"options"`
}

type TaskManager interface {
	// getTasks will get all Task objects in database, convert them to JSON, put them in a variable, and return that variable #TODO: make this comment better
	getTasks() []Task // TODO: do map of uuid: Task
	createTask(task Task) error
}

type LocalTaskStore struct {
	tasks []Task
}

func NewLocalTaskStore() *LocalTaskStore {
	lts := &LocalTaskStore{
		tasks: make([]Task, 0),
	}

	lts.fillWithMockData()

	return lts
}

func (lts *LocalTaskStore) fillWithMockData() {
}

// TODO: using method pointer receiver correctly?
func (lts *LocalTaskStore) getTasks() []Task {
	return lts.tasks
}

func (lts *LocalTaskStore) createTask(task Task) error {
	lts.tasks = append(lts.tasks, task)
	return nil
}
