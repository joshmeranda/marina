package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joshmeranda/marina/gateway/api/auth"
	"github.com/joshmeranda/marina/gateway/api/core"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	"github.com/joshmeranda/marina/gateway/api/user"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

type Client struct {
	terminalClient terminal.TerminalServiceClient
	userClient     user.UserServiceClient
	authClient     auth.AuthServiceClient

	health healthgrpc.HealthClient
	logger *slog.Logger
}

func NewClient(conn grpc.ClientConnInterface, logger *slog.Logger) *Client {
	return &Client{
		terminalClient: terminal.NewTerminalServiceClient(conn),
		userClient:     user.NewUserServiceClient(conn),
		authClient:     auth.NewAuthServiceClient(conn),

		health: healthgrpc.NewHealthClient(conn),
		logger: logger,
	}
}

func (c *Client) Exec(ctx context.Context, config *rest.Config, pod *core.NamespacedName, terminal *core.NamespacedName) error {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	req := client.CoreV1().RESTClient().Post().Resource("pods").Namespace(pod.Namespace).Name(pod.Name).SubResource("exec")

	execOpts := &corev1.PodExecOptions{
		Command: []string{"sh"},
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	req.VersionedParams(execOpts, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return fmt.Errorf("failed to stream: %w", err)
	}

	return nil
}
