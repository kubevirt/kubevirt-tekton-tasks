package test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/clients"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/utils"
	. "github.com/onsi/ginkgo/v2"
	ginkgo_reporters "github.com/onsi/ginkgo/v2/reporters"
	v1 "github.com/openshift/api/template/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1reporter "kubevirt.io/client-go/reporter"
	qe_reporters "kubevirt.io/qe-tools/pkg/ginkgo-reporters"

	. "github.com/onsi/gomega"
)

var afterSuiteReporters []Reporter

var _ = ReportAfterSuite("Tests", func(report Report) {
	for _, reporter := range afterSuiteReporters {
		ginkgo_reporters.ReportViaDeprecatedReporter(reporter, report)
	}
})

func TestExit(t *testing.T) {
	RegisterFailHandler(Fail)
	BuildTestSuite()

	if qe_reporters.JunitOutput != "" {
		afterSuiteReporters = append(afterSuiteReporters, v1reporter.NewV1JUnitReporter(qe_reporters.JunitOutput))
	}

	if qe_reporters.Polarion.Run {
		afterSuiteReporters = append(afterSuiteReporters, &qe_reporters.Polarion)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Functional test suite")
}

func BuildTestSuite() {
	BeforeSuite(func() {
		err := testoptions.InitTestOptions(framework.TestOptionsInstance)
		noErr(err)
		err = clients.InitClients(framework.ClientsInstance, framework.TestOptionsInstance)
		noErr(err)

		if framework.TestOptionsInstance.EnvScope == constants.OKDEnvScope {
			templateList, err := framework.ClientsInstance.TemplateClient.Templates("openshift").List(context.Background(), metav1.ListOptions{
				LabelSelector: "template.kubevirt.io/type=base",
			})
			noErr(err)

			framework.TestOptionsInstance.CommonTemplatesVersion = getCommonTemplatesVersion(templateList)
		}
	})
}

func getCommonTemplatesVersion(templateList *v1.TemplateList) string {
	var commonTemplatesVersion []int
	found := false
	requiredTemplate := "fedora-server-medium"

	for _, template := range templateList.Items {
		if strings.HasPrefix(template.Name, requiredTemplate) {
			found = true
			parts := strings.Split(template.Name, fmt.Sprintf("%v-v", requiredTemplate))
			if len(parts) == 2 {
				nextVersion, err := utils.ConvertStringSliceToInt(strings.Split(parts[1], "."))
				noErr(err)
				if utils.IsBVersionHigher(commonTemplatesVersion, nextVersion) {
					commonTemplatesVersion = nextVersion
				}
			} else {
				// no version suffix
				commonTemplatesVersion = nil
				break
			}
		}
	}

	if len(commonTemplatesVersion) == 0 {
		if found {
			return "" // no version suffix
		}
		Expect(templateList).ShouldNot(BeNil())
		Fail(fmt.Sprintf("Could not compute common templates version. Number of found templates = %v", len(templateList.Items)))
	}

	return fmt.Sprintf("-v%v", utils.JoinIntSlice(commonTemplatesVersion, "."))
}

func noErr(err error) {
	if err != nil {
		Fail(err.Error())
	}
}
