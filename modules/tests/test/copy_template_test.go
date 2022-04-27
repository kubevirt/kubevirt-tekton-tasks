package test

import (
	"context"

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
			Entry("no service account", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "cannot get resource \"templates\" in API group \"template.openshift.io\"",
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
