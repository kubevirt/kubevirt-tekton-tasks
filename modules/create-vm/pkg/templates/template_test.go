package templates_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
)

var _ = Describe("Template", func() {

	It("GetFlagLabelByPrefix", func() {
		key, value := templates.GetFlagLabelByPrefix(template.NewFedoraServerTinyTemplate().Build(), "workload.template.kubevirt.io")
		Expect(key).To(Equal("workload.template.kubevirt.io/server"))
		Expect(value).To(Equal("true"))
	})

	It("GetTemplateValidations", func() {
		validations, err := templates.GetTemplateValidations(template.NewFedoraServerTinyTemplate().Build())
		Expect(err).Should(Succeed())
		Expect(validations.IsEmpty()).To(BeFalse())
		Expect(validations.GetDefaultDiskBus()).To(Equal(testconstants.Virtio))
	})

	It("GetOs", func() {
		osID, osName := templates.GetOs(template.NewFedoraServerTinyTemplate().Build())
		Expect(osID).To(Equal("fedora29"))
		Expect(osName).To(Equal("Fedora 27 or higher"))
	})
})
