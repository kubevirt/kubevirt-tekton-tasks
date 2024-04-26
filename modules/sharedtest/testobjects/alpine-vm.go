package testobjects

import v1 "kubevirt.io/api/core/v1"

func NewTestAlpineVM(name string) *TestVM {
	containerDiskName := "containerdisk"

	vmi := newRandomVMI()
	vmi.Name = name

	vmi.Spec.Domain.Devices.Disks = []v1.Disk{
		{
			Name: containerDiskName,
			DiskDevice: v1.DiskDevice{
				Disk: &v1.DiskTarget{
					Bus: "virtio",
				},
			},
		},
	}

	vmi.Spec.Volumes = []v1.Volume{
		{
			Name: containerDiskName,
			VolumeSource: v1.VolumeSource{
				ContainerDisk: &v1.ContainerDiskSource{
					Image: "quay.io/kubevirt/alpine-container-disk-demo:20240426_ca94b81c6",
				},
			},
		},
	}
	return (&TestVM{
		Data: newRandomVirtualMachine(vmi, false),
	}).WithMemory("128Mi")
}
