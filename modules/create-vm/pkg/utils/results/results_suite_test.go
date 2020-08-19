package results_test

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResults(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Results Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
var _ = AfterSuite(utilstest.TearDownSuite)
