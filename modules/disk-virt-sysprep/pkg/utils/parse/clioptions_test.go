package parse_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-sysprep/pkg/utils/parse"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"reflect"
)

const (
	sysprepCommands = `update
install make,ansible
operations firewall-rules,defaults
`
)

var _ = Describe("CLIOptions", func() {
	table.DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init().Error()).To(ContainSubstring(expectedErrMessage))
	},
		table.Entry("no sysprep commands", "sysprep-commands option or SYSPREP_COMMANDS env variable is required", &parse.CLIOptions{}),
	)
	table.DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		table.Entry("returns valid defaults", &parse.CLIOptions{
			SysprepCommands: "test",
		}, map[string]interface{}{
			"GetSysprepCommands":              "test",
			"GetAdditionalVirtSysprepOptions": "",
			"GetDebugLevel":                   zapcore.InfoLevel,
			"IsVerbose":                       false,
		}),
		table.Entry("handles cli arguments", &parse.CLIOptions{
			SysprepCommands:              sysprepCommands,
			AdditionalVirtSysprepOptions: "--network --dry-run -q -v",
			Verbose:                      "true",
		}, map[string]interface{}{
			"GetSysprepCommands":              sysprepCommands,
			"GetAdditionalVirtSysprepOptions": "--network --dry-run -q -v",
			"GetDebugLevel":                   zapcore.DebugLevel,
			"IsVerbose":                       true,
		}),
	)
})
