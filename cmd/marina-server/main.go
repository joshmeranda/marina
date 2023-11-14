package main

import (
	"fmt"
	"net"

	"github.com/joshmeranda/marina/pkg/gateway"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var Version string

func getClusterClient(ctx *cli.Context) (client.Client, error) {
	var config *rest.Config
	var err error

	if kubeconfig := ctx.String("kubeconfig"); kubeconfig != "" {
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return nil, fmt.Errorf("failed to get config from kubeconfig: %w", err)
		}
	} else if config, err = rest.InClusterConfig(); err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	client, err := client.New(config, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

func Start(ctx *cli.Context) error {
	port := ctx.Int("port")
	addr := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	client, err := getClusterClient(ctx)
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	gateway := gateway.NewGateway(client)
	gateway.Register(server)

	if err := server.Serve(listener); err != nil {
		return err
	}

	return nil
}

func main() {

}
