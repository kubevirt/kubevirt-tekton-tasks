package zutils_test

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/internal/intutilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = BeforeSuite(intutilstest.SetupTestSuite)
var _ = AfterSuite(intutilstest.TearDownSuite)
