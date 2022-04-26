package env_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	. "github.com/onsi/ginkgo/v2"
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

		DescribeTable("should identify false flag", func(tested string) {
			SetEnv(existingVar, tested)
			Expect(env.IsEnvVarTrue(existingVar)).To(BeFalse())
		},
			Entry("False", "false"),
			Entry("Bad", "falzee"),
			Entry("UpperCase", "FALSE"),
			Entry("Partially UpperCase", "FAlsE"),
		)

		DescribeTable("should identify true flag", func(tested string) {
			SetEnv(existingVar, tested)
			Expect(env.IsEnvVarTrue(existingVar)).To(BeTrue())
		},
			Entry("True", "true"),
			Entry("UpperCase", "TRUE"),
			Entry("Partially UpperCase", "True"),
		)
	})

	It("should return tekton dir", func() {
		Expect(env.GetTektonResultsDir()).To(Equal("/tekton/results"))
	})
})
