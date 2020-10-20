package constants

const CreateVMFromTemplateClusterTaskName = "create-vm-from-template"
const CreateVMFromTemplateServiceAccountName = "create-vm-task"

type createVMFromTemplateParams struct {
	TemplateName              string
	TemplateNamespace         string
	TemplateParams            string
	VmNamespace               string
	DataVolumes               string
	OwnDataVolumes            string
	PersistentVolumeClaims    string
	OwnPersistentVolumeClaims string
}

var CreateVMFromTemplateParams = createVMFromTemplateParams{
	TemplateName:              "templateName",
	TemplateNamespace:         "templateNamespace",
	TemplateParams:            "templateParams",
	VmNamespace:               "vmNamespace",
	DataVolumes:               "dataVolumes",
	OwnDataVolumes:            "ownDataVolumes",
	PersistentVolumeClaims:    "persistentVolumeClaims",
	OwnPersistentVolumeClaims: "ownPersistentVolumeClaims",
}
