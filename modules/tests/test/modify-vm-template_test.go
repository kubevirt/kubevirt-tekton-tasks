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
	kubevirtv1 "kubevirt.io/client-go/api/v1"

	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ExpectedResults struct {
	VM       *kubevirtv1.VirtualMachine
	Template *templatev1.Template

	ExpectedCPUSocketsTopologyNumber uint32
	ExpectedCPUCoresTopologyNumber   uint32
	ExpectedCPUThreadsTopologyNumber uint32
	ExpectedTemplateLabels           map[string]string
	ExpectedTemplateAnnotations      map[string]string
	ExpectedVMLabels                 map[string]string
	ExpectedVMAnnotations            map[string]string
	ExpectedVMMemory                 resource.Quantity
}

func (e ExpectedResults) ExpectCPUTopology() {
	Expect(e.VM.Spec.Template.Spec.Domain.CPU.Sockets).To(Equal(e.ExpectedCPUSocketsTopologyNumber), "cpu sockets should equal")
	Expect(e.VM.Spec.Template.Spec.Domain.CPU.Cores).To(Equal(e.ExpectedCPUCoresTopologyNumber), "cpu cores should equal")
	Expect(e.VM.Spec.Template.Spec.Domain.CPU.Threads).To(Equal(e.ExpectedCPUThreadsTopologyNumber), "cpu threads should equal")
}

func (e ExpectedResults) ExpectTemplateLabels() {
	checkValuesInMap(e.Template.Labels, e.ExpectedTemplateLabels, "template labels")
}

func (e ExpectedResults) ExpectTemplateAnnotations() {
	checkValuesInMap(e.Template.Annotations, e.ExpectedTemplateAnnotations, "template annotations")
}

func (e ExpectedResults) ExpectVMLabels() {
	checkValuesInMap(e.VM.Labels, e.ExpectedVMLabels, "vm labels")
}

func (e ExpectedResults) ExpectVMAnnotations() {
	checkValuesInMap(e.VM.Annotations, e.ExpectedVMAnnotations, "vm annotations")
}

func (e ExpectedResults) ExpectVMMemory() {
	Expect(e.VM.Spec.Template.Spec.Domain.Resources.Requests.Memory()).To(Equal(&e.ExpectedVMMemory), "memory should equal")
}

func checkValuesInMap(updatedMap, expectedResult map[string]string, message string) {
	for key, value := range expectedResult {
		Expect(updatedMap[key]).To(Equal(value), "value should equal in "+message)
	}
}

var _ = Describe("Copy template task", func() {
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

			var err error
			expectedResults.Template, err = f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.TODO(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(expectedResults.Template).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(expectedResults.Template)

			expectedResults.VM, err = zutils.DecodeVM(expectedResults.Template)
			Expect(err).ShouldNot(HaveOccurred(), "decode VM")

			expectedResults.ExpectCPUTopology()
			expectedResults.ExpectTemplateLabels()
			expectedResults.ExpectTemplateAnnotations()
			expectedResults.ExpectVMLabels()
			expectedResults.ExpectVMAnnotations()
			expectedResults.ExpectVMMemory()

		},
			table.Entry("should update template", &testconfigs.ModifyTemplateTestConfig{
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
			}),
		)
	})
})
