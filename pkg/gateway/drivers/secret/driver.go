package secret

import "context"

type Driver interface {
	Get(ctx context.Context, name string, key string) ([]byte, error)

	Set(ctx context.Context, name string, key string, value []byte) error
}
