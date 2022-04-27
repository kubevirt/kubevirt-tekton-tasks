package parse_test

import (
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/utils/parse"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

const (
	testStringSourceName      = "source-name-test"
	testStringSourceNamespace = "source-namespace-test"
	testStringTargetName      = "target-name-test"
	testStringTargetNamespace = "target-namespace-test"
)

var _ = Describe("CLIOptions", func() {
	Context("invalid cli options", func() {
		DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
			err := options.Init()
			Expect(err).Should(HaveOccurred())
			fmt.Println(err.Error())
			Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
		},
			Entry("no source-template-name", "source-template-name param has to be specified", &parse.CLIOptions{}),
			Entry("wrong output type", "non-existing is not a valid output type",
				&parse.CLIOptions{
					SourceTemplateName:      testStringSourceName,
					SourceTemplateNamespace: testStringSourceNamespace,
					TargetTemplateName:      testStringTargetName,
					TargetTemplateNamespace: testStringTargetNamespace,
					Output:                  "non-existing",
				}),
		)
	})
	Context("correct cli options", func() {
		DescribeTable("Init should succeed", func(options *parse.CLIOptions) {
			Expect(options.Init()).To(Succeed())
		},
			Entry("with yaml output", &parse.CLIOptions{
				SourceTemplateName:      testStringSourceName,
				SourceTemplateNamespace: testStringSourceNamespace,
				TargetTemplateName:      testStringTargetName,
				TargetTemplateNamespace: testStringTargetNamespace,
				Output:                  "yaml",
				Debug:                   true,
			}),
			Entry("with json output", &parse.CLIOptions{
				SourceTemplateName:      testStringSourceName,
				SourceTemplateNamespace: testStringSourceNamespace,
				TargetTemplateName:      testStringTargetName,
				TargetTemplateNamespace: testStringTargetNamespace,
				Output:                  "json",
				Debug:                   true,
			}),
			Entry("with AllowReplace", &parse.CLIOptions{
				SourceTemplateName:      testStringSourceName,
				SourceTemplateNamespace: testStringSourceNamespace,
				TargetTemplateName:      testStringTargetName,
				TargetTemplateNamespace: testStringTargetNamespace,
				AllowReplace:            "false",
				Output:                  "json",
				Debug:                   true,
			}),
			Entry("with AllowReplace", &parse.CLIOptions{
				SourceTemplateName:      testStringSourceName,
				SourceTemplateNamespace: testStringSourceNamespace,
				TargetTemplateName:      testStringTargetName,
				TargetTemplateNamespace: testStringTargetNamespace,
				AllowReplace:            "true",
				Output:                  "json",
				Debug:                   true,
			}),
			Entry("no source-template-namespace", &parse.CLIOptions{SourceTemplateName: testStringSourceName}),
			Entry("no target-template-name", &parse.CLIOptions{
				SourceTemplateName:      testStringSourceName,
				SourceTemplateNamespace: testStringSourceNamespace,
			}),
			Entry("no target-template-namespace", &parse.CLIOptions{
				SourceTemplateName:      testStringSourceName,
				SourceTemplateNamespace: testStringSourceNamespace,
				TargetTemplateName:      testStringTargetName,
			}),
		)

		It("Init should trim spaces", func() {
			options := &parse.CLIOptions{
				SourceTemplateName:      " " + testStringSourceName + " ",
				SourceTemplateNamespace: " " + testStringSourceNamespace + " ",
				TargetTemplateName:      " " + testStringTargetName + " ",
				TargetTemplateNamespace: " " + testStringTargetNamespace + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.SourceTemplateName).To(Equal(testStringSourceName), "SourceTemplateName should equal")
			Expect(options.SourceTemplateNamespace).To(Equal(testStringSourceNamespace), "SourceTemplateNamespace should equal")
			Expect(options.TargetTemplateName).To(Equal(testStringTargetName), "TargetTemplateName should equal")
			Expect(options.TargetTemplateNamespace).To(Equal(testStringTargetNamespace), "TargetTemplateNamespace should equal")

		})

		DescribeTable("CLI options should return correct values", func(fnToCall func() string, result string) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("GetSourceTemplateName should return correct value", (&parse.CLIOptions{SourceTemplateName: testStringSourceName}).GetSourceTemplateName, testStringSourceName),
			Entry("GetSourceTemplateNamespace should return correct value", (&parse.CLIOptions{SourceTemplateNamespace: testStringSourceNamespace}).GetSourceTemplateNamespace, testStringSourceNamespace),
			Entry("GetTargetTemplateNamespace should return correct value", (&parse.CLIOptions{TargetTemplateNamespace: testStringTargetNamespace}).GetTargetTemplateNamespace, testStringTargetNamespace),
			Entry("GetTargetTemplateName should return correct value", (&parse.CLIOptions{TargetTemplateName: testStringTargetName}).GetTargetTemplateName, testStringTargetName),
		)

		DescribeTable("GetAllowReplaceValue should return correct values", func(fnToCall func() bool, result bool) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("should return correct true", (&parse.CLIOptions{AllowReplace: "true"}).GetAllowReplaceValue, true),
			Entry("should return correct false", (&parse.CLIOptions{AllowReplace: "false"}).GetAllowReplaceValue, false),
			Entry("should return correct false, when wrong string", (&parse.CLIOptions{AllowReplace: "notAValue"}).GetAllowReplaceValue, false),
		)

		DescribeTable("CLI options should return correct log level", func(options *parse.CLIOptions, level zapcore.Level) {
			Expect(options.GetDebugLevel()).To(Equal(level), "level should equal")
		},
			Entry("GetDebugLevel should return correct debug level", (&parse.CLIOptions{Debug: true}), zapcore.DebugLevel),
			Entry("GetDebugLevel should return correct info level", (&parse.CLIOptions{Debug: false}), zapcore.InfoLevel),
		)
	})
})
