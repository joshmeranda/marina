package gateway

import (
	"context"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
)

var _ terminal.TerminalServer = &Gateway{}

func (g *Gateway) CreateTerminal(ctx context.Context, req *terminal.TerminalCreateRequest) (*empty.Empty, error) {
	terminal := marinav1.Terminal{}

	if err := g.kubeClient.Create(ctx, &terminal); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
