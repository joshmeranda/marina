package client

import (
	"context"

	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var _ healthgrpc.HealthClient = &Client{}

func (c *Client) Check(ctx context.Context, in *healthgrpc.HealthCheckRequest, opts ...grpc.CallOption) (*healthgrpc.HealthCheckResponse, error) {
	resp, err := c.health.Check(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Watch implements grpc_health_v1.HealthClient.
func (c *Client) Watch(ctx context.Context, in *healthgrpc.HealthCheckRequest, opts ...grpc.CallOption) (healthgrpc.Health_WatchClient, error) {
	wc, err := c.health.Watch(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return wc, nil
}
