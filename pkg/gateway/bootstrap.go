package gateway

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/joshmeranda/marina/pkg/apis/user"
	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	charSet = "abcdefghijklmnopqrstuvwxyz0123456789-_"

	bootstrapSecretName = "marina-bootstrap-password"
	bootstrapSecretKey  = "password"
)

func generateRandomPassword(length int) (string, error) {
	bigLength := big.NewInt(int64(len(charSet)))
	passwordRaw := make([]byte, length)
	for i := range passwordRaw {
		r, err := rand.Int(rand.Reader, bigLength)
		if err != nil {
			return "", err
		}
		passwordRaw[i] = charSet[r.Int64()]
	}
	return string(passwordRaw), nil
}

func (g *Gateway) ensureRole(ctx context.Context) error {
	// todo: should create this role in the helm chart
	adminRole := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "admin",
			Namespace: g.namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{rbacv1.VerbAll},
				APIGroups: []string{"terminal.marina.io"},
				Resources: []string{rbacv1.ResourceAll},
			},
		},
	}

	if err := g.kubeClient.Create(ctx, &adminRole); errors.IsAlreadyExists(err) {
		g.logger.Debug("admin role already exists")
		return nil
	} else if err != nil {
		return err
	}

	g.logger.Info("created admin role")

	return nil
}

func (g *Gateway) ensureUser(ctx context.Context) error {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bootstrapSecretName,
			Namespace: g.namespace,
		},
		Data: map[string][]byte{},
	}

	var bootstrapHash []byte

	if err := g.kubeClient.Get(ctx, client.ObjectKeyFromObject(&secret), &secret); errors.IsNotFound(err) {
		// todo: allow users to specify bootstrap password
		bootstrapPassword, err := generateRandomPassword(20)
		if err != nil {
			return fmt.Errorf("failed to generate bootstrap password: %w", err)
		}

		bootstrapHash, err = bcrypt.GenerateFromPassword([]byte(bootstrapPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to generate bootstrap hash: %w", err)
		}

		secret.Data[bootstrapSecretKey] = bootstrapHash

		if err := g.kubeClient.Create(ctx, &secret); errors.IsAlreadyExists(err) {
			g.logger.Debug("bootstrap secret already exists")
		} else if err != nil {
			return fmt.Errorf("failed to create bootstrap secret: %w", err)
		}
	}

	req := &user.UserCreateRequest{
		User: &user.User{
			Name:     "admin",
			Password: string(bootstrapHash),
			Roles:    []string{"admin"},
		},
	}

	if _, err := g.CreateUser(ctx, req); errors.IsAlreadyExists(err) {
		g.logger.Debug("admin user already exists")
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (g *Gateway) Bootstrap(ctx context.Context) error {
	g.logger.Info("bootstrapping cluster")

	if err := g.ensureRole(ctx); err != nil {
		return fmt.Errorf("failed to create admin role: %w", err)
	}

	if err := g.ensureUser(ctx); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
