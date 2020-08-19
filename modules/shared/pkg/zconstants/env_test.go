package zconstants_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"os"
)

const (
	existingVar    = "EXISTING_VAR"
	nonExistingVar = "NON_EXISTING_VAR"
)

var _ = Describe("Env", func() {
	var defaultOutOfClusterENV string

	BeforeEach(func() {
		defaultOutOfClusterENV = os.Getenv(zconstants.OutOfClusterENV)
	})
	AfterEach(func() {
		UnSetEnv(existingVar)
		UnSetEnv(nonExistingVar)
		SetEnv(zconstants.OutOfClusterENV, defaultOutOfClusterENV)
	})

	Describe("Identifies ENV flags", func() {
		It("should identify missing flag", func() {
			Expect(zconstants.IsEnvVarTrue(nonExistingVar)).To(BeFalse())
		})

		table.DescribeTable("should identify false flag", func(tested string) {
			SetEnv(existingVar, tested)
			Expect(zconstants.IsEnvVarTrue(existingVar)).To(BeFalse())
		},
			table.Entry("False", "false"),
			table.Entry("Bad", "falzee"),
			table.Entry("UpperCase", "FALSE"),
			table.Entry("Partially UpperCase", "FAlsE"),
		)

		table.DescribeTable("should identify true flag", func(tested string) {
			SetEnv(existingVar, tested)
			Expect(zconstants.IsEnvVarTrue(existingVar)).To(BeTrue())
		},
			table.Entry("True", "true"),
			table.Entry("UpperCase", "TRUE"),
			table.Entry("Partially UpperCase", "True"),
		)
	})
	Describe("should lookup active namespace", func() {
		It("should fail out of cluster", func() {
			ns, err := zconstants.GetActiveNamespace()
			Expect(ns).To(BeEmpty())
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("return results tekton dir", func() {
		It("should return supplemental dir if out of cluster", func() {
			SetEnv(zconstants.OutOfClusterENV, "true")
			Expect(zconstants.GetTektonResultsDir()).To(BeADirectory())
		})
		It("should return tekton dir in a cluster", func() {
			SetEnv(zconstants.OutOfClusterENV, "false")
			Expect(zconstants.GetTektonResultsDir()).To(Equal(zconstants.TektonResultsDirPath))
		})
	})

})
