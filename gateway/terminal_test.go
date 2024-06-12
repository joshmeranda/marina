package gateway_test

import (
	"context"
	"log/slog"
	"os"

	terminalv1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/gateway"
	"github.com/joshmeranda/marina/gateway/api/core"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gateway Terminal Service", Ordered, func() {
	var logger *slog.Logger
	var g *gateway.Gateway
	var namespace string

	BeforeAll(func() {
		var err error

		namespace, err = generateNamespaceName()
		Expect(err).ToNot(HaveOccurred())

		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

		g, err = gateway.NewGateway(gateway.WithLogger(logger), gateway.WithKubeConfig(cfg), gateway.WithNamespace(namespace))
		Expect(err).ToNot(HaveOccurred())
		k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
	})

	AfterAll(func() {
		err := k8sClient.Delete(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	It("can create a terminal", func(ctx context.Context) {
		req := &terminal.TerminalCreateRequest{
			Name: &core.NamespacedName{
				Name:      "terminal-test",
				Namespace: namespace,
			},
		}
		_, err := g.CreateTerminal(ctx, req)
		Expect(err).NotTo(HaveOccurred())

		var foundTerminal terminalv1.Terminal
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      req.Name.Name,
			Namespace: req.Name.Namespace,
		}, &foundTerminal)
		Expect(err).ToNot(HaveOccurred())
	})

	It("can delete a terminal", func(ctx context.Context) {
		req := &terminal.TerminalDeleteRequest{
			Name: &core.NamespacedName{
				Name:      "terminal-test",
				Namespace: namespace,
			},
		}
		_, err := g.DeleteTerminal(ctx, req)
		Expect(err).NotTo(HaveOccurred())

		var foundTerminal terminalv1.Terminal
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      req.Name.Name,
			Namespace: req.Name.Namespace,
		}, &foundTerminal)
		Expect(errors.IsNotFound(err)).To(BeTrue())
		Expect(foundTerminal).To(BeZero())
	})
})
