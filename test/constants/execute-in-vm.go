package constants

const ExecuteInVMTaskName = "execute-in-vm"

const CleanupVMTaskName = "cleanup-vm"

type executeOrCleanupVMParams struct {
	VMName      string
	VMNamespace string
	SecretName  string
	Command     string
	Args        string
	Script      string
	Stop        string
	Delete      string
	Timeout     string
}

var ExecuteOrCleanupVMParams = executeOrCleanupVMParams{
	VMName:      "vmName",
	VMNamespace: "vmNamespace",
	SecretName:  "secretName",
	Command:     "command",
	Args:        "args",
	Script:      "script",
	Stop:        "stop",
	Delete:      "delete",
	Timeout:     "timeout",
}

type ExecInVMMode string

const (
	ExecuteInVMMode ExecInVMMode = "execute-in-vm"
	CleanupVMMode   ExecInVMMode = "cleanup-vm"
)
