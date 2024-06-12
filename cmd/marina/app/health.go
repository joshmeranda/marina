package app

import (
	"context"
	"fmt"

	"github.com/joshmeranda/marina/gateway/api/auth"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	"github.com/joshmeranda/marina/gateway/api/user"
	"github.com/urfave/cli/v2"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func health(ctx *cli.Context) error {
	var services []string
	if ctx.NArg() == 0 {
		services = []string{
			terminal.TerminalService_ServiceDesc.ServiceName,
			user.UserService_ServiceDesc.ServiceName,
			auth.AuthService_ServiceDesc.ServiceName,
		}
	} else {
		services = ctx.Args().Slice()
	}

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	for _, service := range services {
		req := &healthgrpc.HealthCheckRequest{
			Service: service,
		}

		resp, err := client.Check(context.Background(), req)
		if err != nil {
			return err
		}

		fmt.Printf("%s: %s\n", service, resp.Status.String())
	}

	return nil
}

var (
	healthCheckCommand = &cli.Command{
		Name:   "health",
		Usage:  "check the health of the gateway",
		Action: health,
	}
)
