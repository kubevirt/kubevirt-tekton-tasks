package vm

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	templatev1 "github.com/openshift/api/template/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"

	lab "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants/labels"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
)

func AddMetadata(vm *kubevirtv1.VirtualMachine, template *templatev1.Template) {
	tempLabels := k8s.EnsureLabels(&vm.Spec.Template.ObjectMeta)

	if template != nil {
		labels := k8s.EnsureLabels(&vm.ObjectMeta)

		// reference origin template
		labels[lab.TemplateNameLabel] = template.GetName()
		labels[lab.TemplateNamespace] = template.GetNamespace()

		// set template flavor
		if flavorKey, flavorValue := templates.GetFlagLabelByPrefix(template, lab.TemplateFlavorLabel); flavorKey != "" {
			labels[flavorKey] = flavorValue
			tempLabels[flavorKey] = flavorValue
		}

		// set template workload
		if workloadKey, workloadValue := templates.GetFlagLabelByPrefix(template, lab.TemplateWorkloadLabel); workloadKey != "" {
			labels[workloadKey] = workloadValue
			tempLabels[workloadKey] = workloadValue
		}

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

func AddVolumes(vm *kubevirtv1.VirtualMachine, templateValidations *validations.TemplateValidations, cliParams *parse.CLIOptions) {
	if templateValidations == nil {
		templateValidations = validations.NewTemplateValidations(nil)
	}
	defaultBus := templateValidations.GetDefaultDiskBus()
	for _, diskName := range cliParams.GetAllDiskNames() {
		disk := kubevirtv1.Disk{
			Name: diskName,
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{Bus: defaultBus},
			},
		}

		vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, disk)
	}

	for _, pvcName := range cliParams.GetAllPVCNames() {
		volume := kubevirtv1.Volume{
			Name: pvcName,
			VolumeSource: kubevirtv1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: pvcName},
			},
		}

		vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volume)
	}

	for _, dvName := range cliParams.GetAllDVNames() {
		volume := kubevirtv1.Volume{
			Name: dvName,
			VolumeSource: kubevirtv1.VolumeSource{
				DataVolume: &kubevirtv1.DataVolumeSource{Name: dvName},
			},
		}

		vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volume)
	}
}

func AsVMOwnerReference(vm *kubevirtv1.VirtualMachine) metav1.OwnerReference {
	blockOwnerDeletion := true
	isController := false
	return metav1.OwnerReference{
		APIVersion:         vm.GroupVersionKind().GroupVersion().String(),
		Kind:               vm.GetObjectKind().GroupVersionKind().Kind,
		Name:               vm.GetName(),
		UID:                vm.GetUID(),
		BlockOwnerDeletion: &blockOwnerDeletion,
		Controller:         &isController,
	}
}
