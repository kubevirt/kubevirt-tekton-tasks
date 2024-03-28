package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

type ExecuteOrCleanupVMTaskData struct {
	Secret *corev1.Secret
	VM     *kubevirtv1.VirtualMachine

	VMTargetNamespace TargetNamespace

	UseDefaultVMNamespacesInTaskParams bool
	ShouldStartVM                      bool
	// supplied
	ExecInVMMode ExecInVMMode

	// Params
	// these three are set if VM is not nil
	VMName      string
	VMNamespace string
	SecretName  string

	Script      string
	Command     []string
	CommandArgs []string
	// cleanup VM
	Stop    bool
	Delete  bool
	Timeout *metav1.Duration
}

type ExecuteOrCleanupVMTestConfig struct {
	TaskRunTestConfig
	TaskData ExecuteOrCleanupVMTaskData

	deploymentNamespace string
}

func (c *ExecuteOrCleanupVMTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
	c.TaskData.VMNamespace = options.GetDeployNamespace()

	if vm := c.TaskData.VM; vm != nil {
		if vm.Name != "" {
			vm.Name = E2ETestsRandomName(vm.Name + "-" + string(c.TaskData.ExecInVMMode))
			vm.Spec.Template.ObjectMeta.Name = vm.Name
		}
		vm.Namespace = c.TaskData.VMNamespace

		c.TaskData.VMName = vm.Name
	}

	if secret := c.TaskData.Secret; secret != nil {
		if secret.Name != "" {
			if vm := c.TaskData.VM; vm != nil {
				secret.Name = vm.Name
			} else {
				secret.Name = E2ETestsRandomName(secret.Name + "-" + string(c.TaskData.ExecInVMMode))
			}
		}
		secret.Namespace = options.DeployNamespace

		c.TaskData.SecretName = secret.Name
	}
}

func (c *ExecuteOrCleanupVMTestConfig) GetTaskRun() *pipev1.TaskRun {
	var taskName, vmNamespace string

	if !c.TaskData.UseDefaultVMNamespacesInTaskParams {
		vmNamespace = c.TaskData.VMNamespace
	}

	params := []pipev1.Param{
		{
			Name: ExecuteOrCleanupVMParams.VMName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.VMName,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.VMNamespace,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: vmNamespace,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.SecretName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.SecretName,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.Script,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.Script,
			},
		},
	}

	if len(c.TaskData.Command) > 0 {
		params = append(params, pipev1.Param{
			Name: ExecuteOrCleanupVMParams.Command,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: c.TaskData.Command,
			},
		})
	}

	if len(c.TaskData.CommandArgs) > 0 {
		params = append(params, pipev1.Param{
			Name: ExecuteOrCleanupVMParams.Args,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: c.TaskData.CommandArgs,
			},
		})
	}

	if c.TaskData.ExecInVMMode == CleanupVMMode {
		taskName = CleanupVMTaskName

		params = append(params,
			pipev1.Param{
				Name: ExecuteOrCleanupVMParams.Stop,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: ToStringBoolean(c.TaskData.Stop),
				},
			},
			pipev1.Param{
				Name: ExecuteOrCleanupVMParams.Delete,
				Value: pipev1.ParamValue{
					Type:      pipev1.ParamTypeString,
					StringVal: ToStringBoolean(c.TaskData.Delete),
				},
			})
		if c.TaskData.Timeout != nil {
			params = append(params,
				pipev1.Param{
					Name: ExecuteOrCleanupVMParams.Timeout,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.Timeout.Duration.String(),
					},
				})
		}
	} else {
		taskName = ExecuteInVMTaskName
	}

	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-" + string(c.TaskData.ExecInVMMode)),
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
