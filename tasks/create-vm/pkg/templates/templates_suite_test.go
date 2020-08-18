package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utilstest"
)

func TestTemplates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Templates Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
var _ = AfterSuite(utilstest.TearDownSuite)
