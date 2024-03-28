package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineRunData struct {
	Name         string
	Params       []pipev1.Param
	TaskRunSpecs []pipev1.PipelineTaskRunSpec
	PipelineRef  *pipev1.PipelineRef
}

type PipelineTestConfig struct {
	TaskRunTestConfig
	Pipeline *pipev1.Pipeline

	deploymentNamespace string
	PipelineRunData     PipelineRunData
	PipelineRun         *pipev1.PipelineRun
}

func (c *PipelineTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
}

func (c *PipelineTestConfig) GetPipelineRun() *pipev1.PipelineRun {
	pipelineRun := &pipev1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(c.PipelineRunData.Name),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.PipelineRunSpec{
			PipelineRef:  c.PipelineRunData.PipelineRef,
			Timeouts:     &pipev1.TimeoutFields{Pipeline: &metav1.Duration{Duration: c.GetTaskRunTimeout()}},
			TaskRunSpecs: c.PipelineRunData.TaskRunSpecs,
			Params:       c.PipelineRunData.Params,
		},
	}
	c.PipelineRun = pipelineRun
	return pipelineRun
}
