package store

import "sync"

type memoryStore[K comparable, V any] struct {
	innerMu sync.RWMutex
	inner   map[K]V
}

func NewMemoryStore[K comparable, V any]() *memoryStore[K, V] {
	return &memoryStore[K, V]{
		inner: make(map[K]V),
	}
}

func (m *memoryStore[K, V]) Get(key K) (V, error) {
	m.innerMu.RLock()
	defer m.innerMu.RUnlock()

	value, ok := m.inner[key]
	if !ok {
		return value, ErrNotFound
	}

	return value, nil
}

func (m *memoryStore[K, V]) Set(key K, value V) error {
	m.innerMu.Lock()
	defer m.innerMu.Unlock()

	m.inner[key] = value

	return nil
}

func (m *memoryStore[K, V]) Delete(key K, value V) error {
	m.innerMu.Lock()
	defer m.innerMu.Unlock()

	delete(m.inner, key)

	return nil
}
