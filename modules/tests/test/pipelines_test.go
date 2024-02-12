package test

import (
	"context"

	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

var _ = Describe("Pipelines tests", func() {
	f := framework.NewFramework()
	It("DV is created, disk-virt-sysprep, create dv, delete dvs", func() {
		config := &testconfigs.PipelineTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				Timeout: Timeouts.PipelineRunExtraWaitDelay,
			},
			Pipeline: &pipev1.Pipeline{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pipeline-dvs",
				},
				Spec: pipev1.PipelineSpec{
					Tasks: []pipev1.PipelineTask{
						{
							Name: "create-dv",
							TaskRef: &pipev1.TaskRef{
								Kind: pipev1.NamespacedTaskKind,
								Name: "modify-data-object",
							},
							Params: []pipev1.Param{
								{
									Name: "waitForSuccess",
									Value: pipev1.ParamValue{
										Type:      pipev1.ParamTypeString,
										StringVal: "true",
									},
								}, {
									Name: "manifest",
									Value: pipev1.ParamValue{
										Type: pipev1.ParamTypeString,
										StringVal: `
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: test-dv
  annotations:
    cdi.kubevirt.io/storage.bind.immediate.requested: "true"
spec:
  pvc:
    resources:
      requests:
        storage: 13Gi
    accessModes:
    - ReadWriteOnce
  source:
    registry:
      url: "docker://quay.io/containerdisks/centos-stream:9"
											`,
									},
								},
							},
						}, {
							Name: "sysprep-dv",
							TaskRef: &pipev1.TaskRef{
								Kind: pipev1.NamespacedTaskKind,
								Name: "disk-virt-sysprep",
							},
							Params: []pipev1.Param{
								{
									Name: "sysprepCommands",
									Value: pipev1.ParamValue{
										Type:      pipev1.ParamTypeString,
										StringVal: "run-command echo 'krtek' > new",
									},
								}, {
									Name: "pvc",
									Value: pipev1.ParamValue{
										Type:      pipev1.ParamTypeString,
										StringVal: "$(tasks.create-dv.results.name)",
									},
								},
							},
							RunAfter: []string{"create-dv"},
						}, {
							Name: "create-updated-dv",
							TaskRef: &pipev1.TaskRef{
								Kind: pipev1.NamespacedTaskKind,
								Name: "modify-data-object",
							},
							Params: []pipev1.Param{
								{

									Name: "manifest",
									Value: pipev1.ParamValue{
										Type: pipev1.ParamTypeString,
										StringVal: `
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: test-dv-updated
  annotations:
    cdi.kubevirt.io/storage.bind.immediate.requested: "true"
spec:
  pvc:
    resources:
      requests:
        storage: 13Gi
    accessModes:
    - ReadWriteOnce
  source:
    pvc:
      name: "$(tasks.create-dv.results.name)"
      namespace: "$(tasks.create-dv.results.namespace)"

											`,
									},
								},
							},
							RunAfter: []string{"sysprep-dv"},
						},
					},
				},
			},
			PipelineRunData: testconfigs.PipelineRunData{
				Name:   "test-dv-hardening",
				Params: []pipev1.Param{},
				PipelineRef: &pipev1.PipelineRef{
					Name: "test-pipeline-dvs",
				},
			},
		}

		f.TestSetup(config)
		pipelineRun := config.GetPipelineRun()
		_, err := f.TknClient.Pipelines(config.PipelineRun.Namespace).Create(context.Background(), config.Pipeline, metav1.CreateOptions{})
		Expect(err).ShouldNot(HaveOccurred())

		f.ManagePipelines(config.Pipeline)
		f.ManagePipelineRuns(config.PipelineRun)

		runner.NewPipelineRunRunner(f, pipelineRun).
			CreatePipelineRun().
			ExpectSuccess()

		pr, err := f.TknClient.PipelineRuns(config.PipelineRun.Namespace).Get(context.Background(), config.PipelineRun.Name, metav1.GetOptions{})
		Expect(err).ShouldNot(HaveOccurred())

		succeededConditionFound := false
		for _, condition := range pr.Status.Conditions {
			if condition.Type == apis.ConditionSucceeded && condition.Status == corev1.ConditionTrue {
				succeededConditionFound = true
			}
		}
		Expect(succeededConditionFound).To(BeTrue(), "pipelineRun should succeed")
	})
})
