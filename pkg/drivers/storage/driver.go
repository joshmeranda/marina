package storage

import "context"

type KeyValueStore[K any, V any] interface {
	Get(ctx context.Context, key K) (V, error)

	Set(ctx context.Context, key K, value V) error

	Delete(ctx context.Context, key K) error
}
