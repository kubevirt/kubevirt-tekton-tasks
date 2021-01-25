package constants

const CreateVMFromTemplateClusterTaskName = "create-vm-from-template"
const CreateVMFromTemplateServiceAccountName = "create-vm-from-template-task"

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

type createVMFromManifestResults struct {
	Name      string
	Namespace string
}

var CreateVMFromManifestResults = createVMFromManifestResults{
	Name:      "name",
	Namespace: "namespace",
}

const ExpectedSuccessfulVMCreation = "apiVersion: kubevirt.io/v1alpha3\nkind: VirtualMachine\n"
