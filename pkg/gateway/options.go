package gateway

import (
	"log/slog"

	marina "github.com/joshmeranda/marina/pkg"
	"github.com/joshmeranda/marina/pkg/gateway/drivers/auth"
	"github.com/joshmeranda/marina/pkg/gateway/drivers/secret"
	"github.com/joshmeranda/marina/pkg/gateway/drivers/storage"
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

func WithSecretDriver(driver secret.Driver) Option {
	return func(g *Gateway) {
		g.secretDriver = driver
	}
}

func WithAccessListStore(store storage.KeyValueStore[string, marina.UserAccessList]) Option {
	return func(g *Gateway) {
		g.accessListStore = store
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
