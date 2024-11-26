package parse_test

import (
	"reflect"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

var (
	defaultNS      = "default"
	testVMManifest = testobjects.NewTestVM().ToString()
)

var _ = Describe("CLIOptions", func() {
	DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		err := options.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
	},
		Entry("no mode", "only one of vm-manifest or virtctl should be specified", &parse.CLIOptions{}),
		Entry("useless virtctl param", "only one of vm-manifest or virtctl should be specified", &parse.CLIOptions{
			VirtualMachineManifest: testVMManifest,
			Virtctl:                "K1:V1",
		}),
		Entry("invalidManifest", "could not read VM manifest", &parse.CLIOptions{
			VirtualMachineManifest: "blabla",
		}),
		Entry("invalid output", "not a valid output type", &parse.CLIOptions{
			VirtualMachineManifest: testVMManifest,
			Output:                 "incorrect-fmt",
		}),
	)

	DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		Entry("handles virtctl cli arguments", &parse.CLIOptions{
			Virtctl:                 "test",
			VirtualMachineNamespace: defaultNS,
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtctl":                 "test",
			"GetDebugLevel":              zapcore.InfoLevel,
			"GetCreationMode":            constants.VirtctlCreatingMode,
			"GetStartVMFlag":             false,
			"GetRunStrategy":             "",
		}),
		Entry("handles vm cli arguments", &parse.CLIOptions{
			VirtualMachineManifest:  testVMManifest,
			VirtualMachineNamespace: defaultNS,
			Output:                  output.YamlOutput, // check if passes validation
			Debug:                   true,
			StartVM:                 "false",
			RunStrategy:             "Always",
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtualMachineManifest":  testVMManifest,
			"GetDebugLevel":              zapcore.DebugLevel,
			"GetCreationMode":            constants.VMManifestCreationMode,
			"GetStartVMFlag":             false,
			"GetRunStrategy":             "Always",
		}),
		Entry("handles trim", &parse.CLIOptions{
			VirtualMachineManifest:  testVMManifest,
			VirtualMachineNamespace: defaultNS + "  ",
		}, map[string]interface{}{
			"GetVirtualMachineNamespace": defaultNS,
		}),
	)

})
