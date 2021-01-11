package constants

import "time"

// Exit codes
// reserve 0+ numbers for the exit code of the command
const (
	InvalidArguments       = -1 // same as go-arg invalid args exit
	ExecutorInitialization = -2
	ExecutorActionsFailed  = -3
	CommandTimeout         = -4
)

const PollVMIInterval = 3 * time.Second
const PollValidConnectionInterval = 3 * time.Second
const CheckSSHConnectionTimeout = 3 * time.Second
const PollVMtoDeleteInterval = 1 * time.Second
const PollVMItoStopInterval = 1 * time.Second
const SetupConnectionDelay = 2 * time.Second

const EmptyConnectionSecretName = "__empty__"
