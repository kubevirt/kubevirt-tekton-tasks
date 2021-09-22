package test

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo"
	templatev1 "github.com/openshift/api/template/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"

	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Copy template task", func() {
	f := framework.NewFramework().LimitEnvScope(OKDEnvScope)
	mockArray := []string{"newKey: value", "test: true"}
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
					ServiceAccount: CopyTemplateServiceAccountName,
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
					Memory:       "nonsense value",
				},
			}),
			table.Entry("wrong number of CPU cores value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"nonsense value\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUCores:     "nonsense value",
				},
			}),
			table.Entry("wrong number of CPU Sockets value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"nonsense value\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUSockets:   "nonsense value",
				},
			}),
			table.Entry("wrong number of CPU Threads value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"nonsense value\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUThreads:   "nonsense value",
				},
			}),
		)
	})
	Context("modify template sucess", func() {
		table.DescribeTable("taskrun succeded and template is updated", func(config *testconfigs.ModifyTemplateTestConfig, checkResults func(*templatev1.Template)) {
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

			updatedTemplate, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.TODO(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(updatedTemplate).ToNot(Equal(nil), "new template should exists")
			f.ManageTemplates(updatedTemplate)

			checkResults(updatedTemplate)
		},
			table.Entry("should update template", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:            testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUTopologyNumberStr,
					CPUSockets:          CPUTopologyNumberStr,
					CPUThreads:          CPUTopologyNumberStr,
					Memory:              MemoryValue,
					TemplateLabels:      mockArray,
					TemplateAnnotations: mockArray,
					VMAnnotations:       mockArray,
					VMLabels:            mockArray,
				},
			}, func(t *templatev1.Template) {
				checkValuesInMap(t.Labels)
				checkValuesInMap(t.Annotations)

				vms, err := decodeVM(t)
				Expect(err).ShouldNot(HaveOccurred())
				for _, vm := range vms {
					checkValuesInMap(vm.Labels)
					checkValuesInMap(vm.Annotations)

					Expect(vm.Spec.Template.Spec.Domain.CPU.Sockets).To(Equal(CPUTopologyNumber), "cpu sockets should equal")
					Expect(vm.Spec.Template.Spec.Domain.CPU.Cores).To(Equal(CPUTopologyNumber), "cpu cores should equal")
					Expect(vm.Spec.Template.Spec.Domain.CPU.Threads).To(Equal(CPUTopologyNumber), "cpu threads should equal")
					memory := resource.MustParse(MemoryValue)
					Expect(vm.Spec.Template.Spec.Domain.Resources.Requests.Memory()).To(Equal(&memory), "memory should equal")
				}
			}),
		)
	})
})

func checkValuesInMap(mapToCheck map[string]string) {
	Expect(mapToCheck["newKey"]).To(Equal("value"))
	Expect(mapToCheck["test"]).To(Equal("true"))
}

func decodeVM(template *templatev1.Template) ([]*kubevirtv1.VirtualMachine, error) {
	var vms []*kubevirtv1.VirtualMachine

	for _, obj := range template.Objects {
		decoder := kubevirtv1.Codecs.UniversalDecoder(kubevirtv1.GroupVersion)
		decoded, err := runtime.Decode(decoder, obj.Raw)
		if err != nil {
			return nil, err
		}
		vm, ok := decoded.(*kubevirtv1.VirtualMachine)
		if ok {
			vms = append(vms, vm)
			break
		}
	}
	if len(vms) == 0 {
		return nil, zerrors.NewMissingRequiredError("no VM object found in the template")
	}
	return vms, nil
}
