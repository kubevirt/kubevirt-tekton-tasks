package parse_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-data-object/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datasource"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	testStrDataObjectNamespace1 = "data-object-namespace-test-1"
	testStrDataObjectNamespace2 = "data-object-namespace-test-2"
	testStrTrue                 = "true"
)

var (
	testDvManifest1 = strings.TrimSpace(datavolume.NewBlankDataVolume("testDv1").ToString())
	testDvManifest2 = strings.TrimSpace(datavolume.NewBlankDataVolume("testDv2").WithNamespace(testStrDataObjectNamespace2).ToString())
	testDsManifest  = strings.TrimSpace(datasource.NewDataSource("testDs").ToString())
)

var _ = Describe("CLIOptions", func() {
	Describe("invalid cli options", func() {
		DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
			err := options.Init()
			Expect(err).Should(HaveOccurred())
			fmt.Println(err.Error())
			Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
		},
			Entry("no data-object-manifest", "data-object-manifest param has to be specified", &parse.CLIOptions{}),
			Entry("wrong output type", "non-existing is not a valid output type",
				&parse.CLIOptions{
					DataObjectManifest: testDvManifest1,
					Output:             "non-existing",
				}),
		)
	})

	Describe("correct cli options", func() {
		DescribeTable("Init should succeed with DataVolume", func(options *parse.CLIOptions) {
			Expect(options.Init()).To(Succeed())
		},
			Entry("with yaml output", &parse.CLIOptions{
				DataObjectManifest:  testDvManifest1,
				DataObjectNamespace: testStrDataObjectNamespace1,
				Output:              "yaml",
			}),
			Entry("with json output", &parse.CLIOptions{
				DataObjectManifest:  testDvManifest1,
				DataObjectNamespace: testStrDataObjectNamespace1,
				Output:              "json",
			}),
			Entry("with debug loglevel", &parse.CLIOptions{
				DataObjectManifest:  testDvManifest1,
				DataObjectNamespace: testStrDataObjectNamespace1,
				Debug:               true,
			}),
			Entry("with WaitForSuccess", &parse.CLIOptions{
				DataObjectManifest:  testDvManifest1,
				DataObjectNamespace: testStrDataObjectNamespace1,
				WaitForSuccess:      testStrTrue,
			}),
			Entry("with AllowReplace", &parse.CLIOptions{
				DataObjectManifest:  testDvManifest1,
				DataObjectNamespace: testStrDataObjectNamespace1,
				AllowReplace:        testStrTrue,
			}),
		)

		DescribeTable("Init should succeed with DataSource", func(options *parse.CLIOptions) {
			Expect(options.Init()).To(Succeed())
		},
			Entry("with yaml output", &parse.CLIOptions{
				DataObjectManifest:  testDsManifest,
				DataObjectNamespace: testStrDataObjectNamespace1,
				Output:              "yaml",
			}),
			Entry("with json output", &parse.CLIOptions{
				DataObjectManifest:  testDsManifest,
				DataObjectNamespace: testStrDataObjectNamespace1,
				Output:              "json",
			}),
			Entry("with debug loglevel", &parse.CLIOptions{
				DataObjectManifest:  testDsManifest,
				DataObjectNamespace: testStrDataObjectNamespace1,
				Debug:               true,
			}),
			Entry("with WaitForSuccess", &parse.CLIOptions{
				DataObjectManifest:  testDsManifest,
				DataObjectNamespace: testStrDataObjectNamespace1,
				WaitForSuccess:      testStrTrue,
			}),
			Entry("with AllowReplace", &parse.CLIOptions{
				DataObjectManifest:  testDsManifest,
				DataObjectNamespace: testStrDataObjectNamespace1,
				AllowReplace:        testStrTrue,
			}),
		)

		It("Init should trim spaces", func() {
			options := &parse.CLIOptions{
				DataObjectManifest:  " " + testDvManifest1 + " ",
				DataObjectNamespace: " " + testStrDataObjectNamespace1 + " ",
				WaitForSuccess:      " " + testStrTrue + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.DataObjectManifest).To(Equal(testDvManifest1), "DataObjectManifest should equal")
			Expect(options.DataObjectNamespace).To(Equal(testStrDataObjectNamespace1), "DataObjectNamespace should equal")
			Expect(options.WaitForSuccess).To(Equal(testStrTrue), "WaitForSuccess should equal")
		})

		DescribeTable("CLI options should return correct values", func(fnToCall func() string, result string) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("GetDataObjectManifest should return correct value", (&parse.CLIOptions{DataObjectManifest: testDvManifest1}).GetDataObjectManifest, testDvManifest1),
			Entry("GetSourceTemplateNamespace should return correct value", (&parse.CLIOptions{DataObjectNamespace: testStrDataObjectNamespace1}).GetDataObjectNamespace, testStrDataObjectNamespace1),
		)

		DescribeTable("GetWaitForSuccess should return correct values", func(fnToCall func() bool, result bool) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("should return correct true", (&parse.CLIOptions{WaitForSuccess: "true"}).GetWaitForSuccess, true),
			Entry("should return correct false", (&parse.CLIOptions{WaitForSuccess: "false"}).GetWaitForSuccess, false),
			Entry("should return correct false, when wrong string", (&parse.CLIOptions{WaitForSuccess: "notAValue"}).GetWaitForSuccess, false),
		)

		DescribeTable("GetAllowReplace should return correct values", func(fnToCall func() bool, result bool) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("should return correct true", (&parse.CLIOptions{AllowReplace: "true"}).GetAllowReplace, true),
			Entry("should return correct false", (&parse.CLIOptions{AllowReplace: "false"}).GetAllowReplace, false),
			Entry("should return correct false, when wrong string", (&parse.CLIOptions{AllowReplace: "notAValue"}).GetAllowReplace, false),
		)

		DescribeTable("CLI options should return correct log level", func(options *parse.CLIOptions, level zapcore.Level) {
			Expect(options.GetDebugLevel()).To(Equal(level), "level should equal")
		},
			Entry("GetDebugLevel should return correct debug level", (&parse.CLIOptions{Debug: true}), zapcore.DebugLevel),
			Entry("GetDebugLevel should return correct info level", (&parse.CLIOptions{Debug: false}), zapcore.InfoLevel),
		)

		It("Init should read the namespace from the manifest", func() {
			options := &parse.CLIOptions{
				DataObjectManifest: " " + testDvManifest2 + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.DataObjectNamespace).To(Equal(testStrDataObjectNamespace2), "DataObjectNamespace should equal")
		})

		It("Init should try to get the active namespace", func() {
			options := &parse.CLIOptions{
				DataObjectManifest: " " + testDvManifest1 + " ",
			}

			err := options.Init()
			if err == nil {
				Expect(options.DataObjectNamespace).ToNot(BeEmpty())
			} else {
				Expect(err).To(MatchError("can't get active namespace: could not detect active namespace"))
			}
		})

		It("GetUnstructuredDataObject should return correct value", func() {
			c := &parse.CLIOptions{
				DataObjectManifest:  testDsManifest,
				DataObjectNamespace: testStrDataObjectNamespace1,
			}
			err := c.Init()
			Expect(err).ShouldNot(HaveOccurred())

			unstructuredDo := unstructured.Unstructured{}
			err = yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(testDsManifest)), 1024).Decode(&unstructuredDo)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(c.GetUnstructuredDataObject()).To(Equal(unstructuredDo), "result should equal")
		})
	})
})
