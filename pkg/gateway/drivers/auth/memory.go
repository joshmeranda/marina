package auth

import (
	"context"
	"fmt"
	"slices"

	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/joshmeranda/marina/pkg/gateway/drivers/storage"
)

type Memory struct {
	store storage.KeyValueStore[string, []byte]
}

func NewMemory() Driver {
	return &Memory{
		store: storage.NewMemoryStore[string, []byte](),
	}
}

func (d *Memory) Authenticate(ctx context.Context, req *auth.LoginRequest) error {
	password, err := d.store.Get(ctx, req.User)
	if err != nil {
		return err
	}

	if slices.Equal(password, req.Secret) {
		return fmt.Errorf("password is not valid for user %s", req.User)
	}

	return nil
}
