package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joshmeranda/marina/pkg/apis"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Version string

func getClient(ctx *cli.Context) (apis.MarinaClient, error) {
	conn, err := grpc.Dial(ctx.String("address"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := apis.NewMarinaClient(conn)

	_, err = client.CreateTerminal(context.Background(), &apis.TerminalCreateRequest{})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Create(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	req := apis.TerminalCreateRequest{
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
		Name:           "marina",
		Version:        Version,
		Description:    "interact with the marina gateway",
		DefaultCommand: "",
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
