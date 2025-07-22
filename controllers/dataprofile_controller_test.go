package controllers

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
phantomdbv1alpha1 "github.com/vandal-db/vandal-db/apis/v1alpha1"
	//+kubebuilder:scaffold:imports
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = phantomdbv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("DataProfile controller", func() {
	Context("When creating a DataProfile", func() {
		It("Should create a new DataProfile object", func() {
			By("Creating a new DataProfile")
			ctx := context.Background()
			dataProfile := &phantomdbv1alpha1.DataProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-dataprofile",
					Namespace: "default",
				},
				Spec: phantomdbv1alpha1.DataProfileSpec{
					Target: phantomdbv1alpha1.DatabaseTarget{
						SecretName: "test-secret",
					},
				},
			}
			Expect(k8sClient.Create(ctx, dataProfile)).Should(Succeed())
		})
	})
})
