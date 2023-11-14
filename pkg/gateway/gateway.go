package gateway

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/pkg/apis"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ apis.MarinaServer = &Gateway{}

type Gateway struct {
	apis.UnimplementedMarinaServer
	client client.Client
}

func NewGateway(client client.Client) *Gateway {
	return &Gateway{
		client: client,
	}
}

func (g *Gateway) CreateTerminal(ctx context.Context, req *apis.TerminalCreateRequest) (*empty.Empty, error) {
	terminal := marinav1.Terminal{}

	if err := g.client.Create(ctx, &terminal); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
