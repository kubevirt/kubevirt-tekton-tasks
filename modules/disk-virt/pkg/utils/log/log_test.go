package log_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/log"
)

var _ = Describe("Log", func() {

	It("creates and returns logger", func() {
		Expect(log.GetLogger()).To(BeNil())
		first := log.InitLogger(zapcore.DebugLevel)
		Expect(first).ToNot(BeNil())
		Expect(log.InitLogger(zapcore.DebugLevel)).To(Equal(first))
		Expect(log.GetLogger()).ToNot(BeNil())
		Expect(log.GetLogger()).To(Equal(first))
	})
})
