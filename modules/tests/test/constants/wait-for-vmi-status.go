package constants

const (
	WaitForVMIStatusClusterTaskName    = "wait-for-vmi-status"
	WaitForVMIStatusServiceAccountName = "wait-for-vmi-status-task"
)

type waitForVMIStatusTasksParams struct {
	VMIName          string
	VMINamespace     string
	SuccessCondition string
	FailureCondition string
}

var WaitForVMIStatusTasksParams = waitForVMIStatusTasksParams{
	VMIName:          "vmiName",
	VMINamespace:     "vmiNamespace",
	SuccessCondition: "successCondition",
	FailureCondition: "failureCondition",
}
