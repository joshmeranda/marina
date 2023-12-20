package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v57/github"
	marina "github.com/joshmeranda/marina/pkg"
	"github.com/joshmeranda/marina/pkg/apis/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// todo: device flow with refres
// todo: admin user login crednetials

var _ auth.AuthServiceServer = &Gateway{}

const (
	TokenSigningSecretName  = "jet-signing-key"
	TokenSigningSecretField = "value"

	UserAccessFieldKeyName = "user-access-list"
)

type customDataClaims struct {
	jwt.RegisteredClaims

	User string `json:"user,omitempty"`
}

func (g *Gateway) TokenAuthInterceptor() grpc.UnaryServerInterceptor {
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

		client := github.NewClient(nil)

		if isUserAllowed, err := g.isUserAllowed(ctx, client, customClaim.User); err != nil {
			return nil, fmt.Errorf("error checking for user access: %w", err)
		} else if !isUserAllowed {
			return nil, fmt.Errorf("user '%s' is not allowed", customClaim.User)
		}

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func (g *Gateway) isUserAllowed(ctx context.Context, ghClient *github.Client, username string) (bool, error) {
	orgs, _, err := ghClient.Organizations.List(context.Background(), "", nil)
	if err != nil {
		return false, fmt.Errorf("could not retrieve user organizations: %w", err)
	}

	orgNames := make([]string, len(orgs))
	for i, org := range orgs {
		orgNames[i] = org.GetLogin()
	}

	list, err := g.accessListStore.Get(ctx, UserAccessFieldKeyName)
	if err != nil {
		return false, fmt.Errorf("could not retrieve user access list: %w", err)
	}

	switch accessType := list.GetAccessFor(username, orgNames); accessType {
	case marina.AccessTypeAllow:
		return true, nil
	case marina.AccessTypeDeny | marina.AccessTypeUnknown:
		return false, nil
	default:
		panic(fmt.Sprintf("bug: encountered unsupported accesss type %d", accessType))
	}
}

func (g *Gateway) generateTokenForUser(ctx context.Context, user string) (string, error) {
	g.logger.Info("generating token for user", "user", user)

	claims := customDataClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "marina",
			ExpiresAt: &jwt.NumericDate{
				// token is valid for roughly 1 week
				Time: time.Now().Add(24 * time.Hour * 7),
			},
		},
		User: user,
	}

	signingKey, err := g.secretDriver.Get(ctx, TokenSigningSecretName, TokenSigningSecretField)
	if err != nil {
		return "", fmt.Errorf("could not get data from secret '%s' at field '%s'", TokenSigningSecretName, TokenSigningSecretField)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bearerToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return bearerToken, nil
}

func (g *Gateway) githubLogin(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := g.authDriver.Authenticate(ctx, req); err != nil {
		return nil, fmt.Errorf("could not authenticate user '%s': %w", req.User, err)
	}

	// isUserAllowed, err := g.isUserAllowed(ctx, client, user.GetLogin())
	// if err != nil {
	// 	return nil, fmt.Errorf("error checking for user access: %w", err)
	// }

	// if !isUserAllowed {
	// 	return nil, fmt.Errorf("user '%s' is not allowed", user.GetLogin())
	// }

	bearerToken, err := g.generateTokenForUser(ctx, req.User)
	if err != nil {
		return nil, fmt.Errorf("could not generate token for user '%s': %w", req.User, err)
	}

	return &auth.LoginResponse{
		Token: bearerToken,
	}, nil
}

func (g *Gateway) marinaLogin(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := g.authDriver.Authenticate(ctx, req); err != nil {
		return nil, fmt.Errorf("could not authenticate user '%s': %w", req.User, err)
	}

	// todo: ideally we'd be able to prevent users from abusing the costliness of this operation (cache?)
	bearerToken, err := g.generateTokenForUser(ctx, req.User)
	if err != nil {
		return nil, fmt.Errorf("could not generate token for user '%s': %w", req.User, err)
	}

	return &auth.LoginResponse{
		Token: bearerToken,
	}, nil
}

func (g *Gateway) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	switch req.SecretType {
	case auth.SecretType_Github:
		return g.githubLogin(ctx, req)
	case auth.SecretType_Password:
		return g.marinaLogin(ctx, req)
	default:
		return nil, fmt.Errorf("recevied unknown token kind: %s", req.SecretType)
	}
}
