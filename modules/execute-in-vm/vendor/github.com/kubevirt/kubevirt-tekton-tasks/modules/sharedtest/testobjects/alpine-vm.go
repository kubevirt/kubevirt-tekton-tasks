package testobjects

import (
	"kubevirt.io/client-go/api/v1"
)

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
					Image: "kubevirt/alpine-container-disk-demo:latest",
				},
			},
		},
	}
	return (&TestVM{
		Data: newRandomVirtualMachine(vmi, false),
	}).WithMemory("64Mi")
}
