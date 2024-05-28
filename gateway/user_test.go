package gateway_test

import (
	"context"
	"log/slog"
	"os"

	terminalv1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/apis/user"
	"github.com/joshmeranda/marina/gateway"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gateway User Service", Ordered, func() {
	var logger *slog.Logger
	var g *gateway.Gateway
	var namespace string

	BeforeAll(func() {
		var err error

		namespace, err = generateNamespaceName()
		Expect(err).ToNot(HaveOccurred())

		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

		g, err = gateway.NewGateway(gateway.WithLogger(logger), gateway.WithKubeClient(k8sClient), gateway.WithNamespace(namespace))
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.Background(), &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "SomeRole",
				Namespace: namespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		err = k8sClient.Create(context.Background(), &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "AnotherRole",
				Namespace: namespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())
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

	When("user roles exist", func() {
		It("can create a user", func(ctx context.Context) {
			req := &user.UserCreateRequest{
				User: &user.User{
					Name:  "bbaggins",
					Roles: []string{"SomeRole"},
				},
			}

			_, err := g.CreateUser(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			var foundUser terminalv1.User
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      req.User.Name,
				Namespace: namespace,
			}, &foundUser)
			Expect(err).ToNot(HaveOccurred())
		})

		It("can list users", func(ctx context.Context) {
			req := &user.UserListRequest{
				Query: &user.UserQuery{},
			}

			resp, err := g.ListUser(ctx, req)
			Expect(err).ToNot(HaveOccurred())

			type internalUser struct {
				Name     string
				Password []byte
				Roles    []string
			}

			userList := make([]internalUser, len(resp.Users))
			for i, u := range resp.Users {
				userList[i] = internalUser{
					Name:     u.Name,
					Password: u.Password,
					Roles:    u.Roles,
				}
			}

			expected := []internalUser{
				{
					Name:  "bbaggins",
					Roles: []string{"SomeRole"},
				},
			}
			Expect(userList).To(ConsistOf(expected))
		})

		It("can edit a user", func(ctx context.Context) {
			req := &user.UserUpdateRequest{
				User: &user.User{
					Name:  "bbaggins",
					Roles: []string{"AnotherRole"},
				},
			}

			_, err := g.UpdateUser(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			var foundUser terminalv1.User
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      req.User.Name,
				Namespace: namespace,
			}, &foundUser)
			Expect(err).ToNot(HaveOccurred())
			Expect(foundUser.Spec.Roles).To(Equal(req.User.Roles))
		})

		It("can delete a user", func(ctx context.Context) {
			req := &user.UserDeleteRequest{
				Name: "bbaggins",
			}

			_, err := g.DeleteUser(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			var foundUser terminalv1.User
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      req.Name,
				Namespace: namespace,
			}, &foundUser)
			Expect(errors.IsNotFound(err)).To(BeTrue())
			Expect(foundUser).To(BeZero())
		})
	})

	When("user roles do not exist", func() {
		It("returns an error", func(ctx context.Context) {
			req := &user.UserCreateRequest{
				User: &user.User{
					Name:  "bbaggins",
					Roles: []string{"NonExistantRole"},
				},
			}

			_, err := g.CreateUser(ctx, req)
			Expect(err).To(HaveOccurred())

			var foundUser terminalv1.User
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      req.User.Name,
				Namespace: namespace,
			}, &foundUser)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})
})
