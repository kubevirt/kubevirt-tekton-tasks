package test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dv"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	customizeCommand = `run-command date
mkdir /test
touch /test/whale
append-line /test/whale:hello blue world
run-command date
no-logfile
`
)

var _ = Describe("Run disk virt-customize", func() {
	f := framework.NewFramework()

	table.DescribeTable("taskrun fails", func(config *testconfigs.DiskVirtCustomizeTestConfig) {
		f.TestSetup(config)

		if dataVolume := config.TaskData.Datavolume; dataVolume != nil {
			dataVolume, err := f.CdiClient.DataVolumes(dataVolume.Namespace).Create(dataVolume)
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dataVolume)
			err = dv.WaitForSuccessfulDataVolume(f.CdiClient, dataVolume.Namespace, dataVolume.Name, constants.Timeouts.SmallDVCreation.Duration)
			Expect(err).ShouldNot(HaveOccurred())
		}

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectFailure().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(nil)
	},
		table.Entry("no pvc", &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
			TaskData:          testconfigs.DiskVirtCustomizeTaskData{},
		}),
		table.Entry("invalid pvc", &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				Timeout: &metav1.Duration{time.Minute},
			},
			TaskData: testconfigs.DiskVirtCustomizeTaskData{
				PVCName: "illegal-test-pvc-1548748",
			},
		}),
		table.Entry("no customize commands", &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "customize-commands option or CUSTOMIZE_COMMANDS env variable is required",
			},
			TaskData: testconfigs.DiskVirtCustomizeTaskData{
				Datavolume: dv.NewBlankDataVolume("no-customize-commands").Build(),
			},
		}),
		table.Entry("wrong customize commands", &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "command 'illegal-operation' not valid",
			},
			TaskData: testconfigs.DiskVirtCustomizeTaskData{
				Datavolume:        dv.NewBlankDataVolume("wrong-customize-commands").Build(),
				CustomizeCommands: "illegal-operation",
			},
		}),
		table.Entry("wrong additional options", &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "unrecognized option '--illegal-command'",
			},
			TaskData: testconfigs.DiskVirtCustomizeTaskData{
				Datavolume:        dv.NewBlankDataVolume("wrong-additional-options").Build(),
				CustomizeCommands: "update",
				AdditionalOptions: "--verbose --illegal-command illegal args",
			},
		}),
		table.Entry("empty disk fails", &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogsList: []string{
					"no operating systems were found in the guest image",
				},
			},
			TaskData: testconfigs.DiskVirtCustomizeTaskData{
				Datavolume:        dv.NewBlankDataVolume("empty-disk").Build(),
				CustomizeCommands: "update",
			},
		}),
	)

	It("virt-customize works", func() {
		// main test
		testConfig := &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogsList: []string{
					"Running: date",
					"Making directory: /test",
					"Running touch: /test/whale",
					"Appending line to /test/whale",
					"Scrubbing the log file",
					"Finishing off",
				},
			},
			TaskData: testconfigs.DiskVirtCustomizeTaskData{
				Datavolume: dv.NewBlankDataVolume("basic-functionality").
					WithSize(5, resource.Giga).
					WithRegistrySource("docker://kubevirt/fedora-cloud-registry-disk-demo").
					Build(),
				CustomizeCommands: customizeCommand,
			},
		}
		f.TestSetup(testConfig)

		dataVolume := testConfig.TaskData.Datavolume
		// prepare DataVolume
		dataVolume, err := f.CdiClient.DataVolumes(dataVolume.Namespace).Create(dataVolume)
		Expect(err).ShouldNot(HaveOccurred())
		f.ManageDataVolumes(dataVolume)
		err = dv.WaitForSuccessfulDataVolume(f.CdiClient, dataVolume.Namespace, dataVolume.Name, constants.Timeouts.SmallDVCreation.Duration)
		Expect(err).ShouldNot(HaveOccurred())

		// run main test
		runner.NewTaskRunRunner(f, testConfig.GetTaskRunWithName("-test")).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(testConfig.GetAllExpectedLogs()...).
			ExpectResults(nil)

		// test the result data is preserved on the PVC

		expectConfig := &testconfigs.DiskVirtCustomizeTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogsList: []string{
					"\"level\":\"debug\"",
					"-x",
					"Running: cat /test/whale",
					"hello blue world",
					"Finishing off",
				},
			},
			TaskData: testconfigs.DiskVirtCustomizeTaskData{
				PVCName:           dataVolume.Name,
				CustomizeCommands: "run-command cat /test/whale",
				Verbose:           true,
				AdditionalOptions: "--smp 2 --verbose",
			},
		}

		f.TestSetup(expectConfig)

		runner.NewTaskRunRunner(f, expectConfig.GetTaskRunWithName("-expect")).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(expectConfig.GetAllExpectedLogs()...).
			ExpectResults(nil)
	})
})
