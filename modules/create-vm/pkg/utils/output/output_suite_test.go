package output_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOutput(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Output Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
