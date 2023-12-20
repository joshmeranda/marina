package auth

import (
	"context"

	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/joshmeranda/marina/pkg/gateway/drivers/secret"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordFieldName = "value"
)

type Local struct {
	secretStore secret.Driver
}

func NewLocal(secretStore secret.Driver) Driver {
	return &Local{
		secretStore: secretStore,
	}
}

func (d *Local) Authenticate(ctx context.Context, req *auth.LoginRequest) error {
	data, err := d.secretStore.Get(ctx, req.User, passwordFieldName)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(data, []byte(req.Secret)); err != nil {
		return err
	}

	return nil
}
