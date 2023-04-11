package constants

const CreateVMFromTemplateTaskName = "create-vm-from-template"
const CreateVMFromTemplateServiceAccountName = "create-vm-from-template-task"
const CreateVMFromTemplateServiceAccountNameNamespaced = "create-vm-from-template-test"

const CreateVMFromManifestTaskName = "create-vm-from-manifest"
const CreateVMFromManifestServiceAccountName = "create-vm-from-manifest-task"
const CreateVMFromManifestServiceAccountNameNamespaced = "create-vm-from-manifest-test"

type createVMParams struct {
	DataVolumes               string
	OwnDataVolumes            string
	PersistentVolumeClaims    string
	OwnPersistentVolumeClaims string
}

type createVMFromManifestParams struct {
	Namespace string
	Manifest  string
}

type createVMFromTemplateParams struct {
	TemplateName      string
	TemplateNamespace string
	TemplateParams    string
	VmNamespace       string
	StartVM           string
	RunStrategy       string
}

var CreateVMParams = createVMParams{
	DataVolumes:               "dataVolumes",
	OwnDataVolumes:            "ownDataVolumes",
	PersistentVolumeClaims:    "persistentVolumeClaims",
	OwnPersistentVolumeClaims: "ownPersistentVolumeClaims",
}

var CreateVMFromManifestParams = createVMFromManifestParams{
	Namespace: "namespace",
	Manifest:  "manifest",
}

var CreateVMFromTemplateParams = createVMFromTemplateParams{
	TemplateName:      "templateName",
	TemplateNamespace: "templateNamespace",
	TemplateParams:    "templateParams",
	VmNamespace:       "vmNamespace",
	StartVM:           "startVM",
	RunStrategy:       "runStrategy",
}

type createVMResults struct {
	Name      string
	Namespace string
}

var CreateVMResults = createVMResults{
	Name:      "name",
	Namespace: "namespace",
}

type CreateVMMode string

const (
	CreateVMTemplateMode   CreateVMMode = "CreateVMTemplateMode"
	CreateVMVMManifestMode CreateVMMode = "CreateVMVMManifestMode"
)

const ExpectedSuccessfulVMCreation = "apiVersion: kubevirt.io/v1\nkind: VirtualMachine\n"
