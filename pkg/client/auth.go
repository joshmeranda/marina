package client

import (
	"context"
	"fmt"

	"github.com/cli/oauth"
	"github.com/cli/oauth/api"
	"github.com/joshmeranda/marina/pkg/apis/auth"
	"google.golang.org/grpc"
)

var _ auth.AuthServiceClient = &Client{}

const (
	ClientId string = "614fad3dd8cd7deb6892"
)

func (g *Client) getAccessToken(ctx context.Context, req *auth.LoginRequest) (*api.AccessToken, error) {
	flow := oauth.Flow{
		Host:     oauth.GitHubHost("https://github.com"),
		Scopes:   []string{"read:user"},
		ClientID: ClientId,
	}

	token, err := flow.DeviceFlow()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (g *Client) Login(ctx context.Context, req *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	if req.Token == "" {
		token, err := g.getAccessToken(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("error getting access token: %w", err)
		}

		req.Token = token.Token
	}

	resp, err := g.authClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: resp.Token,
	}, nil
}
