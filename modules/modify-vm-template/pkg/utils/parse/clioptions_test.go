package parse_test

import (
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-vm-template/pkg/utils/parse"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	testString                = "test"
	testNumberOfCPU           = "2"
	testNumberOfCPUInt uint32 = 2
	testStringMemory          = "124M"
)

var (
	mockArray                 = []string{"newKey: value", "test: true"}
	diskArray                 = []string{"{\"name\": \"test\", \"cdrom\": {\"bus\": \"sata\"}}"}
	volumeArray               = []string{"{\"name\": \"test\", \"containerDisk\": {\"image\": \"URL\"}}"}
	templateParametersArray   = []string{"{\"description\": \"VM name\", \"name\": \"NAME\"}"}
	dataVolumeArray           = []string{"{\"apiVersion\": \"cdi.kubevirt.io/v1beta1\", \"kind\": \"DataVolume\", \"metadata\":{\"name\": \"test1\"}, \"spec\": {\"source\": {\"http\": {\"url\": \"test.somenonexisting\"}}}}"}
	resultMap                 = map[string]string{"newKey": "value", "test": "true"}
	testStringMemoryResource  = resource.MustParse(testStringMemory)
	parsedDisk                = []kubevirtv1.Disk{{Name: "test", DiskDevice: kubevirtv1.DiskDevice{CDRom: &kubevirtv1.CDRomTarget{Bus: "sata"}}}}
	parsedVolume              = []kubevirtv1.Volume{{Name: "test", VolumeSource: kubevirtv1.VolumeSource{ContainerDisk: &kubevirtv1.ContainerDiskSource{Image: "URL"}}}}
	parsedDataVolumeTemplates = []kubevirtv1.DataVolumeTemplateSpec{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DataVolume",
				APIVersion: "cdi.kubevirt.io/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test1",
			},
			Spec: cdiv1.DataVolumeSpec{
				Source: &cdiv1.DataVolumeSource{
					HTTP: &cdiv1.DataVolumeSourceHTTP{
						URL: "test.somenonexisting",
					},
				},
			},
		},
	}
)

var _ = Describe("CLIOptions", func() {
	Context("invalid cli options", func() {
		DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
			err := options.Init()
			Expect(err).Should(HaveOccurred())
			fmt.Println(err.Error())
			Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
		},
			Entry("no template-name", "template-name param has to be specified", &parse.CLIOptions{}),
			Entry("wrong output type", "non-existing is not a valid output type", &parse.CLIOptions{TemplateName: testString, Output: "non-existing"}),
			Entry("wrong cpu sockets", "parsing \"wrong cpu sockets\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: "wrong cpu sockets"}),
			Entry("wrong cpu cores", "parsing \"wrong cpu cores\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: "wrong cpu cores"}),
			Entry("wrong cpu threads", "parsing \"wrong cpu threads\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: "wrong cpu threads"}),
			Entry("negative cpu sockets", "parsing \"-2\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: "-2"}),
			Entry("negative cpu cores", "parsing \"-1\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: "-1"}),
			Entry("negative cpu threads", "parsing \"-3\": invalid syntax", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: "-3"}),
			Entry("wrong template labels", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: []string{"singleKey"}}),
			Entry("wrong template annotations", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: []string{"singleKey"}}),
			Entry("wrong vm labels", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: []string{"singleKey"}}),
			Entry("wrong vm annotations", "pair should be in \"KEY:VAL\" format", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: mockArray, VMAnnotations: []string{"singleKey"}}),
			Entry("wrong disk json", "invalid character 'w'", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: mockArray, Disks: []string{"{wrongJson: value}"}}),
			Entry("wrong volume json", "invalid character 'k'", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: mockArray, Volumes: []string{"{key: value}"}}),
			Entry("wrong dataVolumeTemplate json", "invalid character 'e' in literal true", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: mockArray, Volumes: mockArray, DatavolumeTemplates: []string{"{wrong value}"}}),
			Entry("wrong templateParameters json", "invalid character 'e' in literal true", &parse.CLIOptions{TemplateName: testString, CPUCores: testNumberOfCPU, CPUThreads: testNumberOfCPU, TemplateLabels: mockArray, TemplateAnnotations: mockArray, VMLabels: mockArray, Volumes: mockArray, DatavolumeTemplates: mockArray, TemplateParameters: []string{"{wrong value}"}}),
		)
	})
	Context("correct cli options", func() {
		DescribeTable("Init should succeed", func(options *parse.CLIOptions) {
			Expect(options.Init()).To(Succeed())
		},
			Entry("with yaml output", &parse.CLIOptions{
				TemplateName: testString,
				Output:       "yaml",
				Debug:        true,
			}),
			Entry("with json output", &parse.CLIOptions{
				TemplateName: testString,
				Output:       "json",
				Debug:        true,
			}),
			Entry("should succeed with all options", &parse.CLIOptions{
				TemplateName:             testString,
				CPUCores:                 testNumberOfCPU,
				CPUThreads:               testNumberOfCPU,
				TemplateLabels:           mockArray,
				TemplateAnnotations:      mockArray,
				VMLabels:                 mockArray,
				VMAnnotations:            mockArray,
				Disks:                    diskArray,
				Volumes:                  volumeArray,
				DeleteVolumes:            true,
				DeleteDisks:              true,
				DeleteTemplateParameters: true,
				DeleteDatavolumeTemplate: true,
				DatavolumeTemplates:      dataVolumeArray,
				TemplateParameters:       templateParametersArray,
			}),
		)

		It("Init should trim spaces", func() {
			options := &parse.CLIOptions{
				TemplateName: " " + testString + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.TemplateName).To(Equal(testString), "TemplateName should equal")
		})

		DescribeTable("CLI options should return correct string values", func(fnToCall func() string, result string) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("GetTemplateName should return correct value", (&parse.CLIOptions{TemplateName: testString}).GetTemplateName, testString),
			Entry("GetTemplateNamespace should return correct value", (&parse.CLIOptions{TemplateNamespace: testString}).GetTemplateNamespace, testString),
		)

		DescribeTable("CLI options should return correct int values", func(fnToCall func() uint32, result uint32) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("GetCPUCores should return correct value", (&parse.CLIOptions{CPUCores: testNumberOfCPU}).GetCPUCores, testNumberOfCPUInt),
			Entry("GetCPUSockets should return correct value", (&parse.CLIOptions{CPUSockets: testNumberOfCPU}).GetCPUSockets, testNumberOfCPUInt),
			Entry("GetCPUThreads should return correct value", (&parse.CLIOptions{CPUThreads: testNumberOfCPU}).GetCPUThreads, testNumberOfCPUInt),
		)

		DescribeTable("CLI options should return correct Quantity values", func(fnToCall func() *resource.Quantity, result *resource.Quantity) {
			r := fnToCall()
			Expect(r).To(Equal(result), "result should equal")
		},
			Entry("GetMemory should return correct value", (&parse.CLIOptions{Memory: testStringMemory}).GetMemory, &testStringMemoryResource),
			Entry("GetMemory should return nil", (&parse.CLIOptions{}).GetMemory, nil),
		)

		DescribeTable("CLI options should return correct log level", func(options *parse.CLIOptions, level zapcore.Level) {
			Expect(options.GetDebugLevel()).To(Equal(level), "level should equal")
		},
			Entry("GetDebugLevel should return correct debug level", (&parse.CLIOptions{Debug: true}), zapcore.DebugLevel),
			Entry("GetDebugLevel should return correct info level", (&parse.CLIOptions{Debug: false}), zapcore.InfoLevel),
		)

		cli := &parse.CLIOptions{
			TemplateName:             testString,
			TemplateLabels:           mockArray,
			TemplateAnnotations:      mockArray,
			VMLabels:                 mockArray,
			VMAnnotations:            mockArray,
			Disks:                    diskArray,
			Volumes:                  volumeArray,
			DatavolumeTemplates:      dataVolumeArray,
			DeleteDatavolumeTemplate: true,
			DeleteDisks:              true,
			DeleteVolumes:            true,
			DeleteTemplateParameters: true,
		}
		DescribeTable("CLI options should return correct map of annotations / labels", func(obj *parse.CLIOptions, fnToCall func() map[string]string, result map[string]string) {
			Expect(obj.Init()).To(Succeed(), "should succeeded")
			Expect(fnToCall()).To(Equal(result), "maps should equal")
		},
			Entry("GetTemplateLabels should return correct template labels", cli, cli.GetTemplateLabels, resultMap),
			Entry("GetTemplateAnnotations should return correct template annotations", cli, cli.GetTemplateAnnotations, resultMap),
			Entry("GetVMLabels should return correct VM labels", cli, cli.GetVMLabels, resultMap),
			Entry("GetVMAnnotations should return correct VM annotations", cli, cli.GetVMAnnotations, resultMap),
		)

		DescribeTable("CLI options should return correct Disk values", func(obj *parse.CLIOptions, fnToCall func() []kubevirtv1.Disk, result []kubevirtv1.Disk) {
			Expect(obj.Init()).To(Succeed(), "should succeeded")
			r := fnToCall()
			Expect(r[0].Name).To(Equal(result[0].Name), "disk name should equal")
			Expect(r[0].CDRom.Bus).To(Equal(result[0].CDRom.Bus), "disk bus should equal")
		},
			Entry("GetDisks should return correct value", cli, cli.GetDisks, parsedDisk),
		)

		DescribeTable("CLI options should return correct Volume values", func(obj *parse.CLIOptions, fnToCall func() []kubevirtv1.Volume, result []kubevirtv1.Volume) {
			Expect(obj.Init()).To(Succeed(), "should succeeded")
			r := fnToCall()
			Expect(r[0].Name).To(Equal(result[0].Name), "volume name should equal")
			Expect(r[0].ContainerDisk.Image).To(Equal(result[0].ContainerDisk.Image), "volume image should equal")
		},
			Entry("GetVolumes should return correct value", cli, cli.GetVolumes, parsedVolume),
		)

		DescribeTable("CLI options should return correct dataVolume templates values", func(obj *parse.CLIOptions, fnToCall func() []kubevirtv1.DataVolumeTemplateSpec, result []kubevirtv1.DataVolumeTemplateSpec) {
			Expect(obj.Init()).To(Succeed(), "should succeeded")
			r := fnToCall()
			Expect(r[0].Name).To(Equal(result[0].Name), "volume name should equal")
			Expect(r[0].Spec.Source.HTTP.URL).To(Equal(result[0].Spec.Source.HTTP.URL), "URL should equal")
		},
			Entry("GetVolumes should return correct value", cli, cli.GetDatavolumeTemplates, parsedDataVolumeTemplates),
		)

		DescribeTable("CLI options should return correct value for bool functions", func(obj *parse.CLIOptions, fnToCall func() bool, result bool) {
			Expect(obj.Init()).To(Succeed(), "should succeeded")
			Expect(fnToCall()).To(Equal(result), "bool should equal")
		},
			Entry("GetDeleteDatavolumeTemplate should return correct value", cli, cli.GetDeleteDatavolumeTemplate, true),
			Entry("GetDeleteDisks should return correct value", cli, cli.GetDeleteDisks, true),
			Entry("GetDeleteVolumes should return correct value", cli, cli.GetDeleteVolumes, true),
			Entry("GetDeleteTemplateParameters should return correct value", cli, cli.GetDeleteTemplateParameters, true),
		)
	})
})
