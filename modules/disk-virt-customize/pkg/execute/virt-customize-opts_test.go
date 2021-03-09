package execute_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetupVirtCustomizeOptions", func() {
	table.DescribeTable("sets options correctly", func(inputCliOptions *parse.CLIOptions, expected []string) {
		Expect(inputCliOptions.Init()).Should(Succeed())

		opts, err := options.NewCommandOptions(inputCliOptions.AdditionalVirtCustomizeOptions)
		Expect(err).Should(Succeed())

		execute.SetupVirtCustomizeOptions(opts, inputCliOptions)
		Expect(opts.GetAll()).Should(Equal(expected))
	},
		table.Entry("empty", &parse.CLIOptions{
			CustomizeCommands: "update",
		},
			[]string{},
		),
		table.Entry("verbose false does not change anything verbose cli arguments", &parse.CLIOptions{
			CustomizeCommands:              "update",
			AdditionalVirtCustomizeOptions: "--smp 4",
			Verbose:                        "false",
		}, []string{
			"--smp", "4",
		}),
		table.Entry("verbose adds only one argument", &parse.CLIOptions{
			CustomizeCommands:              "update",
			AdditionalVirtCustomizeOptions: "--smp 4 --memsize 2048 -q -v",
			Verbose:                        "true",
		}, []string{
			"--smp", "4",
			"--memsize", "2048",
			"-q",
			"-v",
			"-x",
		}),
		table.Entry("verbose adds both arguments", &parse.CLIOptions{
			CustomizeCommands: "update",
			Verbose:           "true",
		}, []string{
			"--verbose",
			"-x",
		}),
	)
})
