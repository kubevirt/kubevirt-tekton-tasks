package testconfigs

import (
	"time"

	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/onsi/ginkgo/v2"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/yaml"
)

type CreateDataObjectTaskData struct {
	DataVolume     *cdiv1beta1.DataVolume
	DataSource     *cdiv1beta1.DataSource
	RawManifest    string
	WaitForSuccess bool
	AllowReplace   bool
	Namespace      TargetNamespace
}

type CreateDataObjectTestConfig struct {
	TaskRunTestConfig
	TaskData CreateDataObjectTaskData

	deploymentNamespace string
}

func (c *CreateDataObjectTestConfig) GetWaitForDataObjectTimeout() time.Duration {
	if c.TaskData.WaitForSuccess {
		return Timeouts.Zero.Duration
	}
	return c.GetTaskRunTimeout()
}

func (c *CreateDataObjectTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace

	count := 0
	if c.TaskData.RawManifest != "" {
		count += 1
	}

	if c.TaskData.DataVolume != nil {
		count += 1

		dv := c.TaskData.DataVolume
		if dv.Name != "" {
			dv.Name = E2ETestsRandomName(dv.Name)
		}
		dv.Namespace = options.ResolveNamespace(c.TaskData.Namespace, "")

		if options.StorageClass != "" {
			dv.Spec.PVC.StorageClassName = &options.StorageClass
		}
	}

	if c.TaskData.DataSource != nil {
		count += 1

		ds := c.TaskData.DataSource
		if ds.Name != "" {
			ds.Name = E2ETestsRandomName(ds.Name)
		}
		ds.Namespace = options.ResolveNamespace(c.TaskData.Namespace, "")
	}

	if count > 1 {
		ginkgo.Fail("Need exactly one of DataVolume, DataSource or RawManifest")
	}

	if c.Timeout == nil || !c.TaskData.WaitForSuccess {
		c.Timeout = Timeouts.DefaultTaskRun
	}
}

func (c *CreateDataObjectTestConfig) GetTaskRun() *v1beta1.TaskRun {
	var do interface{}
	if c.TaskData.DataVolume != nil {
		do = c.TaskData.DataVolume
	} else if c.TaskData.DataSource != nil {
		do = c.TaskData.DataSource
	}

	doStr := c.TaskData.RawManifest
	if do != nil {
		doBytes, err := yaml.Marshal(do)
		if err != nil {
			ginkgo.Fail(err.Error())
		}
		doStr = string(doBytes)
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(CreateDataObjectTaskrunName),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: CreateDataObjectClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: c.ServiceAccount,
			Params: []v1beta1.Param{
				{
					Name: CreateDataObjectParams.Manifest,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: doStr,
					},
				},
				{
					Name: CreateDataObjectParams.WaitForSuccess,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: ToStringBoolean(c.TaskData.WaitForSuccess),
					},
				},
				{
					Name: CreateDataObjectParams.AllowReplace,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: ToStringBoolean(c.TaskData.AllowReplace),
					},
				},
			},
		},
	}
}
