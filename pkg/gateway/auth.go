package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v57/github"
	"github.com/joshmeranda/marina/pkg/apis/auth"
)

var _ auth.AuthServiceServer = &Gateway{}

type customDataClaims struct {
	jwt.RegisteredClaims

	User string `json:"user,omitempty"`
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

	// signingKey := make([]byte, 10)
	// if _, err := rand.Read(signingKey); err != nil {
	// 	return nil, fmt.Errorf("error generating signing key: %w", err)
	// }

	// todo: secret provider
	signingKey := []byte("secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bearerToken, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	if err := g.authStore.Set(user.GetLogin(), bearerToken); err != nil {
		return nil, fmt.Errorf("could not store user credentials: %w", err)
	}

	return &auth.LoginResponse{
		Token: bearerToken,
	}, nil
}
