package dv

import (
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/utils"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
	"sigs.k8s.io/yaml"
)

type CreateDVTestConfig struct {
	Datavolume     *cdiv1beta1.DataVolume
	ServiceAccount string
	WaitForSuccess bool
	Timeout        *metav1.Duration
	Namespace      TargetNamespace
	LimitScope     utils.TestScope
	ExpectedLogs   string
	testConfig     *utils.TestConfig
}

func (c *CreateDVTestConfig) GetTaskRunTimeout() *metav1.Duration {
	if c.Timeout != nil && c.WaitForSuccess {
		return c.Timeout
	}
	return Timeouts.DefaultTaskRun
}

func (c *CreateDVTestConfig) GetWaitForDVTimeout() *metav1.Duration {
	if c.WaitForSuccess {
		return Timeouts.Zero
	}
	return c.Timeout
}

func (c *CreateDVTestConfig) Init(testConfig *utils.TestConfig) *CreateDVTestConfig {
	c.testConfig = testConfig
	if c.Datavolume != nil {
		if c.Datavolume.Name != "" {
			c.Datavolume.Name = E2ETestsName(c.Datavolume.Name)
		}
		c.Datavolume.Namespace = c.testConfig.GetResolvedTestNamespace(c.Namespace)

		if c.testConfig.StorageClass != "" {
			c.Datavolume.Spec.PVC.StorageClassName = &c.testConfig.StorageClass
		}
	}
	return c
}

func (c *CreateDVTestConfig) AsTaskRun() (*v1beta1.TaskRun, error) {
	var dv string
	if c.Datavolume != nil {
		dvbytes, err := yaml.Marshal(c.Datavolume)
		if err != nil {
			return nil, err
		}
		dv = string(dvbytes)
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsName("taskrun"),
			Namespace: c.testConfig.DeployNamespace,
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
						StringVal: ToStringBoolean(c.WaitForSuccess),
					},
				},
			},
		},
	}, nil
}
