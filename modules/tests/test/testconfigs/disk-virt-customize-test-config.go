package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

type DiskVirtCustomizeTaskData struct {
	Datavolume        *v1beta12.DataVolume
	PVCName           string
	CustomizeCommands string
	AdditionalOptions string
	Verbose           bool
}

type DiskVirtCustomizeTestConfig struct {
	TaskRunTestConfig
	TaskData DiskVirtCustomizeTaskData

	deploymentNamespace string
}

func (c *DiskVirtCustomizeTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace

	if dv := c.TaskData.Datavolume; dv != nil {
		if dv.Name != "" {
			dv.Name = E2ETestsRandomName(dv.Name)
		}
		c.TaskData.PVCName = dv.Name

		dv.Namespace = options.DeployNamespace

		if options.StorageClass != "" {
			dv.Spec.PVC.StorageClassName = &options.StorageClass
		}

	}

}

func (c *DiskVirtCustomizeTestConfig) GetTaskRun() *v1beta1.TaskRun {
	return c.GetTaskRunWithName("")
}

func (c *DiskVirtCustomizeTestConfig) GetTaskRunWithName(nameSuffix string) *v1beta1.TaskRun {
	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-disk-virt-customize" + nameSuffix),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: DiskVirtCustomizeClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params: []v1beta1.Param{
				{
					Name: DiskVirtCustomizeParams.PVCName,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.PVCName,
					},
				},
				{
					Name: DiskVirtCustomizeParams.CustomizeCommands,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.CustomizeCommands,
					},
				},
				{
					Name: DiskVirtCustomizeParams.AdditionalOptions,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.AdditionalOptions,
					},
				},
				{
					Name: DiskVirtCustomizeParams.Verbose,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: ToStringBoolean(c.TaskData.Verbose),
					},
				},
			},
		},
	}
}
