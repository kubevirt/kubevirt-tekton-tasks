package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineRunData struct {
	Name         string
	Params       []v1beta1.Param
	TaskRunSpecs []v1beta1.PipelineTaskRunSpec
	PipelineRef  *v1beta1.PipelineRef
}

type PipelineTestConfig struct {
	TaskRunTestConfig
	Pipeline *v1beta1.Pipeline

	deploymentNamespace string
	PipelineRunData     PipelineRunData
	PipelineRun         *v1beta1.PipelineRun
}

func (c *PipelineTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
}

func (c *PipelineTestConfig) GetPipelineRun() *v1beta1.PipelineRun {
	pipelineRun := &v1beta1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(c.PipelineRunData.Name),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef:  c.PipelineRunData.PipelineRef,
			Timeout:      &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			TaskRunSpecs: c.PipelineRunData.TaskRunSpecs,
			Params:       c.PipelineRunData.Params,
		},
	}
	c.PipelineRun = pipelineRun
	return pipelineRun
}
