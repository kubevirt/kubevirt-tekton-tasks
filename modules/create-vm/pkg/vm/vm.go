package vm

import (
	lab "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants/labels"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

func AddMetadata(vm *kubevirtv1.VirtualMachine) {
	var tempLabels map[string]string
	if vm.Spec.Template == nil {
		tempLabels = k8s.EnsureLabels(&vm.ObjectMeta)
	} else {
		tempLabels = k8s.EnsureLabels(&vm.Spec.Template.ObjectMeta)
	}

	// for pairing service-vm (like for RDP)
	if vmName := vm.GetName(); vmName != "" {
		tempLabels[lab.VMNameLabel] = vmName
	}
}
