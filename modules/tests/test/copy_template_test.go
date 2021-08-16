package test

import (
	"context"

	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	openshiftNamespace TargetNamespace = "openshift"
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
					DoNotSetSourceName: true,
					TargetTemplateName: NewTemplateName,
				},
			}),

			table.Entry("no service account", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "cannot get resource \"templates\" in API group \"template.openshift.io\"",
				},
				TaskData: testconfigs.CopyTemplateTaskData{},
			}),
			table.Entry("[NAMESPACE SCOPED] cannot copy template in different namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io is forbidden",
					LimitTestScope: NamespaceTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					ForbiddenNamespace: openshiftNamespace,
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
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.ResultTemplateName,
					"namespace": config.TaskData.TargetNamespace,
				})

			results := r.GetResults()
			resultTemplateName := results["name"]
			resultTemplateNamespace := results["namespace"]

			newTemplate, err := f.TemplateClient.Templates(resultTemplateNamespace).Get(context.TODO(), resultTemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newTemplate).ToNot(Equal(nil), "new template should exists")

			f.ManageTemplates(newTemplate)
		},
			table.Entry("should create template in the same namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{},
			}),
			table.Entry("should create template in different namespace", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
					LimitTestScope: ClusterTestScope,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					TargetTemplateNamespace: openshiftNamespace,
				},
			}),
		)
		table.DescribeTable("taskrun succeded without target name", func(config *testconfigs.CopyTemplateTestConfig) {
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
			Expect(newTemplate).ToNot(Equal(nil), "new template should exists")

			f.ManageTemplates(newTemplate)
		},
			table.Entry("no target template name specified", &testconfigs.CopyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CopyTemplateServiceAccountName,
				},
				TaskData: testconfigs.CopyTemplateTaskData{
					DoNotSetTargetName: true,
				},
			}),
		)
	})
})
