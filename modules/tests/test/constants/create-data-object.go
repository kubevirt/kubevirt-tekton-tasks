package constants

const (
	CreateDataObjectClusterTaskName    = "create-data-object"
	CreateDataObjectServiceAccountName = "create-data-object-task"
	CreateDataObjectTaskrunName        = "taskrun-create-data-object"

	UnusualRestartCountThreshold = 3
	ReasonError                  = "Error"
)

type createDataObjectParams struct {
	Manifest            string
	WaitForSuccess      string
	AllowReplace        string
	DeleteObject        string
	DeleteObjectName    string
	DeleteObjectKind    string
	DataObjectNamespace string
}

var CreateDataObjectParams = createDataObjectParams{
	Manifest:            "manifest",
	WaitForSuccess:      "waitForSuccess",
	AllowReplace:        "allowReplace",
	DeleteObject:        "deleteObject",
	DeleteObjectName:    "deleteObjectName",
	DeleteObjectKind:    "deleteObjectKind",
	DataObjectNamespace: "namespace",
}

type createDataObjectResults struct {
	Name      string
	Namespace string
}

var CreateDataObjectResults = createDataObjectResults{
	Name:      "name",
	Namespace: "namespace",
}
