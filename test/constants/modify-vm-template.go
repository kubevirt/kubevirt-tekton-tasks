package constants

import (
	templatev1 "github.com/openshift/api/template/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	ModifyTemplateTaskName    = "modify-vm-template"
	ModifyTemplateTaskRunName = "taskrun-modify-vm-template"

	TemplateNameOptionName             = "templateName"
	TemplateNamespaceOptionName        = "templateNamespace"
	CPUCoresOptionName                 = "cpuCores"
	CPUSocketsOptionName               = "cpuSockets"
	CPUThreadsOptionName               = "cpuThreads"
	MemoryOptionName                   = "memory"
	TemplateLabelsOptionName           = "templateLabels"
	TemplateAnnotationsOptionName      = "templateAnnotations"
	VMLabelsOptionName                 = "vmLabels"
	VMAnnotationsOptionName            = "vmAnnotations"
	DisksOptionName                    = "disks"
	VolumesOptionName                  = "volumes"
	DataVolumeTemplatesOptionName      = "datavolumeTemplates"
	DeleteDatavolumeTemplateOptionName = "deleteDatavolumeTemplate"
	DeleteDisksOptionName              = "deleteDisks"
	DeleteVolumesOptionName            = "deleteVolumes"
	DeleteTemplateParametersOptionName = "deleteTemplateParameters"
	DeleteTemplateOptionName           = "deleteTemplate"
	TemplateParametersOptionName       = "templateParameters"

	CPUSocketsTopologyNumber    uint32 = 1
	CPUCoresTopologyNumber      uint32 = 2
	CPUThreadsTopologyNumber    uint32 = 3
	CPUSocketsTopologyNumberStr        = "1"
	CPUCoresTopologyNumberStr          = "2"
	CPUThreadsTopologyNumberStr        = "3"

	MemoryValue = "180M"
)

var (
	bootOrder               uint = 2
	MockTemplateAnnotations      = []string{"newKey: TemplateAnnotation", "testTemplateAnnotation: true"}
	MockTemplateLabels           = []string{"newKey: TemplateLabel", "testTemplateLabel: true"}
	MockVMAnnotations            = []string{"newKey: VMAnnotation", "testVMAnnotation: true"}
	MockVMLabels                 = []string{"newKey: VMLabel", "testVMLabel: true"}
	WrongStrSlice                = []string{"wrong value"}
	MockDisk                     = []string{"{\"name\": \"test\", \"cdrom\": {\"bus\": \"sata\"}}"}
	MockDisks                    = []string{"{\"name\": \"test\", \"cdrom\": {\"bus\": \"sata\"}}", "{\"name\": \"containerdisk\", \"disk\": {\"bus\": \"sata\"}, \"bootOrder\": 2}"}
	MockVolume                   = []string{"{\"name\": \"test3\"}"}
	MockVolumes                  = []string{"{\"name\": \"containerdisk\", \"containerDisk\": {\"image\": \"URL\"}}", "{\"name\": \"cloudinitdisk\"}", "{\"name\": \"test3\"}"}
	MockDataVolumeTemplates      = []string{"{\"apiVersion\": \"cdi.kubevirt.io/v1beta1\", \"kind\": \"DataVolume\", \"metadata\":{\"name\": \"test1\"}, \"spec\": {\"source\": {\"http\": {\"url\": \"test.somenonexisting\"}}}}"}
	MockTemplateParameter        = []string{"{\"name\": \"test\", \"value\": \"test\"}"}

	TemplateAnnotationsMap = map[string]string{
		"newKey":                 "TemplateAnnotation",
		"testTemplateAnnotation": "true",
	}
	TemplateLabelsMap = map[string]string{
		"newKey":            "TemplateLabel",
		"testTemplateLabel": "true",
	}
	VMAnnotationsMap = map[string]string{
		"newKey":           "VMAnnotation",
		"testVMAnnotation": "true",
	}
	VMLabelsMap = map[string]string{
		"newKey":      "VMLabel",
		"testVMLabel": "true",
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
	Disk = kubevirtv1.Disk{
		Name: "test",
		DiskDevice: kubevirtv1.DiskDevice{
			CDRom: &kubevirtv1.CDRomTarget{
				Bus: "sata",
			},
		},
	}
	Volume = kubevirtv1.Volume{
		Name: "test3",
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

	TemplateParameters = []templatev1.Parameter{
		{
			Name:  "test",
			Value: "test",
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
