package execute_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestParse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Execute Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
