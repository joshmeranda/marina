package client

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joshmeranda/marina/apis/terminal"
	"google.golang.org/grpc"
)

var _ terminal.TerminalServiceClient = &Client{}

func (c *Client) CreateTerminal(ctx context.Context, req *terminal.TerminalCreateRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	_, err := c.terminalClient.CreateTerminal(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (c *Client) DeleteTerminal(ctx context.Context, req *terminal.TerminalDeleteRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	_, err := c.terminalClient.DeleteTerminal(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
