package secret_test

import (
	"context"

	"github.com/joshmeranda/marina/pkg/gateway/drivers/secret"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var _ = Describe("Kube Secret Driver", Ordered, func() {
	var cfg *rest.Config
	var k8sClient client.Client
	var testEnv *envtest.Environment

	var driver secret.Driver
	var namespace string

	BeforeAll(func() {
		logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

		By("bootstrapping test environment")
		testEnv = &envtest.Environment{
			ErrorIfCRDPathMissing: true,
		}

		var err error
		// cfg is defined in this file globally.
		cfg, err = testEnv.Start()
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())

		//+kubebuilder:scaffold:scheme

		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
		Expect(k8sClient).NotTo(BeNil())

		namespace = "marina-system"
		driver = secret.NewKubeDriver(k8sClient, namespace)

		err = k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace,
				Namespace: namespace,
			},
		})
		Expect(err).ToNot(HaveOccurred())

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"test-key": []byte("test-value"),
			},
		}

		err = k8sClient.Create(context.Background(), &secret)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		By("tearing down the test environment")
		err := testEnv.Stop()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can pull secret", func(ctx context.Context) {
		data, err := driver.Get(ctx, "test-secret", "test-key")
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal([]byte("test-value")))
	})

	It("errors on non-existent secret", func(ctx context.Context) {
		data, err := driver.Get(ctx, "test-secret-2", "test-key")
		Expect(err).To(HaveOccurred())
		Expect(data).To(BeNil())
	})
})
