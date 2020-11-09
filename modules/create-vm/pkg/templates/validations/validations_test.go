package validations_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest/testconstants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest/testobjects"
)

var _ = Describe("Validations", func() {
	table.DescribeTable("gets default bus", func(templateValidations []validations.CommonTemplateValidation, expectedBus string) {
		Expect(validations.NewTemplateValidations(templateValidations).GetDefaultDiskBus()).To(Equal(expectedBus))
	},
		table.Entry("nil", nil, Virtio),
		table.Entry("empty", testobjects.NewTestCommonTemplateValidations(), Virtio),
		table.Entry("one", testobjects.NewTestCommonTemplateValidations(Scsi), Scsi),
		table.Entry("two with virtio", testobjects.NewTestCommonTemplateValidations(Scsi, Virtio), Virtio),
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
