package env_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Env Suite")
}

var _ = BeforeSuite(func() {
	utilstest.SetEnv(zconstants.OutOfClusterENV, "true")
})
var _ = AfterSuite(func() {
	utilstest.UnSetEnv(zconstants.OutOfClusterENV)
})
