package output_test

import (
	"io/ioutil"
	"os"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestStruct struct {
	Data string
}

var _ = Describe("Output", func() {
	DescribeTable("Checks output type", func(value string, expectedValue bool) {
		Expect(output.IsOutputType(value)).To(Equal(expectedValue))

	},
		Entry("empty", "", true),
		Entry("yaml", "yaml", true),
		Entry("json", "json", true),
		Entry("invalid", "invalid", false),
	)

	DescribeTable("pretty prints", func(toPrint interface{}, outputType output.OutputType, expectedValue string) {
		r, w, _ := os.Pipe()
		tmp := os.Stdout
		defer func() {
			os.Stdout = tmp
		}()
		os.Stdout = w
		go func() {
			output.PrettyPrint(toPrint, outputType)
			w.Close()
		}()
		stdout, _ := ioutil.ReadAll(r)
		Expect(string(stdout)).To(Equal(expectedValue))

	},
		Entry("empty", []string{"test"}, output.OutputType(""), ""),
		Entry("array to yaml", []string{"test"}, output.YamlOutput, "- test\n"),
		Entry("array to json", []string{"test"}, output.JsonOutput, "[\n    \"test\"\n]\n"),
		Entry("struct to yaml", TestStruct{Data: "test"}, output.YamlOutput, "Data: test\n"),
		Entry("struct to json", TestStruct{Data: "test"}, output.JsonOutput, "{\n    \"Data\": \"test\"\n}\n"),
	)
})
