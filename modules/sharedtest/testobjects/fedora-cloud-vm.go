package testobjects

import v1 "kubevirt.io/api/core/v1"

func NewTestFedoraCloudVM(name string) *TestVM {
	cloudConfig := &CloudConfig{
		Password: "fedora",
	}
	containerDiskName, cloudinitDiskName := "containerdisk", "cloudinitdisk"

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
		{
			Name: cloudinitDiskName,
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
					Image: "kubevirt/fedora-cloud-container-disk-demo:latest",
				},
			},
		},
		{
			Name: cloudinitDiskName,
			VolumeSource: v1.VolumeSource{
				CloudInitNoCloud: &v1.CloudInitNoCloudSource{
					UserData: cloudConfig.ToString(),
				},
			},
		},
	}
	return (&TestVM{
		Data: newRandomVirtualMachine(vmi, false),
	}).WithMemory("1Gi")
}
