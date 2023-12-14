package store

import (
	"errors"
)

var (
	ErrNotFound = errors.New("key not found")
)

type KeyValueStore[K any, V any] interface {
	Get(key K) (V, error)

	Set(key K, value V) error

	Delete(key K, value V) error
}
