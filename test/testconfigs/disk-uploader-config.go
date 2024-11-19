package testconfigs

import (
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
)

type DiskUploaderTaskData struct {
	ExportSourceKind string
	ExportSourceName string
	VolumeName       string
	ImageDestination string
	PushTimeout      string
	SecretName       string
}

type DiskUploaderTestConfig struct {
	TaskRunTestConfig
	TaskData DiskUploaderTaskData

	deploymentNamespace string
}

func (c *DiskUploaderTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
}

func (c *DiskUploaderTestConfig) GetTaskRun() *pipev1.TaskRun {
	params := []pipev1.Param{
		{
			Name: "EXPORT_SOURCE_KIND",
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.ExportSourceKind,
			},
		},
		{
			Name: "EXPORT_SOURCE_NAME",
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.ExportSourceName,
			},
		},
		{
			Name: "VOLUME_NAME",
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.VolumeName,
			},
		},
		{
			Name: "IMAGE_DESTINATION",
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.ImageDestination,
			},
		},
		{
			Name: "PUSH_TIMEOUT",
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.PushTimeout,
			},
		},
		{
			Name: "SECRET_NAME",
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.SecretName,
			},
		},
	}

	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.E2ETestsRandomName("taskrun-disk-uploader"),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: "disk-uploader",
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params:  params,
		},
	}
}
