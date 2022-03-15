package test

import (
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Wait for VMI Status", func() {
	f := framework.NewFramework()

	table.DescribeTable("executes correctly", func(config *testconfigs.WaitForVMIStatusTestConfig) {
		f.TestSetup(config)

		if vm := config.TaskData.VM; vm != nil {
			vm, err := f.KubevirtClient.VirtualMachine(vm.Namespace).Create(vm)
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageVMs(vm)
			if config.TaskData.ShouldStartVM {
				err := f.KubevirtClient.VirtualMachine(vm.Namespace).Start(vm.Name, &v1.StartOptions{})
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
		table.Entry("no vmi", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "vmi-name should not be empty",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{},
		}),
		table.Entry("invalid vmi name", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "invalid vmi-name value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName: "name with spaces",
			},
		}),
		table.Entry("invalid vmi namespace", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "invalid vmi-namespace value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName:      "test",
				VMINamespace: "namespace with spaces",
			},
		}),
		table.Entry("invalid success condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "success-condition: invalid condition: cannot parse jsonpath",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName:          "test",
				SuccessCondition: "test.....test",
			},
		}),
		table.Entry("invalid failure condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "failure-condition: could not parse condition",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName:          "test",
				FailureCondition: "invalid#$%^$&",
			},
		}),
		table.Entry("no service account", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "cannot list resource \"virtualmachineinstances\" in API group \"kubevirt.io\"",
				Timeout:      &metav1.Duration{1 * time.Minute},
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("no-serviceaccount").Build(),
				SuccessCondition: "status.phase == Running",
			},
		}),
		table.Entry("[NAMESPACE SCOPED] cannot check status for VMI in different namespace", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				LimitTestScope: constants.NamespaceTestScope,
				ExpectedLogs:   "cannot list resource \"virtualmachineinstances\" in API group \"kubevirt.io\"",
				Timeout:        &metav1.Duration{1 * time.Minute},
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:                testobjects.NewTestAlpineVM("wait-for-vmi-status-in-different-ns").Build(),
				VMTargetNamespace: constants.SystemTargetNS,
				SuccessCondition:  "status.phase == Running",
			},
		}),
		// negative cases
		table.Entry("fulfills failure condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
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
		table.Entry("no conditions report success immediately", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectSuccess:  true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName: "test",
			},
		}),
		table.Entry("fulfills success condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectSuccess:  true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-success-condition").Build(),
				SuccessCondition: "status.phase == Running",
				FailureCondition: "status.phase == Failed",
				ShouldStartVM:    true,
			},
		}),
		table.Entry("fulfills complex success condition", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectSuccess:  true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-complex-success-condition").Build(),
				SuccessCondition: "status.phase in (Scheduling, Running), status.phase == Running",
				FailureCondition: "status.phase == Failed, spec.running == false",
				ShouldStartVM:    true,
			},
		}),
		table.Entry("[CLUSTER SCOPED] fulfills success condition in the same namespace as deploy", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectSuccess:  true,
				LimitTestScope: constants.ClusterTestScope,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:                testobjects.NewTestAlpineVM("fulfills-success-condition-in-same-ns").Build(),
				SuccessCondition:  "status.phase == Scheduled",
				ShouldStartVM:     true,
				VMTargetNamespace: constants.DeployTargetNS,
			},
		}),
	)
})
