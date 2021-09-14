package test

import (
	"context"

	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Copy template task", func() {
	f := framework.NewFramework().LimitEnvScope(OKDEnvScope)

	Context("copy template fail", func() {

		table.DescribeTable("taskrun fails and no template is created", func(config *testconfigs.CopyTemplateTestConfig) {
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
			table.Entry("no source template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "source-template-name param has to be specified",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					TargetTemplateName: NewTemplateName,
				},
			}),
			table.Entry("source template doesn't exist", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io \"cirros-vm-template\" not found",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
				},
			}),
			table.Entry("no service account", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "cannot get resource \"templates\" in API group \"template.openshift.io\"",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
				},
			}),
			table.Entry("[NAMESPACE SCOPED] cannot copy template in different namespace", &testconfigs.CopyTemplateTestConfig{
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
		table.DescribeTable("taskrun succeded and template is created", func(config *testconfigs.CopyTemplateTestConfig) {
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
			table.Entry("should create template in the same namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					TargetTemplateName: NewTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("should create template in different namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					TargetTemplateName:      NewTemplateName,
					TargetTemplateNamespace: TestTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("no target template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName: testtemplate.CirrosTemplateName,
					Template:           testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("no target namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					TargetTemplateNamespace: EmptyTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("no source namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: EmptyTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			table.Entry("no namespaces specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: EmptyTargetNS,
					TargetTemplateNamespace: EmptyTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
		)
	})
})
