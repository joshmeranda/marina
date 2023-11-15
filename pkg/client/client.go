package client

import (
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type Client struct {
	terminalClient terminal.TerminalClient
	health         healthgrpc.HealthClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		terminalClient: terminal.NewTerminalClient(conn),
		health:         healthgrpc.NewHealthClient(conn),
	}
}
