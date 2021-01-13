package constants

const ExecuteInVMClusterTaskName = "execute-in-vm"
const ExecuteInVMServiceAccountName = "execute-in-vm-task"

const CleanupVMClusterTaskName = "cleanup-vm"
const CleanupVMServiceAccountName = "cleanup-vm-task"

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
