package client

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joshmeranda/marina/pkg/apis"
	"google.golang.org/grpc"
)

var _ apis.MarinaClient = &Client{}

type Client struct{}

func (c *Client) CreateTerminal(ctx context.Context, req *apis.TerminalCreateRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
