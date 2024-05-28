package e2e_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/joshmeranda/marina/apis/auth"
	"github.com/joshmeranda/marina/apis/terminal"
	"github.com/joshmeranda/marina/apis/user"
	marinaclient "github.com/joshmeranda/marina/client"
	"github.com/phayes/freeport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Admin Login", Ordered, func() {
	var ctx context.Context
	var cancel context.CancelCauseFunc
	var namespace string
	var port int
	var err error
	var configDir string

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
		client := marinaclient.NewClient(conn, logger)

		services := []string{
			terminal.TerminalService_ServiceDesc.ServiceName,
			user.UserService_ServiceDesc.ServiceName,
			auth.AuthService_ServiceDesc.ServiceName,
		}

		Eventually(func() int {
			var n int

			for _, serviceName := range services {
				status, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
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

	It("fails to access endpoints requiring auth", func() {
		args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "list"}
		err := clientApp.RunContext(ctx, args)
		Expect(err).To(HaveOccurred())
	})

	It("receives credentials", func() {
		passwordSecret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "marina-bootstrap-password",
				Namespace: namespace,
			},
		}
		err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&passwordSecret), &passwordSecret)
		Expect(err).ToNot(HaveOccurred())

		password := passwordSecret.Data["password"]

		args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "auth", "login", "password", "admin", string(password)}
		err = clientApp.RunContext(ctx, args)
		Expect(err).ToNot(HaveOccurred())

		Expect(path.Join(configDir, "config.yaml")).To(BeAnExistingFile())
	})

	It("is able to access endpoints requiring auth", func() {
		args := []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "--config-dir", configDir, "user", "list"}
		err := clientApp.RunContext(ctx, args)
		Expect(err).ToNot(HaveOccurred())
	})
})
