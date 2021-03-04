package execattributes

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"strconv"
)

const (
	defaultSSHPort = 22
)

func parsePort(sshOptions *options.CommandOptions) (int, error) {
	if !sshOptions.IncludesOption("-p") {
		return defaultSSHPort, nil
	}

	portStr := sshOptions.GetOptionValue("-p")

	if portStr == "" {
		return 0, zerrors.NewMissingRequiredError("ssh option requires an argument -- p")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, zerrors.NewMissingRequiredError("Bad port '%v'", portStr)
	}
	return port, nil
}
