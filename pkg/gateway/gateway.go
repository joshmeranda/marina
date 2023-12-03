package gateway

import (
	"context"
	"log/slog"

	"github.com/joshmeranda/marina/pkg/apis/terminal"
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
	terminal.UnimplementedTerminalServer

	health     healthgrpc.HealthServer
	kubeClient client.Client

	logger *slog.Logger
}

func NewGateway(client client.Client, logger *slog.Logger) *Gateway {
	return &Gateway{
		kubeClient: client,
		health:     health.NewServer(),
		logger:     logger,
	}
}

func (g *Gateway) Register(s *grpc.Server) {
	healthUpdater := g.health.(*health.Server)

	terminal.RegisterTerminalServer(s, g)
	healthUpdater.SetServingStatus(terminal.Terminal_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	healthgrpc.RegisterHealthServer(s, g)
	healthUpdater.SetServingStatus(healthgrpc.Health_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)
}
