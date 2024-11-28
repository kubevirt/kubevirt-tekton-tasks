package vm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kubevirtv1 "kubevirt.io/api/core/v1"

	vm2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vm"
	shtestobjects "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
)

var _ = Describe("VM", func() {
	var vm *kubevirtv1.VirtualMachine

	BeforeEach(func() {
		vm = shtestobjects.NewTestVM().Build()
	})

	It("Adds correct default metadata", func() {
		vm2.AddMetadata(vm)

		Expect(vm.Spec.Template.ObjectMeta.Labels).To(Equal(map[string]string{
			"vm.kubevirt.io/name": vm.Name,
			"name":                vm.Name,
		}))

	})
})
