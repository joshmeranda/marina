package secret

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type kubeDriver struct {
	namespace  string
	kubeClient client.Client
}

func NewKubeDriver(kubeClient client.Client, namespace string) Driver {
	return &kubeDriver{
		namespace:  namespace,
		kubeClient: kubeClient,
	}
}

func (d *kubeDriver) Get(ctx context.Context, name string, key string) ([]byte, error) {
	objectKey := types.NamespacedName{
		Name:      name,
		Namespace: d.namespace,
	}

	var secret corev1.Secret
	if err := d.kubeClient.Get(ctx, objectKey, &secret); err != nil {
		return nil, fmt.Errorf("error fetching secret: %w", err)
	}

	value, ok := secret.Data[key]
	if !ok {
		return nil, fmt.Errorf("no such field")
	}

	return value, nil
}

func (d *kubeDriver) Set(ctx context.Context, name string, key string, value []byte) error {
	objectKey := types.NamespacedName{
		Name:      name,
		Namespace: d.namespace,
	}

	var secret corev1.Secret
	if err := d.kubeClient.Get(ctx, objectKey, &secret); errors.IsNotFound(err) {
		secret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: d.namespace,
			},
			Data: map[string][]byte{
				key: value,
			},
		}

		if err := d.kubeClient.Create(ctx, &secret); err != nil {
			return fmt.Errorf("error creating secret: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error fetching secret: %w", err)
	}

	secret.Data[key] = value

	if err := d.kubeClient.Update(ctx, &secret); err != nil {
		return fmt.Errorf("error updating secret: %w", err)
	}

	return nil
}
