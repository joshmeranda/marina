package e2e_test

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	terminalv1 "github.com/joshmeranda/marina-operator/api/v1"
	marinacmd "github.com/joshmeranda/marina-operator/cmd"
	marinaclient "github.com/joshmeranda/marina/client"
	gatewayapp "github.com/joshmeranda/marina/cmd/gateway/app"
	marinapp "github.com/joshmeranda/marina/cmd/marina/app"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment

	clientApp cli.App
	serverApp cli.App
	marinaApp cli.App

	testDir        string
	kubeconfigPath string
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)

	suiteConfig, reporterConfig := GinkgoConfiguration()

	if focusFiles := os.Getenv("FOCUS_FILES"); focusFiles != "" {
		suiteConfig.FocusFiles = strings.Split(focusFiles, ",")
	}

	if focusStrings := os.Getenv("FOCUS_STRINGS"); focusStrings != "" {
		suiteConfig.FocusFiles = strings.Split(focusStrings, ",")
	}

	RunSpecs(t, "E2E Suite", suiteConfig, reporterConfig)
}

var _ = BeforeSuite(func() {
	clientApp = marinapp.App()
	serverApp = gatewayapp.App()
	marinaApp = marinacmd.App()

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "crds")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = terminalv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	testDir, err := os.MkdirTemp("", "marina-e2e")
	Expect(err).NotTo(HaveOccurred())

	kubeconfigPath = filepath.Join(testDir, "kubeconfig")
	err = storeKubeconfig(cfg, kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())

	err = os.RemoveAll(testDir)
	Expect(err).NotTo(HaveOccurred())
})

func generateRandomSuffix() string {
	length := 10
	suffixMin := int(math.Pow10(length))
	suffixMax := int(math.Pow10(length+1) - 1)

	suffix := rand.Intn(suffixMax-suffixMin) + suffixMin

	return fmt.Sprintf("%d", suffix)
}

func generateNamespaceName() (string, error) {
	suffix := generateRandomSuffix()
	return fmt.Sprintf("marina-system-%s", suffix), nil
}

func storeKubeconfig(cfg *rest.Config, kubeconfig string) error {
	clientConfig := clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters: map[string]*clientcmdapi.Cluster{
			"default": {
				Server:                   cfg.Host,
				CertificateAuthorityData: cfg.CAData,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"default": {
				Cluster:  "default",
				AuthInfo: "default",
			},
		},
		CurrentContext: "default",
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"default": {
				ClientCertificateData: cfg.CertData,
				ClientKeyData:         cfg.KeyData,
			},
		},
	}

	err := clientcmd.WriteToFile(clientConfig, kubeconfig)
	if err != nil {
		return err
	}

	return nil
}

func runServerWithArgs(ctx context.Context, namespace string, port int, additionalArgs []string) {
	GinkgoHelper()
	defer GinkgoRecover()

	args := []string{"marina-server",
		"--etcd", testEnv.ControlPlane.Etcd.URL.String(),
		"--namespace", namespace,
		"--port", fmt.Sprintf("%d", port),
		"--kubeconfig", kubeconfigPath,
		// "--silent",
	}
	args = append(args, additionalArgs...)

	By("by starting marina server")
	err := serverApp.RunContext(ctx, args)
	Expect(err).ToNot(HaveOccurred())
}

func runOperatorWithArgs(ctx context.Context, args []string) {
	GinkgoHelper()
	defer GinkgoRecover()

	args = append([]string{"marina-operator",
		"--kubeconfig", kubeconfigPath,
		"--metrics-bind-address", "0",
		"--health-probe-bind-address", "0",
		"--webhook-port", "0",
		// "--silent",
	}, args...)

	By("by starting marina operator")
	err := marinaApp.RunContext(ctx, args)
	Expect(err).ToNot(HaveOccurred())
}

func clientCall(ctx context.Context, address string, bearerToken string, f func(ctx context.Context, client *marinaclient.Client) error) error {
	conn, err := grpc.Dial(address,
		grpc.WithUnaryInterceptor(marinaclient.TokenAuthInterceptor(bearerToken)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	c := marinaclient.NewClient(conn, logger)

	err = f(ctx, c)

	return err
}
