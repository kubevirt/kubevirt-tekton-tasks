package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type DiskVirtLibguestfsTaskData struct {
	Datavolume        *v1beta12.DataVolume
	PVCName           string
	Commands          string
	AdditionalOptions string
	Verbose           bool

	// supplied
	LibguestfsTaskType LibguestfsTaskType
}

type DiskVirtLibguestfsTestConfig struct {
	TaskRunTestConfig
	TaskData DiskVirtLibguestfsTaskData

	deploymentNamespace string
}

func (c *DiskVirtLibguestfsTestConfig) Init(options *testoptions.TestOptions) {
	if c.TaskData.LibguestfsTaskType == "" {
		panic("unknow libguestfs type")
	}

	c.deploymentNamespace = options.DeployNamespace

	if dv := c.TaskData.Datavolume; dv != nil {
		if dv.Name != "" {
			dv.Name = E2ETestsRandomName(dv.Name + "-" + string(c.TaskData.LibguestfsTaskType))
		}
		c.TaskData.PVCName = dv.Name

		dv.Namespace = options.DeployNamespace

		if options.StorageClass != "" {
			dv.Spec.PVC.StorageClassName = &options.StorageClass
		}

	}

}

func (c *DiskVirtLibguestfsTestConfig) GetTaskRun() *v1beta1.TaskRun {
	return c.GetTaskRunWithName("")
}

func (c *DiskVirtLibguestfsTestConfig) GetTaskRunWithName(nameSuffix string) *v1beta1.TaskRun {
	var taskName string

	params := []v1beta1.Param{
		{
			Name: DiskVirtLibguestfsTasksParams.PVCName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.PVCName,
			},
		},
		{
			Name: DiskVirtLibguestfsTasksParams.Verbose,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: ToStringBoolean(c.TaskData.Verbose),
			},
		},
		{
			Name: DiskVirtLibguestfsTasksParams.AdditionalOptions,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.AdditionalOptions,
			},
		},
	}
	if c.TaskData.LibguestfsTaskType == VirtSysPrepTaskType {
		params = append(params, v1beta1.Param{
			Name: DiskVirtLibguestfsTasksParams.SysprepCommands,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.Commands,
			},
		})
		taskName = DiskVirtSysprepClusterTaskName
	} else {
		params = append(params, v1beta1.Param{
			Name: DiskVirtLibguestfsTasksParams.CustomizeCommands,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.Commands,
			},
		})
		taskName = DiskVirtCustomizeClusterTaskName
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-disk-" + string(c.TaskData.LibguestfsTaskType) + nameSuffix),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: taskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params:  params,
		},
	}
}
