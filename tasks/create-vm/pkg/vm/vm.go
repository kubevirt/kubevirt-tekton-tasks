package vm

import (
	templatev1 "github.com/openshift/api/template/v1"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/templates"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/parse"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

const (
	// templateOsLabel is a label that specifies the OS of the template
	templateOsLabel = "os.template.kubevirt.io"

	// templateWorkloadLabel is a label that specifies the workload of the template
	templateWorkloadLabel = "workload.template.kubevirt.io"

	// templateFlavorLabel is a label that specifies the flavor of the template
	templateFlavorLabel = "flavor.template.kubevirt.io"

	// templateNameOsAnnotation is an annotation that specifies human readable os name
	templateNameOsAnnotation = "name.os.template.kubevirt.io"

	// templateNameLabel defines a label of the template name which was used to created the VM
	templateNameLabel = "vm.kubevirt.io/template"

	// templateNamespace defines a label of the template namespace which was used to create the VM
	templateNamespace = "vm.kubevirt.io/template.namespace"

	// vmNameLabel defines a label of virtual machine name which was used to create the VM
	vmNameLabel = "vm.kubevirt.io/name"
)

func AddMetadata(vm *kubevirtv1.VirtualMachine, template *templatev1.Template) {
	labels := vm.ObjectMeta.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
		vm.ObjectMeta.SetLabels(labels)
	}

	tempLabels := vm.Spec.Template.ObjectMeta.GetLabels()

	// reference origin template
	labels[templateNameLabel] = template.GetName()
	labels[templateNamespace] = template.GetNamespace()

	// set template flavor
	if flavorKey, flavorValue := templates.GetFlagLabelByPrefix(template, templateFlavorLabel); flavorKey != "" {
		labels[flavorKey] = flavorValue
		tempLabels[flavorKey] = flavorValue
	}

	// set template workload
	if workloadKey, workloadValue := templates.GetFlagLabelByPrefix(template, templateWorkloadLabel); workloadKey != "" {
		labels[workloadKey] = workloadValue
		tempLabels[workloadKey] = workloadValue
	}

	// TODO search for correct os label and annotation from template and set here

	// for pairing service-vm (like for RDP)
	if vmName := vm.GetName(); vmName != "" {
		tempLabels[vmNameLabel] = vmName
	}
}

func AddVolumes(vm *kubevirtv1.VirtualMachine, template *templatev1.Template, cliParams *parse.CLIParams) {
	for _, diskName := range cliParams.GetAllDiskNames() {
		disk := kubevirtv1.Disk{
			Name: diskName,
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{Bus: "virtio"}, // TODO get from template validations or default to virtio
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
