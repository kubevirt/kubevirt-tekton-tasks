package output_test

import (
	"testing"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOutput(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Output Suite")
}

func SetupTestSuite() {
	log.InitLogger(zap.InfoLevel)
}

var _ = BeforeSuite(SetupTestSuite)
