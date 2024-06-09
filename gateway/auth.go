package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joshmeranda/marina/apis/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// todo: device flow with refresh tokens

var _ auth.AuthServiceServer = &Gateway{}

const (
	TokenSigningSecretName  = "jwt-signing-key"
	TokenSigningSecretField = "value"

	UserMetadataFieldName = "username"
)

// todo: add serviceaccount token?
type customDataClaims struct {
	jwt.RegisteredClaims

	User string `json:"user,omitempty"`
}

func (g *Gateway) TokenAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		switch info.FullMethod {
		case "/auth.AuthService/Login", "/grpc.health.v1.Health/Check":
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
			secret := corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      TokenSigningSecretName,
					Namespace: g.namespace,
				},
			}
			if err := g.kubeClient.Get(ctx, client.ObjectKeyFromObject(&secret), &secret); err != nil {
				return "", fmt.Errorf("could not get signing secret: %w", err)
			}

			signingKey := secret.Data[TokenSigningSecretField]

			return signingKey, nil
		})
		if err != nil {
			return resp, fmt.Errorf("could not parse token: %w", err)
		}

		customClaim, ok := token.Claims.(*customDataClaims)
		if !ok {
			return nil, fmt.Errorf("unsupported token claim type: %t", token.Claims)
		}

		// todo: we aren't actually checking for authentication here
		// todo: if the claim is expired, we should return an error

		md.Set(UserMetadataFieldName, customClaim.User)
		ctx = metadata.NewIncomingContext(ctx, md)

		resp, err = handler(ctx, req)

		return resp, err
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

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenSigningSecretName,
			Namespace: g.namespace,
		},
	}
	if err := g.kubeClient.Get(ctx, client.ObjectKeyFromObject(&secret), &secret); err != nil {
		return "", fmt.Errorf("could not get signing secret: %w", err)
	}

	signingKey := secret.Data[TokenSigningSecretField]

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bearerToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return bearerToken, nil
}

func (g *Gateway) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := g.authDriver.Authenticate(ctx, req); err != nil {
		return nil, fmt.Errorf("could not authenticate user '%s': %w", req.User, err)
	}

	// todo: deal with access list stuff (might not really need it)

	bearerToken, err := g.generateTokenForUser(ctx, req.User)
	if err != nil {
		return nil, fmt.Errorf("could not generate token for user '%s': %w", req.User, err)
	}

	return &auth.LoginResponse{
		Token: bearerToken,
	}, nil
}
