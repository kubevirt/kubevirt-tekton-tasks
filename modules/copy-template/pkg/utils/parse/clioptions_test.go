package parse_test

import (
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/utils/parse"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

const (
	testString = "test"
)

var _ = Describe("CLIOptions", func() {
	Context("invalid cli options", func() {
		table.DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
			err := options.Init()
			Expect(err).Should(HaveOccurred())
			fmt.Println(err.Error())
			Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
		},
			table.Entry("no source-template-name", "source-template-name param has to be specified", &parse.CLIOptions{}),
			table.Entry("no source-template-namespace", "source-template-namespace param has to be specified", &parse.CLIOptions{SourceTemplateName: testString}),
			table.Entry("no target-template-name", "target-template-name param has to be specified", &parse.CLIOptions{SourceTemplateName: testString, SourceTemplateNamespace: testString}),
			table.Entry("no target-template-namespace", "target-template-namespace param has to be specified", &parse.CLIOptions{SourceTemplateName: testString, SourceTemplateNamespace: testString, TargetTemplateName: testString}),
			table.Entry("wrong output type", "non-existing is not a valid output type", &parse.CLIOptions{SourceTemplateName: testString, SourceTemplateNamespace: testString, TargetTemplateName: testString, TargetTemplateNamespace: testString, Output: "non-existing"}),
		)
	})
	Context("correct cli options", func() {
		table.DescribeTable("Init should succeed", func(options *parse.CLIOptions) {
			Expect(options.Init()).To(Succeed())
		},
			table.Entry("with yaml output", &parse.CLIOptions{
				SourceTemplateName:      testString,
				SourceTemplateNamespace: testString,
				TargetTemplateName:      testString,
				TargetTemplateNamespace: testString,
				Output:                  "yaml",
				Debug:                   true,
			}),
			table.Entry("with json output", &parse.CLIOptions{
				SourceTemplateName:      testString,
				SourceTemplateNamespace: testString,
				TargetTemplateName:      testString,
				TargetTemplateNamespace: testString,
				Output:                  "json",
				Debug:                   true,
			}),
		)

		It("Init should trim spaces", func() {
			options := &parse.CLIOptions{
				SourceTemplateName:      " " + testString + " ",
				SourceTemplateNamespace: " " + testString + " ",
				TargetTemplateName:      " " + testString + " ",
				TargetTemplateNamespace: " " + testString + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.SourceTemplateName).To(Equal(testString), "SourceTemplateName should equal")
			Expect(options.SourceTemplateNamespace).To(Equal(testString), "SourceTemplateNamespace should equal")
			Expect(options.TargetTemplateName).To(Equal(testString), "TargetTemplateName should equal")
			Expect(options.TargetTemplateNamespace).To(Equal(testString), "TargetTemplateNamespace should equal")

		})

		table.DescribeTable("CLI options should return correct values", func(fnToCall func() string, result string) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			table.Entry("GetSourceTemplateName should return correct value", (&parse.CLIOptions{SourceTemplateName: testString}).GetSourceTemplateName, testString),
			table.Entry("GetSourceTemplateNamespace should return correct value", (&parse.CLIOptions{SourceTemplateNamespace: testString}).GetSourceTemplateNamespace, testString),
			table.Entry("GetTargetTemplateNamespace should return correct value", (&parse.CLIOptions{TargetTemplateNamespace: testString}).GetTargetTemplateNamespace, testString),
			table.Entry("GetTargetTemplateName should return correct value", (&parse.CLIOptions{TargetTemplateName: testString}).GetTargetTemplateName, testString),
		)

		table.DescribeTable("CLI options should return correct log level", func(options *parse.CLIOptions, level zapcore.Level) {
			Expect(options.GetDebugLevel()).To(Equal(level), "level should equal")
		},
			table.Entry("GetDebugLevel should return correct debug level", (&parse.CLIOptions{Debug: true}), zapcore.DebugLevel),
			table.Entry("GetDebugLevel should return correct info level", (&parse.CLIOptions{Debug: false}), zapcore.InfoLevel),
		)
	})
})
