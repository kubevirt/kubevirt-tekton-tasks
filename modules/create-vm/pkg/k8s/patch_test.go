package k8s_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
)

const expectedPatch1 = `[
  {
    "op": "add",
    "path": "/metadata/labels",
    "value": {
      "app": "test"
    }
  }
]`

var _ = Describe("Patch", func() {
	It("Creates Patch", func() {
		before := v1.Pod{}
		after := v1.Pod{}
		after.Labels = map[string]string{"app": "test"}

		patch, err := k8s.CreatePatch(before, after)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(patch)).To(Equal(expectedPatch1))
	})

})
