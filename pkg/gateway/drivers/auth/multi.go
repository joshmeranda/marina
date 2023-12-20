package auth

import (
	"context"
	"fmt"

	"github.com/joshmeranda/marina/pkg/apis/auth"
)

type MultiAuth struct {
	Drivers map[auth.SecretType]Driver
}

func (d *MultiAuth) Authenticate(ctx context.Context, req *auth.LoginRequest) error {
	driver, ok := d.Drivers[req.SecretType]
	if !ok {
		return fmt.Errorf("unsupported secret type: %s", req.SecretType)
	}

	return driver.Authenticate(ctx, req)
}
