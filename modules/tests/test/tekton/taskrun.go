package tekton

import (
	"bytes"
	"fmt"
	"github.com/onsi/gomega"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	tkntest "github.com/tektoncd/pipeline/test"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/yaml"
	"time"
)

func WaitForTaskRunState(client pipev1beta1.TaskRunInterface, name string, timeout time.Duration, inState tkntest.ConditionAccessorFn) *v1beta1.TaskRun {
	err := wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		r, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(&r.Status)
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	taskRun, err := client.Get(name, metav1.GetOptions{})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	return taskRun
}

func CreateTaskRun(taskRunClient pipev1beta1.TaskRunInterface, taskRun *v1beta1.TaskRun) (*v1beta1.TaskRun, string) {
	taskRun, err := taskRunClient.Create(taskRun)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return taskRun, taskRun.Name
}

func GetTaskRunLogs(podClient clientv1.PodInterface, taskRun *v1beta1.TaskRun) string {
	if taskRun.Status.PodName == "" {
		return ""
	}

	// print logs
	req := podClient.GetLogs(taskRun.Status.PodName, &v1.PodLogOptions{})

	podLogs, err := req.Stream()
	if err != nil {
		return ""
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return ""
	}
	return buf.String()
}

func DeleteTaskRun(taskRunClient pipev1beta1.TaskRunInterface, podClient clientv1.PodInterface, taskRunName string, debug bool) {
	originalError := recover()
	failed := originalError != nil

	if failed {
		defer panic(originalError)
	}

	if !debug || failed {
		taskRun, err := taskRunClient.Get(taskRunName, metav1.GetOptions{})
		if err == nil {
			if !debug {
				defer taskRunClient.Delete(taskRunName, &metav1.DeleteOptions{})
			}

			if failed {
				// print conditions
				conditions, _ := yaml.Marshal(taskRun.Status.Conditions)
				fmt.Printf("taskrun conditions:\n%v\n", string(conditions))

				if taskRun.Status.PodName == "" {
					return
				}
				fmt.Printf("%v pod logs:\n%v\n", taskRun.Status.PodName, GetTaskRunLogs(podClient, taskRun))
			}
		}
	}
}
