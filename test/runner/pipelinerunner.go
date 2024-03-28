package runner

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	framework2 "github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/tekton"
	. "github.com/onsi/gomega"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tkntest "github.com/tektoncd/pipeline/test"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

type PipelineRunRunner struct {
	framework   *framework2.Framework
	pipelineRun *pipev1.PipelineRun
	logs        string
}

func NewPipelineRunRunner(framework *framework2.Framework, pipelineRun *pipev1.PipelineRun) *PipelineRunRunner {
	Expect(pipelineRun).ShouldNot(BeNil())
	return &PipelineRunRunner{
		framework:   framework,
		pipelineRun: pipelineRun,
	}
}

func (r *PipelineRunRunner) GetPipelineRun() *pipev1.PipelineRun {
	return r.pipelineRun
}

func (r *PipelineRunRunner) CreatePipelineRun() *PipelineRunRunner {
	pipelineRun, err := r.framework.TknClient.PipelineRuns(r.pipelineRun.Namespace).Create(context.Background(), r.pipelineRun, v1.CreateOptions{})
	Expect(err).ShouldNot(HaveOccurred())
	r.pipelineRun = pipelineRun
	r.framework.ManagePipelineRuns(r.pipelineRun)
	return r
}

func (r *PipelineRunRunner) ExpectFailure() *PipelineRunRunner {
	r.pipelineRun, r.logs = tekton.WaitForPipelineRunState(r.framework.Clients, r.pipelineRun.Namespace, r.pipelineRun.Name,
		r.pipelineRun.PipelineTimeout(context.Background())+constants.Timeouts.PipelineRunExtraWaitDelay.Duration,
		tkntest.PipelineRunFailed(r.pipelineRun.Name))
	return r
}

func (r *PipelineRunRunner) WaitForPipelineRunFinish() *PipelineRunRunner {
	r.pipelineRun, r.logs = tekton.WaitForPipelineRunState(r.framework.Clients, r.pipelineRun.Namespace, r.pipelineRun.Name,
		r.pipelineRun.PipelineTimeout(context.Background())+constants.Timeouts.PipelineRunExtraWaitDelay.Duration,
		func(accessor apis.ConditionAccessor) (bool, error) {
			succeeded, _ := tkntest.PipelineRunSucceed(r.pipelineRun.Name)(accessor)
			return succeeded, nil
		})
	return r
}

func (r *PipelineRunRunner) ExpectSuccess() *PipelineRunRunner {
	r.pipelineRun, r.logs = tekton.WaitForPipelineRunState(r.framework.Clients, r.pipelineRun.Namespace, r.pipelineRun.Name,
		r.pipelineRun.PipelineTimeout(context.Background())+constants.Timeouts.PipelineRunExtraWaitDelay.Duration,
		tkntest.PipelineRunSucceed(r.pipelineRun.Name))
	return r
}

func (r *PipelineRunRunner) ExpectSuccessOrFailure(expectSuccess bool) *PipelineRunRunner {
	if expectSuccess {
		r.ExpectSuccess()
	} else {
		r.ExpectFailure()
	}
	return r
}

func (r *PipelineRunRunner) ExpectLogs(logs ...string) *PipelineRunRunner {
	if len(logs) > 0 {
		for _, snippet := range logs {
			Expect(r.logs).Should(ContainSubstring(snippet))
		}
	}
	return r
}

func (r *PipelineRunRunner) GetResults() map[string]string {
	return tekton.PipelineResultsToMap(r.pipelineRun.Status.Results)
}

func (r *PipelineRunRunner) ExpectResults(results map[string]string) *PipelineRunRunner {
	return r.ExpectResultsWithLen(results, len(results))
}

func (r *PipelineRunRunner) ExpectResultsWithLen(results map[string]string, expectedLen int) *PipelineRunRunner {
	receivedResults := r.GetResults()

	Expect(receivedResults).Should(HaveLen(expectedLen))

	for resultKey, resultValue := range results {
		Expect(receivedResults[resultKey]).To(Equal(resultValue))
	}
	return r
}
