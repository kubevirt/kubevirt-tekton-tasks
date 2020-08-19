package vm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest"
)

func TestVm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vm Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
var _ = AfterSuite(utilstest.TearDownSuite)
