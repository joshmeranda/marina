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

func (g *Gateway) clientFromContext(ctx context.Context, opts client.Options) (client.Client, error) {
	foundUsers := metadata.ValueFromIncomingContext(ctx, UserMetadataFieldName)
	config := rest.CopyConfig(g.kubeConfig)

	if len(foundUsers) > 1 {
		return nil, fmt.Errorf("bug: could not get user from context: multiple authenticated users")
	}

	if len(foundUsers) == 0 {
		return g.kubeClient, nil
	}

	config.Impersonate.UserName = fmt.Sprintf("system:serviceaccount:%s:%s", g.namespace, foundUsers[0])

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
