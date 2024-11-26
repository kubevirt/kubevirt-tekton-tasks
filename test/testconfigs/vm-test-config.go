package testconfigs

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type CreateVMTaskData struct {
	VM *kubevirtv1.VirtualMachine

	StartVM                   string
	RunStrategy               string
	ExpectedAdditionalDiskBus string

	// this is set if VM is not nil
	VMManifest string

	SetOwnerReference string
	VMNamespace       string
	Virtctl           string
}

func (c *CreateVMTaskData) GetExpectedVM() (*kubevirtv1.VirtualMachine, error) {
	var vm *kubevirtv1.VirtualMachine

	err := yaml.Unmarshal([]byte(c.VMManifest), &vm)
	if err != nil {
		return nil, err
	}

	if vm.Namespace == "" {
		vm.Namespace = c.VMNamespace
	}

	return vm, err
}

type CreateVMTestConfig struct {
	TaskRunTestConfig
	TaskData CreateVMTaskData

	deploymentNamespace string
}

func (c *CreateVMTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
	c.initCreateVMManifest(options)
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

func (c *CreateVMTestConfig) GetTaskRun() *pipev1.TaskRun {
	vmNamespace := c.TaskData.VMNamespace
	taskName := CreateVMFromManifestTaskName
	taskRunName := "taskrun-vm-create-from-manifest"

	params := []pipev1.Param{
		{
			Name: CreateVMFromManifestParams.StartVM,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.StartVM,
			},
		}, {
			Name: CreateVMFromManifestParams.RunStrategy,
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
		}, {
			Name: CreateVMFromManifestParams.Manifest,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.VMManifest,
			},
		}, {
			Name: CreateVMFromManifestParams.Namespace,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: vmNamespace,
			},
		}, {
			Name: CreateVMFromManifestParams.Virtctl,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.Virtctl,
			},
		},
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
