package test

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Wait for VMI Status", func() {
	f := framework.NewFramework()

	DescribeTable("executes correctly", func(config *testconfigs.WaitForVMIStatusTestConfig) {
		f.TestSetup(config)

		if vm := config.TaskData.VM; vm != nil {
			vm, err := f.KubevirtClient.VirtualMachine(vm.Namespace).Create(context.Background(), vm, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageVMs(vm)
			if config.TaskData.ShouldStartVM {
				err := f.KubevirtClient.VirtualMachine(vm.Namespace).Start(context.Background(), vm.Name, &v1.StartOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			}
		}

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccessOrFailure(config.ExpectSuccess).
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectTermination(config.ExpectedTermination).
			ExpectResults(nil)
	},
		// negative validation cases
		Entry("no vmi", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "vmi-name should not be empty",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{},
		}),
		Entry("invalid vmi name", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "invalid vmi-name value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName: "name with spaces",
			},
		}),
		Entry("invalid success condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "success-condition: invalid condition: cannot parse jsonpath",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName:          "test",
				SuccessCondition: "test.....test",
			},
		}),
		Entry("invalid failure condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "failure-condition: could not parse condition",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName:          "test",
				FailureCondition: "invalid#$%^$&",
			},
		}),
		// negative cases
		Entry("fulfills failure condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedTermination: &testconfigs.TaskRunExpectedTermination{
					ExitCode: 2,
				},
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-failure-condition").Build(),
				FailureCondition: "status.phase == Running",
				ShouldStartVM:    true,
			},
		}),
		// positive cases
		Entry("no conditions report success immediately", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectSuccess: true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName: "test",
			},
		}),
		Entry("fulfills success condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectSuccess: true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-success-condition").Build(),
				SuccessCondition: "status.phase == Running",
				FailureCondition: "status.phase == Failed",
				ShouldStartVM:    true,
			},
		}),
		Entry("fulfills complex success condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectSuccess: true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-complex-success-condition").Build(),
				SuccessCondition: "status.phase in (Scheduling, Running), status.phase == Running",
				FailureCondition: "status.phase == Failed, spec.running == false",
				ShouldStartVM:    true,
			},
		}),
		Entry("fulfills success condition in the same namespace as deploy", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectSuccess: true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-success-condition-in-same-ns").Build(),
				SuccessCondition: "status.phase == Scheduled",
				ShouldStartVM:    true,
			},
		}),
	)
})
