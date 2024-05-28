package app

import (
	"fmt"

	"github.com/joshmeranda/marina/apis/user"
	"github.com/urfave/cli/v2"
)

func userDelete(ctx *cli.Context) error {
	var name string

	if narg := ctx.NArg(); narg == 1 {
		name = ctx.Args().Get(0)
	} else {
		return fmt.Errorf("expected exactly 1 arguments, received %d", narg)
	}

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	req := &user.UserDeleteRequest{
		Name: name,
	}

	_, err = client.DeleteUser(ctx.Context, req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

var (
	userDeleteCommand = &cli.Command{
		Name:      "delete",
		Usage:     "dewlete a user",
		ArgsUsage: "<name>",
		Action:    userDelete,
	}
)
