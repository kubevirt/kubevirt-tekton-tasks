package requirements_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	req "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/requirements"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utilstest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	v1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Reruirements", func() {
	Context("GetLabelRequirement", func() {
		DescribeTable("works correctly", func(selector string, expectedResult labels.Requirements) {
			reqs, err := req.GetLabelRequirement(selector)
			Expect(err).Should(Succeed())
			Expect(reqs).To(Equal(expectedResult))
		},
			Entry("empty", "  ", nil),
			Entry("basic", testSelector, labels.Requirements{
				utilstest.GetRequirement("invalid.path", selection.NotIn, []string{"1", "2", "3"}),
				utilstest.GetRequirement("metadata", selection.Exists, []string{}),
				utilstest.GetRequirement("metadata.name", selection.In, []string{"fedora", "ubuntu"}),
				utilstest.GetRequirement("spec.runStrategy", selection.NotEquals, []string{"true"}),
			}),
			Entry("with spaces", "  metadata.name   ", labels.Requirements{
				utilstest.GetRequirement("metadata.name", selection.Exists, []string{}),
			}),
		)

		DescribeTable("fails", func(selector string, expectedError string) {
			reqs, err := req.GetLabelRequirement(selector)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring(expectedError))
			Expect(reqs).Should(BeNil())
		},
			Entry("invalid jsonpath", "test.....test", "invalid condition: cannot parse jsonpath"),
			Entry("invalid condition", "invalid#$%^$&", "could not parse condition"),
		)
	})

	Context("MatchesRequirements", func() {
		DescribeTable("works correctly", func(vm *v1.VirtualMachine, requirements labels.Requirements, expectedResult bool) {
			Expect(req.MatchesRequirements(vm, requirements)).Should(Equal(expectedResult))
		},
			Entry("nil vm and no requirements", nil, nil, true),
			Entry("vm and no requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), nil, true),
			Entry("vm and empty requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{}, true),
			Entry("matches requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{
				utilstest.GetRequirement("metadata.name", selection.In, []string{"fedora", "ubuntu"}),
				utilstest.GetRequirement("spec.runStrategy", selection.NotEquals, []string{"true"}),
			}, true),
			Entry("does not match requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{
				utilstest.GetRequirement("metadata.name", selection.In, []string{"ubuntu", "arch"}),
			}, false),
			Entry("does not match multiple requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{
				utilstest.GetRequirement("invalid.path", selection.In, []string{"1", "2", "3"}),
				utilstest.GetRequirement("metadata", selection.Exists, []string{}),
				utilstest.GetRequirement("metadata.name", selection.In, []string{"ubuntu", "arch"}),
				utilstest.GetRequirement("spec.runStrategy", selection.NotEquals, []string{"true"}),
			}, false),
		)
	})
})
