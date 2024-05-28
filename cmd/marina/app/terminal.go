package app

import "github.com/urfave/cli/v2"

var (
	terminalCommand = &cli.Command{
		Name:  "terminal",
		Usage: "interact with the marina terminal api",
		Subcommands: []*cli.Command{
			terminalCreateCommand,
		},
	}
)
