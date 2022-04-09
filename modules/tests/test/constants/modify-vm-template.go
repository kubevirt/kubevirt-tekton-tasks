package constants

import (
	kubevirtv1 "kubevirt.io/api/core/v1"
)

const (
	ModifyTemplateClusterTaskName    = "modify-vm-template"
	ModifyTemplateServiceAccountName = "modify-vm-template-task"
	ModifyTemplateTaskRunName        = "taskrun-modify-vm-template"

	TemplateNameOptionName        = "templateName"
	TemplateNamespaceOptionName   = "templateNamespace"
	CPUCoresOptionName            = "cpuCores"
	CPUSocketsOptionName          = "cpuSockets"
	CPUThreadsOptionName          = "cpuThreads"
	MemoryOptionName              = "memory"
	TemplateLabelsOptionName      = "templateLabels"
	TemplateAnnotationsOptionName = "templateAnnotations"
	VMLabelsOptionName            = "vmLabels"
	VMAnnotationsOptionName       = "vmAnnotations"
	DisksOptionName               = "disks"
	VolumesOptionName             = "volumes"
	DeleteDatavolumeTemplateName  = "deleteDatavolumeTemplate"

	CPUSocketsTopologyNumber    uint32 = 1
	CPUCoresTopologyNumber      uint32 = 2
	CPUThreadsTopologyNumber    uint32 = 3
	CPUSocketsTopologyNumberStr        = "1"
	CPUCoresTopologyNumberStr          = "2"
	CPUThreadsTopologyNumberStr        = "3"

	MemoryValue = "180M"
)

var (
	bootOrder     uint = 2
	MockArray          = []string{"newKey: value", "test: true"}
	WrongStrSlice      = []string{"wrong vaue"}

	MockDisks   = []string{"{\"name\": \"test\", \"cdrom\": {\"bus\": \"sata\"}}", "{\"name\": \"containerdisk\", \"disk\": {\"bus\": \"sata\"}, \"bootOrder\": 2}"}
	MockVolumes = []string{"{\"name\": \"containerdisk\", \"containerDisk\": {\"image\": \"URL\"}}", "{\"name\": \"cloudinitdisk\"}", "{\"name\": \"test3\"}"}

	LabelsAnnotationsMap = map[string]string{
		"newKey": "value",
		"test":   "true",
	}
	Disks = []kubevirtv1.Disk{
		{
			Name: "containerdisk",
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{
					Bus: "sata",
				},
			},
			BootOrder: &bootOrder,
		}, {
			Name: "cloudinitdisk",
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{
					Bus: "virtio",
				},
			},
		}, {
			Name: "test",
			DiskDevice: kubevirtv1.DiskDevice{
				CDRom: &kubevirtv1.CDRomTarget{
					Bus: "sata",
				},
			},
		},
	}
	Volumes = []kubevirtv1.Volume{
		{
			Name: "containerdisk",
			VolumeSource: kubevirtv1.VolumeSource{
				ContainerDisk: &kubevirtv1.ContainerDiskSource{
					Image: "URL",
				},
			},
		},
		{
			Name: "cloudinitdisk",
		},
		{
			Name: "test3",
		},
	}
)
