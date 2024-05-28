package app

import "github.com/urfave/cli/v2"

var (
	authCommand = &cli.Command{
		Name:  "auth",
		Usage: "manage authentication",
		Subcommands: []*cli.Command{
			{
				Name:  "login",
				Usage: "get credentials fro mmarina gteway",
				Subcommands: []*cli.Command{
					githubLoginCommand,
					passwordLoginCommand,
				},
			},
		},
	}
)
