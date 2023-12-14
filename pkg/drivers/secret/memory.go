package secret

import (
	"context"
	"fmt"
)

type Secret map[string][]byte

type memoryDriver struct {
	data map[string]Secret
}

func NewMemoryDriver(data map[string]Secret) Driver {
	return &memoryDriver{
		data: data,
	}
}

func (d *memoryDriver) Get(_ context.Context, name string, key string) ([]byte, error) {
	secret, ok := d.data[name]
	if !ok {
		return nil, fmt.Errorf("no such secret")
	}

	value, ok := secret[key]
	if !ok {
		return nil, fmt.Errorf("no such field")
	}

	return value, nil
}
