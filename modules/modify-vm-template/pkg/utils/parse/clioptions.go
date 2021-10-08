package parse

import (
	"strconv"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	templateNameOptionName      = "template-name"
	templateNamespaceOptionName = "template-namespace"
	colonSeparator              = ":"
)

type CLIOptions struct {
	TemplateName        string            `arg:"--template-name,env:TEMPLATE_NAME,required" placeholder:"NAME" help:"Name of a template"`
	TemplateNamespace   string            `arg:"--template-namespace,env:TEMPLATE_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace of a template"`
	CPUSockets          string            `arg:"--cpu-sockets,env:CPU_SOCKETS" placeholder:"CPU_SOCKETS" help:"Number of CPU sockets"`
	CPUCores            string            `arg:"--cpu-cores,env:CPU_CORES" placeholder:"CPU_CORES" help:"Number of CPU cores"`
	CPUThreads          string            `arg:"--cpu-threads,env:CPU_THREADS" placeholder:"CPU_THREADS" help:"Number of CPU threads"`
	Memory              string            `arg:"--memory,env:MEMORY" placeholder:"MEMORY" help:"Memory of the vm, format 1M, 1G"`
	TemplateLabels      []string          `arg:"--template-labels" placeholder:"KEY: VALUE KEY: VALUE" help:"Adds labels to template"`
	TemplateAnnotations []string          `arg:"--template-annotations" placeholder:"KEY: VALUE KEY: VALUE" help:"Adds annotations to template"`
	VMLabels            []string          `arg:"--vm-labels" placeholder:"KEY: VALUE KEY: VALUE" help:"Adds labels to VMs"`
	VMAnnotations       []string          `arg:"--vm-annotations" placeholder:"KEY: VALUE KEY: VALUE" help:"Adds annotations to VMs"`
	Output              output.OutputType `arg:"-o" placeholder:"FORMAT" help:"Output format. One of: yaml|json"`
	Debug               bool              `arg:"--debug" help:"Sets DEBUG log level"`

	templateLabels      map[string]string
	templateAnnotations map[string]string
	vmLabels            map[string]string
	vmAnnotations       map[string]string
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetCPUSockets() int {
	res, _ := strconv.Atoi(c.CPUSockets)
	return res
}

func (c *CLIOptions) GetCPUCores() int {
	res, _ := strconv.Atoi(c.CPUCores)
	return res
}

func (c *CLIOptions) GetCPUThreads() int {
	res, _ := strconv.Atoi(c.CPUThreads)
	return res
}

func (c *CLIOptions) GetMemory() *resource.Quantity {
	if c.Memory == "" {
		return nil
	}
	q := resource.MustParse(c.Memory)
	return &q
}

func (c *CLIOptions) GetTemplateName() string {
	return c.TemplateName
}

func (c *CLIOptions) GetTemplateNamespace() string {
	return c.TemplateNamespace
}

func (c *CLIOptions) GetTemplateLabels() map[string]string {
	return c.templateLabels
}

func (c *CLIOptions) GetTemplateAnnotations() map[string]string {
	return c.templateAnnotations
}

func (c *CLIOptions) GetVMAnnotations() map[string]string {
	return c.vmAnnotations
}

func (c *CLIOptions) GetVMLabels() map[string]string {
	return c.vmLabels
}

func (c *CLIOptions) Init() error {
	c.trimSpaces()

	if err := c.assertValidParams(); err != nil {
		return err
	}

	if err := c.assertValidTypes(); err != nil {
		return err
	}

	c.setDefaultValues()

	return nil
}

func (c *CLIOptions) trimSpaces() {
	for _, strVariablePtr := range []*string{&c.TemplateName, &c.TemplateNamespace} {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}

	for i, value := range c.TemplateLabels {
		value = strings.ReplaceAll(value, " ", "")
		c.TemplateLabels[i] = value
	}
	for i, value := range c.TemplateAnnotations {
		value = strings.ReplaceAll(value, " ", "")
		c.TemplateAnnotations[i] = value
	}
	for i, value := range c.VMLabels {
		value = strings.ReplaceAll(value, " ", "")
		c.VMLabels[i] = value
	}
	for i, value := range c.VMAnnotations {
		value = strings.ReplaceAll(value, " ", "")
		c.VMAnnotations[i] = value
	}
}

func (c *CLIOptions) setDefaultValues() error {
	if c.TemplateNamespace == "" {
		activeNamespace, err := env.GetActiveNamespace()
		if err != nil {
			return zerrors.NewMissingRequiredError("can't get active namespace: %v", err.Error())
		}

		c.TemplateNamespace = activeNamespace
	}
	return nil
}

func checkCorrectInt(value string) error {
	if value == "" {
		return nil
	}
	_, err := strconv.Atoi(value)
	return err
}

func (c *CLIOptions) assertValidParams() error {
	if c.TemplateName == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", templateNameOptionName)
	}

	if c.Memory != "" {
		_, err := resource.ParseQuantity(c.Memory)
		if err != nil {
			return err
		}
	}

	err := checkCorrectInt(c.CPUCores)
	if err != nil {
		return err
	}

	err = checkCorrectInt(c.CPUThreads)
	if err != nil {
		return err
	}

	err = checkCorrectInt(c.CPUSockets)
	if err != nil {
		return err
	}

	c.templateLabels, err = zutils.ExtractKeysAndValuesByLastKnownKey(c.TemplateLabels, colonSeparator)
	if err != nil {
		return err
	}

	c.templateAnnotations, err = zutils.ExtractKeysAndValuesByLastKnownKey(c.TemplateAnnotations, colonSeparator)
	if err != nil {
		return err
	}

	c.vmLabels, err = zutils.ExtractKeysAndValuesByLastKnownKey(c.VMLabels, colonSeparator)
	if err != nil {
		return err
	}

	c.vmAnnotations, err = zutils.ExtractKeysAndValuesByLastKnownKey(c.VMAnnotations, colonSeparator)
	if err != nil {
		return err
	}
	return nil
}

func (c *CLIOptions) assertValidTypes() error {
	if !output.IsOutputType(string(c.Output)) {
		return zerrors.NewMissingRequiredError("%v is not a valid output type", c.Output)
	}
	return nil
}
