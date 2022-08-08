package vm

import (
	lab "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants/labels"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	templatev1 "github.com/openshift/api/template/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
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

// returns transient pointer to the Disk struct in array
func getDisk(vm *kubevirtv1.VirtualMachine, name string) *kubevirtv1.Disk {
	for i := 0; i < len(vm.Spec.Template.Spec.Domain.Devices.Disks); i++ {
		if vm.Spec.Template.Spec.Domain.Devices.Disks[i].Name == name {
			return &vm.Spec.Template.Spec.Domain.Devices.Disks[i]
		}
	}

	return nil
}

// returns transient pointer to the Volume struct in array
func getVolume(vm *kubevirtv1.VirtualMachine, name string) *kubevirtv1.Volume {
	for i := 0; i < len(vm.Spec.Template.Spec.Volumes); i++ {
		if vm.Spec.Template.Spec.Volumes[i].Name == name {
			return &vm.Spec.Template.Spec.Volumes[i]
		}
	}

	return nil
}

func AddVolumes(vm *kubevirtv1.VirtualMachine, templateValidations *validations.TemplateValidations, cliParams *parse.CLIOptions) {
	if templateValidations == nil {
		templateValidations = validations.NewTemplateValidations(nil)
	}
	defaultBus := templateValidations.GetDefaultDiskBus()

	ensureDisk := func(diskName string) *kubevirtv1.Disk {
		if disk := getDisk(vm, diskName); disk != nil {
			return disk
		}
		disk := kubevirtv1.Disk{
			Name: diskName,
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{Bus: kubevirtv1.DiskBus(defaultBus)},
			},
		}

		vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, disk)

		return getDisk(vm, diskName)
	}

	ensureVolume := func(volumeName string) *kubevirtv1.Volume {
		if volume := getVolume(vm, volumeName); volume != nil {
			return volume
		}
		volume := kubevirtv1.Volume{
			Name: volumeName,
		}

		vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volume)

		return getVolume(vm, volumeName)
	}

	for volumeName, pvcName := range cliParams.GetPVCDiskNamesMap() {
		ensureDisk(volumeName)
		volume := ensureVolume(volumeName)

		if volume.PersistentVolumeClaim == nil {
			volume.VolumeSource = kubevirtv1.VolumeSource{
				PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			}
		} else {
			volume.PersistentVolumeClaim.ClaimName = pvcName
		}
	}

	for volumeName, dvName := range cliParams.GetDVDiskNamesMap() {
		ensureDisk(volumeName)
		volume := ensureVolume(volumeName)

		if volume.DataVolume == nil {
			volume.VolumeSource = kubevirtv1.VolumeSource{
				DataVolume: &kubevirtv1.DataVolumeSource{Name: dvName},
			}
		} else {
			volume.DataVolume.Name = dvName
		}
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
