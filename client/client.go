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

// todo: https://docs.armory.io/continuous-deployment/armory-admin/manual-service-account/
func (c *Client) Exec(ctx context.Context, pod *core.NamespacedName, terminal *core.NamespacedName) error {
	// todo: receive kubeconfig from gateway rather than relying on local kubeconfig
	// todo: access with pod exec rather than sub-process
	//
	// kubeconfigPath := ctx.String("kubeconfig")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to build rest config: %w", err)
	// }
	//
	// kubeClient, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	return err
	// }
	//
	// opts := exec.ExecOptions{
	// 	StreamOptions: exec.StreamOptions{
	// 		Namespace:       terminal.Namespace,
	// 		PodName:         terminal.Name,
	// 		ContainerName:   "",
	// 		Stdin:           true,
	// 		TTY:             true,
	// 		Quiet:           false,
	// 		InterruptParent: &interrupt.Handler{},
	// 		IOStreams:       genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	// 	},
	// 	// FilenameOptions:  resource.FilenameOptions{},
	// 	// ResourceName:     "",
	// 	Command: []string{"sh"},
	// 	// EnforceNamespace: false,
	// 	// Builder: func() *resource.Builder {
	// 	// },
	// 	ExecutablePodFn: polymorphichelpers.AttachablePodForObjectFn,
	// 	// Pod:           &v1.Pod{},
	// 	Executor:      &exec.DefaultRemoteExecutor{},
	// 	PodClient:     kubeClient.CoreV1(),
	// 	GetPodTimeout: time.Second * 5,
	// 	Config:        restConfig,
	// }
	//
	// if err := opts.Validate(); err != nil {
	// 	return fmt.Errorf("failed to validate exec options: %w", err)
	// }
	//
	// if err := opts.Run(); err != nil {
	// 	return fmt.Errorf("failed to run exec: %w", err)
	// }
	//
	// assumes we have a working kubeconfig for the cluster
	// cmd := exec.Command("kubectl", "exec", "--stdin", "--tty", "--namespace", terminal.Namespace, terminal.Name, "--", "sh")
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	//
	// if err := cmd.Run(); err != nil {
	// 	return err
	// }
	//
	// return nil

	var bearerToken string
	var host string

	config := &rest.Config{
		BearerToken: bearerToken,
		Host:        host,
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create pod from bearer token: %w", err)
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
