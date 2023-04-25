package constants

// Exit codes
const (
	GenericExitCode           = 1
	InvalidCLIInputExitCode   = 2
	VolumesNotPresentExitCode = 3
	CreateVMErrorExitCode     = 4
	OwnVolumesErrorExitCode   = 5
	WriteResultsExitCode      = 6
	StartVMErrorExitCode      = 7
)

// Result names
const (
	NameResultName      = "name"
	NamespaceResultName = "namespace"
)

type CreationMode string

const (
	TemplateCreationMode   CreationMode = "TemplateCreationMode"
	VMManifestCreationMode CreationMode = "VMManifestCreationMode"
	VirtctlCreatingMode    CreationMode = "VirtctlCreatingMode"
)
