package constants

const (
	CreateVMFromManifestTaskName = "create-vm-from-manifest"
	SetOwnerReference            = "setOwnerReference"
)

type createVMFromManifestParams struct {
	Namespace   string
	Manifest    string
	Virtctl     string
	StartVM     string
	RunStrategy string
}

var CreateVMFromManifestParams = createVMFromManifestParams{
	Namespace:   "namespace",
	Manifest:    "manifest",
	Virtctl:     "virtctl",
	StartVM:     "startVM",
	RunStrategy: "runStrategy",
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
	CreateVMVMManifestMode CreateVMMode = "CreateVMVMManifestMode"
	CreateVMVirtctlMode    CreateVMMode = "CreateVMVirtctlMode"
)

const ExpectedSuccessfulVMCreation = "apiVersion: kubevirt.io/v1\nkind: VirtualMachine\n"
