package parse_test

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/parse"
)

const (
	sysprepCommands = `update
install make,ansible
operations firewall-rules,defaults
`
)

var _ = Describe("CLIOptions", func() {
	DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init().Error()).To(ContainSubstring(expectedErrMessage))
	},
		Entry("no sysprep commands", "virt-commands option or VIRT_COMMANDS env variable is required", &parse.CLIOptions{}),
	)
	DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		Entry("returns valid defaults", &parse.CLIOptions{
			Commands: "test",
		}, map[string]interface{}{
			"GetCommands":              "test",
			"GetAdditionalVirtOptions": "",
			"GetDebugLevel":            zapcore.InfoLevel,
			"IsVerbose":                false,
		}),
		Entry("handles cli arguments", &parse.CLIOptions{
			Commands:              sysprepCommands,
			AdditionalVirtOptions: "--network --dry-run -q -v",
			Verbose:               "true",
		}, map[string]interface{}{
			"GetCommands":              sysprepCommands,
			"GetAdditionalVirtOptions": "--network --dry-run -q -v",
			"GetDebugLevel":            zapcore.DebugLevel,
			"IsVerbose":                true,
		}),
	)
})
