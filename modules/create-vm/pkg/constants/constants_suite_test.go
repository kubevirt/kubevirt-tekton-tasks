package constants_test

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConstants(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Constants Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
var _ = AfterSuite(utilstest.TearDownSuite)
