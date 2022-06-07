package test

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
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
			Entry("no source template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "source-template-name param has to be specified",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					TargetTemplateName: NewTemplateName,
				},
			}),
			Entry("source template doesn't exist", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io \"cirros-vm-template\" not found",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
				},
			}),
			Entry("[NAMESPACE SCOPED] cannot copy template in different namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io is forbidden",
					LimitTestScope: NamespaceTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					TargetTemplateNamespace: SystemTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
		)
	})
	Context("copy template sucess", func() {
		DescribeTable("taskrun succeded and template is created", func(config *testconfigs.CopyTemplateTestConfig) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
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

			newTemplate, err := f.TemplateClient.Templates(resultTemplateNamespace).Get(context.TODO(), resultTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(BeNil(), "new template should exists")

			f.ManageTemplates(newTemplate)
		},
			Entry("should create template in the same namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					TargetTemplateName: NewTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("should create template in different namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					TargetTemplateName:      NewTemplateName,
					TargetTemplateNamespace: DeployTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("no target template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("no target namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					TargetTemplateNamespace: EmptyTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("no source namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: EmptyTargetNS,
					TemplateNamespace:       DeployTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("no namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: EmptyTargetNS,
					TargetTemplateNamespace: EmptyTargetNS,
					TemplateNamespace:       DeployTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
		)
	})

	Context("Remove common templates labels / annotations", func() {
		It("taskrun succeeds and template is updated", func() {
			config := &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.RhelTemplateName,
					TargetTemplateName:      NewTemplateName,
					TargetTemplateNamespace: TestTargetNS,
					Template:                testtemplate.NewRhelDesktopTinyTemplate().Build(),
				},
			}
			f.TestSetup(config)

			t, err := f.TemplateClient.Templates(config.TaskData.SourceNamespace).Create(context.TODO(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			r := runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...)

			results := r.GetResults()
			resultTemplateName := results["name"]
			resultTemplateNamespace := results["namespace"]

			newTemplate, err := f.TemplateClient.Templates(resultTemplateNamespace).Get(context.TODO(), resultTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(BeNil(), " template should exists")
			f.ManageTemplates(newTemplate)
			//check template type
			Expect(newTemplate.Labels[templates.TemplateTypeLabel]).To(Equal("vm"), "template type should equal VM")

			checkRemovedRecordsTemplate(newTemplate.Labels)
			checkRemovedRecordsTemplate(newTemplate.Annotations)

			vm, _, err := zutils.DecodeVM(newTemplate)
			Expect(err).ToNot(HaveOccurred())

			checkRemovedRecordsVM(vm.Spec.Template.ObjectMeta.Labels)
			checkRemovedRecordsVM(vm.Spec.Template.ObjectMeta.Annotations)
			Expect(vm.Labels[templates.VMTemplateNameLabel]).To(Equal(newTemplate.Name), "template name should be changed")
		})
	})

	Context("Allow replace", func() {
		It("taskrun fails and new template is not created", func() {
			config := &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
					ExpectedLogs:   "templates.template.openshift.io \"test-template\" already exists",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					TargetTemplateName:      NewTemplateName,
					TargetTemplateNamespace: TestTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}
			f.TestSetup(config)

			t, err := f.TemplateClient.Templates(config.TaskData.SourceNamespace).Create(context.TODO(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			//create template which has the same name as template which will be created
			config.TaskData.Template.Name = NewTemplateName
			t, err = f.TemplateClient.Templates(string(config.TaskData.SourceNamespace)).Create(context.TODO(), config.TaskData.Template, v1.CreateOptions{})
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
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					TargetTemplateName:      NewTemplateName,
					TargetTemplateNamespace: TestTargetNS,
					AllowReplace:            "true",
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}
			f.TestSetup(config)

			t, err := f.TemplateClient.Templates(config.TaskData.SourceNamespace).Create(context.TODO(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			//create template which has the same name as template which will be created
			config.TaskData.Template.Name = NewTemplateName
			config.TaskData.Template.Objects = []runtime.RawExtension{}
			t, err = f.TemplateClient.Templates(string(config.TaskData.SourceNamespace)).Create(context.TODO(), config.TaskData.Template, v1.CreateOptions{})

			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(t)

			r := runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...)

			results := r.GetResults()
			resultTemplateName := results["name"]
			resultTemplateNamespace := results["namespace"]

			newTemplate, err := f.TemplateClient.Templates(resultTemplateNamespace).Get(context.TODO(), resultTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(BeNil(), " template should exists")
			Expect(len(newTemplate.Objects)).To(Equal(1), "template should be updated")

			f.ManageTemplates(newTemplate)
		})
	})
})

func checkRemovedRecordsTemplate(obj map[string]string) {

	for record, _ := range obj {

		if strings.HasPrefix(record, templates.TemplateOsLabelPrefix) {
			Expect(record).To(Equal(""), fmt.Sprintf("there should be no %s labels", templates.TemplateOsLabelPrefix))
		}

		if strings.HasPrefix(record, templates.TemplateFlavorLabelPrefix) {
			Expect(record).To(Equal(""), fmt.Sprintf("there should be no %s labels", templates.TemplateFlavorLabelPrefix))
		}

		if strings.HasPrefix(record, templates.TemplateWorkloadLabelPrefix) {
			Expect(record).To(Equal(""), fmt.Sprintf("there should be no %s labels", templates.TemplateWorkloadLabelPrefix))
		}
	}
	Expect(obj[templates.TemplateVersionLabel]).To(Equal(""))
	Expect(obj[templates.TemplateDeprecatedAnnotation]).To(Equal(""))
	Expect(obj[templates.KubevirtDefaultOSVariant]).To(Equal(""))

	Expect(obj[templates.OpenshiftDocURL]).To(Equal(""))
	Expect(obj[templates.OpenshiftProviderDisplayName]).To(Equal(""))
	Expect(obj[templates.OpenshiftSupportURL]).To(Equal(""))

	Expect(obj[templates.TemplateKubevirtProvider]).To(Equal(""))
	Expect(obj[templates.TemplateKubevirtProviderSupportLevel]).To(Equal(""))
	Expect(obj[templates.TemplateKubevirtProviderURL]).To(Equal(""))

	Expect(obj[templates.OperatorSDKPrimaryResource]).To(Equal(""))
	Expect(obj[templates.OperatorSDKPrimaryResourceType]).To(Equal(""))

	Expect(obj[templates.AppKubernetesComponent]).To(Equal(""))
	Expect(obj[templates.AppKubernetesName]).To(Equal(""))
	Expect(obj[templates.AppKubernetesPartOf]).To(Equal(""))
	Expect(obj[templates.AppKubernetesVersion]).To(Equal(""))
	Expect(obj[templates.AppKubernetesManagedBy]).To(Equal(""))
}

func checkRemovedRecordsVM(obj map[string]string) {
	Expect(obj[templates.VMFlavorAnnotation]).To(Equal(""))
	Expect(obj[templates.VMOSAnnotation]).To(Equal(""))
	Expect(obj[templates.VMWorkloadAnnotation]).To(Equal(""))
	Expect(obj[templates.VMDomainLabel]).To(Equal(""))
	Expect(obj[templates.VMSizeLabel]).To(Equal(""))
	Expect(obj[templates.VMTemplateRevisionLabel]).To(Equal(""))
	Expect(obj[templates.VMTemplateVersionLabel]).To(Equal(""))
}
