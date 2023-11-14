package gateway

import (
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Gateway struct {
	terminal.UnimplementedTerminalServer
	client client.Client
}

func NewGateway(client client.Client) *Gateway {
	return &Gateway{
		client: client,
	}
}

func (g *Gateway) Register(s *grpc.Server) {
	terminal.RegisterTerminalServer(s, g)
}
