package client

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"google.golang.org/grpc"
)

var _ terminal.TerminalClient = &Client{}

func (c *Client) CreateTerminal(ctx context.Context, req *terminal.TerminalCreateRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return c.terminalClient.CreateTerminal(ctx, req, opts...)
}
