package e2e_test

import (
	"context"
	"fmt"

	"github.com/phayes/freeport"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marina Check", func() {
	var ctx context.Context
	var cancel context.CancelCauseFunc
	var namespace string
	var port int
	var err error

	BeforeEach(func() {
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
	})

	AfterEach(func() {
		cancel(fmt.Errorf("AfterAll test complete"))

		err := k8sClient.Delete(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	It("receives client health status", func() {
		Eventually(func() error {
			err := clientApp.RunContext(ctx, []string{"marina", "--address", fmt.Sprintf("127.0.0.1:%d", port), "health"})
			return err
		}, "5s").Should(Succeed())
	})
})
