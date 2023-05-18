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
		Entry("no mode", "only one of vm-manifest, template-name or virtctl should be specified", &parse.CLIOptions{}),
		Entry("multiple modes", "only one of vm-manifest, template-name or virtctl should be specified", &parse.CLIOptions{
			TemplateName:           "test",
			VirtualMachineManifest: testVMManifest,
		}),
		Entry("useless template ns", "template-namespace, template-params options are not applicable for vm-manifest", &parse.CLIOptions{
			VirtualMachineManifest: testVMManifest,
			TemplateNamespace:      defaultNS,
		}),
		Entry("useless template params", "template-namespace, template-params options are not applicable for vm-manifest", &parse.CLIOptions{
			VirtualMachineManifest: testVMManifest,
			TemplateParams:         []string{"K1:V1"},
		}),
		Entry("invalidManifest", "could not read VM manifest", &parse.CLIOptions{
			VirtualMachineManifest: "blabla",
		}),
		Entry("invalid output", "not a valid output type", &parse.CLIOptions{
			TemplateName: "test",
			Output:       "incorrect-fmt",
		}),
		Entry("invalid template params 1", "invalid template-params: no key found before \"V1\"; pair should be in \"KEY:VAL\" format", &parse.CLIOptions{
			TemplateName:   "test",
			TemplateParams: []string{"V1", "K2=V2"},
		}),
		Entry("invalid template params 2", "invalid template-params: no key found before \":V1\"; pair should be in \"KEY:VAL\" format", &parse.CLIOptions{
			TemplateName:   "test",
			TemplateParams: []string{":V1"},
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
			TemplateName:            "test",
			TemplateNamespace:       defaultNS,
			VirtualMachineNamespace: defaultNS,
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtualMachineManifest":  "",
			"GetTemplateParams":          map[string]string{},
			"GetDebugLevel":              zapcore.InfoLevel,
			"GetCreationMode":            constants.TemplateCreationMode,
			"GetStartVMFlag":             false,
			"GetRunStrategy":             "",
		}),
		Entry("handles template cli arguments", &parse.CLIOptions{
			TemplateName:            "test",
			TemplateNamespace:       defaultNS,
			TemplateParams:          []string{"K1:V1", "with", "space", "K2:V2"},
			VirtualMachineNamespace: defaultNS,
			Output:                  output.YamlOutput, // check if passes validation
			Debug:                   true,
			StartVM:                 "true",
			RunStrategy:             "Always",
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtualMachineManifest":  "",
			"GetTemplateParams": map[string]string{
				"K1": "V1 with space",
				"K2": "V2",
			},
			"GetDebugLevel":   zapcore.DebugLevel,
			"GetCreationMode": constants.TemplateCreationMode,
			"GetStartVMFlag":  true,
			"GetRunStrategy":  "Always",
		}),
		Entry("handles vm cli arguments", &parse.CLIOptions{
			VirtualMachineManifest:  testVMManifest,
			VirtualMachineNamespace: defaultNS,
			Output:                  output.YamlOutput, // check if passes validation
			Debug:                   true,
			StartVM:                 "false",
			RunStrategy:             "Always",
		}, map[string]interface{}{
			"GetTemplateNamespace":       "",
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtualMachineManifest":  testVMManifest,
			"GetTemplateParams":          map[string]string{},
			"GetDebugLevel":              zapcore.DebugLevel,
			"GetCreationMode":            constants.VMManifestCreationMode,
			"GetStartVMFlag":             false,
			"GetRunStrategy":             "Always",
		}),
		Entry("handles trim", &parse.CLIOptions{
			TemplateName:            "test",
			TemplateNamespace:       "  " + defaultNS + " ",
			VirtualMachineNamespace: defaultNS + "  ",
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
		}),
	)

})
