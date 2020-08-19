package k8s_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	v1core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
)

var _ = Describe("Ownerereferences", func() {
	table.DescribeTable("Append OwnerReferences", func(refs []v1.OwnerReference) {
		pod := v1core.Pod{}
		pod.OwnerReferences = k8s.AppendOwnerReferences(pod.OwnerReferences, refs)
		Expect(pod.OwnerReferences).To(HaveLen(len(refs)))

	},
		table.Entry("nil", nil),
		table.Entry("empty", []v1.OwnerReference{}),
		table.Entry("one", []v1.OwnerReference{{APIVersion: "v1", Kind: "Pod", Name: "first", UID: "", Controller: nil, BlockOwnerDeletion: nil}}),
	)
})
