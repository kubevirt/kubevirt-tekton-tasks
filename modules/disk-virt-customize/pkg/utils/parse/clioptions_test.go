package parse_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/parse"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"reflect"
)

const (
	customizeCommands = `update
install make,ansible
delete /var/cache/dnf
`
)

var _ = Describe("CLIOptions", func() {
	table.DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init().Error()).To(ContainSubstring(expectedErrMessage))
	},
		table.Entry("no customize commands", "customize-commands option or CUSTOMIZE_COMMANDS env variable is required", &parse.CLIOptions{}),
	)
	table.DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		table.Entry("returns valid defaults", &parse.CLIOptions{
			CustomizeCommands: "test",
		}, map[string]interface{}{
			"GetCustomizeCommands":              "test",
			"GetAdditionalVirtCustomizeOptions": "",
			"GetDebugLevel":                     zapcore.InfoLevel,
			"IsVerbose":                         false,
		}),
		table.Entry("handles cli arguments", &parse.CLIOptions{
			CustomizeCommands:              customizeCommands,
			AdditionalVirtCustomizeOptions: "--smp 4 --memsize 2048 -q -v",
			Verbose:                        "true",
		}, map[string]interface{}{
			"GetCustomizeCommands":              customizeCommands,
			"GetAdditionalVirtCustomizeOptions": "--smp 4 --memsize 2048 -q -v -x",
			"GetDebugLevel":                     zapcore.DebugLevel,
			"IsVerbose":                         true,
		}),
		table.Entry("adds verbose cli arguments", &parse.CLIOptions{
			CustomizeCommands: customizeCommands,
			Verbose:           "true",
		}, map[string]interface{}{
			"GetCustomizeCommands":              customizeCommands,
			"GetAdditionalVirtCustomizeOptions": "--verbose -x",
			"GetDebugLevel":                     zapcore.DebugLevel,
			"IsVerbose":                         true,
		}),
	)

	It("does common operations correctly", func() {
		cliOptions := &parse.CLIOptions{
			CustomizeCommands:              customizeCommands,
			AdditionalVirtCustomizeOptions: "-v --let-x 4",
			Verbose:                        "false",
		}

		err := cliOptions.Init()
		Expect(err).Should(Succeed())

		// IncludesSSHOption
		Expect(cliOptions.IncludesVirtCustomizeOption("-v")).To(BeTrue())
		Expect(cliOptions.IncludesVirtCustomizeOption("-x")).To(BeFalse())

		// IncludesSSHOption and AddAdditionalSSHOption
		cliOptions.AddAdditionalVirtCustomizeOption("-x")
		Expect(cliOptions.IncludesVirtCustomizeOption("-v")).To(BeTrue())
		Expect(cliOptions.IncludesVirtCustomizeOption("-x")).To(BeTrue())
		Expect(cliOptions.GetAdditionalVirtCustomizeOptions()).Should(Equal("-v --let-x 4 -x"))
	})

})
