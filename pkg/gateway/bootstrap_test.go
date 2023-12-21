package gateway_test

import (
	"context"
	"log/slog"
	"os"
	"sync"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/pkg/gateway"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Gateway Bootstrap", Ordered, func() {
	var logger *slog.Logger
	var g *gateway.Gateway
	var namespace string

	var user marinav1.User
	var role rbacv1.Role
	var secret corev1.Secret

	BeforeAll(func() {
		var err error

		namespace, err = generateNamespaceName()
		Expect(err).ToNot(HaveOccurred())

		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		user = marinav1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "admin",
				Namespace: namespace,
			},
		}

		role = rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "admin",
				Namespace: namespace,
			},
		}

		secret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "marina-bootstrap-password",
				Namespace: namespace,
			},
		}

		g, err = gateway.NewGateway(gateway.WithLogger(logger), gateway.WithKubeClient(k8sClient), gateway.WithNamespace(namespace))
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())
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

	AfterEach(func() {
		err := k8sClient.Delete(context.TODO(), &user)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Delete(context.TODO(), &role)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Delete(context.Background(), &secret)
		Expect(err).ToNot(HaveOccurred())
	})

	When("bootstraping in sequence", Ordered, func() {
		When("initially bootstrapping", func() {
			It("successfully creates admin role and user", func(ctx context.Context) {
				err := g.Bootstrap(ctx)
				Expect(err).ToNot(HaveOccurred())

				user := user
				role := role
				secret := secret

				err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&user), &user)
				Expect(err).ToNot(HaveOccurred())

				err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&role), &role)
				Expect(err).ToNot(HaveOccurred())

				err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&secret), &secret)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("successive bootstraps", func() {
			It("successfully does nothing", func(ctx context.Context) {
				err := g.Bootstrap(ctx)
				Expect(err).ToNot(HaveOccurred())

				user := user
				role := role
				secret := secret

				err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&user), &user)
				Expect(err).ToNot(HaveOccurred())

				err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&role), &role)
				Expect(err).ToNot(HaveOccurred())

				err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&secret), &secret)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	When("bootstrapping in parallel", func() {
		It("successfully creates admin role and user", func(ctx context.Context) {
			n := 10
			wg := sync.WaitGroup{}
			wg.Add(n)

			for i := 0; i < n; i++ {
				go func() {
					defer GinkgoRecover()

					ctx := context.WithoutCancel(ctx)
					err := g.Bootstrap(ctx)
					wg.Done()
					Expect(err).ToNot(HaveOccurred())
				}()
			}

			wg.Wait()

			user := user
			role := role
			secret := secret

			err := k8sClient.Get(ctx, client.ObjectKeyFromObject(&user), &user)
			Expect(err).ToNot(HaveOccurred())

			err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&role), &role)
			Expect(err).ToNot(HaveOccurred())

			err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&secret), &secret)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
