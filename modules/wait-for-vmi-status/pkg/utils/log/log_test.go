package log_test

import (
	log2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

var _ = Describe("Log", func() {

	It("creates and returns logger", func() {
		Expect(log2.Logger()).To(BeNil())
		first := log2.InitLogger(zapcore.DebugLevel)
		Expect(first).ToNot(BeNil())
		Expect(log2.InitLogger(zapcore.DebugLevel)).To(Equal(first))
		Expect(log2.Logger()).ToNot(BeNil())
		Expect(log2.Logger()).To(Equal(first))
	})
})
