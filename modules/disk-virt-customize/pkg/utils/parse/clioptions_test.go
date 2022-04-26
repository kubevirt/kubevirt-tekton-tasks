package parse_test

import (
	"reflect"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/parse"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

const (
	customizeCommands = `update
install make,ansible
delete /var/cache/dnf
`
)

var _ = Describe("CLIOptions", func() {
	DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init().Error()).To(ContainSubstring(expectedErrMessage))
	},
		Entry("no customize commands", "customize-commands option or CUSTOMIZE_COMMANDS env variable is required", &parse.CLIOptions{}),
	)
	DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		Entry("returns valid defaults", &parse.CLIOptions{
			CustomizeCommands: "test",
		}, map[string]interface{}{
			"GetCustomizeCommands":              "test",
			"GetAdditionalVirtCustomizeOptions": "",
			"GetDebugLevel":                     zapcore.InfoLevel,
			"IsVerbose":                         false,
		}),
		Entry("handles cli arguments", &parse.CLIOptions{
			CustomizeCommands:              customizeCommands,
			AdditionalVirtCustomizeOptions: "--smp 4 --memsize 2048 -q -v",
			Verbose:                        "true",
		}, map[string]interface{}{
			"GetCustomizeCommands":              customizeCommands,
			"GetAdditionalVirtCustomizeOptions": "--smp 4 --memsize 2048 -q -v",
			"GetDebugLevel":                     zapcore.DebugLevel,
			"IsVerbose":                         true,
		}),
	)
})
