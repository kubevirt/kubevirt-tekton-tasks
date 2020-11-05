package testconfigs

import (
	v1 "github.com/openshift/api/template/v1"
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/dv"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/template"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	"strings"
)

type CreateVMFromTemplateTaskData struct {
	Template *v1.Template

	TemplateTargetNamespace TargetNamespace
	VMTargetNamespace       TargetNamespace

	DataVolumesToCreate                      []*dv.TestDataVolume
	PVCsAreNotDataVolumes                    bool
	IsCommonTemplate                         bool
	UseDefaultTemplateNamespacesInTaskParams bool
	UseDefaultVMNamespacesInTaskParams       bool
	ExpectedAdditionalDiskBus                string

	// Params
	// these two are set if Template is not nil
	TemplateName      string
	TemplateNamespace string

	TemplateParams            []string
	VMNamespace               string
	DataVolumes               []string
	OwnDataVolumes            []string
	PersistentVolumeClaims    []string
	OwnPersistentVolumeClaims []string
}

func (c *CreateVMFromTemplateTaskData) GetTemplateParam(key string) string {
	for _, param := range c.TemplateParams {
		fragments := strings.SplitN(param, ":", 2)
		if len(fragments) == 2 && fragments[0] == key {
			return fragments[1]
		}
	}
	return ""
}

func (c *CreateVMFromTemplateTaskData) ArePVCsDataVolumes() bool {
	return !c.PVCsAreNotDataVolumes
}

func (c *CreateVMFromTemplateTaskData) GetExpectedVMStubMeta() *kubevirtv1.VirtualMachine {
	var disks []kubevirtv1.Disk
	var volumes []kubevirtv1.Volume

	if c.Template != nil && c.Template.Objects != nil {
		disks = append(disks, template.GetVM(c.Template).Spec.Template.Spec.Domain.Devices.Disks...)
		volumes = append(volumes, template.GetVM(c.Template).Spec.Template.Spec.Volumes...)
	}

	for _, dataVolume := range c.DataVolumesToCreate {
		name := dataVolume.Data.Name
		disk := kubevirtv1.Disk{
			Name: name,
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{Bus: c.ExpectedAdditionalDiskBus},
			},
		}
		volume := kubevirtv1.Volume{
			Name: name,
		}

		switch dataVolume.AttachmentType {
		case dv.DV, dv.OwnedDV:
			volume.DataVolume = &kubevirtv1.DataVolumeSource{
				Name: name,
			}
		case dv.PVC, dv.OwnedPVC:
			volume.PersistentVolumeClaim = &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: name,
			}
		}

		disks = append(disks, disk)
		volumes = append(volumes, volume)
	}

	return &kubevirtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GetTemplateParam(template.NameParam),
			Namespace: c.VMNamespace,
		},
		Spec: kubevirtv1.VirtualMachineSpec{
			Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{
				Spec: kubevirtv1.VirtualMachineInstanceSpec{
					Volumes: volumes,
					Domain: kubevirtv1.DomainSpec{
						Devices: kubevirtv1.Devices{
							Disks: disks,
						},
					},
				},
			},
		},
	}
}

func (c *CreateVMFromTemplateTaskData) SetDVorPVC(name string, attachmentType dv.TestDataVolumeAttachmentType) {
	switch attachmentType {
	case dv.DV:
		c.DataVolumes = append(c.DataVolumes, name)
	case dv.OwnedDV:
		c.OwnDataVolumes = append(c.OwnDataVolumes, name)
	case dv.PVC:
		c.PersistentVolumeClaims = append(c.PersistentVolumeClaims, name)
	case dv.OwnedPVC:
		c.OwnPersistentVolumeClaims = append(c.OwnPersistentVolumeClaims, name)
	}
}

type CreateVMFromTemplateTestConfig struct {
	TaskRunTestConfig
	TaskData CreateVMFromTemplateTaskData

	deploymentNamespace string
}

func (c *CreateVMFromTemplateTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
	c.TaskData.VMNamespace = options.ResolveNamespace(c.TaskData.VMTargetNamespace)
	if c.TaskData.Template != nil {
		tmpl := c.TaskData.Template
		if tmpl.Name != "" {
			tmpl.Name = E2ETestsRandomName(tmpl.Name)
		}
		tmpl.Namespace = options.ResolveNamespace(c.TaskData.TemplateTargetNamespace)

		c.TaskData.TemplateName = tmpl.Name
		c.TaskData.TemplateNamespace = tmpl.Namespace
	} else {
		if c.TaskData.TemplateTargetNamespace != "" {
			// for negative cases
			c.TaskData.TemplateNamespace = options.ResolveNamespace(c.TaskData.TemplateTargetNamespace)
		}
		if c.TaskData.TemplateName != "" && c.TaskData.IsCommonTemplate {
			c.TaskData.TemplateName += "-" + options.CommonTemplatesVersion
		}
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

func (c *CreateVMFromTemplateTestConfig) GetTaskRun() *v1beta1.TaskRun {
	var templateNamespace, vmNamespace string

	if !c.TaskData.UseDefaultTemplateNamespacesInTaskParams {
		templateNamespace = c.TaskData.TemplateNamespace
	}

	if !c.TaskData.UseDefaultVMNamespacesInTaskParams {
		vmNamespace = c.TaskData.VMNamespace
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-vm-create-from-template"),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: CreateVMFromTemplateClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: c.ServiceAccount,
			Params: []v1beta1.Param{
				{
					Name: CreateVMFromTemplateParams.TemplateName,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.TemplateName,
					},
				},
				{
					Name: CreateVMFromTemplateParams.TemplateNamespace,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: templateNamespace,
					},
				},
				{
					Name: CreateVMFromTemplateParams.TemplateParams,
					Value: v1beta1.ArrayOrString{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: c.TaskData.TemplateParams,
					},
				},
				{
					Name: CreateVMFromTemplateParams.VmNamespace,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: vmNamespace,
					},
				},
				{
					Name: CreateVMFromTemplateParams.DataVolumes,
					Value: v1beta1.ArrayOrString{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: c.TaskData.DataVolumes,
					},
				},
				{
					Name: CreateVMFromTemplateParams.OwnDataVolumes,
					Value: v1beta1.ArrayOrString{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: c.TaskData.OwnDataVolumes,
					},
				},
				{
					Name: CreateVMFromTemplateParams.PersistentVolumeClaims,
					Value: v1beta1.ArrayOrString{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: c.TaskData.PersistentVolumeClaims,
					},
				},
				{
					Name: CreateVMFromTemplateParams.OwnPersistentVolumeClaims,
					Value: v1beta1.ArrayOrString{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: c.TaskData.OwnPersistentVolumeClaims,
					},
				},
			},
		},
	}
}
