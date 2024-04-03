package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

type WaitForVMIStatusTaskData struct {
	VM                *kubevirtv1.VirtualMachine
	VMTargetNamespace TargetNamespace

	ShouldStartVM bool

	// Params
	// these two are set if Template is not nil
	VMIName      string
	VMINamespace string

	SuccessCondition string
	FailureCondition string
}

type WaitForVMIStatusTestConfig struct {
	TaskRunTestConfig
	TaskData WaitForVMIStatusTaskData

	deploymentNamespace string
}

func (c *WaitForVMIStatusTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
	c.TaskData.VMINamespace = options.GetDeployNamespace()

	if vm := c.TaskData.VM; vm != nil {
		if vm.Name != "" {
			vm.Name = E2ETestsRandomName(vm.Name)
			vm.Spec.Template.ObjectMeta.Name = vm.Name
		}
		vm.Spec.Template.ObjectMeta.Namespace = ""
		vm.Namespace = c.TaskData.VMINamespace

		c.TaskData.VMIName = vm.Name
	}
}

func (c *WaitForVMIStatusTestConfig) GetTaskRun() *pipev1.TaskRun {
	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-" + WaitForVMIStatusTaskName),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: WaitForVMIStatusTaskName,
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params: []pipev1.Param{
				{
					Name: WaitForVMIStatusTasksParams.VMINamespace,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.VMINamespace,
					},
				},
				{
					Name: WaitForVMIStatusTasksParams.VMIName,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.VMIName,
					},
				}, {
					Name: WaitForVMIStatusTasksParams.SuccessCondition,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.SuccessCondition,
					},
				}, {
					Name: WaitForVMIStatusTasksParams.FailureCondition,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.FailureCondition,
					},
				},
			},
		},
	}
}
