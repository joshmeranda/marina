package secret

import (
	"context"
	"fmt"
	"sync"
)

type Secret map[string][]byte

type memoryDriver struct {
	dataMut sync.RWMutex
	data    map[string]Secret
}

func NewMemoryDriver(data map[string]Secret) Driver {
	return &memoryDriver{
		data: data,
	}
}

func (d *memoryDriver) Get(_ context.Context, name string, key string) ([]byte, error) {
	d.dataMut.RLock()
	defer d.dataMut.RUnlock()

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

func (d *memoryDriver) Set(_ context.Context, name string, key string, value []byte) error {
	d.dataMut.Lock()
	defer d.dataMut.Unlock()

	secret, ok := d.data[name]
	if !ok {
		return fmt.Errorf("no such secret")
	}

	secret[key] = value

	return nil
}
