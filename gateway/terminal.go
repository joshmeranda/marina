package gateway

import (
	"context"
	"fmt"
	"time"

	marinav1 "github.com/joshmeranda/marina/api/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joshmeranda/marina/gateway/api/core"
	"github.com/joshmeranda/marina/gateway/api/terminal"
)

var _ terminal.TerminalServiceServer = &Gateway{}

const (
	LabelKeyTerminalName = "terminal.marina.io/name"
	LabelKeyUsername     = "user.marina.io/username"
)

func (g *Gateway) getPodForTerminal(ctx context.Context, kubeClient client.Client, t *marinav1.Terminal) (*corev1.Pod, error) {
	backoff := wait.Backoff{
		Duration: time.Second,
		Factor:   1.1,
		Steps:    20,
		Cap:      time.Minute * 5,
	}

	podList := corev1.PodList{}

	labels := client.MatchingLabels{
		LabelKeyTerminalName: t.Name,
		LabelKeyUsername:     t.Spec.User,
	}

	condition := func(ctx context.Context) (bool, error) {
		if err := kubeClient.List(ctx, &podList, client.InNamespace(t.Namespace), labels); err != nil {
			g.logger.Warn("could not list pods for terminal, retrying", "err", err)
			return false, nil
		}

		if len(podList.Items) == 0 {
			g.logger.Warn("no pods found for terminal, retrying...")
			return false, nil
		}

		return true, nil
	}

	if err := wait.ExponentialBackoffWithContext(ctx, backoff, condition); err != nil {
		return nil, fmt.Errorf("failed to get pod for terminal: %w", err)
	}

	if len(podList.Items) != 1 {
		return nil, fmt.Errorf("expected 1 pod, got %d", len(podList.Items))
	}

	return &podList.Items[0], nil
}

func (g *Gateway) getExecTokenForTerminal(ctx context.Context, kubeClient client.Client, t *marinav1.Terminal) (string, error) {
	tokenRequest := &authenticationv1.TokenRequest{
		// todo: we will to add either expiration or bind this to a secret to limit the token's lifetime
		// Spec: authenticationv1.TokenRequestSpec{
		// 	Audiences: []string{"system:serviceaccount:" + terminal.Namespace + ":" + terminal.Name},
		// 	BoundObjectRef: &authenticationv1.BoundObjectReference{
		// 		Kind:       terminal.Kind,
		// 		APIVersion: terminal.APIVersion,
		// 		Name:       terminal.Name,
		// 		UID:        terminal.UID,
		// 	},
		// },
	}

	sa := &corev1.ServiceAccount{}
	if err := kubeClient.Get(ctx, client.ObjectKey{Namespace: t.Namespace, Name: t.Spec.User}, sa); err != nil {
		return "", fmt.Errorf("could not fetch service account: %w", err)
	}

	if err := g.kubeClient.SubResource("token").Create(ctx, sa, tokenRequest); err != nil {
		return "", fmt.Errorf("could not create token: %w", err)
	}

	return tokenRequest.Status.Token, nil
}

func (g *Gateway) CreateTerminal(ctx context.Context, req *terminal.TerminalCreateRequest) (*terminal.TerminalCreateResponse, error) {
	kubeClient, err := g.clientFromContext(ctx, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to impersonate user: %w", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from context: %w", err)
	}

	t := marinav1.Terminal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name.Name,
			Namespace: req.Name.Namespace,
		},
		Spec: marinav1.TerminalSpec{
			Image: req.Spec.Image,
			User:  user,
		},
	}

	if err := kubeClient.Create(ctx, &t); err != nil {
		return nil, err
	}

	pod, err := g.getPodForTerminal(ctx, kubeClient, &t)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pod for terminal: %w", err)
	}

	token, err := g.getExecTokenForTerminal(ctx, kubeClient, &t)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exec token for terminal: %w", err)
	}

	return &terminal.TerminalCreateResponse{
		Pod: &core.NamespacedName{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		Token: token,
	}, nil
}

func (g *Gateway) DeleteTerminal(ctx context.Context, req *terminal.TerminalDeleteRequest) (*empty.Empty, error) {
	kubeClient, err := g.clientFromContext(ctx, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to impersonate user: %w", err)
	}

	t := marinav1.Terminal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name.Name,
			Namespace: req.Name.Namespace,
		},
	}

	if err := kubeClient.Delete(ctx, &t); err != nil {
		return nil, fmt.Errorf("could not delete terminal '%s': %w", req.Name, err)
	}

	return &empty.Empty{}, nil
}
