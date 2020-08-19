package constants

// ENV variables
const (
	OutOfClusterENV = "OUT_OF_CLUSTER"
)

// Exit codes
const (
	GenericExitCode           = 1
	InvalidNamespacesExitCode = 2
	VolumesNotPresentExitCode = 3
	CreateVMErrorExitCode     = 4
	OwnVolumesErrorExitCode   = 5
	WriteResultsExitCode      = 6
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

// misc
const (
	True  = "true"
	False = "false"
)
