package app

import (
	"fmt"
	"slices"

	"github.com/joshmeranda/marina/apis/user"
	"github.com/urfave/cli/v2"
)

func userUpdate(ctx *cli.Context) error {
	if narg := ctx.NArg(); narg != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %d", narg)
	}

	username := ctx.Args().First()

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	u, err := client.GetUser(ctx.Context, &user.UserGetRequest{
		Name: username,
	})
	if err != nil {
		return fmt.Errorf("failed to get exisitng user: %w", err)
	}

	addRoles := ctx.StringSlice("add-role")
	removeRoles := ctx.StringSlice("remove-role")
	password := ctx.String("password")

	u.Roles = slices.DeleteFunc(u.Roles, func(role string) bool {
		return slices.Contains(removeRoles, role)
	})

	u.Roles = append(u.Roles, addRoles...)

	if password != "" {
		u.Password = []byte(password)
	}

	_, err = client.UpdateUser(ctx.Context, &user.UserUpdateRequest{
		User: u,
	})

	return err
}

var (
	userUpdateCommand = &cli.Command{
		Name:      "update",
		Usage:     "update a user's fields",
		ArgsUsage: "<name>",
		Action:    userUpdate,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "add-role",
				Usage: "add role to user",
			},
			&cli.StringSliceFlag{
				Name:  "remove-role",
				Usage: "remove role from user",
			},

			&cli.StringFlag{
				Name:  "password",
				Usage: "change user's password",
			},
		},
	}
)
