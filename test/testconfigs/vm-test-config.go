package testconfigs

import (
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	template2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	v1 "github.com/openshift/api/template/v1"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type CreateVMTaskData struct {
	CreateMode CreateVMMode

	Template                *v1.Template
	TemplateTargetNamespace TargetNamespace

	VM *kubevirtv1.VirtualMachine

	IsCommonTemplate          bool
	StartVM                   string
	RunStrategy               string
	ExpectedAdditionalDiskBus string

	// Params
	// these two are set if Template is not nil
	TemplateName      string
	TemplateNamespace string

	// this is set if VM is not nil
	VMManifest string

	SetOwnerReference string
	TemplateParams    []string
	VMNamespace       string
	Virtctl           string
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

	return &kubevirtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vmName,
			Namespace: vmNamespace,
		},
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
		if c.TaskData.VM != nil {
			c.initCreateVMManifest(options)
		}
	case CreateVMTemplateMode:
		c.initCreateVMTemplate(options)
	default:
		panic("unknown VM create mode")
	}
}

func (c *CreateVMTestConfig) initCreateVMManifest(options *testoptions.TestOptions) {
	vm := c.TaskData.VM
	if vm.Name != "" {
		vm.Name = E2ETestsRandomName(vm.Name)
		vm.Spec.Template.ObjectMeta.Name = vm.Name
	}

	vm.Spec.Template.ObjectMeta.Namespace = ""

	vm.Namespace = ""
	c.TaskData.VMNamespace = options.GetDeployNamespace()

	c.TaskData.VMManifest = (&testobjects.TestVM{Data: vm}).ToString()

}

func (c *CreateVMTestConfig) initCreateVMTemplate(options *testoptions.TestOptions) {
	c.TaskData.VMNamespace = options.GetDeployNamespace()

	if tmpl := c.TaskData.Template; tmpl != nil {
		if tmpl.Name != "" {
			tmpl.Name = E2ETestsRandomName(tmpl.Name)
		}
		tmpl.Namespace = options.GetDeployNamespace()

		c.TaskData.TemplateName = tmpl.Name
	} else {
		if c.TaskData.TemplateName != "" && c.TaskData.IsCommonTemplate {
			c.TaskData.TemplateName += options.CommonTemplatesVersion
		}
	}
}

func (c *CreateVMTestConfig) GetTaskRun() *pipev1.TaskRun {
	var taskName, taskRunName string

	params := []pipev1.Param{
		{
			Name: CreateVMFromTemplateParams.StartVM,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.StartVM,
			},
		},
		{
			Name: CreateVMFromTemplateParams.RunStrategy,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.RunStrategy,
			},
		}, {
			Name: SetOwnerReference,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.SetOwnerReference,
			},
		},
	}

	vmNamespace := c.TaskData.VMNamespace

	switch c.TaskData.CreateMode {
	case CreateVMVMManifestMode:
		taskName = CreateVMFromManifestTaskName
		taskRunName = "taskrun-vm-create-from-manifest"

		params = append(params,
			pipev1.Param{
				Name: CreateVMFromManifestParams.Manifest,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: c.TaskData.VMManifest,
				},
			},
			pipev1.Param{
				Name: CreateVMFromManifestParams.Namespace,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: vmNamespace,
				},
			},
			pipev1.Param{
				Name: CreateVMFromManifestParams.Virtctl,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: c.TaskData.Virtctl,
				},
			},
		)
	case CreateVMTemplateMode:
		taskName = CreateVMFromTemplateTaskName
		taskRunName = "taskrun-vm-create-from-template"

		templateNamespace := c.TaskData.TemplateNamespace

		params = append(params,
			pipev1.Param{
				Name: CreateVMFromTemplateParams.TemplateName,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: c.TaskData.TemplateName,
				},
			},
			pipev1.Param{
				Name: CreateVMFromTemplateParams.TemplateNamespace,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: templateNamespace,
				},
			},

			pipev1.Param{
				Name: CreateVMFromTemplateParams.VmNamespace,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: vmNamespace,
				},
			},
		)

		if len(c.TaskData.TemplateParams) > 0 {
			params = append(params, pipev1.Param{
				Name: CreateVMFromTemplateParams.TemplateParams,
				Value: pipev1.ParamValue{
					Type:     pipev1.ParamTypeArray,
					ArrayVal: c.TaskData.TemplateParams,
				},
			})
		}
	}

	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(taskRunName),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: taskName,
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params:  params,
		},
	}
}
