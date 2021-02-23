package results_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResults(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Results Suite")
}

var _ = BeforeSuite(func() {
	utilstest.SetEnv(zconstants.OutOfClusterENV, "true")
})
var _ = AfterSuite(func() {
	utilstest.UnSetEnv(zconstants.OutOfClusterENV)
})
