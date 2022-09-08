package test

import (
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
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
		Entry("invalid vmi namespace", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "invalid vmi-namespace value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName:      "test",
				VMINamespace: "namespace with spaces",
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
		Entry("[NAMESPACE SCOPED] cannot check status for VMI in different namespace", &testconfigs.WaitForVMIStatusTestConfig{
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
		Entry("fulfills failure condition", &testconfigs.WaitForVMIStatusTestConfig{
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
		Entry("no conditions report success immediately", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectSuccess:  true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VMIName: "test",
			},
		}),
		Entry("fulfills success condition", &testconfigs.WaitForVMIStatusTestConfig{
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
		Entry("fulfills complex success condition", &testconfigs.WaitForVMIStatusTestConfig{
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
		Entry("[CLUSTER SCOPED] fulfills success condition in the same namespace as deploy", &testconfigs.WaitForVMIStatusTestConfig{
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
		// JSONPath tests
		Entry("fulfills success condition - jsonPath", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectSuccess:  true,
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-success-condition").Build(),
				SuccessCondition: "jsonpath='{.status.phase}' == Running",
				FailureCondition: "jsonpath='{.status.phase}' == Failed",
				ShouldStartVM:    true,
			},
		}),
		Entry("invalid success condition - jsonPath", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectedLogs:   "success-condition: could not parse condition {.status.phase} == Running: unable to parse requirement: <nil>: Invalid value: \"{.status.phase}\": name part must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyName',  or 'my.name',  or '123-abc', regex used for validation is '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]')",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-success-condition").Build(),
				SuccessCondition: "{.status.phase} == Running",
				FailureCondition: "jsonpath='{.status.phase}' == Failed",
				ShouldStartVM:    true,
			},
		}),
		Entry("invalid failure condition - jsonPath", &testconfigs.WaitForVMIStatusTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: constants.WaitForVMIStatusServiceAccountName,
				ExpectedLogs:   "failure-condition: valid jsonpath format is jsonpath='{.status.phase}' == Success",
			},
			TaskData: testconfigs.WaitForVMIStatusTaskData{
				VM:               testobjects.NewTestAlpineVM("fulfills-success-condition").Build(),
				SuccessCondition: "jsonpath='{.status.phase}' == Running",
				FailureCondition: "jsonpath={.status.phase}' == Failed",
				ShouldStartVM:    true,
			},
		}),
	)
})
