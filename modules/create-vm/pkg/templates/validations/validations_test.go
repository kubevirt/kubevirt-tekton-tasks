package validations_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest/testobjects"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
)

var _ = Describe("Validations", func() {
	DescribeTable("gets default bus", func(templateValidations []validations.CommonTemplateValidation, expectedBus string) {
		Expect(validations.NewTemplateValidations(templateValidations).GetDefaultDiskBus()).To(Equal(expectedBus))
	},
		Entry("nil", nil, Virtio),
		Entry("empty", testobjects.NewTestCommonTemplateValidations(), Virtio),
		Entry("one", testobjects.NewTestCommonTemplateValidations(Scsi), Scsi),
		Entry("two with virtio", testobjects.NewTestCommonTemplateValidations(Scsi, Virtio), Virtio),
	)

	It("gets prefered bus", func() {
		allowed := testobjects.NewTestCommonTemplateValidations(Scsi, Sata, Virtio)
		otherAllowed := testobjects.NewTestCommonTemplateValidations(Virtio)
		prefered := testobjects.NewTestCommonTemplateValidations(Sata)
		prefered[0].JustWarning = true

		var finalTemplateValidations []validations.CommonTemplateValidation

		finalTemplateValidations = append(finalTemplateValidations, allowed...)
		finalTemplateValidations = append(finalTemplateValidations, otherAllowed...)
		finalTemplateValidations = append(finalTemplateValidations, prefered...)

		Expect(validations.NewTemplateValidations(finalTemplateValidations).GetDefaultDiskBus()).To(Equal(Sata))
	})
})
