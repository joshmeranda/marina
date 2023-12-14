package gateway

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"github.com/joshmeranda/marina/pkg/apis/user"
	"github.com/joshmeranda/marina/pkg/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
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

func TokenAuthInterceptor(l *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if info.FullMethod == "/auth.AuthService/Login" {
			resp, err := handler(ctx, req)
			return resp, err
		}

		md, found := metadata.FromIncomingContext(ctx)
		if !found {
			return nil, fmt.Errorf("could not get tokens from context: missing metadata")
		}

		tokens, ok := md["token"]
		if !ok {
			return nil, fmt.Errorf("could not get tokens from context: missing token")
		}

		token, err := jwt.ParseWithClaims(tokens[0], &customDataClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			return resp, fmt.Errorf("could not parse token: %w", err)
		}

		customClaim, ok := token.Claims.(*customDataClaims)
		if !ok {
			return nil, fmt.Errorf("unsupported token claim type: %t", token.Claims)
		}

		// todo: check that the user exists
		_ = customClaim

		resp, err = handler(ctx, req)

		return resp, err
	}
}

type Gateway struct {
	terminal.UnimplementedTerminalServiceServer
	user.UnimplementedUserServiceServer
	auth.UnimplementedAuthServiceServer

	health healthgrpc.HealthServer

	kubeClient client.Client
	logger     *slog.Logger
	authStore  store.KeyValueStore[string, string]
}

func NewGateway(client client.Client, logger *slog.Logger) *Gateway {
	return &Gateway{
		kubeClient: client,
		health:     health.NewServer(),
		logger:     logger,
		authStore:  store.NewMemoryStore[string, string](),
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
