package store

import "errors"

var (
	NotFoundErr = errors.New("not found")
)

type Stoerer[T any] interface {
	Get(id uint32) (T, error)
	List() (map[uint32]T, error)
	Add(id uint32, newT T) error
	Update(id uint32, updatedT T) error
	Delete(id uint32) error
}
