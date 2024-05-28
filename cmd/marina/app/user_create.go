package app

import (
	"fmt"

	"github.com/joshmeranda/marina/apis/user"
	"github.com/urfave/cli/v2"
)

func userCreate(ctx *cli.Context) error {
	var name string
	var password string

	if narg := ctx.NArg(); narg == 2 {
		name = ctx.Args().Get(0)
		password = ctx.Args().Get(1)
	} else {
		return fmt.Errorf("expected exactly 2 arguments, received %d", narg)
	}

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	req := &user.UserCreateRequest{
		User: &user.User{
			Name:     name,
			Password: []byte(password),
			Roles:    ctx.StringSlice("add-role"),
		},
	}

	_, err = client.CreateUser(ctx.Context, req)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

var (
	userCreateCommand = &cli.Command{
		Name:      "create",
		Usage:     "create a user",
		ArgsUsage: "<name> <password>",
		Action:    userCreate,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "add-role",
				Usage:   "the additional roles to assign to the user",
				Aliases: []string{"r"},
			},
		},
	}
)
