package parse

import (
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap/zapcore"
)

const (
	vmManifestOptionName        = "vm-manifest"
	vmNamespaceOptionName       = "vm-namespace"
	templateNameOptionName      = "template-name"
	templateNamespaceOptionName = "template-namespace"
	templateParamsOptionName    = "template-params"
	virtctlOptionName           = "virtctl"
)

const templateParamSep = ":"
const volumesSep = ":"

type CLIOptions struct {
	TemplateName            string            `arg:"--template-name,env:TEMPLATE_NAME" placeholder:"NAME" help:"Name of a template to create VM from"`
	TemplateNamespace       string            `arg:"--template-namespace,env:TEMPLATE_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace of a template to create VM from"`
	TemplateParams          []string          `arg:"--template-params" placeholder:"KEY1:VAL1 KEY2:VAL2" help:"Template params to pass when processing the template manifest"`
	VirtualMachineManifest  string            `arg:"--vm-manifest,env:VM_MANIFEST" placeholder:"MANIFEST" help:"YAML manifest of a VirtualMachine resource to be created (can be set by VM_MANIFEST env variable)."`
	VirtualMachineNamespace string            `arg:"--vm-namespace,env:VM_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace where to create the VM"`
	StartVM                 string            `arg:"--start-vm,env:START_VM" help:"Start vm after creation"`
	RunStrategy             string            `arg:"--run-strategy,env:RUN_STRATEGY" help:"Set run strategy to vm"`
	SetOwnerReference       string            `arg:"--set-owner-reference,env:SET_OWNER_REFERENCE" placeholder:"false" help:"Set owner reference to the new object created by the task run pod. Allowed values true/false"`
	Output                  output.OutputType `arg:"-o" placeholder:"FORMAT" help:"Output format. One of: yaml|json"`
	Debug                   bool              `arg:"--debug" help:"Sets DEBUG log level"`
	Virtctl                 string            `arg:"--virtctl,env:VIRTCTL" placeholder:"VIRTCTL" help:"Specifies the parameters for virtctl create vm command that will be used to create VirtualMachine."`
}

func (c *CLIOptions) GetStartVMFlag() bool {
	return c.StartVM == "true"
}

func (c *CLIOptions) GetRunStrategy() string {
	return c.RunStrategy
}

func (c *CLIOptions) GetVirtctl() string {
	return c.Virtctl
}

func (c *CLIOptions) GetSetOwnerReferenceValue() bool {
	return c.SetOwnerReference == "true"
}

func (c *CLIOptions) GetTemplateParams() map[string]string {
	result, err := zutils.ExtractKeysAndValuesByLastKnownKey(c.TemplateParams, templateParamSep)

	if err != nil {
		panic(fmt.Errorf("init was not called: %v", err.Error()))
	}
	return result
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetCreationMode() constants.CreationMode {
	// Input validation is done in Init
	if c.VirtualMachineManifest != "" {
		return constants.VMManifestCreationMode
	}

	if c.TemplateName != "" {
		return constants.TemplateCreationMode
	}

	if c.Virtctl != "" {
		return constants.VirtctlCreatingMode
	}

	return ""
}

func (c *CLIOptions) GetTemplateNamespace() string {
	return c.TemplateNamespace
}

func (c *CLIOptions) GetVirtualMachineManifest() string {
	return c.VirtualMachineManifest
}

func (c *CLIOptions) GetVirtualMachineNamespace() string {
	return c.VirtualMachineNamespace
}

func (c *CLIOptions) Init() error {
	if err := c.assertValidMode(); err != nil {
		return err
	}

	if err := c.assertValidTypes(); err != nil {
		return err
	}

	if _, err := zutils.ExtractKeysAndValuesByLastKnownKey(c.TemplateParams, templateParamSep); err != nil {
		return zerrors.NewMissingRequiredError("invalid %v: %v", templateParamsOptionName, err.Error())
	}

	if err := c.resolveDefaultNamespacesAndManifests(); err != nil {
		return err
	}

	c.trimSpaces()

	return nil
}
