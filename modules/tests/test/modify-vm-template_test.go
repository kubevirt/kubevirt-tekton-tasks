package test

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo"
	kubevirtv1 "kubevirt.io/api/core/v1"

	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ExpectedResults struct {
	ExpectedCPUSocketsTopologyNumber uint32
	ExpectedCPUCoresTopologyNumber   uint32
	ExpectedCPUThreadsTopologyNumber uint32
	ExpectedTemplateLabels           map[string]string
	ExpectedTemplateAnnotations      map[string]string
	ExpectedVMLabels                 map[string]string
	ExpectedVMAnnotations            map[string]string
	ExpectedVMMemory                 resource.Quantity
	ExpectedDisks                    []kubevirtv1.Disk
	ExpectedVolumes                  []kubevirtv1.Volume
}

func (e ExpectedResults) ExpectCPUTopology(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Domain.CPU.Sockets).To(Equal(e.ExpectedCPUSocketsTopologyNumber), "cpu sockets should equal")
	Expect(vm.Spec.Template.Spec.Domain.CPU.Cores).To(Equal(e.ExpectedCPUCoresTopologyNumber), "cpu cores should equal")
	Expect(vm.Spec.Template.Spec.Domain.CPU.Threads).To(Equal(e.ExpectedCPUThreadsTopologyNumber), "cpu threads should equal")
}

func (e ExpectedResults) ExpectTemplateLabels(template *templatev1.Template) {
	checkValuesInMap(template.Labels, e.ExpectedTemplateLabels, "template labels")
}

func (e ExpectedResults) ExpectTemplateAnnotations(template *templatev1.Template) {
	checkValuesInMap(template.Annotations, e.ExpectedTemplateAnnotations, "template annotations")
}

func (e ExpectedResults) ExpectVMLabels(vm *kubevirtv1.VirtualMachine) {
	checkValuesInMap(vm.Labels, e.ExpectedVMLabels, "vm labels")
}

func (e ExpectedResults) ExpectVMAnnotations(vm *kubevirtv1.VirtualMachine) {
	checkValuesInMap(vm.Annotations, e.ExpectedVMAnnotations, "vm annotations")
}

func (e ExpectedResults) ExpectVMMemory(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Domain.Resources.Requests.Memory()).To(Equal(&e.ExpectedVMMemory), "memory should equal")
}

func (e ExpectedResults) ExpectVMDisks(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Domain.Devices.Disks).To(Equal(e.ExpectedDisks), "vm disks should equal")
}

func (e ExpectedResults) ExpectVMVolumes(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Volumes).To(Equal(e.ExpectedVolumes), "vm volumes should equal")
}

func checkValuesInMap(updatedMap, expectedResult map[string]string, message string) {
	for key, value := range expectedResult {
		Expect(updatedMap[key]).To(Equal(value), "value should equal in "+message)
	}
}

var _ = Describe("Modify template task", func() {
	f := framework.NewFramework().LimitEnvScope(OKDEnvScope)
	Context("modify template fail", func() {
		table.DescribeTable("taskrun fails and no template is created", func(config *testconfigs.ModifyTemplateTestConfig) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}
			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)
		},
			table.Entry("no source template name specified", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "template-name param has to be specified",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("no template specified", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io \"cirros-vm-template\" not found",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
				},
			}),
			table.Entry("no service account", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "cannot get resource \"templates\" in API group \"template.openshift.io\"",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					Template:     testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("[NAMESPACE SCOPED] cannot updated template in different namespace", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "cannot get resource \"templates\" in API group \"template.openshift.io\"",
					LimitTestScope: NamespaceTestScope,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:            testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: SystemTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("wrong memory value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "quantities must match the regular expression",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					Memory:       "wrong memory value",
				},
			}),
			table.Entry("wrong number of CPU cores value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"wrong cpu cores\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUCores:     "wrong cpu cores",
				},
			}),
			table.Entry("wrong number of CPU Sockets value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"wrong cpu sockets\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUCores:     CPUCoresTopologyNumberStr,
					CPUThreads:   CPUThreadsTopologyNumberStr,
					CPUSockets:   "wrong cpu sockets",
				},
			}),
			table.Entry("wrong number of CPU Threads value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"wrong cpu threads\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUCores:     CPUCoresTopologyNumberStr,
					CPUThreads:   "wrong cpu threads",
				},
			}),
			table.Entry("wrong template labels value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:   testtemplate.CirrosTemplateName,
					CPUCores:       CPUCoresTopologyNumberStr,
					CPUThreads:     CPUThreadsTopologyNumberStr,
					CPUSockets:     CPUSocketsTopologyNumberStr,
					TemplateLabels: []string{"singleLabel"},
				},
			}),
			table.Entry("wrong template annotations value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					TemplateLabels:      MockArray,
					TemplateAnnotations: []string{"singleAnnotation"},
				},
			}),
			table.Entry("wrong vm labels value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					TemplateLabels:      MockArray,
					TemplateAnnotations: MockArray,
					VMLabels:            []string{"singleLabel"},
				},
			}),
			table.Entry("wrong template annotations value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					TemplateLabels:      MockArray,
					TemplateAnnotations: MockArray,
					VMLabels:            MockArray,
					VMAnnotations:       []string{"singleAnnotation"},
				},
			}),
			table.Entry("wrong disks value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "invalid character 'w' looking for beginning of value",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					Disks:        WrongStrSlice,
				},
			}),
			table.Entry("wrong volumes value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "invalid character 'w' looking for beginning of value",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					Volumes:      WrongStrSlice,
				},
			}),
		)
	})
	Context("modify template sucess", func() {
		table.DescribeTable("taskrun succeded and template is updated", func(config *testconfigs.ModifyTemplateTestConfig, expectedResults ExpectedResults) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			template, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.TODO(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(template).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(template)

			vm, _, err := zutils.DecodeVM(template)
			Expect(err).ShouldNot(HaveOccurred(), "decode VM")

			expectedResults.ExpectCPUTopology(vm)
			expectedResults.ExpectTemplateLabels(template)
			expectedResults.ExpectTemplateAnnotations(template)
			expectedResults.ExpectVMLabels(vm)
			expectedResults.ExpectVMAnnotations(vm)
			expectedResults.ExpectVMMemory(vm)
			expectedResults.ExpectVMDisks(vm)
			expectedResults.ExpectVMVolumes(vm)

		},
			table.Entry("should update template in the same namespace", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:            testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					Memory:              MemoryValue,
					TemplateLabels:      MockArray,
					TemplateAnnotations: MockArray,
					VMAnnotations:       MockArray,
					VMLabels:            MockArray,
					Disks:               MockDisks,
					Volumes:             MockVolumes,
				},
			}, ExpectedResults{
				ExpectedCPUSocketsTopologyNumber: CPUSocketsTopologyNumber,
				ExpectedCPUCoresTopologyNumber:   CPUCoresTopologyNumber,
				ExpectedCPUThreadsTopologyNumber: CPUThreadsTopologyNumber,
				ExpectedVMMemory:                 resource.MustParse(MemoryValue),
				ExpectedTemplateLabels:           LabelsAnnotationsMap,
				ExpectedTemplateAnnotations:      LabelsAnnotationsMap,
				ExpectedVMLabels:                 LabelsAnnotationsMap,
				ExpectedVMAnnotations:            LabelsAnnotationsMap,
				ExpectedDisks:                    Disks,
				ExpectedVolumes:                  Volumes,
			}),
			table.Entry("should update template across namespaces", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateName:            testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: DeployTargetNS,
					CPUCores:                CPUCoresTopologyNumberStr,
					CPUSockets:              CPUSocketsTopologyNumberStr,
					CPUThreads:              CPUThreadsTopologyNumberStr,
					Memory:                  MemoryValue,
					TemplateLabels:          MockArray,
					TemplateAnnotations:     MockArray,
					VMAnnotations:           MockArray,
					VMLabels:                MockArray,
					Disks:                   MockDisks,
					Volumes:                 MockVolumes,
				},
			}, ExpectedResults{
				ExpectedCPUSocketsTopologyNumber: CPUSocketsTopologyNumber,
				ExpectedCPUCoresTopologyNumber:   CPUCoresTopologyNumber,
				ExpectedCPUThreadsTopologyNumber: CPUThreadsTopologyNumber,
				ExpectedVMMemory:                 resource.MustParse(MemoryValue),
				ExpectedTemplateLabels:           LabelsAnnotationsMap,
				ExpectedTemplateAnnotations:      LabelsAnnotationsMap,
				ExpectedVMLabels:                 LabelsAnnotationsMap,
				ExpectedVMAnnotations:            LabelsAnnotationsMap,
				ExpectedDisks:                    Disks,
				ExpectedVolumes:                  Volumes,
			}),
		)

		It("taskrun succeded and template datavolume is removed and volumes, disks are updated", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                 testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:             testtemplate.RhelTemplateName,
					SourceTemplateNamespace:  DeployTargetNS,
					DeleteDatavolumeTemplate: true,
				},
			}
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			template, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.TODO(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(template).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(template)

			vm, _, err := zutils.DecodeVM(template)
			Expect(err).ShouldNot(HaveOccurred(), "decode VM")
			Expect(len(vm.Spec.DataVolumeTemplates)).To(Equal(0), "datavolume template should be empty")

			Expect(len(vm.Spec.Template.Spec.Volumes)).To(Equal(1), "there should be only 1 volume")
			emptyVolume := kubevirtv1.Volume{}
			Expect(vm.Spec.Template.Spec.Volumes[0].DataVolume).To(Equal(emptyVolume.DataVolume), "data volume should be nil")

			Expect(len(vm.Spec.Template.Spec.Domain.Devices.Disks)).To(Equal(1), "there should be only 1 disk")
			Expect(vm.Spec.Template.Spec.Domain.Devices.Disks[0].Name).ToNot(Equal("rootdisk"), "disk should not have name rootdisk")
		})
	})
})
