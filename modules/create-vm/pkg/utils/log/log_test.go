package log_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/log"
	"go.uber.org/zap/zapcore"
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
