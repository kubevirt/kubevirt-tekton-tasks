package constants

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
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
	DataVolumeTemplatesName       = "datavolumeTemplates"
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

	MockDisks               = []string{"{\"name\": \"test\", \"cdrom\": {\"bus\": \"sata\"}}", "{\"name\": \"containerdisk\", \"disk\": {\"bus\": \"sata\"}, \"bootOrder\": 2}"}
	MockVolumes             = []string{"{\"name\": \"containerdisk\", \"containerDisk\": {\"image\": \"URL\"}}", "{\"name\": \"cloudinitdisk\"}", "{\"name\": \"test3\"}"}
	MockDataVolumeTemplates = []string{"{\"apiVersion\": \"cdi.kubevirt.io/v1beta1\", \"kind\": \"DataVolume\", \"metadata\":{\"name\": \"test1\"}, \"spec\": {\"source\": {\"http\": {\"url\": \"test.somenonexisting\"}}}}"}

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

	DataVolumeTemplates = []kubevirtv1.DataVolumeTemplateSpec{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DataVolume",
				APIVersion: "cdi.kubevirt.io/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test1",
			},
			Spec: cdiv1.DataVolumeSpec{
				Source: &cdiv1.DataVolumeSource{
					HTTP: &cdiv1.DataVolumeSourceHTTP{
						URL: "test.somenonexisting",
					},
				},
			},
		},
	}
)
