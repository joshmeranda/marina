package auth

import (
	"context"

	"github.com/joshmeranda/marina/pkg/apis/auth"
)

type Driver interface {
	Authenticate(ctx context.Context, req *auth.LoginRequest) error
}
