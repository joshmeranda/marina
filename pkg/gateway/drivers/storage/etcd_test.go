package storage_test

import (
	"context"

	"github.com/joshmeranda/marina/pkg/gateway/drivers/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var _ = Describe("Etcd Storage", Ordered, func() {
	var testEtcd envtest.Etcd

	BeforeAll(func() {
		logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

		By("bootstrapping test etcd")

		testEtcd = envtest.Etcd{}
		err := testEtcd.Start()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		err := testEtcd.Stop()
		Expect(err).NotTo(HaveOccurred())
	})

	When("storing primitives", Ordered, func() {
		var etcdStore storage.KeyValueStore[string, string]
		var sampleKey string
		var sampleData string

		BeforeAll(func() {
			config := clientv3.Config{
				Endpoints: []string{testEtcd.URL.String()},
			}
			etcdClient, err := clientv3.New(config)
			Expect(err).ToNot(HaveOccurred())

			etcdStore = storage.NewEtcdStore[string](etcdClient, yaml.Marshal, yaml.Unmarshal)
			sampleKey = "test-key"
			sampleData = "test-value"
		})

		It("can store data in etcd", func(ctx context.Context) {
			err := etcdStore.Set(ctx, sampleKey, sampleData)
			Expect(err).ToNot(HaveOccurred())
		})

		It("can get data", func() {
			value, err := etcdStore.Get(context.Background(), sampleKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal(sampleData))
		})

		It("can delete data", func() {
			err := etcdStore.Delete(context.Background(), sampleKey)
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails to get deleted data", func() {
			value, err := etcdStore.Get(context.Background(), sampleKey)
			Expect(err).To(HaveOccurred())
			Expect(value).To(BeZero())
		})
	})

	When("storing non-primitive data", Ordered, func() {
		type SampleData struct {
			Name string
			Age  int
		}

		var etcdStore storage.KeyValueStore[string, SampleData]
		var sampleKey string
		var sampleData SampleData

		BeforeAll(func() {
			config := clientv3.Config{
				Endpoints: []string{testEtcd.URL.String()},
			}
			etcdClient, err := clientv3.New(config)
			Expect(err).ToNot(HaveOccurred())

			etcdStore = storage.NewEtcdStore[SampleData](etcdClient, yaml.Marshal, yaml.Unmarshal)
			sampleData = SampleData{
				Name: "Bilbo Baggins",
				Age:  111,
			}
			sampleKey = "test-key"
		})

		It("can store data in etcd", func(ctx context.Context) {
			err := etcdStore.Set(ctx, sampleKey, sampleData)
			Expect(err).ToNot(HaveOccurred())
		})

		It("can get data", func() {
			value, err := etcdStore.Get(context.Background(), sampleKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal(sampleData))
		})

		It("can delete data", func() {
			err := etcdStore.Delete(context.Background(), sampleKey)
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails to get deleted data", func() {
			value, err := etcdStore.Get(context.Background(), sampleKey)
			Expect(err).To(HaveOccurred())
			Expect(value).To(BeZero())
		})
	})
})
