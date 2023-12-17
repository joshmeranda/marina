package gateway

import (
	"context"
	"fmt"

	"github.com/joshmeranda/marina/pkg/apis/user"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// todo: make namespace configurable

func (g *Gateway) Bootstrap(ctx context.Context) error {
	g.logger.Info("bootstrapping cluster")

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
		g.logger.Info("admin role already exists")
	} else if err != nil {
		return fmt.Errorf("failed to create admin role: %w", err)
	}

	req := &user.UserCreateRequest{
		User: &user.User{
			Name:  "admin",
			Roles: []string{"admin"},
		},
	}

	if _, err := g.CreateUser(ctx, req); err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("could not create admin user: %w", err)
	}

	return nil
}
