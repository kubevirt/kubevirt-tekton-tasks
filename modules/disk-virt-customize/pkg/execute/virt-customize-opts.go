package execute

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
)

func SetupVirtCustomizeOptions(opts *options.CommandOptions, inputCliOptions *parse.CLIOptions) {
	if inputCliOptions.IsVerbose() {
		if !opts.IncludesOption("-v") && !opts.IncludesOption("--verbose") {
			opts.AddFlag("--verbose")
		}

		if !opts.IncludesOption("-x") {
			opts.AddFlag("-x")
		}
	}
}
