package dv

import (
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/utils"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type CreateDVTaskData struct {
	Datavolume     *TestDataVolume
	WaitForSuccess bool
}

type CreateDVTestConfig struct {
	utils.TaskRunTestConfig
	TaskData CreateDVTaskData
}

func (c *CreateDVTestConfig) GetTaskRunTimeout() *metav1.Duration {
	if c.Timeout != nil && c.TaskData.WaitForSuccess {
		return c.Timeout
	}
	return Timeouts.DefaultTaskRun
}

func (c *CreateDVTestConfig) GetWaitForDVTimeout() *metav1.Duration {
	if c.TaskData.WaitForSuccess {
		return Timeouts.Zero
	}
	return c.Timeout
}

func (c *CreateDVTestConfig) Init(testConfig *utils.TestConfig) *CreateDVTestConfig {
	c.SetTestConfig(testConfig)
	if c.TaskData.Datavolume != nil {
		dv := c.TaskData.Datavolume.Data
		if dv.Name != "" {
			dv.Name = E2ETestsName(dv.Name)
		}
		dv.Namespace = testConfig.GetResolvedTestNamespace(c.Namespace)

		if testConfig.StorageClass != "" {
			dv.Spec.PVC.StorageClassName = &testConfig.StorageClass
		}
	}
	return c
}

func (c *CreateDVTestConfig) AsTaskRun() (*v1beta1.TaskRun, error) {
	var dv string
	if c.TaskData.Datavolume != nil {
		dvbytes, err := yaml.Marshal(c.TaskData.Datavolume.Data)
		if err != nil {
			return nil, err
		}
		dv = string(dvbytes)
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsName("taskrun-dv-create"),
			Namespace: c.GetTestConfig().DeployNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: CreateDataVolumeFromManifestClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            c.Timeout,
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
	}, nil
}
