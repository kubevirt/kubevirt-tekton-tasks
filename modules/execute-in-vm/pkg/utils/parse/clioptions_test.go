package parse_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"time"

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
		table.Entry("no script or command", "no action was specified: at least one of the following options is required: command|script|stop|delete", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
		}),
		table.Entry("script and command", "one of command|script options is allowed", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			Command:                  commandArr,
		}),
		table.Entry("no connection secret", "connection secret should not be empty", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
		}),
		table.Entry("empty connection secret", "connection secret should not be empty", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			ConnectionSecretName:     "__empty__",
		}),
		table.Entry("invalid connection secret", "connection secret does not have a valid name", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			ConnectionSecretName:     "secret!",
		}),
		table.Entry("invalid timeout", "could not parse timeout: time: unknown unit", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			Timeout:                  "1h5q",
			ConnectionSecretName:     "my-secret",
		}),
		table.Entry("invalid stop", "invalid option stop stahp, only true|false is allowed", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			Stop:                     "stahp",
			ConnectionSecretName:     "my-secret",
		}),
		table.Entry("invalid delete", "invalid option delete yes, only true|false is allowed", &parse.CLIOptions{
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			Delete:                   "yes",
			ConnectionSecretName:     "my-secret",
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
			ConnectionSecretName:     "my-secret",
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  script,
			"GetDebugLevel":              zapcore.InfoLevel,
			"GetScriptTimeout":           0 * time.Second,
			"ShouldStop":                 false,
			"ShouldDelete":               false,
		}),
		table.Entry("handles multiple ns from cli", &parse.CLIOptions{
			VirtualMachineNamespaces: multipleNSArr,
			Script:                   script,
			ConnectionSecretName:     "my-secret",
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
		}),
		table.Entry("handles Script cli arguments", &parse.CLIOptions{
			VirtualMachineName:       "vm",
			VirtualMachineNamespaces: defaultNSArr,
			Script:                   script,
			Debug:                    true,
			Timeout:                  "5m10s",
			Stop:                     "true",
			Delete:                   "false",
			ConnectionSecretName:     "my-secret",
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  script,
			"GetDebugLevel":              zapcore.DebugLevel,
			"GetScriptTimeout":           5*time.Minute + 10*time.Second,
			"ShouldStop":                 true,
			"ShouldDelete":               false,
		}),
		table.Entry("handles simple Command cli arguments", &parse.CLIOptions{
			VirtualMachineName:       "vm",
			VirtualMachineNamespaces: defaultNSArr,
			Command:                  []string{"ls"},
			ConnectionSecretName:     "my-secret",
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  "ls",
		}),
		table.Entry("handles Command cli arguments", &parse.CLIOptions{
			VirtualMachineName:       "vm",
			VirtualMachineNamespaces: defaultNSArr,
			Command:                  commandArr,
			Debug:                    true,
			Timeout:                  "12h5m10s",
			Stop:                     "true",
			Delete:                   "true",
			ConnectionSecretName:     "my-secret",
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetScript":                  expectedCommand,
			"GetDebugLevel":              zapcore.DebugLevel,
			"GetScriptTimeout":           12*time.Hour + 5*time.Minute + 10*time.Second,
			"ShouldStop":                 true,
			"ShouldDelete":               true,
		}),
	)

})
