package parse_test

import (
	"reflect"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utilstest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

var (
	defaultNS = "default"
)

var _ = Describe("CLIOptions", func() {
	DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init().Error()).To(ContainSubstring(expectedErrMessage))
	},
		Entry("empty vmi name", "vmi-name should not be empty", &parse.CLIOptions{}),
		Entry("invalid vmi name", "invalid vmi-name value: a lowercase RFC 1123 subdomain must consist of", &parse.CLIOptions{
			VirtualMachineInstanceName: "invalid name",
		}),
		Entry("invalid vm namespace", "invalid vmi-namespace value: a lowercase RFC 1123 subdomain must consist of", &parse.CLIOptions{
			VirtualMachineInstanceName:      "test",
			VirtualMachineInstanceNamespace: "@ns",
		}),
		Entry("invalid success condition", "success-condition: could not parse condition", &parse.CLIOptions{
			VirtualMachineInstanceName:      "test",
			VirtualMachineInstanceNamespace: defaultNS,
			SuccessCondition:                "invalid#$%^$&",
		}),
		Entry("invalid success condition jsonpath", "success-condition: invalid condition: cannot parse jsonpath", &parse.CLIOptions{
			VirtualMachineInstanceName:      "test",
			VirtualMachineInstanceNamespace: defaultNS,
			SuccessCondition:                "test.....test",
		}),
		Entry("invalid failure condition", "failure-condition: could not parse condition", &parse.CLIOptions{
			VirtualMachineInstanceName:      "test",
			VirtualMachineInstanceNamespace: defaultNS,
			FailureCondition:                "invalid#$%^$&",
		}),
		Entry("invalid failure condition jsonpath", "failure-condition: invalid condition: cannot parse jsonpath", &parse.CLIOptions{
			VirtualMachineInstanceName:      "test",
			VirtualMachineInstanceNamespace: defaultNS,
			FailureCondition:                "test.....test",
		}),
	)

	DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		Entry("returns valid defaults", &parse.CLIOptions{
			VirtualMachineInstanceName:      "test",
			VirtualMachineInstanceNamespace: defaultNS,
		}, map[string]interface{}{
			"GetVirtualMachineInstanceName":      "test",
			"GetVirtualMachineInstanceNamespace": defaultNS,
			"GetSuccessCondition":                "",
			"GetFailureCondition":                "",
			"GetSuccessRequirements":             labels.Requirements(nil),
			"GetFailureRequirements":             labels.Requirements(nil),
			"GetDebugLevel":                      zapcore.InfoLevel,
		}),
		Entry("handles cli arguments + trim", &parse.CLIOptions{
			VirtualMachineInstanceName:      " test  ",
			VirtualMachineInstanceNamespace: "  " + defaultNS,
			SuccessCondition:                " metadata.name in (fedora, ubuntu), status.phase == Succeeded  ",
			FailureCondition:                " status.phase in (Failed, Unknown)",
			Debug:                           true,
		}, map[string]interface{}{
			"GetVirtualMachineInstanceName":      "test",
			"GetVirtualMachineInstanceNamespace": defaultNS,
			"GetSuccessCondition":                "metadata.name in (fedora, ubuntu), status.phase == Succeeded",
			"GetFailureCondition":                "status.phase in (Failed, Unknown)",
			"GetSuccessRequirements": labels.Requirements{
				utilstest.GetRequirement("metadata.name", selection.In, []string{"fedora", "ubuntu"}),
				utilstest.GetRequirement("status.phase", selection.DoubleEquals, []string{"Succeeded"}),
			},
			"GetFailureRequirements": labels.Requirements{
				utilstest.GetRequirement("status.phase", selection.In, []string{"Failed", "Unknown"}),
			},
			"GetDebugLevel": zapcore.DebugLevel,
		}),
	)
})
