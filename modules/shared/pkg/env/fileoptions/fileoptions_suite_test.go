package fileoptions_test

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/internal/intutilstest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFileoptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fileoptions Suite")
}

var _ = BeforeSuite(intutilstest.SetupTestSuite)
var _ = AfterSuite(intutilstest.TearDownSuite)
