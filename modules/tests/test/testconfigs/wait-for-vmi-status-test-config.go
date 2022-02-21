package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
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
	c.TaskData.VMINamespace = options.ResolveNamespace(c.TaskData.VMTargetNamespace, c.TaskData.VMINamespace)

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

func (c *WaitForVMIStatusTestConfig) GetTaskRun() *v1beta1.TaskRun {
	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-" + WaitForVMIStatusClusterTaskName),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: WaitForVMIStatusClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: c.ServiceAccount,
			Params: []v1beta1.Param{
				{
					Name: WaitForVMIStatusTasksParams.VMINamespace,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.VMINamespace,
					},
				},
				{
					Name: WaitForVMIStatusTasksParams.VMIName,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.VMIName,
					},
				}, {
					Name: WaitForVMIStatusTasksParams.SuccessCondition,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.SuccessCondition,
					},
				}, {
					Name: WaitForVMIStatusTasksParams.FailureCondition,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.FailureCondition,
					},
				},
			},
		},
	}
}
