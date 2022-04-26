package zutils_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("IsTrue", func() {
		DescribeTable("returns false", func(tested string) {
			Expect(zutils.IsTrue(tested)).To(BeFalse())
		},
			Entry("False", "false"),
			Entry("Bad", "falzee"),
			Entry("UpperCase", "FALSE"),
			Entry("Partially UpperCase", "FAlsE"),
		)

		DescribeTable("returns true", func(tested string) {
			Expect(zutils.IsTrue(tested)).To(BeTrue())
		},
			Entry("True", "true"),
			Entry("UpperCase", "TRUE"),
			Entry("Partially UpperCase", "True"),
		)
	})

	Describe("DecodeVM", func() {
		It("DecodeVM", func() {
			vm, vmIndex, err := zutils.DecodeVM(template.NewFedoraServerTinyTemplate().Build())
			Expect(err).Should(Succeed())
			Expect(vm.Kind).To(Equal("VirtualMachine"))
			Expect(vmIndex).To(Equal(0))
			Expect(vm.Name).To(Equal("${NAME}"))
			Expect(vm.Spec.Template.Spec.Domain.Devices.Interfaces[0].Name).To(Equal("default"))
		})

		It("DecodeVM fails", func() {
			template := template.NewFedoraServerTinyTemplate().Build()
			template.Objects = nil
			vm, _, err := zutils.DecodeVM(template)
			Expect(err).Should(HaveOccurred())
			Expect(vm).To(BeNil())
		})
	})
})
