package zutils_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
)

var _ = Describe("Utils", func() {
	Describe("IsTrue", func() {
		table.DescribeTable("returns false", func(tested string) {
			Expect(zutils.IsTrue(tested)).To(BeFalse())
		},
			table.Entry("False", "false"),
			table.Entry("Bad", "falzee"),
			table.Entry("UpperCase", "FALSE"),
			table.Entry("Partially UpperCase", "FAlsE"),
		)

		table.DescribeTable("returns true", func(tested string) {
			Expect(zutils.IsTrue(tested)).To(BeTrue())
		},
			table.Entry("True", "true"),
			table.Entry("UpperCase", "TRUE"),
			table.Entry("Partially UpperCase", "True"),
		)
	})
})
