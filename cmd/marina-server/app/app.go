package app

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	authapis "github.com/joshmeranda/marina/apis/auth"
	marinagateway "github.com/joshmeranda/marina/gateway"
	"github.com/joshmeranda/marina/gateway/drivers/auth"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes/scheme"
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

	schema := scheme.Scheme
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

func getLogger(ctx *cli.Context) *slog.Logger {
	var out io.Writer = os.Stdout
	opts := &slog.HandlerOptions{}

	if ctx.Bool("quiet") {
		opts.Level = slog.LevelWarn
	}

	if ctx.Bool("silent") {
		out = io.Discard
	}

	if ctx.Bool("verbose") {
		opts.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(out, opts))

	return logger
}

func start(ctx *cli.Context) error {
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

	namespace := ctx.String("namespace")

	logger := getLogger(ctx)

	authDriver := auth.MultiAuth{
		Drivers: map[authapis.SecretType]auth.Driver{
			authapis.SecretType_Github:   auth.NewGithub(),
			authapis.SecretType_Password: auth.NewLocal(client, namespace),
		},
	}

	gateway, err := marinagateway.NewGateway(
		marinagateway.WithLogger(logger),
		marinagateway.WithKubeClient(client),
		marinagateway.WithNamespace(namespace),
		marinagateway.WithAuthDriver(&authDriver),
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

func App() cli.App {
	return cli.App{
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
			&cli.StringSliceFlag{
				Name:    "etcd",
				Usage:   "the endpoints for the etcd cluster to use for storing access lists",
				Aliases: []string{"e"},
			},

			&cli.BoolFlag{
				Name:    "quiet",
				Usage:   "suppress all output except for warnings and errors",
				Aliases: []string{"q"},
			},
			&cli.BoolFlag{
				Name:    "silent",
				Usage:   "suppress all output",
				Aliases: []string{"s"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "run verbosely",
				Aliases: []string{"v"},
			},
		},
		Action: start,
	}
}
