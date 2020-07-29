package constants

// Option names
const (
	TemplateNameOptionName              = "templateName"
	TemplateNamespaceOptionName         = "templateNamespace"
	TemplateParamsOptionName            = "templateParams"
	OutputParamOptionName               = "output"
	DataVolumesOptionName               = "dataVolumes"
	OwnDataVolumesOptionName            = "ownDataVolumes"
	PersistentVolumeClaimsOptionName    = "persistentVolumeClaims"
	OwnPersistentVolumeClaimsOptionName = "ownPersistentVolumeClaims"
)

// Result names
const (
	NameResultName      = "name"
	NamespaceResultName = "namespace"
)

// Exit codes
const (
	WrongArgsExitCode                         = 2
	VolumesNotPresentExitCode                 = 3
	MissingRequiredTemplateParametersExitCode = 4
	CreateVMErrorExitCode                     = 5
	OwnVolumesErrorExitCode                   = 6
)

// Env related constants
const (
	OutOfClusterENV             = "OUT_OF_CLUSTER"
	PodNamespaceENV             = "POD_NAMESPACE"
	serviceAccountNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	TektonResultsDirPath        = "/tekton/results"
)
