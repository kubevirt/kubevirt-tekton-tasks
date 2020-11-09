package test

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/utils"
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
	RunSpecs(t, "E2E Tests Suite")
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
		if strings.HasPrefix(template.Name, requiredTemplate) {
			parts := strings.Split(template.Name, fmt.Sprintf("%v-v", requiredTemplate))
			if len(parts) == 2 {
				nextVersion, err := utils.ConvertStringSliceToInt(strings.Split(parts[1], "."))
				noErr(err)
				if utils.IsBVersionHigher(commonTemplatesVersion, nextVersion) {
					commonTemplatesVersion = nextVersion
				}
			}
		}
	}

	if len(commonTemplatesVersion) == 0 {
		Fail("Could not compute common templates version")
	}

	return fmt.Sprintf("v%v", utils.JoinIntSlice(commonTemplatesVersion, "."))
}

func noErr(err error) {
	if err != nil {
		Fail(err.Error())
	}
}
