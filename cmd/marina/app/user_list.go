package app

import (
	"fmt"

	"github.com/joshmeranda/marina/gateway/api/user"
	"github.com/urfave/cli/v2"
)

func userList(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	req := &user.UserListRequest{
		Query: &user.UserQuery{},
	}

	resp, err := client.ListUser(ctx.Context, req)
	if err != nil {
		return fmt.Errorf("could not list users: %w", err)
	}

	for _, user := range resp.Users {
		fmt.Printf("%-10s %v\n", user.Name, user.Roles)
	}

	return nil
}

var (
	userListCommand = &cli.Command{
		Name:   "list",
		Usage:  "filter and list users in the marina cluster",
		Action: userList,
	}
)
