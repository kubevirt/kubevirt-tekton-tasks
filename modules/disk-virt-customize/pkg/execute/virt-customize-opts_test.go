package execute_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetupVirtCustomizeOptions", func() {
	DescribeTable("sets options correctly", func(inputCliOptions *parse.CLIOptions, expected []string) {
		Expect(inputCliOptions.Init()).Should(Succeed())

		opts, err := options.NewCommandOptions(inputCliOptions.AdditionalVirtCustomizeOptions)
		Expect(err).Should(Succeed())

		execute.SetupVirtCustomizeOptions(opts, inputCliOptions)
		Expect(opts.GetAll()).Should(Equal(expected))
	},
		Entry("empty", &parse.CLIOptions{
			CustomizeCommands: "update",
		},
			[]string{},
		),
		Entry("verbose false does not change anything verbose cli arguments", &parse.CLIOptions{
			CustomizeCommands:              "update",
			AdditionalVirtCustomizeOptions: "--smp 4",
			Verbose:                        "false",
		}, []string{
			"--smp", "4",
		}),
		Entry("verbose adds only one argument", &parse.CLIOptions{
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
		Entry("verbose adds both arguments", &parse.CLIOptions{
			CustomizeCommands: "update",
			Verbose:           "true",
		}, []string{
			"--verbose",
			"-x",
		}),
	)
})
