package auth

import (
	"context"

	"github.com/joshmeranda/marina/apis/auth"
)

type Driver interface {
	Authenticate(ctx context.Context, req *auth.LoginRequest) error
}
