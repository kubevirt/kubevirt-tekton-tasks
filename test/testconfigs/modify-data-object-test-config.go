package testconfigs

import (
	"time"

	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	"github.com/onsi/ginkgo/v2"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/yaml"
)

type ModifyDataObjectTaskData struct {
	DataVolume          *cdiv1beta1.DataVolume
	DataSource          *cdiv1beta1.DataSource
	RawManifest         string
	WaitForSuccess      bool
	AllowReplace        bool
	DeleteObject        bool
	DeleteObjectName    string
	DeleteObjectKind    string
	Namespace           TargetNamespace
	dataObjectNamespace string
	SetOwnerReference   string
}

type ModifyDataObjectTestConfig struct {
	TaskRunTestConfig
	TaskData ModifyDataObjectTaskData

	deploymentNamespace string
}

func (c *ModifyDataObjectTestConfig) GetWaitForDataObjectTimeout() time.Duration {
	if c.TaskData.WaitForSuccess {
		return Timeouts.Zero.Duration
	}
	return c.GetTaskRunTimeout()
}

func (c *ModifyDataObjectTestConfig) Init(options *testoptions.TestOptions) {
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
		dv.Namespace = options.GetDeployNamespace()

		if options.StorageClass != "" {
			dv.Spec.PVC.StorageClassName = &options.StorageClass
		}
		if c.TaskData.DeleteObjectName != "" {
			c.TaskData.DeleteObjectName = dv.Name
		}

		c.TaskData.dataObjectNamespace = dv.Namespace
	}

	if c.TaskData.DataSource != nil {
		count += 1

		ds := c.TaskData.DataSource
		if ds.Name != "" {
			ds.Name = E2ETestsRandomName(ds.Name)
		}
		ds.Namespace = options.GetDeployNamespace()
		if c.TaskData.DeleteObjectName != "" {
			c.TaskData.DeleteObjectName = ds.Name
		}
		c.TaskData.dataObjectNamespace = ds.Namespace
	}

	if count > 1 {
		ginkgo.Fail("Need exactly one of DataVolume, DataSource or RawManifest")
	}

	if c.Timeout == nil || !c.TaskData.WaitForSuccess {
		c.Timeout = Timeouts.DefaultTaskRun
	}
}

func (c *ModifyDataObjectTestConfig) GetTaskRun() *pipev1.TaskRun {
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

	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(ModifyDataObjectTaskrunName),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: ModifyDataObjectTaskName,
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params: []pipev1.Param{
				{
					Name: ModifyDataObjectParams.Manifest,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: doStr,
					},
				},
				{
					Name: ModifyDataObjectParams.WaitForSuccess,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: ToStringBoolean(c.TaskData.WaitForSuccess),
					},
				},
				{
					Name: ModifyDataObjectParams.AllowReplace,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: ToStringBoolean(c.TaskData.AllowReplace),
					},
				},
				{
					Name: ModifyDataObjectParams.DeleteObject,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: ToStringBoolean(c.TaskData.DeleteObject),
					},
				},
				{
					Name: ModifyDataObjectParams.DeleteObjectName,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.DeleteObjectName,
					},
				},
				{
					Name: ModifyDataObjectParams.DeleteObjectKind,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.DeleteObjectKind,
					},
				},
				{
					Name: ModifyDataObjectParams.DataObjectNamespace,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.dataObjectNamespace,
					},
				}, {
					Name: SetOwnerReference,
					Value: pipev1.ParamValue{
						Type:      pipev1.ParamTypeString,
						StringVal: c.TaskData.SetOwnerReference,
					},
				},
			},
		},
	}
}
