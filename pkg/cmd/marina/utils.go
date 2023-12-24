package marina

import (
	"log/slog"
	"os"

	"github.com/joshmeranda/marina/pkg/client"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getClient(ctx *cli.Context) (*client.Client, error) {
	// todo: use real transport credentials
	conn, err := grpc.Dial(ctx.String("address"),
		grpc.WithUnaryInterceptor(client.TokenAuthInterceptor(cm.Config.BearerToken)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	return client.NewClient(conn, logger), nil
}

func ToPtr[T any](t T) *T {
	return &t
}
