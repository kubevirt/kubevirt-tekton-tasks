package parse_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"reflect"
)

var (
	defaultNS     = "default"
	defaultNSArr  = []string{defaultNS}
	multipleNSArr = []string{"overriden-ns", defaultNS}
)

var _ = Describe("CLIOptions", func() {
	table.DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		err := options.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
	},
		table.Entry("invalid output", "not a valid output type", &parse.CLIOptions{
			Output: "incorrect-fmt",
		}),
		table.Entry("invalid template params 1", "parameters have incorrect format", &parse.CLIOptions{
			TemplateParams: []string{"K1:V1", "K2=V2"},
		}),
		table.Entry("invalid template params 2", "parameters have incorrect format", &parse.CLIOptions{
			TemplateParams: []string{":V1"},
		}),
	)

	table.DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		table.Entry("returns valid defaults", &parse.CLIOptions{
			TemplateNamespaces:       defaultNSArr,
			VirtualMachineNamespaces: defaultNSArr,
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
			"GetAllPVCNames":             []string(nil),
			"GetAllDVNames":              []string(nil),
			"GetAllDiskNames":            []string(nil),
			"GetTemplateParams":          map[string]string{},
			"GetDebugLevel":              zapcore.InfoLevel,
		}),
		table.Entry("handles multiple ns from cli", &parse.CLIOptions{
			TemplateNamespaces:       multipleNSArr,
			VirtualMachineNamespaces: multipleNSArr,
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
		}),
		table.Entry("handles cli arguments", &parse.CLIOptions{
			TemplateName:              "test",
			TemplateNamespaces:        defaultNSArr,
			TemplateParams:            []string{"K1:V1", "K2:V2"},
			VirtualMachineNamespaces:  defaultNSArr,
			Output:                    output.YamlOutput, // check if passes validation
			OwnDataVolumes:            []string{"dv1"},
			DataVolumes:               []string{"dv2", "dv3"},
			OwnPersistentVolumeClaims: []string{"pvc1", "pvc2"},
			PersistentVolumeClaims:    []string{"pvc3"},
			Debug:                     true,
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
			"GetAllPVCNames":             []string{"pvc1", "pvc2", "pvc3"},
			"GetAllDVNames":              []string{"dv1", "dv2", "dv3"},
			"GetAllDiskNames":            []string{"pvc1", "pvc2", "pvc3", "dv1", "dv2", "dv3"},
			"GetTemplateParams": map[string]string{
				"K1": "V1",
				"K2": "V2",
			},
			"GetDebugLevel": zapcore.DebugLevel,
		}),
	)

})
