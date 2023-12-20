package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	marina "github.com/joshmeranda/marina/pkg"
	marinagateway "github.com/joshmeranda/marina/pkg/gateway"
	"github.com/joshmeranda/marina/pkg/gateway/drivers/secret"
	"github.com/joshmeranda/marina/pkg/gateway/drivers/storage"
	"github.com/urfave/cli/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/runtime"
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

	schema := runtime.NewScheme()
	if err := marinav1.AddToScheme(schema); err != nil {
		return nil, fmt.Errorf("failed to add marina scheme: %w", err)
	}

	opts := client.Options{
		Scheme: schema,
	}

	client, err := client.New(config, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

func getEtcdClient(ctx *cli.Context) (*clientv3.Client, error) {
	config := clientv3.Config{}
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	secretDriver := secret.NewKubeDriver(client, "marina-system")

	etcdClient, err := getEtcdClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create etcd client: %w", err)
	}
	storageDriver := storage.NewEtcdStore[marina.UserAccessList](etcdClient, json.Marshal, json.Unmarshal)

	gateway, err := marinagateway.NewGateway(
		marinagateway.WithLogger(logger),
		marinagateway.WithKubeClient(client),
		marinagateway.WithSecretDriver(secretDriver),
		marinagateway.WithAccessListStore(storageDriver),
		marinagateway.WithNamespace(ctx.String("namespace")),
	)
	if err != nil {
		return fmt.Errorf("failed to create gateway: %w", err)
	}

	err = gateway.Bootstrap(ctx.Context)
	if err != nil {
		return fmt.Errorf("failed to bootstrap marina gateway: %w", err)
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(marinagateway.LoggingInterceptor(logger), gateway.TokenAuthInterceptor()))

	gateway.Register(server)

	logger.Info("starting server", "addr", addr)

	if err := server.Serve(listener); err != nil {
		return err
	}

	return nil
}

func main() {
	app := cli.App{
		Name:        "marina-server",
		Version:     Version,
		Description: "run the marina gateay server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "namespace",
				Aliases: []string{"n"},
				Usage:   "The name space to use for managin k8s resources",
				Value:   marinagateway.DefaultNamespace,
			},
			&cli.IntFlag{
				Name:    "port",
				Usage:   "the port for the gateway to listen on",
				Aliases: []string{"p"},
				EnvVars: []string{"MARINA_GATEWAY_PORT"},
				Value:   8081, // todo: estalish default port
			},
			&cli.StringFlag{
				Name:    "kubeconfig",
				Usage:   "the path to the kubeconfig file to use for the terminal",
				EnvVars: []string{"KUBECONFIG"},
				Aliases: []string{"f"},
			},
		},
		Action: Start,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
