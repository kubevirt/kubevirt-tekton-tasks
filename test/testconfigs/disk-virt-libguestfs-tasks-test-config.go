package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	pod "github.com/tektoncd/pipeline/pkg/apis/pipeline/pod"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type DiskVirtLibguestfsTaskData struct {
	Datavolume        *cdiv1beta1.DataVolume
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

func (c *DiskVirtLibguestfsTestConfig) GetTaskRun() *pipev1.TaskRun {
	return c.GetTaskRunWithName("")
}

func (c *DiskVirtLibguestfsTestConfig) GetTaskRunWithName(nameSuffix string) *pipev1.TaskRun {
	var taskName string

	params := []pipev1.Param{
		{
			Name: DiskVirtLibguestfsTasksParams.PVCName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.PVCName,
			},
		},
		{
			Name: DiskVirtLibguestfsTasksParams.Verbose,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: ToStringBoolean(c.TaskData.Verbose),
			},
		},
		{
			Name: DiskVirtLibguestfsTasksParams.AdditionalOptions,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.AdditionalOptions,
			},
		},
		{
			Name: DiskVirtLibguestfsTasksParams.VirtCommandsFileNameCommands,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.Commands,
			},
		},
	}
	if c.TaskData.LibguestfsTaskType == VirtSysPrepTaskType {
		taskName = DiskVirtSysprepTaskName
	} else {
		taskName = DiskVirtCustomizeTaskName
	}
	var qemuUser int64 = 107
	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-disk-" + string(c.TaskData.LibguestfsTaskType) + nameSuffix),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: taskName,
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params:  params,
			PodTemplate: &pod.PodTemplate{
				SecurityContext: &corev1.PodSecurityContext{
					RunAsUser: &qemuUser,
					FSGroup:   &qemuUser,
				},
			},
		},
	}
}
