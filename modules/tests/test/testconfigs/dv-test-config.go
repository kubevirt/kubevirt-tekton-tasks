package testconfigs

import (
	"github.com/onsi/ginkgo"
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
	"sigs.k8s.io/yaml"
	"time"
)

type CreateDVTaskData struct {
	Datavolume     *v1beta12.DataVolume
	WaitForSuccess bool
	Namespace      TargetNamespace
}

type CreateDVTestConfig struct {
	TaskRunTestConfig
	TaskData CreateDVTaskData

	deploymentNamespace string
}

func (c *CreateDVTestConfig) GetWaitForDVTimeout() time.Duration {
	if c.TaskData.WaitForSuccess {
		return Timeouts.Zero.Duration
	}

	if c.Timeout != nil {
		return c.Timeout.Duration

	}
	return Timeouts.DefaultTaskRun.Duration
}

func (c *CreateDVTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
	if c.TaskData.Datavolume != nil {
		dv := c.TaskData.Datavolume
		if dv.Name != "" {
			dv.Name = E2ETestsRandomName(dv.Name)
		}
		dv.Namespace = options.ResolveNamespace(c.TaskData.Namespace)

		if options.StorageClass != "" {
			dv.Spec.PVC.StorageClassName = &options.StorageClass
		}
	}

	if c.Timeout == nil || !c.TaskData.WaitForSuccess {
		c.Timeout = Timeouts.DefaultTaskRun
	}
}

func (c *CreateDVTestConfig) GetTaskRun() *v1beta1.TaskRun {
	var dv string
	if c.TaskData.Datavolume != nil {
		dvbytes, err := yaml.Marshal(c.TaskData.Datavolume)
		if err != nil {
			ginkgo.Fail(err.Error())
		}
		dv = string(dvbytes)
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-dv-create"),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: CreateDataVolumeFromManifestClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: c.ServiceAccount,
			Params: []v1beta1.Param{
				{
					Name: CreateDataVolumeFromManifestParams.Manifest,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: dv,
					},
				},
				{
					Name: CreateDataVolumeFromManifestParams.WaitForSuccess,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: ToStringBoolean(c.TaskData.WaitForSuccess),
					},
				},
			},
		},
	}
}
