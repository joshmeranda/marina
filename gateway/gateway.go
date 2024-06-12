package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/gateway/api/auth"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	"github.com/joshmeranda/marina/gateway/api/user"
	authdriver "github.com/joshmeranda/marina/gateway/drivers/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DefaultNamespace string = "marina-system"
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

	kubeClient client.Client
	kubeConfig *rest.Config

	logger     *slog.Logger
	namespace  string
	authDriver authdriver.Driver
}

func NewGateway(opts ...Option) (*Gateway, error) {
	gateway := &Gateway{
		health: health.NewServer(),
	}

	for _, opt := range opts {
		opt(gateway)
	}

	var err error

	if gateway.kubeConfig == nil {
		gateway.kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
		}
	}

	if gateway.kubeClient == nil {
		schema := scheme.Scheme
		if err := marinav1.AddToScheme(schema); err != nil {
			return nil, fmt.Errorf("failed to add marina scheme: %w", err)
		}

		opts := client.Options{
			Scheme: schema,
		}

		gateway.kubeClient, err = client.New(gateway.kubeConfig, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}
	}

	if gateway.logger == nil {
		gateway.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	if gateway.namespace == "" {
		gateway.namespace = DefaultNamespace
	}

	if gateway.authDriver == nil {
		gateway.authDriver = authdriver.NewMemory()
	}

	return gateway, nil
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
