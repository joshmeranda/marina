package gateway_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	terminalv1 "github.com/joshmeranda/marina/api/v1"
	"github.com/joshmeranda/marina/gateway"
	"github.com/joshmeranda/marina/gateway/api/core"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	"github.com/joshmeranda/marina/gateway/images"
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

		imagesAccessList := images.ImagesAccessList{
			Blocked: []images.ImageMatcher{
				{
					Tag: "latest",
				},
			},
		}

		g, err = gateway.NewGateway(gateway.WithLogger(logger), gateway.WithKubeConfig(cfg), gateway.WithNamespace(namespace), gateway.WithAccessList(imagesAccessList))
		Expect(err).ToNot(HaveOccurred())
		k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})

		err = g.Bootstrap(context.Background())
		Expect(err).ToNot(HaveOccurred())

		terminalPod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "terminal-test",
				Namespace: namespace,
				Labels: map[string]string{
					gateway.LabelKeyTerminalName: "terminal-test",
					gateway.LabelKeyUsername:     "admin",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "exec-shell",
						Image:   "library/busybox:1.36.1",
						Command: []string{"/bin/sh", "-ec", "trap : TERM INT; sleep infinity & wait"},
					},
				},
			},
		}
		err = k8sClient.Create(context.Background(), terminalPod)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() corev1.PodPhase {
			err = k8sClient.Get(context.Background(), types.NamespacedName{
				Name:      terminalPod.Name,
				Namespace: terminalPod.Namespace,
			}, terminalPod)
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("=== [GatewayTerminalService BerforeAll] 000 '%+v' ===\n", terminalPod.Status)

			return terminalPod.Status.Phase
		}).Should(Equal(corev1.PodRunning), "10s")
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
			Spec: &terminal.TerminalSpec{
				Image: "repository/image:tag",
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

	It("does not create a terminal with a blocked image", func(ctx context.Context) {
		req := &terminal.TerminalCreateRequest{
			Name: &core.NamespacedName{
				Name:      "terminal-test",
				Namespace: namespace,
			},
			Spec: &terminal.TerminalSpec{
				Image: "repository/image:latest",
			},
		}
		_, err := g.CreateTerminal(ctx, req)
		Expect(err).To(HaveOccurred())

		var foundTerminal terminalv1.Terminal
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      req.Name.Name,
			Namespace: req.Name.Namespace,
		}, &foundTerminal)
		Expect(err).To(HaveOccurred())
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
