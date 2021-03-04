package generate

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants/connectionsecret"
)

func ensureComment(opts *options.CommandOptions, cliOptions *parse.CLIOptions) {
	// comment
	if !opts.IncludesOption("-C") {
		connectionOptions := cliOptions.GetPrivateKeyConnectionOptions()
		user := "default"

		if u := connectionOptions[connectionsecret.SSHConnectionSecretKeys.User]; u != "" {
			user = u
		}

		opts.AddOption("-C", fmt.Sprintf("%v@generated", user))
	}
}
func setDefaultOptions(opts *options.CommandOptions) {
	// type of key
	if !opts.IncludesOption("-t") {
		opts.AddOption("-t", "rsa")
	}

	// number of bits in the key
	if opts.GetOptionValue("-t") == "rsa" && !opts.IncludesOption("-b") {
		opts.AddOption("-b", "4096")
	}

	// new passphrase
	if !opts.IncludesOption("-N") {
		opts.AddOption("-N", "")
	}
}
