package marina

import "github.com/urfave/cli/v2"

var (
	userCommand = &cli.Command{
		Name:  "user",
		Usage: "interact with the user api",
		Subcommands: []*cli.Command{
			userListCommand,
			userCreateCommand,
			userDeleteCommand,
			userUpdateCommand,
		},
	}
)
