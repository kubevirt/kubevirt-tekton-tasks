package parse

import "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils"

type CLIParams struct {
	TemplateName              string
	TemplateNamespace         string
	TemplateParams            map[string]string
	DataVolumes               []string
	OwnDataVolumes            []string
	PersistentVolumeClaims    []string
	OwnPersistentVolumeClaims []string
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
