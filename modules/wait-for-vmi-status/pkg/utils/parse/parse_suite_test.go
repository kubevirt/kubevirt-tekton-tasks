package parse_test

import (
	"testing"

	log2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/log"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestParse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parse Suite")
}

var _ = BeforeSuite(func() {
	log2.InitLogger(zap.DebugLevel)
})
