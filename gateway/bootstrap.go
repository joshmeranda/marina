package gateway

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/joshmeranda/marina/gateway/api/user"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	charSet = "abcdefghijklmnopqrstuvwxyz0123456789-_"

	bootstrapSecretName = "marina-bootstrap-password"
	bootstrapSecretKey  = "password"
)

func generateRandomPassword(length int) ([]byte, error) {
	bigLength := big.NewInt(int64(len(charSet)))
	passwordRaw := make([]byte, length)
	for i := range passwordRaw {
		r, err := rand.Int(rand.Reader, bigLength)
		if err != nil {
			return nil, err
		}
		passwordRaw[i] = charSet[r.Int64()]
	}
	return passwordRaw, nil
}

func (g *Gateway) ensureAdminRole(ctx context.Context) error {
	adminRole := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "admin",
			Namespace: g.namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"core.marina.io"},
				Resources: []string{"users", "terminals"},
				Verbs:     []string{"create", "delete", "get", "list", "update", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create", "delete", "get"},
			},
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"roles", "rolebindings"},
				Verbs:     []string{"create", "delete", "get"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments"},
				Verbs:     []string{"create", "delete", "get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts/token"},
				Verbs:     []string{"get"},
			},
		},
	}

	// todo: could be configurable
	backoff := wait.Backoff{
		Duration: time.Second,
		Factor:   1.1,
		Steps:    20,
		Cap:      time.Minute * 5,
	}

	condition := func(ctx context.Context) (bool, error) {
		err := g.kubeClient.Create(ctx, &adminRole)

		switch {
		case err == nil:
			g.logger.Info("created admin role")
			return true, nil
		case errors.IsAlreadyExists(err):
			g.logger.Debug("admin role already exists")
			return true, nil
		default:
			g.logger.Warn("failed to create admin role, retrying", "error", err)
			return false, nil
		}
	}

	if err := wait.ExponentialBackoffWithContext(ctx, backoff, condition); err != nil {
		return fmt.Errorf("failed to create admin role: %w", err)
	}

	return nil
}

func (g *Gateway) ensureTokenSigningSecret(ctx context.Context) error {
	signingKey := make([]byte, 512)
	if _, err := rand.Read(signingKey); err != nil {
		return fmt.Errorf("failed to generate signing key: %w", err)
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenSigningSecretName,
			Namespace: g.namespace,
		},
		Data: map[string][]byte{
			TokenSigningSecretField: signingKey,
		},
	}

	if err := g.kubeClient.Create(ctx, &secret); errors.IsAlreadyExists(err) {
		g.logger.Debug("signing secret role already exists")
		return nil
	} else if err != nil {
		return err
	}

	g.logger.Info("created signing secret")

	return nil
}

func (g *Gateway) ensureAdminUser(ctx context.Context) error {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bootstrapSecretName,
			Namespace: g.namespace,
		},
		Data: map[string][]byte{},
	}

	var bootstrapPassword []byte

	if err := g.kubeClient.Get(ctx, client.ObjectKeyFromObject(&secret), &secret); errors.IsNotFound(err) {
		g.logger.Debug("bootstrap secret does not exist, creating")

		// todo: allow users to specify bootstrap password
		bootstrapPassword, err = generateRandomPassword(20)
		if err != nil {
			return fmt.Errorf("failed to generate bootstrap password: %w", err)
		}

		secret.Data[bootstrapSecretKey] = bootstrapPassword

		if err := g.kubeClient.Create(ctx, &secret); errors.IsAlreadyExists(err) {
			g.logger.Debug("bootstrap secret already exists")
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to create bootstrap secret: %w", err)
		}
	}

	req := &user.UserCreateRequest{
		User: &user.User{
			Name:     "admin",
			Password: bootstrapPassword,
			Roles:    []string{"admin"},
		},
	}

	if _, err := g.CreateUser(ctx, req); errors.IsAlreadyExists(err) {
		g.logger.Debug("admin user already exists")
		return nil
	} else if err != nil {
		return err
	}

	g.logger.Info("created admin user")

	return nil
}

func (g *Gateway) Bootstrap(ctx context.Context) error {
	g.logger.Info("bootstrapping cluster")

	if err := g.ensureAdminRole(ctx); err != nil {
		return fmt.Errorf("failed to create admin role: %w", err)
	}

	if err := g.ensureAdminUser(ctx); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := g.ensureTokenSigningSecret(ctx); err != nil {
		return fmt.Errorf("failed to create signing secret: %w", err)
	}

	return nil
}
