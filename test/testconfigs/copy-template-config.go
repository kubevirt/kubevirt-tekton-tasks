package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	v1 "github.com/openshift/api/template/v1"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CopyTemplateTaskData struct {
	Template                   *v1.Template
	TemplateNamespace          string
	SourceTemplateName         string
	TargetTemplateName         string
	AllowReplace               string
	UsePlainTargetTemplateName bool
}

type CopyTemplateTestConfig struct {
	TaskRunTestConfig
	TaskData CopyTemplateTaskData

	deploymentNamespace string
}

func (c *CopyTemplateTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
	if !c.TaskData.UsePlainTargetTemplateName && c.TaskData.AllowReplace != "true" {
		c.TaskData.TargetTemplateName = E2ETestsRandomName(c.TaskData.TargetTemplateName)
	}
	if c.TaskData.Template != nil {
		c.TaskData.Template.Namespace = options.GetDeployNamespace()
		c.TaskData.TemplateNamespace = options.GetDeployNamespace()
		originalTemplateName := c.TaskData.Template.Name
		c.TaskData.Template.Name = E2ETestsRandomName(c.TaskData.Template.Name)
		if c.TaskData.SourceTemplateName == originalTemplateName {
			c.TaskData.SourceTemplateName = c.TaskData.Template.Name
		}
	}
}

func (c *CopyTemplateTestConfig) GetTaskRun() *pipev1.TaskRun {
	params := []pipev1.Param{
		{
			Name: SourceTemplateNameOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.SourceTemplateName,
			},
		},
		{
			Name: SourceTemplateNamespaceOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.TemplateNamespace,
			},
		},
		{
			Name: TargetTemplateNameOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.TargetTemplateName,
			},
		},
		{
			Name: TargetTemplateNamespaceOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.TemplateNamespace,
			},
		}, {
			Name: AllowReplaceOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.AllowReplace,
			},
		},
	}

	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(CopyTemplateTaskRunName),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: CopyTemplateTaskName,
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params:  params,
		},
	}
}
