package parse_test

import (
	"reflect"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

var (
	defaultNS      = "default"
	testVMManifest = testobjects.NewTestVM().ToString()
	trueVar        = true
	falseVar       = false
)

var _ = Describe("CLIOptions", func() {
	table.DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		err := options.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
	},
		table.Entry("no mode", "one of vm-manifest, template-name should be specified", &parse.CLIOptions{}),
		table.Entry("both modes", "only one of vm-manifest, template-name should be specified", &parse.CLIOptions{
			TemplateName:           "test",
			VirtualMachineManifest: testVMManifest,
		}),
		table.Entry("useless template ns", "template-namespace, template-params options are not applicable for vm-manifest", &parse.CLIOptions{
			VirtualMachineManifest: testVMManifest,
			TemplateNamespace:      defaultNS,
		}),
		table.Entry("useless template params", "template-namespace, template-params options are not applicable for vm-manifest", &parse.CLIOptions{
			VirtualMachineManifest: testVMManifest,
			TemplateParams:         []string{"K1:V1"},
		}),
		table.Entry("invalidManifest", "could not read VM manifest", &parse.CLIOptions{
			VirtualMachineManifest: "blabla",
		}),
		table.Entry("invalid output", "not a valid output type", &parse.CLIOptions{
			TemplateName: "test",
			Output:       "incorrect-fmt",
		}),
		table.Entry("invalid template params 1", "invalid template-params: no key found before \"V1\"; pair should be in \"KEY:VAL\" format", &parse.CLIOptions{
			TemplateName:   "test",
			TemplateParams: []string{"V1", "K2=V2"},
		}),
		table.Entry("invalid template params 2", "invalid template-params: no key found before \":V1\"; pair should be in \"KEY:VAL\" format", &parse.CLIOptions{
			TemplateName:   "test",
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
			TemplateName:            "test",
			TemplateNamespace:       defaultNS,
			VirtualMachineNamespace: defaultNS,
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtualMachineManifest":  "",
			"GetPVCNames":                []string(nil),
			"GetOwnPVCNames":             []string(nil),
			"GetDVNames":                 []string(nil),
			"GetOwnDVNames":              []string(nil),
			"GetPVCDiskNamesMap":         map[string]string{},
			"GetDVDiskNamesMap":          map[string]string{},
			"GetTemplateParams":          map[string]string{},
			"GetDebugLevel":              zapcore.InfoLevel,
			"GetCreationMode":            constants.TemplateCreationMode,
			"GetStartVMFlag":             falseVar,
		}),
		table.Entry("handles template cli arguments", &parse.CLIOptions{
			TemplateName:              "test",
			TemplateNamespace:         defaultNS,
			TemplateParams:            []string{"K1:V1", "with", "space", "K2:V2"},
			VirtualMachineNamespace:   defaultNS,
			Output:                    output.YamlOutput, // check if passes validation
			PersistentVolumeClaims:    []string{"pvc1"},
			OwnPersistentVolumeClaims: []string{"mydisk1:pvc2", "pvc3"},
			DataVolumes:               []string{"dv1", "mydisk2:dv2"},
			OwnDataVolumes:            []string{"mydisk3:dv3", "dv4", "mydisk4:dv5"},
			Debug:                     true,
			StartVM:                   "true",
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtualMachineManifest":  "",
			"GetPVCNames":                []string{"pvc1"},
			"GetOwnPVCNames":             []string{"pvc2", "pvc3"},
			"GetDVNames":                 []string{"dv1", "dv2"},
			"GetOwnDVNames":              []string{"dv3", "dv4", "dv5"},
			"GetPVCDiskNamesMap": map[string]string{
				"pvc1":    "pvc1",
				"mydisk1": "pvc2",
				"pvc3":    "pvc3",
			},
			"GetDVDiskNamesMap": map[string]string{
				"dv1":     "dv1",
				"mydisk2": "dv2",
				"mydisk3": "dv3",
				"dv4":     "dv4",
				"mydisk4": "dv5",
			},
			"GetTemplateParams": map[string]string{
				"K1": "V1 with space",
				"K2": "V2",
			},
			"GetDebugLevel":   zapcore.DebugLevel,
			"GetCreationMode": constants.TemplateCreationMode,
			"GetStartVMFlag":  trueVar,
		}),
		table.Entry("handles vm cli arguments", &parse.CLIOptions{
			VirtualMachineManifest:    testVMManifest,
			VirtualMachineNamespace:   defaultNS,
			Output:                    output.YamlOutput, // check if passes validation
			PersistentVolumeClaims:    []string{"mydisk1:pvc1"},
			OwnPersistentVolumeClaims: []string{"pvc2", "pvc3"},
			DataVolumes:               []string{"mydisk2:dv1", ":dv2"},
			OwnDataVolumes:            []string{"dv3"},
			Debug:                     true,
			StartVM:                   "false",
		}, map[string]interface{}{
			"GetTemplateNamespace":       "",
			"GetVirtualMachineNamespace": defaultNS,
			"GetVirtualMachineManifest":  testVMManifest,
			"GetPVCNames":                []string{"pvc1"},
			"GetOwnPVCNames":             []string{"pvc2", "pvc3"},
			"GetDVNames":                 []string{"dv1", "dv2"},
			"GetOwnDVNames":              []string{"dv3"},
			"GetPVCDiskNamesMap": map[string]string{
				"mydisk1": "pvc1",
				"pvc2":    "pvc2",
				"pvc3":    "pvc3",
			},
			"GetDVDiskNamesMap": map[string]string{
				"mydisk2": "dv1",
				"dv2":     "dv2",
				"dv3":     "dv3",
			},
			"GetTemplateParams": map[string]string{},
			"GetDebugLevel":     zapcore.DebugLevel,
			"GetCreationMode":   constants.VMManifestCreationMode,
			"GetStartVMFlag":    falseVar,
		}),
		table.Entry("handles trim", &parse.CLIOptions{
			TemplateName:              "test",
			TemplateNamespace:         "  " + defaultNS + " ",
			VirtualMachineNamespace:   defaultNS + "  ",
			PersistentVolumeClaims:    []string{"pvc1 "},
			OwnPersistentVolumeClaims: []string{" pvc2", " pvc3  "},
			DataVolumes:               []string{" dv1", "dv2"},
			OwnDataVolumes:            []string{" dv3     "},
		}, map[string]interface{}{
			"GetTemplateNamespace":       defaultNS,
			"GetVirtualMachineNamespace": defaultNS,
			"GetPVCNames":                []string{"pvc1"},
			"GetOwnPVCNames":             []string{"pvc2", "pvc3"},
			"GetDVNames":                 []string{"dv1", "dv2"},
			"GetOwnDVNames":              []string{"dv3"},
			"GetPVCDiskNamesMap": map[string]string{
				"pvc1": "pvc1",
				"pvc2": "pvc2",
				"pvc3": "pvc3",
			},
			"GetDVDiskNamesMap": map[string]string{
				"dv1": "dv1",
				"dv2": "dv2",
				"dv3": "dv3",
			},
		}),
	)

})
