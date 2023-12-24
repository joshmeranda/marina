package marina

import (
	"context"
	"fmt"

	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"github.com/joshmeranda/marina/pkg/apis/user"
	"github.com/urfave/cli/v2"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func healthCheck(ctx *cli.Context) error {
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
		Name:   "check",
		Usage:  "check the health of the gateway",
		Action: healthCheck,
	}
)
