package store

type KeyValueStore[K any, V any] interface {
	Get(key K) (V, error)

	Set(key K, value V) error

	Delete(key K, value V) error
}
