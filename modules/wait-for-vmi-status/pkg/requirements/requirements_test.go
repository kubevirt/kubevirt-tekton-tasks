package requirements_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	req "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/requirements"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utilstest"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	v1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Reruirements", func() {
	Context("GetLabelRequirement", func() {
		table.DescribeTable("works correctly", func(selector string, expectedResult labels.Requirements) {
			reqs, err := req.GetLabelRequirement(selector)
			Expect(err).Should(Succeed())
			Expect(reqs).To(Equal(expectedResult))
		},
			table.Entry("empty", "  ", nil),
			table.Entry("basic", testSelector, labels.Requirements{
				utilstest.GetRequirement("invalid.path", selection.NotIn, []string{"1", "2", "3"}),
				utilstest.GetRequirement("metadata", selection.Exists, []string{}),
				utilstest.GetRequirement("metadata.name", selection.In, []string{"fedora", "ubuntu"}),
				utilstest.GetRequirement("spec.running", selection.NotEquals, []string{"true"}),
			}),
			table.Entry("with spaces", "  metadata.name   ", labels.Requirements{
				utilstest.GetRequirement("metadata.name", selection.Exists, []string{}),
			}),
		)

		table.DescribeTable("fails", func(selector string, expectedError string) {
			reqs, err := req.GetLabelRequirement(selector)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring(expectedError))
			Expect(reqs).Should(BeNil())
		},
			table.Entry("invalid jsonpath", "test.....test", "invalid condition: cannot parse jsonpath"),
			table.Entry("invalid condition", "invalid#$%^$&", "could not parse condition"),
		)
	})

	Context("MatchesRequirements", func() {
		table.DescribeTable("works correctly", func(vm *v1.VirtualMachine, requirements labels.Requirements, expectedResult bool) {
			Expect(req.MatchesRequirements(vm, requirements)).Should(Equal(expectedResult))
		},
			table.Entry("nil vm and no requirements", nil, nil, true),
			table.Entry("vm and no requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), nil, true),
			table.Entry("vm and empty requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{}, true),
			table.Entry("matches requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{
				utilstest.GetRequirement("metadata.name", selection.In, []string{"fedora", "ubuntu"}),
				utilstest.GetRequirement("spec.running", selection.NotEquals, []string{"true"}),
			}, true),
			table.Entry("does not match requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{
				utilstest.GetRequirement("metadata.name", selection.In, []string{"ubuntu", "arch"}),
			}, false),
			table.Entry("does not match multiple requirements", testobjects.NewTestFedoraCloudVM("fedora").Build(), labels.Requirements{
				utilstest.GetRequirement("invalid.path", selection.In, []string{"1", "2", "3"}),
				utilstest.GetRequirement("metadata", selection.Exists, []string{}),
				utilstest.GetRequirement("metadata.name", selection.In, []string{"ubuntu", "arch"}),
				utilstest.GetRequirement("spec.running", selection.NotEquals, []string{"true"}),
			}, false),
		)
	})
})
