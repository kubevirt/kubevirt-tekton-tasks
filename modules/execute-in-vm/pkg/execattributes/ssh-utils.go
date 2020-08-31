package execattributes

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"regexp"
	"strconv"
)

const (
	defaultSSHPort = 22
)

var portRegex = regexp.MustCompile(`(^|\s)-p\s*([^\s]*)`)

func parsePort(sshOptions string) (int, error) {
	groups := portRegex.FindStringSubmatch(sshOptions)

	if groups == nil {
		return defaultSSHPort, nil
	}
	portStr := groups[2]

	if portStr == "" {
		return 0, zerrors.NewMissingRequiredError("ssh option requires an argument -- p")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, zerrors.NewMissingRequiredError("Bad port '%v'", portStr)
	}
	return port, nil
}
