package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/requirements"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	vmiNameOptionName          = "vmi-name"
	vmiNamespaceOptionName     = "vmi-namespace"
	successConditionOptionName = "success-condition"
	failureConditionOptionName = "failure-condition"
)

type CLIOptions struct {
	VirtualMachineInstanceName      string `arg:"--vmi-name,env:VMI_NAME" placeholder:"NAME" help:"Name of a VMI to wait for."`
	VirtualMachineInstanceNamespace string `arg:"--vmi-namespace,env:VMI_NAMESPACE" placeholder:"NAME" help:"Namespace of a VMI to wait for."`
	SuccessCondition                string `arg:"--success-condition,env:SUCCESS_CONDITION" placeholder:"CONDITION" help:" A label selector expression to decide if the VirtualMachineInstance (VMI) is in a success state. Eg. \"status.phase == Succeeded\". It is evaluated on each VMI update and will result in this task succeeding if true."`
	FailureCondition                string `arg:"--failure-condition,env:FAILURE_CONDITION" placeholder:"CONDITION" help:"A label selector expression to decide if the VirtualMachineInstance (VMI) is in a failed state. Eg. \"status.phase in (Failed, Unknown)\". It is evaluated on each VMI update and will result in this task failing if true."`
	Debug                           bool   `arg:"--debug" help:"Sets DEBUG log level"`
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetVirtualMachineInstanceName() string {
	return c.VirtualMachineInstanceName
}

func (c *CLIOptions) GetVirtualMachineInstanceNamespace() string {
	return c.VirtualMachineInstanceNamespace
}

func (c *CLIOptions) GetSuccessCondition() string {
	return c.SuccessCondition
}

func (c *CLIOptions) GetFailureCondition() string {
	return c.FailureCondition
}

func (c *CLIOptions) GetSuccessRequirements() labels.Requirements {
	reqs, err := requirements.GetLabelRequirement(c.SuccessCondition)
	if err != nil {
		panic("Init should be called first to validate the SuccessCondition")
	}
	return reqs
}

func (c *CLIOptions) GetFailureRequirements() labels.Requirements {
	reqs, err := requirements.GetLabelRequirement(c.FailureCondition)
	if err != nil {
		panic("Init should be called first to validate the FailureCondition")
	}
	return reqs
}

func (c *CLIOptions) Init() error {
	c.trimSpaces()

	if err := c.validateNames(); err != nil {
		return err
	}

	if err := c.resolveDefaultNamespaces(); err != nil {
		return err
	}

	if err := c.validateConditions(); err != nil {
		return err
	}

	return nil
}
