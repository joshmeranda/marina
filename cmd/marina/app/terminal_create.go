package app

import (
	"context"
	"fmt"

	"github.com/joshmeranda/marina/gateway/api/core"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	"github.com/urfave/cli/v2"
)

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

	if _, err := client.CreateTerminal(ctx.Context, &createReq); err != nil {
		return fmt.Errorf("could not create terminal: %w", err)
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

var (
	terminalCreateCommand = &cli.Command{
		Name:  "create",
		Usage: "create a terminal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Usage:   "the name of the terminal",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:    "image",
				Usage:   "the image to use for the terminal",
				Aliases: []string{"i"},
			},
		},
		Action: create,
	}
)
