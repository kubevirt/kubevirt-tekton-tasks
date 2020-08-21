package results_test

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/internal/intutilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResults(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Results Suite")
}

var _ = BeforeSuite(intutilstest.SetupTestSuite)
var _ = AfterSuite(intutilstest.TearDownSuite)
