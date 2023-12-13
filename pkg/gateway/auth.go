package gateway

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v57/github"
	"github.com/joshmeranda/marina/pkg/apis/auth"
)

var _ auth.AuthServiceServer = &Gateway{}

type customDataClaims struct {
	jwt.RegisteredClaims

	Data map[string]any `json:"data,omitempty"`
}

func (g *Gateway) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	client := github.NewClient(nil).WithAuthToken(req.Token)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	claims := customDataClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "marina",
			ExpiresAt: &jwt.NumericDate{
				// token is valid for roughly 1 week
				Time: time.Now().Add(24 * time.Hour * 7),
			},
		},
		Data: map[string]any{
			"username": user.GetLogin(),
		},
	}

	mySigningKey := []byte("AllYourBase")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bearerToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: bearerToken,
	}, nil
}
