package test

import (
	"context"

	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Copy template task", func() {
	f := framework.NewFramework().LimitEnvScope(OKDEnvScope)

	testoptions.InitTestOptions(f.TestOptions)

	cirrosTemplate := testtemplate.NewCirrosServerTinyTemplate().Build()
	cirrosTemplate.Namespace = f.TestOptions.TestNamespace

	newTemplateName := "test-template"
	openshiftNamespace := "openshift"

	Context("copy template fail", func() {

		table.DescribeTable("taskrun fails and no template is created", func(config *testconfigs.CopyTemplateTestConfig) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				template, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(template)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)
		},
			table.Entry("no target template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "target-template-name param has to be specified",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      cirrosTemplate.Name,
					SourceTemplateNamespace: cirrosTemplate.Namespace,
					TargetTemplateNamespace: cirrosTemplate.Namespace,
				},
			}),
			table.Entry("no target template namespace specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "target-template-namespace param has to be specified",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      cirrosTemplate.Name,
					SourceTemplateNamespace: cirrosTemplate.Namespace,
					TargetTemplateName:      newTemplateName,
				},
			}),
			table.Entry("no service account", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "cannot get resource \"templates\" in API group \"template.openshift.io\"",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      cirrosTemplate.Name,
					SourceTemplateNamespace: cirrosTemplate.Namespace,
					TargetTemplateName:      newTemplateName,
					TargetTemplateNamespace: cirrosTemplate.Namespace,
				},
			}),
			table.Entry("no source template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "source-template-name param has to be specified",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateNamespace: cirrosTemplate.Namespace,
					TargetTemplateName:      newTemplateName,
					TargetTemplateNamespace: cirrosTemplate.Namespace,
				},
			}),
			table.Entry("no source template namespace specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "source-template-namespace param has to be specified",
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					SourceTemplateName:      cirrosTemplate.Name,
					TargetTemplateName:      newTemplateName,
					TargetTemplateNamespace: cirrosTemplate.Namespace,
				},
			}),
			table.Entry("[NAMESPACE SCOPED] cannot copy template in different namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io is forbidden",
					LimitTestScope: NamespaceTestScope,
				},

				TaskData: testconfigs.CopyTemplateTaskData{
					Template:                cirrosTemplate,
					SourceTemplateName:      cirrosTemplate.Name,
					SourceTemplateNamespace: cirrosTemplate.Namespace,
					TargetTemplateName:      newTemplateName,
					TargetTemplateNamespace: openshiftNamespace,
				},
			}),
		)
	})
	Context("copy template sucess", func() {
		table.DescribeTable("taskrun succeded and template is created", func(config *testconfigs.CopyTemplateTestConfig, targetNamespace string) {
			f.TestSetup(config)
			if template := config.TaskData.Template; template != nil {
				f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      newTemplateName,
					"namespace": targetNamespace,
				})

			newTemplate, err := f.TemplateClient.Templates(config.TaskData.TargetTemplateNamespace).Get(context.TODO(), newTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(Equal(nil), "new template should exists")

			f.ManageTemplates(newTemplate)
		},
			table.Entry("should create template in the same namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					Template:                cirrosTemplate,
					SourceTemplateName:      cirrosTemplate.Name,
					SourceTemplateNamespace: cirrosTemplate.Namespace,
					TargetTemplateName:      newTemplateName,
					TargetTemplateNamespace: f.TestNamespace,
				},
			}, cirrosTemplate.Namespace,
			),
			table.Entry("should create template in different namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					Template:                cirrosTemplate,
					SourceTemplateName:      cirrosTemplate.Name,
					SourceTemplateNamespace: cirrosTemplate.Namespace,
					TargetTemplateName:      newTemplateName,
					TargetTemplateNamespace: openshiftNamespace,
				},
			}, openshiftNamespace),
		)
	})
})
