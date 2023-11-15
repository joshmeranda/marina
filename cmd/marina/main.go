package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"github.com/joshmeranda/marina/pkg/client"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var Version string

func getClient(ctx *cli.Context) (*client.Client, error) {
	conn, err := grpc.Dial(ctx.String("address"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return client.NewClient(conn), nil
}

func HealthCheck(ctx *cli.Context) error {
	var services []string
	if ctx.NArg() == 0 {
		services = []string{
			terminal.Terminal_ServiceDesc.ServiceName,
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

func Create(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	req := terminal.TerminalCreateRequest{
		Image: ctx.String("image"),
		Shell: ctx.String("shell"),
	}

	if _, err := client.CreateTerminal(context.Background(), &req); err != nil {
		return err
	}

	return nil
}

func main() {
	app := cli.App{
		Name:        "marina",
		Version:     Version,
		Description: "interact with the marina gateway",
		Commands: []*cli.Command{
			{
				Name:        "create",
				Description: "create a new terminal",
				Action:      Create,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "image",
						Usage:   "the name of the image to use for the terminal",
						Aliases: []string{"i"},
					},
					&cli.StringFlag{
						Name:    "shell",
						Usage:   "the shell to use for the terminal",
						Aliases: []string{"s"},
					},
				},
			},
			{
				Name:        "check",
				Description: "check the health of the gateway",
				Action:      HealthCheck,
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
