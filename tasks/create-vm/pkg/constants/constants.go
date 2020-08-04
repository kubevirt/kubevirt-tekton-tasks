package constants

// ENV variables
const (
	OutOfClusterENV = "OUT_OF_CLUSTER"
)

// Exit codes
const (
	GenericExitCode           = 1
	VolumesNotPresentExitCode = 2
	CreateVMErrorExitCode     = 3
	OwnVolumesErrorExitCode   = 4
	WriteResultsExitCode      = 5
)

// Result names
const (
	NameResultName      = "name"
	NamespaceResultName = "namespace"
)

// Env related constants
const (
	serviceAccountNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	TektonResultsDirPath        = "/tekton/results"
)
