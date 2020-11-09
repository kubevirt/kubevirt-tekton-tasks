package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest/testconstants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest/testobjects"
)

var _ = Describe("Template", func() {

	It("GetFlagLabelByPrefix", func() {
		key, value := templates.GetFlagLabelByPrefix(testobjects.NewFedoraServerTinyTemplate(), "workload.template.kubevirt.io")
		Expect(key).To(Equal("workload.template.kubevirt.io/server"))
		Expect(value).To(Equal("true"))
	})

	It("DecodeVM", func() {
		vm, err := templates.DecodeVM(testobjects.NewFedoraServerTinyTemplate())
		Expect(err).Should(Succeed())
		Expect(vm.Kind).To(Equal("VirtualMachine"))
		Expect(vm.Name).To(Equal("${NAME}"))
		Expect(vm.Spec.Template.Spec.Domain.Devices.Interfaces[0].Name).To(Equal("default"))
	})

	It("DecodeVM fails", func() {
		template := testobjects.NewFedoraServerTinyTemplate()
		template.Objects = nil
		vm, err := templates.DecodeVM(template)
		Expect(err).Should(HaveOccurred())
		Expect(vm).To(BeNil())
	})

	It("GetTemplateValidations", func() {
		validations, err := templates.GetTemplateValidations(testobjects.NewFedoraServerTinyTemplate())
		Expect(err).Should(Succeed())
		Expect(validations.IsEmpty()).To(BeFalse())
		Expect(validations.GetDefaultDiskBus()).To(Equal(testconstants.Virtio))
	})

	It("GetOs", func() {
		osID, osName := templates.GetOs(testobjects.NewFedoraServerTinyTemplate())
		Expect(osID).To(Equal("fedora29"))
		Expect(osName).To(Equal("Fedora 27 or higher"))
	})
})
