package controllers

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vandalv1alpha1 "github.com/Oridak771/Vandal/apis/v1alpha1"
)

var _ = Describe("DataClone controller", func() {
	Context("When creating a DataClone", func() {
		It("Should create a new DataClone object", func() {
			By("Creating a new DataClone")
			ctx := context.Background()
			dataClone := &vandalv1alpha1.DataClone{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-dataclone",
					Namespace: "default",
				},
				Spec: vandalv1alpha1.DataCloneSpec{
					SourceProfile: "test-dataprofile",
				},
			}
			Expect(k8sClient.Create(ctx, dataClone)).Should(Succeed())
		})
	})
})
