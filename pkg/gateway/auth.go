package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v57/github"
	"github.com/joshmeranda/marina/pkg/apis/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ auth.AuthServiceServer = &Gateway{}

const (
	TokenSigningSecretName  = "jet-signing-key"
	TokenSigningSecretField = "value"
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

		// todo: check that the user exists
		_ = customClaim

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func (g *Gateway) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// todo: check against white / black listed users
	client := github.NewClient(nil).WithAuthToken(req.Token)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	g.logger.Info("generating token for user", "user", user.GetLogin())

	claims := customDataClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "marina",
			ExpiresAt: &jwt.NumericDate{
				// token is valid for roughly 1 week
				Time: time.Now().Add(24 * time.Hour * 7),
			},
		},
		User: user.GetLogin(),
	}

	signingKey, err := g.secretDriver.Get(ctx, TokenSigningSecretName, TokenSigningSecretField)
	if err != nil {
		return nil, fmt.Errorf("could not get data from secret '%s' at field '%s'", TokenSigningSecretName, TokenSigningSecretField)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bearerToken, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: bearerToken,
	}, nil
}
