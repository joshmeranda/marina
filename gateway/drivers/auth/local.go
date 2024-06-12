package auth

import (
	"context"
	"fmt"

	marinav1 "github.com/joshmeranda/marina/api/v1"
	"github.com/joshmeranda/marina/gateway/api/auth"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Local struct {
	kubeClient client.Client
	namespace  string
}

func NewLocal(kubeClient client.Client, namespace string) Driver {
	return &Local{
		kubeClient: kubeClient,
		namespace:  namespace,
	}
}

func (d *Local) Authenticate(ctx context.Context, req *auth.LoginRequest) error {
	user := marinav1.User{}
	err := d.kubeClient.Get(ctx, types.NamespacedName{
		Name:      req.User,
		Namespace: d.namespace,
	}, &user)
	if err != nil {
		return fmt.Errorf("unable to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Spec.Password, []byte(req.Secret)); err != nil {
		return err
	}

	return nil
}
