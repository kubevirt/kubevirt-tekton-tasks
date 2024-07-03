package execute_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
)

var _ = Describe("SetupVirtOptions", func() {
	DescribeTable("sets options correctly", func(inputCliOptions *parse.CLIOptions, expected []string) {
		Expect(inputCliOptions.Init()).Should(Succeed())

		opts, err := options.NewCommandOptions(inputCliOptions.AdditionalVirtOptions)
		Expect(err).Should(Succeed())

		execute.SetupVirtOptions(opts, inputCliOptions)
		Expect(opts.GetAll()).Should(Equal(expected))
	},
		Entry("empty", &parse.CLIOptions{
			Commands: "update",
		},
			[]string{},
		),
		Entry("verbose false does not change anything verbose cli arguments", &parse.CLIOptions{
			Commands:              "update",
			AdditionalVirtOptions: "--network --dry-run",
			Verbose:               "false",
		}, []string{
			"--network", "--dry-run",
		}),
		Entry("verbose adds only one argument", &parse.CLIOptions{
			Commands:              "update",
			AdditionalVirtOptions: "--network --dry-run -q -v",
			Verbose:               "true",
		}, []string{
			"--network",
			"--dry-run",
			"-q",
			"-v",
			"-x",
		}),
		Entry("verbose adds both arguments", &parse.CLIOptions{
			Commands: "update",
			Verbose:  "true",
		}, []string{
			"--verbose",
			"-x",
		}),
	)
})
