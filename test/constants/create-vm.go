package constants

const CreateVMFromTemplateTaskName = "create-vm-from-template"

const CreateVMFromManifestTaskName = "create-vm-from-manifest"

type createVMFromManifestParams struct {
	Namespace string
	Manifest  string
	Virtctl   string
}

type createVMFromTemplateParams struct {
	TemplateName      string
	TemplateNamespace string
	TemplateParams    string
	VmNamespace       string
	StartVM           string
	RunStrategy       string
}

var CreateVMFromManifestParams = createVMFromManifestParams{
	Namespace: "namespace",
	Manifest:  "manifest",
	Virtctl:   "virtctl",
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
	CreateVMVirtctlMode    CreateVMMode = "CreateVMVirtctlMode"
)

const ExpectedSuccessfulVMCreation = "apiVersion: kubevirt.io/v1\nkind: VirtualMachine\n"
