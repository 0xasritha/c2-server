package task

import (
	"errors"
	"slices"
)

var (
	NotFoundError = errors.New("task not found")
)

type Store interface {
	Add(t Task) error
	AddTasks([]Task) error
	Get(id uint32) (Task, error)
	List() ([]Task, error)
	Update(id uint32, t Task) error
	Remove(id uint32) error
}

type MemStore []Task

func NewMemStore() MemStore {
	return make(MemStore, 0)
}

func (m *MemStore) Add(t Task) error {
	*m = append(*m, t)
	return nil
}

func (m *MemStore) AddTasks(newTasks []Task) error {
	*m = append(*m, newTasks...)
	return nil
}

func (m *MemStore) Get(id uint32) (Task, error) {
	if id > uint32(len(*m)-1) {
		return Task{}, NotFoundError
	}

	return (*m)[id], nil
}

func (m *MemStore) List() ([]Task, error) {
	return *m, nil
}

func (m *MemStore) Update(id uint32, t Task) error {
	if id > uint32(len(*m)-1) {
		return NotFoundError
	}

	(*m)[int(id)] = t
	return nil
}

func (m *MemStore) Remove(id uint32) error {
	if id > uint32(len(*m)-1) {
		return NotFoundError
	}

	*m = slices.Delete(*m, int(id), int(id)+1)
	return nil
}
