package parse

import (
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap/zapcore"
)

const (
	sourceTemplateNameOptionName      = "source-template-name"
	sourceTemplateNamespaceOptionName = "source-template-namespace"
	targetTemplateNameOptionName      = "target-template-name"
	targetTemplateNamespaceOptionName = "target-template-namespace"
)

type CLIOptions struct {
	SourceTemplateName      string `arg:"--source-template-name,env:SOURCE_TEMPLATE_NAME" placeholder:"NAME" help:"Name of a source template"`
	SourceTemplateNamespace string `arg:"--source-template-namespace,env:SOURCE_TEMPLATE_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace of a source template"`
	TargetTemplateName      string `arg:"--target-template-name,env:TARGET_TEMPLATE_NAME" placeholder:"NAME" help:"Name of a target template"`
	TargetTemplateNamespace string `arg:"--target-template-namespace,env:TARGET_TEMPLATE_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace of a target template"`

	Output output.OutputType `arg:"-o" placeholder:"FORMAT" help:"Output format. One of: yaml|json"`
	Debug  bool              `arg:"--debug" help:"Sets DEBUG log level"`
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetSourceTemplateNamespace() string {
	return c.SourceTemplateNamespace
}

func (c *CLIOptions) GetSourceTemplateName() string {
	return c.SourceTemplateName
}

func (c *CLIOptions) GetTargetTemplateNamespace() string {
	return c.TargetTemplateNamespace
}

func (c *CLIOptions) GetTargetTemplateName() string {
	return c.TargetTemplateName
}

func (c *CLIOptions) Init() error {
	if err := c.assertValidMode(); err != nil {
		return err
	}

	if err := c.assertValidTypes(); err != nil {
		return err
	}

	c.trimSpaces()

	return nil
}

func (c *CLIOptions) trimSpaces() {
	for _, strVariablePtr := range []*string{&c.SourceTemplateName, &c.SourceTemplateNamespace, &c.TargetTemplateName, &c.TargetTemplateNamespace} {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}
}

func (c *CLIOptions) assertValidMode() error {
	if c.SourceTemplateName == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", sourceTemplateNameOptionName)
	}

	if c.SourceTemplateNamespace == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", sourceTemplateNamespaceOptionName)
	}

	if c.TargetTemplateName == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", targetTemplateNameOptionName)
	}

	if c.TargetTemplateNamespace == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", targetTemplateNamespaceOptionName)
	}

	return nil
}

func (c *CLIOptions) assertValidTypes() error {
	if !output.IsOutputType(string(c.Output)) {
		return zerrors.NewMissingRequiredError("%v is not a valid output type", c.Output)
	}
	return nil
}
