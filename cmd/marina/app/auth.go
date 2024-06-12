package app

import "github.com/urfave/cli/v2"

// todo: add whoami command
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
