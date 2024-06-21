package gateway

import (
	"context"
	"fmt"

	marinav1 "github.com/joshmeranda/marina/api/v1"
	"google.golang.org/grpc/metadata"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func userFromContext(ctx context.Context) (string, error) {
	foundUsers := metadata.ValueFromIncomingContext(ctx, UserMetadataFieldName)

	if len(foundUsers) > 1 {
		return "", fmt.Errorf("bug: could not get user from context: multiple authenticated users")
	}

	if len(foundUsers) == 0 {
		return "admin", nil
	}

	return foundUsers[0], nil
}

func (g *Gateway) clientFromContext(ctx context.Context, opts client.Options) (client.Client, error) {
	user, err := userFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from context: %w", err)
	}

	if user == "admin" {
		return g.kubeClient, nil
	}

	config := rest.CopyConfig(g.kubeConfig)
	config.Impersonate.UserName = fmt.Sprintf("system:serviceaccount:%s:%s", g.namespace, user)

	if opts.Scheme == nil {
		schema := scheme.Scheme
		if err := marinav1.AddToScheme(schema); err != nil {
			return nil, fmt.Errorf("failed to add marina scheme: %w", err)
		}

		opts.Scheme = schema
	}

	client, err := client.New(config, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
