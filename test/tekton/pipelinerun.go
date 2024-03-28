package tekton

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/clients"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tkntest "github.com/tektoncd/pipeline/test"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/yaml"
)

type taskRunsLogs struct {
	mu   sync.Mutex
	logs map[string]string
}

func (l *taskRunsLogs) getLog(podName string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logs[podName]
}

func (l *taskRunsLogs) setLog(podName, log string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs[podName] = log
}

func (l *taskRunsLogs) getAllLogs() string {
	logs := ""
	l.mu.Lock()
	defer l.mu.Unlock()
	for podName, podLog := range l.logs {
		logs += fmt.Sprintf("\n %s: %s", podName, podLog)
	}
	return logs
}

func WaitForPipelineRunState(clients *clients.Clients, namespace, name string, timeout time.Duration, inState tkntest.ConditionAccessorFn) (*pipev1.PipelineRun, string) {
	pipelinePodsLogs := taskRunsLogs{}
	pipelinePodsLogs.logs = make(map[string]string)
	var wg sync.WaitGroup
	err := wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		pipelineRun, err := clients.TknClient.PipelineRuns(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}

		for _, reference := range pipelineRun.Status.ChildReferences {
			podName := reference.Name
			if pipelinePodsLogs.getLog(podName) == "" && podName != "" {
				wg.Add(1)
				go func() {
					defer wg.Done()
					req := clients.CoreV1Client.Pods(namespace).GetLogs(podName, &v1.PodLogOptions{
						Follow: true,
					})
					podLogs, err := req.Stream(context.Background())
					//when an error occurs, just end the function and do nothing, in next iteration the command will run again to get logs
					if err != nil {
						return
					}

					defer podLogs.Close()
					defer GinkgoRecover()

					result, err := io.ReadAll(podLogs)
					pipelinePodsLogs.setLog(podName, string(result))

					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				}()
			}
		}

		return inState(&pipelineRun.Status)
	})
	wg.Wait()
	logs := pipelinePodsLogs.getAllLogs()

	if err != nil {
		fmt.Println(logs)
	}

	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	pipelineRun, err := clients.TknClient.PipelineRuns(namespace).Get(context.Background(), name, metav1.GetOptions{})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	return pipelineRun, logs
}

func PrintPipelineRunDebugInfo(clients *clients.Clients, pipelineRunNamespace, pipelineRunName string) {
	// print conditions
	pipelineRun, err := clients.TknClient.PipelineRuns(pipelineRunNamespace).Get(context.Background(), pipelineRunName, metav1.GetOptions{})
	if err == nil {
		conditions, _ := yaml.Marshal(pipelineRun.Status.Conditions)
		fmt.Printf("pipelineRun conditions:\n%v\n", string(conditions))
	}
}
