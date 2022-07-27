package constants

import "time"

// Exit codes
const (
	InvalidCLIInputExitCode    = -1
	DataObjectCreatorErrorCode = 1
	CreateDataObjectErrorCode  = 2
	WriteResultsExitCode       = 3
	DeleteObjectExitCode       = 4
)

// apiVersion and kinds
const (
	DataVolumeKind = "DataVolume"
	DataSourceKind = "DataSource"
)

// Result names
const (
	NameResultName      = "name"
	NamespaceResultName = "namespace"
)

// WaitForSuccess
const (
	PollInterval                 = 15 * time.Second
	PollTimeout                  = 600 * time.Second
	UnusualRestartCountThreshold = 3
	ReasonError                  = "Error"
)
