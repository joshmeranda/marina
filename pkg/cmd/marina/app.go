package marina

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/joshmeranda/marina/pkg/apis/core"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"github.com/joshmeranda/marina/pkg/apis/user"
	"github.com/joshmeranda/marina/pkg/client"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var Version string

var cm *configManager

func init() {
	var err error
	cm, err = newDefualtConfigManager()
	if err != nil {
		panic(err)
	}
}

func getClient(ctx *cli.Context) (*client.Client, error) {
	bearerToken, err := cm.GetBearerToken()
	if err != nil {
		return nil, fmt.Errorf("could not get bearer token: %w", err)
	}

	// todo: use real transport credentials
	conn, err := grpc.Dial(ctx.String("address"),
		grpc.WithUnaryInterceptor(client.TokenAuthInterceptor(bearerToken)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	return client.NewClient(conn, logger), nil
}

func create(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	createReq := terminal.TerminalCreateRequest{
		Name: &core.NamespacedName{
			Name:      ctx.String("name"),
			Namespace: "marina-system",
		},
		Spec: &terminal.TerminalSpec{
			Image: ctx.String("image"),
		},
	}

	if err := client.Exec(ctx.Context, createReq.Name); err != nil {
		return fmt.Errorf("could not access terminal: %w", err)
	}

	deleteReq := terminal.TerminalDeleteRequest{
		Name: createReq.Name,
	}

	if _, err := client.DeleteTerminal(context.Background(), &deleteReq); err != nil {
		return fmt.Errorf("could not delete terminal: %w", err)
	}

	return nil
}

func healthCheck(ctx *cli.Context) error {
	var services []string
	if ctx.NArg() == 0 {
		services = []string{
			terminal.TerminalService_ServiceDesc.ServiceName,
			user.UserService_ServiceDesc.ServiceName,
		}
	} else {
		services = ctx.Args().Slice()
	}

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	for _, service := range services {
		req := &healthgrpc.HealthCheckRequest{
			Service: service,
		}

		resp, err := client.Check(context.Background(), req)
		if err != nil {
			return err
		}

		fmt.Printf("%s: %s\n", service, resp.Status.String())
	}

	return nil
}

func login(ctx *cli.Context) error {
	ghAccessToken, err := cm.GetGhAccessToken()
	if err != nil {
		return fmt.Errorf("could not get github access token: %w", err)
	}

	req := &auth.LoginRequest{
		Secret: ghAccessToken,
	}

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Login(context.Background(), req)
	if err != nil {
		return err
	}

	fmt.Printf("token: %s\n", resp.Token)

	return nil
}

func App() cli.App {
	return cli.App{
		Name:        "marina",
		Version:     Version,
		Description: "interact with the marina gateway",
		Commands: []*cli.Command{
			{
				Name:        "create",
				Description: "create a new terminal",
				Action:      create,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "image",
						Usage:   "the name of the image to use for the terminal",
						Aliases: []string{"i"},
					},
					&cli.StringFlag{
						Name:    "name",
						Usage:   "the name of the terminal",
						Aliases: []string{"n"},
					},
				},
			},
			{
				Name:        "check",
				Description: "check the health of the gateway",
				Action:      healthCheck,
			},
			{
				Name: "auth",
				Subcommands: []*cli.Command{
					{
						Name:        "login",
						Description: "get and store gateway authentication credentials",
						Action:      login,
					},
					{
						Name: "store",
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Usage:    "the address of the gateway",
				Required: true,
				Aliases:  []string{"a"},
				EnvVars:  []string{"MARINA_GATEWAY_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "kubeconfig",
				Usage:   "the path to the kubeconfig file to use for the terminal",
				EnvVars: []string{"KUBECONFIG"},
				Aliases: []string{"f"},
			},
		},
	}
}
