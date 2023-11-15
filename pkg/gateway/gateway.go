package gateway

import (
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Gateway struct {
	terminal.UnimplementedTerminalServer
	health     healthgrpc.HealthServer
	kubeClient client.Client
}

func NewGateway(client client.Client) *Gateway {
	return &Gateway{
		kubeClient: client,
		health:     health.NewServer(),
	}
}

func (g *Gateway) Register(s *grpc.Server) {
	healthUpdater := g.health.(*health.Server)

	terminal.RegisterTerminalServer(s, g)
	healthUpdater.SetServingStatus(terminal.Terminal_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	healthgrpc.RegisterHealthServer(s, g)
	healthUpdater.SetServingStatus(healthgrpc.Health_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)
}
