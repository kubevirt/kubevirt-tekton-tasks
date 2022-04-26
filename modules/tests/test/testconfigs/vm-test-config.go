package testconfigs

import (
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	template2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/onsi/ginkgo/v2"
	v1 "github.com/openshift/api/template/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type CreateVMTaskData struct {
	CreateMode CreateVMMode

	Template                *v1.Template
	TemplateTargetNamespace TargetNamespace

	VM                        *kubevirtv1.VirtualMachine
	VMTargetNamespace         TargetNamespace
	VMManifestTargetNamespace TargetNamespace

	DataVolumesToCreate                      []*datavolume.TestDataVolume
	DataSourcesToCreate                      []*datavolume.TestDataVolume
	IsCommonTemplate                         bool
	UseDefaultTemplateNamespacesInTaskParams bool
	UseDefaultVMNamespacesInTaskParams       bool
	StartVM                                  string
	ExpectedAdditionalDiskBus                string

	// Params
	// these two are set if Template is not nil
	TemplateName      string
	TemplateNamespace string

	// this is set if VM is not nil
	VMManifest string

	TemplateParams            []string
	VMNamespace               string
	DataVolumes               []string
	DataSources               []string
	OwnDataVolumes            []string
	PersistentVolumeClaims    []string
	OwnPersistentVolumeClaims []string
}

func (c *CreateVMTaskData) GetTemplateParam(key string) string {
	for _, param := range c.TemplateParams {
		fragments := strings.SplitN(param, ":", 2)
		if len(fragments) == 2 && fragments[0] == key {
			return fragments[1]
		}
	}
	return ""
}

func (c *CreateVMTaskData) GetExpectedVMStubMeta() *kubevirtv1.VirtualMachine {
	var finalDisks []kubevirtv1.Disk
	var finalVolumes []kubevirtv1.Volume
	var vmName, vmNamespace string

	var vm *kubevirtv1.VirtualMachine

	switch c.CreateMode {
	case CreateVMVMManifestMode:
		if err := yaml.Unmarshal([]byte(c.VMManifest), &vm); err != nil || vm == nil {
			vm = nil
		} else {
			if c.VMNamespace != "" {
				vm.Namespace = c.VMNamespace
			}
			vmName = vm.Name
			vmNamespace = vm.Namespace
		}
	case CreateVMTemplateMode:
		if c.Template != nil && c.Template.Objects != nil {
			vm = template2.GetVM(c.Template)
		}

		vmName = c.GetTemplateParam(template2.NameParam)
		vmNamespace = c.VMNamespace
	}
	if vm != nil {
		finalDisks = append(finalDisks, vm.Spec.Template.Spec.Domain.Devices.Disks...)
		finalVolumes = append(finalVolumes, vm.Spec.Template.Spec.Volumes...)
	}

	findDisk := func(name string) *kubevirtv1.Disk {
		for i := 0; i < len(finalDisks); i++ {
			if finalDisks[i].Name == name {
				return &finalDisks[i]
			}
		}
		return nil
	}

	findVolume := func(name string) *kubevirtv1.Volume {
		for i := 0; i < len(finalVolumes); i++ {
			if finalVolumes[i].Name == name {
				return &finalVolumes[i]
			}
		}
		return nil
	}

	for _, dataVolume := range c.DataVolumesToCreate {
		name := dataVolume.Data.Name

		if dataVolume.DiskName == "" || findDisk(dataVolume.DiskName) == nil {
			disk := kubevirtv1.Disk{
				Name: name,
				DiskDevice: kubevirtv1.DiskDevice{
					Disk: &kubevirtv1.DiskTarget{Bus: c.ExpectedAdditionalDiskBus},
				},
			}
			if dataVolume.DiskName != "" {
				disk.Name = dataVolume.DiskName
			}

			finalDisks = append(finalDisks, disk)
		}

		if originalVolume := findVolume(dataVolume.DiskName); dataVolume.DiskName != "" && originalVolume != nil {
			switch dataVolume.AttachmentType {
			case datavolume.DV, datavolume.OwnedDV:
				originalVolume.VolumeSource = kubevirtv1.VolumeSource{
					DataVolume: &kubevirtv1.DataVolumeSource{Name: name},
				}
			case datavolume.PVC, datavolume.OwnedPVC:
				originalVolume.VolumeSource = kubevirtv1.VolumeSource{
					PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
						PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: name,
						},
					},
				}
			}
		} else {
			volume := kubevirtv1.Volume{
				Name: name,
			}

			if dataVolume.DiskName != "" {
				volume.Name = dataVolume.DiskName
			}

			switch dataVolume.AttachmentType {
			case datavolume.DV, datavolume.OwnedDV:
				volume.DataVolume = &kubevirtv1.DataVolumeSource{
					Name: name,
				}
			case datavolume.PVC, datavolume.OwnedPVC:
				volume.PersistentVolumeClaim = &kubevirtv1.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: name,
					},
				}
			}

			finalVolumes = append(finalVolumes, volume)
		}
	}

	return &kubevirtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vmName,
			Namespace: vmNamespace,
		},
		Spec: kubevirtv1.VirtualMachineSpec{
			Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{
				Spec: kubevirtv1.VirtualMachineInstanceSpec{
					Volumes: finalVolumes,
					Domain: kubevirtv1.DomainSpec{
						Devices: kubevirtv1.Devices{
							Disks: finalDisks,
						},
					},
				},
			},
		},
	}
}

func (c *CreateVMTaskData) SetDVorPVC(name string, attachmentType datavolume.TestDataVolumeAttachmentType) {
	switch attachmentType {
	case datavolume.DV:
		c.DataVolumes = append(c.DataVolumes, name)
	case datavolume.OwnedDV:
		c.OwnDataVolumes = append(c.OwnDataVolumes, name)
	case datavolume.PVC:
		c.PersistentVolumeClaims = append(c.PersistentVolumeClaims, name)
	case datavolume.OwnedPVC:
		c.OwnPersistentVolumeClaims = append(c.OwnPersistentVolumeClaims, name)
	}
}

type CreateVMTestConfig struct {
	TaskRunTestConfig
	TaskData CreateVMTaskData

	deploymentNamespace string
}

func (c *CreateVMTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace

	switch c.TaskData.CreateMode {
	case CreateVMVMManifestMode:
		c.initCreateVMManifest(options)
	case CreateVMTemplateMode:
		c.initCreateVMTemplate(options)
	default:
		panic("unknown VM create mode")
	}

	for _, dataVolume := range c.TaskData.DataVolumesToCreate {
		dataVolume.Data.Name = E2ETestsRandomName(dataVolume.Data.Name)
		dataVolume.Data.Namespace = c.TaskData.VMNamespace
		if options.StorageClass != "" {
			dataVolume.Data.Spec.PVC.StorageClassName = &options.StorageClass
		}
	}

	if c.TaskData.ExpectedAdditionalDiskBus == "" {
		c.TaskData.ExpectedAdditionalDiskBus = "virtio"
	}
}

func (c *CreateVMTestConfig) initCreateVMManifest(options *testoptions.TestOptions) {
	if c.TaskData.VMTargetNamespace != "" && c.TaskData.VMManifestTargetNamespace != "" {
		ginkgo.Fail("only one of VMTargetNamespace|VMManifestTargetNamespace can be set")
	}

	if vm := c.TaskData.VM; vm != nil {
		if vm.Name != "" {
			vm.Name = E2ETestsRandomName(vm.Name)
			vm.Spec.Template.ObjectMeta.Name = vm.Name
		}

		vm.Spec.Template.ObjectMeta.Namespace = ""

		if c.TaskData.VMManifestTargetNamespace != "" {
			vm.Namespace = options.ResolveNamespace(c.TaskData.VMManifestTargetNamespace, "")
			c.TaskData.VMNamespace = ""
		} else {
			vm.Namespace = ""
			c.TaskData.VMNamespace = options.ResolveNamespace(c.TaskData.VMTargetNamespace, c.TaskData.VMNamespace)
		}

		c.TaskData.VMManifest = (&testobjects.TestVM{Data: vm}).ToString()
	} else {
		// just for invalid YAMLs - otherwise use TaskData.VM
		if c.TaskData.VMManifestTargetNamespace != "" {
			ginkgo.Fail("VMManifestTargetNamespace cannot be set for manifest")
		}
		if c.TaskData.VMTargetNamespace != "" {
			ginkgo.Fail("VMTargetNamespace cannot be set for manifest")
		}
	}
}

func (c *CreateVMTestConfig) initCreateVMTemplate(options *testoptions.TestOptions) {
	c.TaskData.VMNamespace = options.ResolveNamespace(c.TaskData.VMTargetNamespace, c.TaskData.VMNamespace)
	c.TaskData.TemplateNamespace = options.ResolveNamespace(c.TaskData.TemplateTargetNamespace, c.TaskData.TemplateNamespace)

	if tmpl := c.TaskData.Template; tmpl != nil {
		if tmpl.Name != "" {
			tmpl.Name = E2ETestsRandomName(tmpl.Name)
		}
		tmpl.Namespace = c.TaskData.TemplateNamespace

		c.TaskData.TemplateName = tmpl.Name
	} else {
		if c.TaskData.TemplateName != "" && c.TaskData.IsCommonTemplate {
			c.TaskData.TemplateName += options.CommonTemplatesVersion
		}
	}
}

func (c *CreateVMTestConfig) GetTaskRun() *v1beta1.TaskRun {
	var taskName, taskRunName string

	params := []v1beta1.Param{
		{
			Name: CreateVMParams.DataVolumes,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: c.TaskData.DataVolumes,
			},
		},
		{
			Name: CreateVMParams.OwnDataVolumes,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: c.TaskData.OwnDataVolumes,
			},
		},
		{
			Name: CreateVMParams.PersistentVolumeClaims,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: c.TaskData.PersistentVolumeClaims,
			},
		},
		{
			Name: CreateVMParams.OwnPersistentVolumeClaims,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: c.TaskData.OwnPersistentVolumeClaims,
			},
		},
		{
			Name: CreateVMFromTemplateParams.StartVM,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.StartVM,
			},
		},
	}
	var vmNamespace string
	if !c.TaskData.UseDefaultVMNamespacesInTaskParams {
		vmNamespace = c.TaskData.VMNamespace
	}

	switch c.TaskData.CreateMode {
	case CreateVMVMManifestMode:
		taskName = CreateVMFromManifestClusterTaskName
		taskRunName = "taskrun-vm-create-from-manifest"

		params = append(params,
			v1beta1.Param{
				Name: CreateVMFromManifestParams.Manifest,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: c.TaskData.VMManifest,
				},
			},
			v1beta1.Param{
				Name: CreateVMFromManifestParams.Namespace,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: vmNamespace,
				},
			},
		)
	case CreateVMTemplateMode:
		taskName = CreateVMFromTemplateClusterTaskName
		taskRunName = "taskrun-vm-create-from-template"

		var templateNamespace string

		if !c.TaskData.UseDefaultTemplateNamespacesInTaskParams {
			templateNamespace = c.TaskData.TemplateNamespace
		}

		params = append(params,
			v1beta1.Param{
				Name: CreateVMFromTemplateParams.TemplateName,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: c.TaskData.TemplateName,
				},
			},
			v1beta1.Param{
				Name: CreateVMFromTemplateParams.TemplateNamespace,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: templateNamespace,
				},
			},
			v1beta1.Param{
				Name: CreateVMFromTemplateParams.TemplateParams,
				Value: v1beta1.ArrayOrString{
					Type:     v1beta1.ParamTypeArray,
					ArrayVal: c.TaskData.TemplateParams,
				},
			},
			v1beta1.Param{
				Name: CreateVMFromTemplateParams.VmNamespace,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: vmNamespace,
				},
			},
		)
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(taskRunName),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: taskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: c.ServiceAccount,
			Params:             params,
		},
	}
}
