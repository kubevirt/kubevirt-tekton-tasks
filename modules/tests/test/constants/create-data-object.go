package constants

const (
	CreateDataObjectClusterTaskName    = "create-data-object"
	CreateDataObjectServiceAccountName = "create-data-object-task"
	CreateDataObjectTaskrunName        = "taskrun-create-data-object"

	UnusualRestartCountThreshold = 3
	ReasonError                  = "Error"
)

type createDataObjectParams struct {
	Manifest       string
	WaitForSuccess string
	AllowReplace   string
}

var CreateDataObjectParams = createDataObjectParams{
	Manifest:       "manifest",
	WaitForSuccess: "waitForSuccess",
	AllowReplace:   "allowReplace",
}

type createDataObjectResults struct {
	Name      string
	Namespace string
}

var CreateDataObjectResults = createDataObjectResults{
	Name:      "name",
	Namespace: "namespace",
}
