package parse_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	"go.uber.org/zap/zapcore"

	"reflect"
)

var (
	defaultNS       = "default"
	defaultNSArr    = []string{defaultNS}
	multipleNSArr   = []string{"overriden-ns", defaultNS}
	script          = "#!/bin/bash\necho hello world"
	commandArr      = []string{"echo", "-E", "hello", "world"}
	expectedCommand = "echo -E hello world"
)

var _ = Describe("CLIOptions", func() {
	table.DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init().Error()).To(ContainSubstring(expectedErrMessage))
	},
		table.Entry("namespace missing", "vm-namespace option is empty", &parse.CLIOptions{}),
		table.Entry("no script or command", "command-args|script option is required", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
		}),
		table.Entry("script and command", "one of command|script options is allowed", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			Command:                  commandArr,
		}),
	)
	//
	table.DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		table.Entry("returns valid defaults", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  script,
			"GetDebugLevel":              zapcore.InfoLevel,
		}),
		table.Entry("handles multiple ns from cli", &parse.CLIOptions{
			VirtualMachineNamespaces: multipleNSArr,
			Script:                   script,
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
		}),
		table.Entry("handles Script cli arguments", &parse.CLIOptions{
			VirtualMachineName:       "vm",
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			Debug:                    true,
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  script,
			"GetDebugLevel":              zapcore.DebugLevel,
		}),
		table.Entry("handles simple Command cli arguments", &parse.CLIOptions{
			VirtualMachineName:       "vm",
			VirtualMachineNamespaces: defaultNSArr,
			Command:                  []string{"ls"},
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  "ls",
		}),
		table.Entry("handles Command cli arguments", &parse.CLIOptions{
			VirtualMachineName:       "vm",
			VirtualMachineNamespaces: defaultNSArr,
			Command:                  commandArr,
			Debug:                    true,
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  expectedCommand,
			"GetDebugLevel":              zapcore.DebugLevel,
		}),
	)

})
