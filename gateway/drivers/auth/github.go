package auth

import (
	"context"

	"github.com/google/go-github/v57/github"
	"github.com/joshmeranda/marina/gateway/api/auth"
)

type Github struct{}

func NewGithub() Driver {
	return &Github{}
}

func (d Github) Authenticate(ctx context.Context, req *auth.LoginRequest) error {
	client := github.NewClient(nil).WithAuthToken(string(req.Secret))

	_, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return err
	}

	return nil
}
