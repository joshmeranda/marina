package gateway

import (
	"log/slog"

	"github.com/joshmeranda/marina/gateway/drivers/auth"
	"github.com/joshmeranda/marina/gateway/images"
	"k8s.io/client-go/rest"
)

type Option func(*Gateway)

func WithKubeConfig(config *rest.Config) Option {
	return func(g *Gateway) {
		g.kubeConfig = config
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

func WithAccessList(imagesAccessList images.ImagesAccessList) Option {
	return func(g *Gateway) {
		g.imagesAccessList = imagesAccessList
	}
}
