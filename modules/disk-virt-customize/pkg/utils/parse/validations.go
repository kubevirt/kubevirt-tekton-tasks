package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
)

func (c *CLIOptions) validateCommands() error {
	if c.GetCustomizeCommands() == "" {
		return zerrors.NewMissingRequiredError("%v option or %v env variable is required", customizeCommandsOptionName, customizeCommandsEnvVarName)
	}
	return nil
}
