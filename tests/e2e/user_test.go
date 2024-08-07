package e2e_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	marinaclient "github.com/joshmeranda/marina/client"
	"github.com/joshmeranda/marina/gateway/api/auth"
	"github.com/joshmeranda/marina/gateway/api/terminal"
	"github.com/joshmeranda/marina/gateway/api/user"
	"github.com/phayes/freeport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// todo: allow login for multiple users (see kubeconfig for example)

var _ = Describe("User", Ordered, func() {
	var ctx context.Context
	var cancel context.CancelCauseFunc
	var namespace string
	var port int
	var err error
	var configDir string
	var marinaClient *marinaclient.Client
	var bearerToken string

	BeforeAll(func() {
		configDir = path.Join(testDir, "admin-login-config-"+generateRandomSuffix())
		ctx, cancel = context.WithCancelCause(context.Background())

		port, err = freeport.GetFreePort()
		Expect(err).ToNot(HaveOccurred())

		var err error

		namespace, err = generateNamespaceName()
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())

		go func() {
			GinkgoRecover()

			runOperatorWithArgs(ctx, nil)
		}()

		go func() {
			GinkgoRecover()

			runServerWithArgs(ctx, namespace, port, nil)
		}()

		By("wait for gateway to be up")
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).ToNot(HaveOccurred())

		logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
		marinaClient = marinaclient.NewClient(conn, logger)

		services := []string{
			terminal.TerminalService_ServiceDesc.ServiceName,
			user.UserService_ServiceDesc.ServiceName,
			auth.AuthService_ServiceDesc.ServiceName,
		}

		Eventually(func() int {
			var n int

			for _, serviceName := range services {
				status, err := marinaClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
					Service: serviceName,
				})
				if err != nil {
					// todo: we might want to log any failures here
					continue
				}

				if status.Status == grpc_health_v1.HealthCheckResponse_SERVING {
					n += 1
				}
			}

			return n
		}, "5s").Should(Equal(len(services)))

		By("receiving bearer token from gateway")
		passwordSecret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "marina-bootstrap-password",
				Namespace: namespace,
			},
		}
		err = k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&passwordSecret), &passwordSecret)
		Expect(err).ToNot(HaveOccurred())

		password := passwordSecret.Data["password"]

		err = clientCall(ctx, fmt.Sprintf("127.0.0.1:%d", port), "", func(ctx context.Context, client *marinaclient.Client) error {
			resp, err := client.Login(ctx, &auth.LoginRequest{
				Secret:     password,
				SecretType: auth.SecretType_Password,
				User:       "admin",
			})
			if err != nil {
				return err
			}

			bearerToken = resp.Token

			return nil
		})
		Expect(err).ToNot(HaveOccurred())

		args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "auth", "login", "password", "admin", string(password)}
		err = clientApp.RunContext(ctx, args)
		Expect(err).ToNot(HaveOccurred())

		Expect(path.Join(configDir, "config.yaml")).To(BeAnExistingFile())

		roles := []rbacv1.Role{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hobbit",
					Namespace: namespace,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "burglar",
					Namespace: namespace,
				},
			},
		}

		for _, role := range roles {
			err = k8sClient.Create(context.Background(), &role)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	AfterAll(func() {
		cancel(fmt.Errorf("AfterAll test complete"))

		err := k8sClient.Delete(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())

		err = os.RemoveAll(configDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("only allows 2 args", func() {
		args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "create", "bbaggins"}
		err := clientApp.RunContext(ctx, args)
		Expect(err).To(HaveOccurred())

		args = []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "create", "bbaggins", "secret", "fbaggins", "abother-secret"}
		err = clientApp.RunContext(ctx, args)
		Expect(err).To(HaveOccurred())
	})

	It("can create a user", func() {
		args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "create", "--add-role", "hobbit", "bbaggins", "secret"}
		err := clientApp.RunContext(ctx, args)
		Expect(err).ToNot(HaveOccurred())

		err = clientCall(ctx, fmt.Sprintf("127.0.0.1:%d", port), bearerToken, func(ctx context.Context, client *marinaclient.Client) error {
			_, err = client.GetUser(ctx, &user.UserGetRequest{
				Name: "bbaggins",
			})
			return err
		})
		Expect(err).ToNot(HaveOccurred())
	})

	When("updating users", func() {
		When("updating roles", func() {
			It("fails to add roles that do not exist", func() {
				args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "update", "--add-role", "role-does-not-exist", "bbaggins"}
				err := clientApp.RunContext(ctx, args)
				Expect(err).To(HaveOccurred())
			})

			It("can update roles", func() {
				args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "update", "--add-role", "burglar", "--remove-role", "hobbit", "bbaggins"}
				err := clientApp.RunContext(ctx, args)
				Expect(err).ToNot(HaveOccurred())

				var u *user.User
				err = clientCall(ctx, fmt.Sprintf("127.0.0.1:%d", port), bearerToken, func(ctx context.Context, client *marinaclient.Client) error {
					u, err = client.GetUser(ctx, &user.UserGetRequest{
						Name: "bbaggins",
					})
					return err
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(u.Roles).To(ContainElement("burglar"))
			})
		})

		It("can update user password", func() {
			args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "update", "--password", "new-secret", "bbaggins"}
			err := clientApp.RunContext(ctx, args)
			Expect(err).ToNot(HaveOccurred())

			err = clientCall(ctx, fmt.Sprintf("127.0.0.1:%d", port), bearerToken, func(ctx context.Context, client *marinaclient.Client) error {
				_, err = client.Login(ctx, &auth.LoginRequest{
					Secret:     []byte("secret"),
					SecretType: auth.SecretType_Password,
					User:       "bbaggins",
				})
				return err
			})
			Expect(err).To(HaveOccurred())

			err = clientCall(ctx, fmt.Sprintf("127.0.0.1:%d", port), bearerToken, func(ctx context.Context, client *marinaclient.Client) error {
				_, err = client.Login(ctx, &auth.LoginRequest{
					Secret:     []byte("new-secret"),
					SecretType: auth.SecretType_Password,
					User:       "bbaggins",
				})
				return err
			})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("can delete a user", func() {
		args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "delete", "bbaggins"}
		err := clientApp.RunContext(ctx, args)
		Expect(err).ToNot(HaveOccurred())

		err = clientCall(ctx, fmt.Sprintf("127.0.0.1:%d", port), bearerToken, func(ctx context.Context, client *marinaclient.Client) error {
			_, err = client.GetUser(ctx, &user.UserGetRequest{
				Name: "bbaggins",
			})
			return err
		})
		Expect(err).To(HaveOccurred())
	})
})
