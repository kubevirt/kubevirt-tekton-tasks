package env_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const (
	existingVar    = "EXISTING_VAR"
	nonExistingVar = "NON_EXISTING_VAR"
)

var _ = Describe("Env", func() {
	AfterEach(func() {
		UnSetEnv(existingVar)
		UnSetEnv(nonExistingVar)
	})

	Describe("Identifies ENV flags", func() {
		It("should identify missing flag", func() {
			Expect(env.IsEnvVarTrue(nonExistingVar)).To(BeFalse())
		})

		table.DescribeTable("should identify false flag", func(tested string) {
			SetEnv(existingVar, tested)
			Expect(env.IsEnvVarTrue(existingVar)).To(BeFalse())
		},
			table.Entry("False", "false"),
			table.Entry("Bad", "falzee"),
			table.Entry("UpperCase", "FALSE"),
			table.Entry("Partially UpperCase", "FAlsE"),
		)

		table.DescribeTable("should identify true flag", func(tested string) {
			SetEnv(existingVar, tested)
			Expect(env.IsEnvVarTrue(existingVar)).To(BeTrue())
		},
			table.Entry("True", "true"),
			table.Entry("UpperCase", "TRUE"),
			table.Entry("Partially UpperCase", "True"),
		)
	})

	It("should return tekton dir", func() {
		Expect(env.GetTektonResultsDir()).To(Equal("/tekton/results"))
	})
})
