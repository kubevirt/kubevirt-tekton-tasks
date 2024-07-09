package test

import (
	"context"
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Copy template task", func() {
	f := framework.NewFramework().LimitEnvScope(OKDEnvScope)

	Context("copy template fail", func() {

		DescribeTable("taskrun fails and no template is created", func(config *testconfigs.CopyTemplateTestConfig) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)
		},
			Entry("no source template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "source-template-name param has to be specified",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					TargetTemplateName: NewTemplateName,
				},
			}),
			Entry("source template doesn't exist", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "templates.template.openshift.io \"cirros-vm-template\" not found",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
				},
			}),
		)
	})
	Context("copy template sucess", func() {
		DescribeTable("taskrun succeded and template is created", func(config *testconfigs.CopyTemplateTestConfig) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}

			r := runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...)

			results := r.GetResults()
			resultTemplateName := results["name"]
			resultTemplateNamespace := results["namespace"]

			newTemplate, err := f.TemplateClient.Templates(resultTemplateNamespace).Get(context.Background(), resultTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(BeNil(), "new template should exists")

			if config.TaskData.SetOwnerReference == "true" {
				Expect(newTemplate.OwnerReferences).To(HaveLen(1), "template should has owner reference")
				Expect(newTemplate.OwnerReferences[0].Kind).To(Equal("Pod"), "OwnerReference should have Kind Pod")
				Expect(newTemplate.OwnerReferences[0].Name).To(HavePrefix("e2e-tests-taskrun-copy-template"), "OwnerReference should be binded to correct Pod")
			} else {
				Expect(newTemplate.OwnerReferences).To(BeEmpty(), "template OwnerReference should be empty")
			}

			f.ManageTemplates(newTemplate)
		},
			Entry("should create template in the same namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					TargetTemplateName: NewTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
					SetOwnerReference:  "true",
				},
			}),
			Entry("no target template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
					SetOwnerReference:  "false",
				},
			}),
			Entry("no target namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("no source namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("no namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
		)
	})

	Context("Remove common templates labels / annotations", func() {
		It("taskrun succeeds and template is updated", func() {
			config := &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.RhelTemplateName,
					TargetTemplateName: NewTemplateName,
					Template:           testtemplate.NewRhelDesktopTinyTemplate().Build(),
					SetOwnerReference:  "true",
				},
			}
			f.TestSetup(config)

			t, err := f.TemplateClient.Templates(config.TaskData.TemplateNamespace).Create(context.Background(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			r := runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...)

			results := r.GetResults()
			resultTemplateName := results["name"]
			resultTemplateNamespace := results["namespace"]

			newTemplate, err := f.TemplateClient.Templates(resultTemplateNamespace).Get(context.Background(), resultTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(BeNil(), " template should exists")
			f.ManageTemplates(newTemplate)
			//check template type
			Expect(newTemplate.Labels).To(HaveKeyWithValue(templates.TemplateTypeLabel, "vm"), "template type should equal VM")

			checkRemovedRecordsTemplate(newTemplate.Labels)
			checkRemovedRecordsTemplate(newTemplate.Annotations)
			Expect(newTemplate.OwnerReferences).To(HaveLen(1), "template should has owner reference")

			vm, _, err := zutils.DecodeVM(newTemplate)
			Expect(err).ToNot(HaveOccurred())

			checkRemovedRecordsVM(vm.Spec.Template.ObjectMeta.Labels)
			checkRemovedRecordsVM(vm.Spec.Template.ObjectMeta.Annotations)
			Expect(vm.Labels).To(HaveKeyWithValue(templates.VMTemplateNameLabel, newTemplate.Name), "template name should be changed")
		})
	})

	Context("Allow replace", func() {
		It("taskrun fails and new template is not created", func() {
			config := &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "templates.template.openshift.io \"test-template\" already exists",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:         testtemplate.CirrosTemplateName,
					TargetTemplateName:         NewTemplateName,
					Template:                   testtemplate.NewCirrosServerTinyTemplate().Build(),
					UsePlainTargetTemplateName: true,
				},
			}
			f.TestSetup(config)

			t, err := f.TemplateClient.Templates(config.TaskData.TemplateNamespace).Create(context.Background(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			//create template which has the same name as template which will be created
			config.TaskData.Template.Name = NewTemplateName
			t, err = f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Create(context.Background(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)
		})

		It("taskrun succeeds and template is updated", func() {
			config := &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					TargetTemplateName: NewTemplateName,
					AllowReplace:       "true",
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
					SetOwnerReference:  "true",
				},
			}
			f.TestSetup(config)

			t, err := f.TemplateClient.Templates(config.TaskData.TemplateNamespace).Create(context.Background(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			//create template which has the same name as template which will be created
			config.TaskData.Template.Name = config.TaskData.TargetTemplateName
			config.TaskData.Template.Objects = []runtime.RawExtension{}
			t, err = f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Create(context.Background(), config.TaskData.Template, v1.CreateOptions{})

			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			r := runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...)

			results := r.GetResults()
			resultTemplateName := results["name"]
			resultTemplateNamespace := results["namespace"]

			newTemplate, err := f.TemplateClient.Templates(resultTemplateNamespace).Get(context.Background(), resultTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(BeNil(), " template should exists")
			Expect(newTemplate.Objects).To(HaveLen(1), "template should be updated")
			Expect(newTemplate.OwnerReferences).To(HaveLen(1), "template should has owner reference")
			Expect(newTemplate.OwnerReferences[0].Kind).To(Equal("Pod"), "OwnerReference should have Kind Pod")
			Expect(newTemplate.OwnerReferences[0].Name).To(HavePrefix("e2e-tests-taskrun-copy-template"), "OwnerReference should be binded to correct Pod")
			f.ManageTemplates(newTemplate)
		})
	})
})

func checkRemovedRecordsTemplate(obj map[string]string) {
	for record, _ := range obj {
		Expect(record).ToNot(HavePrefix(templates.TemplateOsLabelPrefix), fmt.Sprintf("there should be no %s labels", templates.TemplateOsLabelPrefix))
		Expect(record).ToNot(HavePrefix(templates.TemplateFlavorLabelPrefix), fmt.Sprintf("there should be no %s labels", templates.TemplateFlavorLabelPrefix))
		Expect(record).ToNot(HavePrefix(templates.TemplateWorkloadLabelPrefix), fmt.Sprintf("there should be no %s labels", templates.TemplateWorkloadLabelPrefix))
	}

	Expect(obj).ToNot(HaveKey(templates.TemplateVersionLabel))
	Expect(obj).ToNot(HaveKey(templates.TemplateDeprecatedAnnotation))
	Expect(obj).ToNot(HaveKey(templates.KubevirtDefaultOSVariant))

	Expect(obj).ToNot(HaveKey(templates.OpenshiftDocURL))
	Expect(obj).ToNot(HaveKey(templates.OpenshiftProviderDisplayName))
	Expect(obj).ToNot(HaveKey(templates.OpenshiftSupportURL))

	Expect(obj).ToNot(HaveKey(templates.TemplateKubevirtProvider))
	Expect(obj).ToNot(HaveKey(templates.TemplateKubevirtProviderSupportLevel))
	Expect(obj).ToNot(HaveKey(templates.TemplateKubevirtProviderURL))

	Expect(obj).ToNot(HaveKey(templates.OperatorSDKPrimaryResource))
	Expect(obj).ToNot(HaveKey(templates.OperatorSDKPrimaryResourceType))

	Expect(obj).ToNot(HaveKey(templates.AppKubernetesComponent))
	Expect(obj).ToNot(HaveKey(templates.AppKubernetesName))
	Expect(obj).ToNot(HaveKey(templates.AppKubernetesPartOf))
	Expect(obj).ToNot(HaveKey(templates.AppKubernetesVersion))
	Expect(obj).ToNot(HaveKey(templates.AppKubernetesManagedBy))
}

func checkRemovedRecordsVM(obj map[string]string) {
	Expect(obj).ToNot(HaveKey(templates.VMFlavorAnnotation))
	Expect(obj).ToNot(HaveKey(templates.VMOSAnnotation))
	Expect(obj).ToNot(HaveKey(templates.VMWorkloadAnnotation))
	Expect(obj).ToNot(HaveKey(templates.VMDomainLabel))
	Expect(obj).ToNot(HaveKey(templates.VMSizeLabel))
	Expect(obj).ToNot(HaveKey(templates.VMTemplateRevisionLabel))
	Expect(obj).ToNot(HaveKey(templates.VMTemplateVersionLabel))
}
