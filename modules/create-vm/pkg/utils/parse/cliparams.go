package parse

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/output"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	vmNamespaceOptionName       = "vm-namespace"
	templateNamespaceOptionName = "template-namespace"
)

const templateParamSep = ":"

// TemplateNamespaces and VirtualMachineNamespaces: arrays allow to have these options without option argument
type CLIOptions struct {
	TemplateName              string            `arg:"--template-name,required" placeholder:"NAME" help:"Name of a template to create VM from"`
	TemplateNamespaces        []string          `arg:"--template-namespace" placeholder:"NAMESPACE" help:"Namespace of a template to create VM from"`
	TemplateParams            []string          `arg:"--template-params" placeholder:"KEY2:VAL1 KEY2:VAL2" help:"Template params to pass when processing the template manifest"`
	VirtualMachineNamespaces  []string          `arg:"--vm-namespace" placeholder:"NAMESPACE" help:"Namespace where to create the VM"`
	DataVolumes               []string          `arg:"--dvs" placeholder:"DV1 DV2" help:"Add DataVolumes to VM Volumes"`
	OwnDataVolumes            []string          `arg:"--own-dvs" placeholder:"DV1 DV2" help:"Add DataVolumes to VM Volumes and add VM to DV ownerReferences. These DVs will be deleted once the created VM gets deleted."`
	PersistentVolumeClaims    []string          `arg:"--pvcs" placeholder:"PVC1 PVC2" help:"Add PersistentVolumeClaims to VM Volumes."`
	OwnPersistentVolumeClaims []string          `arg:"--own-pvcs" placeholder:"PVC1 PVC2" help:"Add PersistentVolumeClaims to VM Volumes and add VM to PVC ownerReferences. These PVCs will be deleted once the created VM gets deleted."`
	Output                    output.OutputType `arg:"-o" placeholder:"FORMAT" help:"Output format. One of: yaml|json"`
	Debug                     bool              `arg:"--debug" help:"Sets DEBUG log level"`
}

func (c *CLIOptions) GetAllPVCNames() []string {
	return zutils.ConcatStringSlices(c.OwnPersistentVolumeClaims, c.PersistentVolumeClaims)
}

func (c *CLIOptions) GetAllDVNames() []string {
	return zutils.ConcatStringSlices(c.OwnDataVolumes, c.DataVolumes)
}

func (c *CLIOptions) GetAllDiskNames() []string {
	return zutils.ConcatStringSlices(c.GetAllPVCNames(), c.GetAllDVNames())
}

func (c *CLIOptions) GetTemplateParams() map[string]string {
	result := make(map[string]string, len(c.TemplateParams))

	for _, keyVal := range c.TemplateParams {
		split := strings.SplitN(keyVal, templateParamSep, 2)
		if len(split) == 2 {
			result[split[0]] = split[1]
		}
	}
	return result
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetTemplateNamespace() string {
	return zutils.GetLast(c.TemplateNamespaces)
}

func (c *CLIOptions) GetVirtualMachineNamespace() string {
	return zutils.GetLast(c.VirtualMachineNamespaces)
}

func (c *CLIOptions) setTemplateNamespace(namespace string) {
	c.TemplateNamespaces = []string{namespace}
}

func (c *CLIOptions) setVirtualMachineNamespace(namespace string) {
	c.VirtualMachineNamespaces = []string{namespace}
}

func (c *CLIOptions) Init() error {
	if err := c.assertValidTypes(); err != nil {
		return err
	}

	if err := c.resolveTemplateParams(); err != nil {
		return err
	}

	if err := c.resolveDefaultNamespaces(); err != nil {
		return err
	}

	c.trimSpaces()

	return nil
}
