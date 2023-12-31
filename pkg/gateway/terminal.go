package gateway

import (
	"context"
	"fmt"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
)

var _ terminal.TerminalServiceServer = &Gateway{}

func (g *Gateway) CreateTerminal(ctx context.Context, req *terminal.TerminalCreateRequest) (*empty.Empty, error) {
	terminal := marinav1.Terminal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name.Name,
			Namespace: req.Name.Namespace,
		},
		Spec: marinav1.TerminalSpec{},
	}

	if err := g.kubeClient.Create(ctx, &terminal); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (g *Gateway) DeleteTerminal(ctx context.Context, req *terminal.TerminalDeleteRequest) (*empty.Empty, error) {
	terminal := marinav1.Terminal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name.Name,
			Namespace: req.Name.Namespace,
		},
	}

	if err := g.kubeClient.Delete(ctx, &terminal); err != nil {
		return nil, fmt.Errorf("could not delete terminal '%s': %w", req.Name, err)
	}

	return &empty.Empty{}, nil
}
