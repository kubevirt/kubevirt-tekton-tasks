package output_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/output"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

type TestStruct struct {
	Data string
}

var _ = Describe("Output", func() {
	table.DescribeTable("Checks output type", func(value string, expectedValue bool) {
		Expect(output.IsOutputType(value)).To(Equal(expectedValue))

	},
		table.Entry("empty", "", true),
		table.Entry("yaml", "yaml", true),
		table.Entry("json", "json", true),
		table.Entry("invalid", "invalid", false),
	)

	table.DescribeTable("pretty prints", func(toPrint interface{}, outputType output.OutputType, expectedValue string) {
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
		table.Entry("empty", []string{"test"}, output.OutputType(""), ""),
		table.Entry("array to yaml", []string{"test"}, output.YamlOutput, "- test\n"),
		table.Entry("array to json", []string{"test"}, output.JsonOutput, "[\n    \"test\"\n]\n"),
		table.Entry("struct to yaml", TestStruct{Data: "test"}, output.YamlOutput, "Data: test\n"),
		table.Entry("struct to json", TestStruct{Data: "test"}, output.JsonOutput, "{\n    \"Data\": \"test\"\n}\n"),
	)
})
