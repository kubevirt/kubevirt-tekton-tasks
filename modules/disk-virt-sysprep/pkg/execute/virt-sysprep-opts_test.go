package execute_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-sysprep/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-sysprep/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetupVirtSysprepOptions", func() {
	table.DescribeTable("sets options correctly", func(inputCliOptions *parse.CLIOptions, expected []string) {
		Expect(inputCliOptions.Init()).Should(Succeed())

		opts, err := options.NewCommandOptions(inputCliOptions.AdditionalVirtSysprepOptions)
		Expect(err).Should(Succeed())

		execute.SetupVirtSysprepOptions(opts, inputCliOptions)
		Expect(opts.GetAll()).Should(Equal(expected))
	},
		table.Entry("empty", &parse.CLIOptions{
			SysprepCommands: "update",
		},
			[]string{},
		),
		table.Entry("verbose false does not change anything verbose cli arguments", &parse.CLIOptions{
			SysprepCommands:              "update",
			AdditionalVirtSysprepOptions: "--network --dry-run",
			Verbose:                      "false",
		}, []string{
			"--network", "--dry-run",
		}),
		table.Entry("verbose adds only one argument", &parse.CLIOptions{
			SysprepCommands:              "update",
			AdditionalVirtSysprepOptions: "--network --dry-run -q -v",
			Verbose:                      "true",
		}, []string{
			"--network",
			"--dry-run",
			"-q",
			"-v",
			"-x",
		}),
		table.Entry("verbose adds both arguments", &parse.CLIOptions{
			SysprepCommands: "update",
			Verbose:         "true",
		}, []string{
			"--verbose",
			"-x",
		}),
	)
})
