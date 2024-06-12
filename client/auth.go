package client

import (
	"context"
	"fmt"

	"github.com/cli/oauth"
	"github.com/cli/oauth/api"
	"github.com/joshmeranda/marina/gateway/api/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ auth.AuthServiceClient = &Client{}

const (
	ClientId string = "614fad3dd8cd7deb6892"
)

func TokenAuthInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "token", token)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

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

func (c *Client) githubLogin(ctx context.Context, req *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	if len(req.Secret) == 0 {
		c.logger.Info("no access token found, must authenticate with oauth provider")
		token, err := c.getAccessToken(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("error getting access token: %w", err)
		}

		req.Secret = []byte(token.Token)
	}

	resp, err := c.authClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: resp.Token,
	}, nil
}

func (c *Client) marinaLogin(ctx context.Context, req *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	resp, err := c.authClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: resp.Token,
	}, nil
}

func (c *Client) Login(ctx context.Context, req *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	switch req.SecretType {
	case auth.SecretType_Github:
		return c.githubLogin(ctx, req, opts...)
	case auth.SecretType_Password:
		return c.marinaLogin(ctx, req, opts...)
	default:
		return nil, fmt.Errorf("recevied unknown token kind: %s", req.SecretType.String())
	}
}
