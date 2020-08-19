package validations_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest"
)

func TestValidations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validations Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
var _ = AfterSuite(utilstest.TearDownSuite)
