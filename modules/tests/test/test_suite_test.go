package test

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/utils"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	v1 "github.com/openshift/api/template/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExit(t *testing.T) {
	RegisterFailHandler(Fail)
	BuildTestSuite()
	junitReporter := reporters.NewJUnitReporter(fmt.Sprintf("../dist/junit_%d.xml", config.GinkgoConfig.ParallelNode))
	RunSpecsWithDefaultAndCustomReporters(t, "E2E Tests Suite", []Reporter{junitReporter})
}

func BuildTestSuite() {
	BeforeSuite(func() {
		err := testoptions.InitTestOptions(framework.TestOptionsInstance)
		noErr(err)
		err = framework.InitClients(framework.ClientsInstance, framework.TestOptionsInstance)
		noErr(err)

		templateList, err := framework.ClientsInstance.TemplateClient.Templates("openshift").List(metav1.ListOptions{})
		noErr(err)

		framework.TestOptionsInstance.CommonTemplatesVersion = getCommonTemplatesVersion(templateList)
	})
}

func getCommonTemplatesVersion(templateList *v1.TemplateList) string {
	var commonTemplatesVersion []int
	requiredTemplate := "fedora-server-tiny"

	for _, template := range templateList.Items {
		fmt.Printf("checking %v\n", template.Name)
		if strings.HasPrefix(template.Name, requiredTemplate) {
			fmt.Println("  hasRightPrefix")
			parts := strings.Split(template.Name, fmt.Sprintf("%v-v", requiredTemplate))
			fmt.Printf("  parts len %v\n", len(parts))
			fmt.Printf("  parts %v\n", parts)
			if len(parts) == 2 {
				nextVersion, err := utils.ConvertStringSliceToInt(strings.Split(parts[1], "."))
				noErr(err)
				fmt.Printf("  nextVersion %v\n", nextVersion)
				if utils.IsBVersionHigher(commonTemplatesVersion, nextVersion) {
					commonTemplatesVersion = nextVersion
				}
				fmt.Printf("  commonTemplatesVersion %v\n", commonTemplatesVersion)
			}
		}
	}
	fmt.Printf("  commonTemplatesVersion %v\n", commonTemplatesVersion)

	if len(commonTemplatesVersion) == 0 {
		Expect(templateList).ShouldNot(BeNil())
		Fail(fmt.Sprintf("Could not compute common templates version. Number of found templates = %v", len(templateList.Items)))
	}

	return fmt.Sprintf("v%v", utils.JoinIntSlice(commonTemplatesVersion, "."))
}

func noErr(err error) {
	if err != nil {
		Fail(err.Error())
	}
}
