package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
)

func (c *CLIOptions) validateCommands() error {
	if c.GetCommands() == "" {
		return zerrors.NewMissingRequiredError("%v option or %v env variable is required", commandsOptionName, commandsEnvVarName)
	}
	return nil
}
