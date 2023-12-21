package gateway_test

import (
	"context"
	"log/slog"
	"os"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/pkg/apis/auth"
	"github.com/joshmeranda/marina/pkg/apis/user"
	"github.com/joshmeranda/marina/pkg/gateway"
	authdriver "github.com/joshmeranda/marina/pkg/gateway/drivers/auth"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gateway Auth", Ordered, func() {
	var logger *slog.Logger
	var g *gateway.Gateway
	var namespace string

	BeforeAll(func() {
		var err error

		namespace, err = generateNamespaceName()
		Expect(err).ToNot(HaveOccurred())

		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

		authDriver := authdriver.NewLocal(k8sClient, namespace)

		g, err = gateway.NewGateway(gateway.WithLogger(logger), gateway.WithKubeClient(k8sClient), gateway.WithNamespace(namespace), gateway.WithAuthDriver(authDriver))
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())

		_, err = g.CreateUser(context.Background(), &user.UserCreateRequest{
			User: &user.User{
				Name:     "test-user",
				Password: "password",
				Roles:    []string{},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		user := marinav1.User{}
		err = k8sClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-user",
			Namespace: namespace,
		}, &user)
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

	When("using correct password", func() {
		It("can login", func(ctx context.Context) {
			resp, err := g.Login(ctx, &auth.LoginRequest{
				Secret:     "password",
				SecretType: auth.SecretType_Password,
				User:       "test-user",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Token).ToNot(BeZero())
		})
	})

	When("using wrong password", func() {
		It("cannot login", func(ctx context.Context) {
			resp, err := g.Login(ctx, &auth.LoginRequest{
				Secret:     "wrong-password",
				SecretType: auth.SecretType_Password,
				User:       "test-user",
			})
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})
})
