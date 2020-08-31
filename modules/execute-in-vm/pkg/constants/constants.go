package constants

import "time"

// Exit codes
// reserve 0+ numbers for the exit code of the command
const (
	InvalidArguments                             = -1 // same as go-arg invalid args exit
	InvalidSecret                                = -2
	ExecutorSetupConnectionFailed                = -3
	EnsureVMRunningFailed                        = -4
	RemoteCommandExecutionFailedFromUnknownError = -5
)

const PollVMIInterval = 5 * time.Second
const PollValidConnectionInterval = 4 * time.Second
const CheckSSHConnectionTimeout = 3 * time.Second
