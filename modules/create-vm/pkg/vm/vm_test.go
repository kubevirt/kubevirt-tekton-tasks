package vm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kubevirtv1 "kubevirt.io/api/core/v1"

	vm2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vm"
	shtestobjects "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	template "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
)

var _ = Describe("VM", func() {
	var vm *kubevirtv1.VirtualMachine

	BeforeEach(func() {
		vm = shtestobjects.NewTestVM().Build()
	})

	It("Adds correct metadata from template", func() {
		vm2.AddMetadata(vm, template.NewFedoraServerTinyTemplate().Build())

		Expect(vm.Labels).To(Equal(map[string]string{
			"vm.kubevirt.io/template":           "fedora-server-tiny-v0.7.0",
			"vm.kubevirt.io/template.namespace": "openshift",
			"os.template.kubevirt.io/fedora29":  "true",
		}))

		Expect(vm.Annotations).To(Equal(map[string]string{
			"name.os.template.kubevirt.io/fedora29": "Fedora 27 or higher",
		}))

		Expect(vm.Spec.Template.ObjectMeta.Labels).To(Equal(map[string]string{
			"vm.kubevirt.io/name":              vm.Name,
			"name":                             vm.Name,
			"os.template.kubevirt.io/fedora29": "true",
		}))

	})

	It("Adds correct default metadata", func() {
		vm2.AddMetadata(vm, nil)

		Expect(vm.Spec.Template.ObjectMeta.Labels).To(Equal(map[string]string{
			"vm.kubevirt.io/name": vm.Name,
			"name":                vm.Name,
		}))

	})
})
