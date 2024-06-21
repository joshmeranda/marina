package app

import (
	"context"
	"fmt"

	"github.com/joshmeranda/marina/gateway/api/core"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	"github.com/joshmeranda/marina/kubeconfig"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func getExecClient(token string) (*rest.Config, error) {
	kubeString, err := kubeconfig.ForTokenBased("marina-exec", "", "https://rancher.local.com:6443", token)
	if err != nil {
		return nil, fmt.Errorf("could not create client from kubeconfig: %w", err)
	}

	getter := func() (*api.Config, error) {
		return clientcmd.Load([]byte(kubeString))
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", getter)
	if err != nil {
		return nil, fmt.Errorf("could not create client from kubeconfig: %w", err)
	}

	return config, nil
}

func create(ctx *cli.Context) error {
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

	createResp, err := client.CreateTerminal(ctx.Context, &createReq)
	if err != nil {
		return fmt.Errorf("could not create terminal: %w", err)
	}

	config, err := getExecClient(createResp.Token)
	if err != nil {
		return fmt.Errorf("could not create kubeconfig: %w", err)
	}

	if err := client.Exec(ctx.Context, config, createResp.Pod, createReq.Name); err != nil {
		return fmt.Errorf("could not access terminal: %w", err)
	}

	deleteReq := terminal.TerminalDeleteRequest{
		Name: createReq.Name,
	}

	if _, err := client.DeleteTerminal(context.Background(), &deleteReq); err != nil {
		return fmt.Errorf("could not delete terminal: %w", err)
	}

	return nil
}

var (
	terminalCreateCommand = &cli.Command{
		Name:  "create",
		Usage: "create a terminal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Usage:   "the name of the terminal",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:    "image",
				Usage:   "the image to use for the terminal",
				Aliases: []string{"i"},
			},
		},
		Action: create,
	}
)
