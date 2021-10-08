package parse_test

import (
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-vm-template/pkg/utils/parse"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	testString         = "test"
	testNumberOfCPU    = "2"
	testNumberOfCPUInt = 2
	testStringMemory   = "124M"
)

var (
	mockArray                = []string{"newKey: value", "test: true"}
	resultMap                = map[string]string{"newKey": "value", "test": "true"}
	testStringMemoryResource = resource.MustParse(testStringMemory)
)

var _ = Describe("CLIOptions", func() {
	Context("invalid cli options", func() {
		table.DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
			err := options.Init()
			Expect(err).Should(HaveOccurred())
			fmt.Println(err.Error())
			Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
		},
			table.Entry("no template-name", "template-name param has to be specified", &parse.CLIOptions{}),
			table.Entry("wrong output type", "non-existing is not a valid output type", &parse.CLIOptions{TemplateName: testString, Output: "non-existing"}),
			table.Entry("wrong cpu sockets", "parsing \"wrong cpu sockets\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: "wrong cpu sockets"}),
			table.Entry("wrong cpu cores", "parsing \"wrong cpu cores\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: "wrong cpu cores"}),
			table.Entry("wrong cpu threads", "parsing \"wrong cpu threads\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: "wrong cpu threads"}),
			table.Entry("wrong template labels", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: []string{"singleKey"}}),
			table.Entry("wrong template annotations", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: []string{"singleKey"}}),
			table.Entry("wrong vm labels", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: []string{"singleKey"}}),
			table.Entry("wrong vm annotations", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: mockArray, VMAnnotations: []string{"singleKey"}}),
		)
	})
	Context("correct cli options", func() {
		table.DescribeTable("Init should succeed", func(options *parse.CLIOptions) {
			Expect(options.Init()).To(Succeed())
		},
			table.Entry("with yaml output", &parse.CLIOptions{
				TemplateName: testString,
				Output:       "yaml",
				Debug:        true,
			}),
			table.Entry("with json output", &parse.CLIOptions{
				TemplateName: testString,
				Output:       "json",
				Debug:        true,
			}),
		)
		It("should succeed with all options", func() {
			options := &parse.CLIOptions{
				TemplateName:        testString,
				CPUCores:            testNumberOfCPU,
				CPUThreads:          testNumberOfCPU,
				TemplateLabels:      mockArray,
				TemplateAnnotations: mockArray,
				VMLabels:            mockArray,
				VMAnnotations:       mockArray,
			}
			Expect(options.Init()).To(Succeed())
		})

		It("Init should trim spaces", func() {
			options := &parse.CLIOptions{
				TemplateName: " " + testString + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.TemplateName).To(Equal(testString), "TemplateName should equal")
		})

		table.DescribeTable("CLI options should return correct string values", func(fnToCall func() string, result string) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			table.Entry("GetTemplateName should return correct value", (&parse.CLIOptions{TemplateName: testString}).GetTemplateName, testString),
			table.Entry("GetTemplateNamespace should return correct value", (&parse.CLIOptions{TemplateNamespace: testString}).GetTemplateNamespace, testString),
		)

		table.DescribeTable("CLI options should return correct int values", func(fnToCall func() int, result int) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			table.Entry("GetCPUCores should return correct value", (&parse.CLIOptions{CPUCores: testNumberOfCPU}).GetCPUCores, testNumberOfCPUInt),
			table.Entry("GetCPUSockets should return correct value", (&parse.CLIOptions{CPUSockets: testNumberOfCPU}).GetCPUSockets, testNumberOfCPUInt),
			table.Entry("GetCPUThreads should return correct value", (&parse.CLIOptions{CPUThreads: testNumberOfCPU}).GetCPUThreads, testNumberOfCPUInt),
		)

		table.DescribeTable("CLI options should return correct Quantity values", func(fnToCall func() *resource.Quantity, result *resource.Quantity) {
			r := fnToCall()
			Expect(r).To(Equal(result), "result should equal")
		},
			table.Entry("GetMemory should return correct value", (&parse.CLIOptions{Memory: testStringMemory}).GetMemory, &testStringMemoryResource),
			table.Entry("GetMemory should return nil", (&parse.CLIOptions{}).GetMemory, nil),
		)

		table.DescribeTable("CLI options should return correct log level", func(options *parse.CLIOptions, level zapcore.Level) {
			Expect(options.GetDebugLevel()).To(Equal(level), "level should equal")
		},
			table.Entry("GetDebugLevel should return correct debug level", (&parse.CLIOptions{Debug: true}), zapcore.DebugLevel),
			table.Entry("GetDebugLevel should return correct info level", (&parse.CLIOptions{Debug: false}), zapcore.InfoLevel),
		)

		cli := &parse.CLIOptions{
			TemplateName:        testString,
			TemplateLabels:      mockArray,
			TemplateAnnotations: mockArray,
			VMLabels:            mockArray,
			VMAnnotations:       mockArray,
		}
		table.DescribeTable("CLI options should return correct map of annotations / labels", func(obj *parse.CLIOptions, fnToCall func() map[string]string, result map[string]string) {
			Expect(obj.Init()).To(Succeed(), "should succeeded")
			Expect(fnToCall()).To(Equal(result), "maps should equal")
		},
			table.Entry("GetTemplateLabels should return correct template labels", cli, cli.GetTemplateLabels, resultMap),
			table.Entry("GetTemplateAnnotations should return correct template annotations", cli, cli.GetTemplateAnnotations, resultMap),
			table.Entry("GetVMLabels should return correct VM labels", cli, cli.GetVMLabels, resultMap),
			table.Entry("GetVMAnnotations should return correct VM annotations", cli, cli.GetVMAnnotations, resultMap),
		)
	})
})
