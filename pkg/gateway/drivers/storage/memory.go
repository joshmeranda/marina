package storage

import (
	"context"
	"sync"
)

type memoryStore[K comparable, V any] struct {
	dataMu sync.RWMutex
	data   map[K]V
}

func NewMemoryStore[K comparable, V any]() KeyValueStore[K, V] {
	return &memoryStore[K, V]{
		data: make(map[K]V),
	}
}

func (s *memoryStore[K, V]) Get(ctx context.Context, key K) (V, error) {
	s.dataMu.RLock()
	defer s.dataMu.RUnlock()
	return s.data[key], nil
}

func (s *memoryStore[K, V]) Set(ctx context.Context, key K, value V) error {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	s.data[key] = value
	return nil
}

func (s *memoryStore[K, V]) Delete(ctx context.Context, key K) error {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	delete(s.data, key)
	return nil
}
