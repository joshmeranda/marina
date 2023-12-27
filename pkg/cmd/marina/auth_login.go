package marina

import (
	"context"
	"fmt"
	"os"

	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/urfave/cli/v2"
)

const (
	loginSecretEnvName = "MARINA_LOGIN_SECRET"
)

func fillLoginRequest(ctx *cli.Context, req *auth.LoginRequest) error {
	args := ctx.Args()
	narg := args.Len()

	var user string
	secret := []byte(os.Getenv(loginSecretEnvName))

	switch narg {
	case 0:
		return fmt.Errorf("expected at least 1 arg, fot %d", narg)
	case 1:
		user = args.First()
	case 2:
		user = args.First()
		secret = []byte(args.Get(1))
	default:
		return fmt.Errorf("expected at most 2 args, got %d", narg)
	}

	if len(secret) == 0 {
		return fmt.Errorf("secret is required but not provided")
	}

	req.User = user
	req.Secret = secret

	return nil
}

func login(ctx *cli.Context, req *auth.LoginRequest) error {
	if err := fillLoginRequest(ctx, req); err != nil {
		return err
	}

	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Login(context.Background(), req)
	if err != nil {
		return err
	}

	cm.Config.BearerToken = resp.Token

	return nil
}

func passwordLogin(ctx *cli.Context) error {
	err := login(ctx, &auth.LoginRequest{
		SecretType: auth.SecretType_Password,
	})
	if err != nil {
		return err
	}

	return nil
}

func githubLogin(ctx *cli.Context) error {
	err := login(ctx, &auth.LoginRequest{
		Secret:     []byte(cm.Config.BearerToken),
		SecretType: auth.SecretType_Github,
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	passwordLoginCommand = &cli.Command{
		Name:      "password",
		Usage:     "authenticate with a password",
		ArgsUsage: "<user> [secret]",
		Action:    passwordLogin,
	}

	githubLoginCommand = &cli.Command{
		Name:      "github",
		Usage:     "authenticate with github",
		ArgsUsage: "<user> [secret]",
		Action:    githubLogin,
	}
)
