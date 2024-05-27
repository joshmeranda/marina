package gateway

import (
	"context"

	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var _ healthgrpc.HealthServer = &Gateway{}

// Check implements grpc_health_v1.HealthServer.
func (g *Gateway) Check(ctx context.Context, req *healthgrpc.HealthCheckRequest) (*healthgrpc.HealthCheckResponse, error) {
	return g.health.Check(ctx, req)
}

// Watch implements grpc_health_v1.HealthServer.
func (g *Gateway) Watch(req *healthgrpc.HealthCheckRequest, server healthgrpc.Health_WatchServer) error {
	return g.health.Watch(req, server)
}
