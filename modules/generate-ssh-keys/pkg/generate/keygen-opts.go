package generate

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
)

func ensureComment(opts *options.CommandOptions, cliOptions *parse.CLIOptions) {
	// comment
	if !opts.Includes("-C") {
		connectionOptions := cliOptions.GetPrivateKeyConnectionOptions()
		user := "default"

		if u := connectionOptions[constants.ConnectionOptions.User]; u != "" {
			user = u
		}

		opts.AddOpt("-C", fmt.Sprintf("%v@generated", user))
	}
}
func setDefaultOptions(opts *options.CommandOptions) {
	// type of key
	if !opts.Includes("-t") {
		opts.AddOpt("-t", "rsa")
	}

	// number of bits in the key
	if opts.GetOptionValue("-t") == "rsa" && !opts.Includes("-b") {
		opts.AddOpt("-b", "4096")
	}

	// new passphrase
	if !opts.Includes("-N") {
		opts.AddOpt("-N", "")
	}
}
