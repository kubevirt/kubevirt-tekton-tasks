package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
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
	if vm := c.TaskData.VM; vm != nil {
		if vm.Name != "" {
			vm.Name = E2ETestsRandomName(vm.Name + "-" + string(c.TaskData.ExecInVMMode))
			vm.Spec.Template.ObjectMeta.Name = vm.Name
		}
		vm.Namespace = options.ResolveNamespace(c.TaskData.VMTargetNamespace)

		c.TaskData.VMName = vm.Name
		c.TaskData.VMNamespace = vm.Namespace
	} else {
		if c.TaskData.VMTargetNamespace != "" {
			// for negative cases
			c.TaskData.VMNamespace = options.ResolveNamespace(c.TaskData.VMTargetNamespace)
		}
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

func (c *ExecuteOrCleanupVMTestConfig) GetTaskRun() *v1beta1.TaskRun {
	var taskName, serviceAccountName, vmNamespace string

	if !c.TaskData.UseDefaultVMNamespacesInTaskParams {
		vmNamespace = c.TaskData.VMNamespace
	}

	params := []v1beta1.Param{
		{
			Name: ExecuteOrCleanupVMParams.VMName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.VMName,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.VMNamespace,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: vmNamespace,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.SecretName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.SecretName,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.Command,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: c.TaskData.Command,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.Args,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: c.TaskData.CommandArgs,
			},
		},
		{
			Name: ExecuteOrCleanupVMParams.Script,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.Script,
			},
		},
	}

	if c.TaskData.ExecInVMMode == CleanupVMMode {
		taskName = CleanupVMClusterTaskName
		if c.ServiceAccount != "" {
			serviceAccountName = CleanupVMServiceAccountName
		}

		params = append(params,
			v1beta1.Param{
				Name: ExecuteOrCleanupVMParams.Stop,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: ToStringBoolean(c.TaskData.Stop),
				},
			},
			v1beta1.Param{
				Name: ExecuteOrCleanupVMParams.Delete,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: ToStringBoolean(c.TaskData.Delete),
				},
			})
		if c.TaskData.Timeout != nil {
			params = append(params,
				v1beta1.Param{
					Name: ExecuteOrCleanupVMParams.Timeout,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.Timeout.Duration.String(),
					},
				})
		}
	} else {
		taskName = ExecuteInVMClusterTaskName
		if c.ServiceAccount != "" {
			serviceAccountName = ExecuteInVMServiceAccountName
		}
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-" + string(c.TaskData.ExecInVMMode)),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: taskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: serviceAccountName,
			Params:             params,
		},
	}
}
