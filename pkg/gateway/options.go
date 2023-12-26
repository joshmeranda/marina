package gateway

import (
	"log/slog"

	"github.com/joshmeranda/marina/pkg/gateway/drivers/auth"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Option func(*Gateway)

func WithKubeClient(client client.Client) Option {
	return func(g *Gateway) {
		g.kubeClient = client
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(g *Gateway) {
		g.logger = logger
	}
}

func WithNamespace(namespace string) Option {
	return func(g *Gateway) {
		g.namespace = namespace
	}
}

func WithAuthDriver(driver auth.Driver) Option {
	return func(g *Gateway) {
		g.authDriver = driver
	}
}
