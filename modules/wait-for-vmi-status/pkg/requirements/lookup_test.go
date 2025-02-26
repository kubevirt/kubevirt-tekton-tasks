package requirements_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	req "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/requirements"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/labels"
	v1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Lookup", func() {
	DescribeTable("lookup works correctly", func(vm *v1.VirtualMachine, selector string, expectedResult labels.Set) {
		reqs, err := labels.Parse(selector)
		Expect(err).Should(Succeed())

		requirements, selectable := reqs.Requirements()
		Expect(selectable).To(BeTrue())

		result, err := req.ObjectToLabelsLookup(vm, requirements)
		Expect(err).Should(Succeed())

		Expect(result).To(Equal(expectedResult))
	},
		Entry("nil vm", nil, testSelector, labels.Set{}),
		Entry("empty requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), "", labels.Set{}),
		Entry("basic", testobjects.NewTestFedoraCloudVM("fedora").Build(), testSelector, labels.Set{
			"metadata.name":    "fedora",
			"spec.runStrategy": "Halted",
			"metadata":         "{\"name\":\"fedora\",\"namespace\":\"default\",\"creationTimestamp\":null}",
		}),
		Entry("with spaces", testobjects.NewTestFedoraCloudVM("fedora").Build(), "  metadata.name   ", labels.Set{
			"metadata.name": "fedora",
		}),
	)

	It("lookup fails on invalid path correctly", func() {
		vm := testobjects.NewTestFedoraCloudVM("fedora").Build()
		reqs, err := labels.Parse("test.....test")
		Expect(err).Should(Succeed())

		requirements, selectable := reqs.Requirements()
		Expect(selectable).To(BeTrue())

		result, err := req.ObjectToLabelsLookup(vm, requirements)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("cannot parse jsonpath"))
		Expect(result).To(BeNil())
	})

})
