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

func (c *Client) getAccessToken(ctx context.Context, req *auth.LoginRequest) (*api.AccessToken, error) {
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

func (c *Client) Login(ctx context.Context, req *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	if req.Token == "" {
		c.logger.Info("no access token found, must authenticate with oauth provider")
		token, err := c.getAccessToken(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("error getting access token: %w", err)
		}

		req.Token = token.Token
	}

	c.logger.Info("requesting token from server")

	resp, err := c.authClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: resp.Token,
	}, nil
}
