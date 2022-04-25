package k8s_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
)

var _ = Describe("ObjectMeta", func() {
	var pod v1.Pod
	var testMap map[string]string

	BeforeEach(func() {
		pod = v1.Pod{}
		testMap = map[string]string{"app": "test"}
	})

	Describe("Ensures Labels", func() {
		It("keeps labels", func() {
			pod.Labels = testMap
			Expect(k8s.EnsureLabels(&pod.ObjectMeta)).ToNot(BeNil())
			Expect(pod.Labels).Should(HaveLen(1))
			Expect(pod.Labels).Should(Equal(testMap))
		})
		It("creates labels", func() {
			Expect(pod.Labels).To(BeNil())
			Expect(k8s.EnsureLabels(&pod.ObjectMeta)).ToNot(BeNil())
			Expect(pod.Labels).ToNot(BeNil())
			Expect(pod.Labels).To(BeEmpty())
		})
	})

	Describe("Ensures Annotations", func() {
		It("keeps annotations", func() {
			pod.Annotations = testMap
			Expect(k8s.EnsureAnnotations(&pod.ObjectMeta)).ToNot(BeNil())
			Expect(pod.Annotations).Should(HaveLen(1))
			Expect(pod.Annotations).Should(Equal(testMap))
		})
		It("creates labels", func() {
			Expect(pod.Annotations).To(BeNil())
			Expect(k8s.EnsureAnnotations(&pod.ObjectMeta)).ToNot(BeNil())
			Expect(pod.Annotations).ToNot(BeNil())
			Expect(pod.Annotations).To(BeEmpty())
		})
	})

})
