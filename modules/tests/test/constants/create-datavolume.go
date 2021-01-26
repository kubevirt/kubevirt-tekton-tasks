package constants

const (
	CreateDataVolumeFromManifestClusterTaskName    = "create-datavolume-from-manifest"
	CreateDataVolumeFromManifestServiceAccountName = "create-datavolume-from-manifest-task"
)

const (
	DataVolumeKind       = "DataVolume"
	DataVolumeApiVersion = "cdi.kubevirt.io/v1beta1"
)

type createDataVolumeFromManifestParams struct {
	Manifest       string
	WaitForSuccess string
}

var CreateDataVolumeFromManifestParams = createDataVolumeFromManifestParams{
	Manifest:       "manifest",
	WaitForSuccess: "waitForSuccess",
}

type createDataVolumeFromManifestResults struct {
	Name      string
	Namespace string
}

var CreateDataVolumeFromManifestResults = createDataVolumeFromManifestResults{
	Name:      "name",
	Namespace: "namespace",
}
