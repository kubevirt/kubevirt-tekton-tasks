package constants

const (
	WaitForVMIStatusTaskName                     = "wait-for-vmi-status"
	WaitForVMIStatusServiceAccountName           = "wait-for-vmi-status-task"
	WaitForVMIStatusServiceAccountNameNamespaced = "wait-for-vmi-status-test"
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
