package parse

import (
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/output"
	"strings"
)

const (
	VMNamespaceOptionName = "vm-namespace"
)

const templateParamSep = ":"

type CLIParams struct {
	TemplateName              string            `arg:"--template-name,required" placeholder:"NAME" help:"Name of a template to create VM from"`
	TemplateNamespace         string            `arg:"--template-namespace" placeholder:"NAMESPACE" help:"Namespace of a template to create VM from"`
	TemplateParams            []string          `arg:"--template-params" placeholder:"KEY2:VAL1 KEY2:VAL2" help:"Template params to pass when processing the template manifest"`
	VirtualMachineNamespace   string            `arg:"--vm-namespace" placeholder:"NAMESPACE" help:"Namespace where to create the VM"`
	DataVolumes               []string          `arg:"--dvs" placeholder:"DV1 DV2" help:"Add DataVolumes to VM Volumes"`
	OwnDataVolumes            []string          `arg:"--own-dvs" placeholder:"DV1 DV2" help:"Add DataVolumes to VM Volumes and add VM to DV ownerReferences. These DVs will be deleted once the created VM gets deleted."`
	PersistentVolumeClaims    []string          `arg:"--pvcs" placeholder:"PVC1 PVC2" help:"Add PersistentVolumeClaims to VM Volumes."`
	OwnPersistentVolumeClaims []string          `arg:"--own-pvcs" placeholder:"PVC1 PVC2" help:"Add PersistentVolumeClaims to VM Volumes and add VM to PVC ownerReferences. These PVCs will be deleted once the created VM gets deleted."`
	Output                    output.OutputType `arg:"-o" placeholder:"FORMAT" help:"Output format. One of: yaml|json"`
}

func (c *CLIParams) GetAllPVCNames() []string {
	return utils.ConcatStringArrays(c.OwnPersistentVolumeClaims, c.PersistentVolumeClaims)
}

func (c *CLIParams) GetAllDVNames() []string {
	return utils.ConcatStringArrays(c.OwnDataVolumes, c.DataVolumes)
}

func (c *CLIParams) GetAllDiskNames() []string {
	return utils.ConcatStringArrays(c.GetAllPVCNames(), c.GetAllDVNames())
}

func (c *CLIParams) GetTemplateParams() map[string]string {
	result := make(map[string]string, len(c.TemplateParams))

	for _, keyVal := range c.TemplateParams {
		split := strings.SplitN(keyVal, templateParamSep, 2)
		if len(split) == 2 {
			result[split[0]] = split[1]
		}
	}
	return result
}
