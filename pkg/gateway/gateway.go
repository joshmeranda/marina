package gateway

import (
	"context"
	"log/slog"

	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"github.com/joshmeranda/marina/pkg/apis/user"
	"github.com/joshmeranda/marina/pkg/drivers/secret"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func LoggingInterceptor(l *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Info("gRPC method call", "method", info.FullMethod)
		resp, err = handler(ctx, req)
		if err != nil {
			l.Error("gRPC method call error",
				"method", info.FullMethod,
				"error", err)
		}
		return resp, err
	}
}

type Gateway struct {
	terminal.UnimplementedTerminalServiceServer
	user.UnimplementedUserServiceServer
	auth.UnimplementedAuthServiceServer

	health healthgrpc.HealthServer

	kubeClient   client.Client
	logger       *slog.Logger
	secretDriver secret.Driver
}

// todo: convert to use options...
func NewGateway(client client.Client, logger *slog.Logger, secretDriver secret.Driver) *Gateway {
	return &Gateway{
		kubeClient:   client,
		health:       health.NewServer(),
		logger:       logger,
		secretDriver: secretDriver,
	}
}

func (g *Gateway) Register(s *grpc.Server) {
	healthUpdater := g.health.(*health.Server)

	terminal.RegisterTerminalServiceServer(s, g)
	healthUpdater.SetServingStatus(terminal.TerminalService_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	user.RegisterUserServiceServer(s, g)
	healthUpdater.SetServingStatus(user.UserService_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	auth.RegisterAuthServiceServer(s, g)
	healthUpdater.SetServingStatus(auth.AuthService_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	healthgrpc.RegisterHealthServer(s, g)
	healthUpdater.SetServingStatus(healthgrpc.Health_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)
}
