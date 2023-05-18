package vm

import (
	lab "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants/labels"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	templatev1 "github.com/openshift/api/template/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

func AddMetadata(vm *kubevirtv1.VirtualMachine, template *templatev1.Template) {
	tempLabels := k8s.EnsureLabels(&vm.Spec.Template.ObjectMeta)

	if template != nil {
		labels := k8s.EnsureLabels(&vm.ObjectMeta)

		// reference origin template
		labels[lab.TemplateNameLabel] = template.GetName()
		labels[lab.TemplateNamespace] = template.GetNamespace()

		if osID, osName := templates.GetOs(template); osID != "" {
			osIDLabel := lab.TemplateOsLabel + "/" + osID
			labels[osIDLabel] = zconstants.True
			tempLabels[osIDLabel] = zconstants.True
			if osName != "" {
				osIdAnnotation := lab.TemplateNameOsAnnotation + "/" + osID
				k8s.EnsureAnnotations(&vm.ObjectMeta)[osIdAnnotation] = osName
			}
		}
	}

	// for pairing service-vm (like for RDP)
	if vmName := vm.GetName(); vmName != "" {
		tempLabels[lab.VMNameLabel] = vmName
	}
}
