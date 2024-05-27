package marina

import (
	"github.com/urfave/cli/v2"
)

var Version string
var cm *configManager

func setup(ctx *cli.Context) error {
	var err error
	cm, err = newConfigManager(ctx.String("config-dir"))
	if err != nil {
		return err
	}

	return nil
}

func teardown(ctx *cli.Context) error {
	if cm != nil {
		if err := cm.Close(); err != nil {
			return err
		}
	}

	return nil
}

func App() cli.App {
	return cli.App{
		Name:        "marina",
		Version:     Version,
		Description: "interact with the marina gateway",
		Commands: []*cli.Command{
			terminalCommand,
			healthCheckCommand,
			authCommand,
			userCommand,
		},
		Before: setup,
		After:  teardown,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-dir",
				Usage:   "the directory to store configuration in",
				Aliases: []string{"c"},
				EnvVars: []string{"MARINA_CONFIG_DIR"},
			},
			&cli.StringFlag{
				Name:     "address",
				Usage:    "the address of the gateway",
				Required: true,
				Aliases:  []string{"a"},
				EnvVars:  []string{"MARINA_GATEWAY_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "kubeconfig",
				Usage:   "the path to the kubeconfig file to use for the terminal",
				EnvVars: []string{"KUBECONFIG"},
				Aliases: []string{"f"},
			},
		},
	}
}
