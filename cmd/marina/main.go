package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joshmeranda/marina/pkg/apis/core"
	"github.com/joshmeranda/marina/pkg/apis/terminal"
	"github.com/joshmeranda/marina/pkg/client"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var Version string

func getClient(ctx *cli.Context) (*client.Client, error) {
	conn, err := grpc.Dial(ctx.String("address"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return client.NewClient(conn), nil
}

// todo: replace with kubeconfig content received from the gateway.
func getKubeClient(ctx *cli.Context) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", ctx.String("KUBECONFIG"))
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func HealthCheck(ctx *cli.Context) error {
	var services []string
	if ctx.NArg() == 0 {
		services = []string{
			terminal.Terminal_ServiceDesc.ServiceName,
		}
	} else {
		services = ctx.Args().Slice()
	}

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	for _, service := range services {
		req := &healthgrpc.HealthCheckRequest{
			Service: service,
		}

		resp, err := client.Check(context.Background(), req)
		if err != nil {
			return err
		}

		fmt.Printf("%s: %s\n", service, resp.Status.String())
	}

	return nil
}

func Create(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	createReq := terminal.TerminalCreateRequest{
		Name: &core.NamespacedName{
			Name:      ctx.String("name"),
			Namespace: "marina-system",
		},
		Spec: &terminal.TerminalSpec{
			Image: ctx.String("image"),
		},
	}

	if _, err := client.CreateTerminal(context.Background(), &createReq); err != nil {
		return err
	}

	kubeClient, err := getKubeClient(ctx)
	if err != nil {
		return err
	}

	opts := &corev1.PodExecOptions{
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Container: "shell",
		Command:   []string{"sh"},
	}

	execReq := kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(ctx.String("name")).
		Namespace("marina-system").
		SubResource("exec").
		VersionedParams(opts, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(&rest.Config{}, "POST", execReq.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(ctx.Context, remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return err
	}

	deleteReq := terminal.TerminalDeleteRequest{
		Name: createReq.Name,
	}
	if _, err := client.DeleteTerminal(context.Background(), &deleteReq); err != nil {
		return err
	}

	return nil
}

func main() {
	app := cli.App{
		Name:        "marina",
		Version:     Version,
		Description: "interact with the marina gateway",
		Commands: []*cli.Command{
			{
				Name:        "create",
				Description: "create a new terminal",
				Action:      Create,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "image",
						Usage:   "the name of the image to use for the terminal",
						Aliases: []string{"i"},
					},
					&cli.StringFlag{
						Name:    "name",
						Usage:   "the name of the terminal",
						Aliases: []string{"n"},
					},
				},
			},
			{
				Name:        "check",
				Description: "check the health of the gateway",
				Action:      HealthCheck,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Usage:    "the address of the gateway",
				Required: true,
				Aliases:  []string{"a"},
				EnvVars:  []string{"MARINA_GATEWAY_ADDRESS"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
