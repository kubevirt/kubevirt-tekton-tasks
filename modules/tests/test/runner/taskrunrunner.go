package runner

import (
	. "github.com/onsi/gomega"
	framework2 "github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/tekton"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tkntest "github.com/tektoncd/pipeline/test"
	"time"
)

type TaskRunRunner struct {
	framework *framework2.Framework
	taskRun   *pipev1beta1.TaskRun
}

const TaskRunExtraWaitDelay = 20 * time.Second

func NewTaskRunRunner(framework *framework2.Framework, taskRun *pipev1beta1.TaskRun) *TaskRunRunner {
	Expect(taskRun).ShouldNot(BeNil())
	return &TaskRunRunner{
		framework: framework,
		taskRun:   taskRun,
	}
}

func (r *TaskRunRunner) GetTaskRun() *pipev1beta1.TaskRun {
	return r.taskRun
}

func (r *TaskRunRunner) CreateTaskRun() *TaskRunRunner {
	taskRun, err := r.framework.TknClient.TaskRuns(r.taskRun.Namespace).Create(r.taskRun)
	Expect(err).ShouldNot(HaveOccurred())
	r.taskRun = taskRun
	r.framework.ManageTaskRun(r.taskRun)
	return r
}

func (r *TaskRunRunner) ExpectFailure() *TaskRunRunner {
	r.taskRun = tekton.WaitForTaskRunState(r.framework.TknClient, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout()+TaskRunExtraWaitDelay,
		tkntest.TaskRunFailed(r.taskRun.Name))
	return r
}

func (r *TaskRunRunner) ExpectSuccess() *TaskRunRunner {
	r.taskRun = tekton.WaitForTaskRunState(r.framework.TknClient, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout()+TaskRunExtraWaitDelay,
		tkntest.TaskRunSucceed(r.taskRun.Name))
	return r
}

func (r *TaskRunRunner) ExpectLogs(logs string) *TaskRunRunner {
	if logs != "" {
		Expect(tekton.GetTaskRunLogs(r.framework.CoreV1Client, r.taskRun)).Should(ContainSubstring(logs))
	}
	return r
}

func (r *TaskRunRunner) ExpectResults(results map[string]string) *TaskRunRunner {
	receivedResults := tekton.TaskResultsToMap(r.taskRun.Status.TaskRunResults)

	Expect(receivedResults).Should(HaveLen(len(results)))

	for resultKey, resultValue := range results {
		Expect(receivedResults[resultKey]).To(Equal(resultValue))
	}
	return r
}
