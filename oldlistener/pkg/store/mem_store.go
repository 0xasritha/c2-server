package store

type MemStore[T any] struct {
	// TODO: rename to data ?? think of a better name
	list map[uint32]T // keys are IDs
}

func NewLocalFileStore[T any]() *MemStore[T] {
	list := make(map[uint32]T)
	return &MemStore[T]{list}
}

func (ms *MemStore[T]) Get(id uint32) (T, error) {
	if val, ok := ms.list[id]; ok {
		return val, nil
	}

	return *new(T), NotFoundErr
}

func (ms *MemStore[T]) List() (map[uint32]T, error) {
	return ms.list, nil
}

func (ms *MemStore[T]) Add(id uint32, newT T) error {
	ms.list[id] = newT
	return nil
}

func (ms *MemStore[T]) Update(id uint32, updatedT T) error {
	if _, ok := ms.list[id]; ok {
		ms.list[id] = updatedT
		return nil
	}

	return NotFoundErr
}

func (ms *MemStore[T]) Delete(id uint32) error {
	delete(ms.list, id)
	return nil
}
