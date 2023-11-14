package client

import (
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"google.golang.org/grpc"
)

type Client struct {
	terminalClient terminal.TerminalClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		terminalClient: terminal.NewTerminalClient(conn),
	}
}
