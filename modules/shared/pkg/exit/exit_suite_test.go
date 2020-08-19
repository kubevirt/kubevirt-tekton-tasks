package exit_test

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Exit Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
var _ = AfterSuite(utilstest.TearDownSuite)
