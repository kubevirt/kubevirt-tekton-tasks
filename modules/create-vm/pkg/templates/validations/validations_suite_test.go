package validations_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest"
)

func TestValidations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validations Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
