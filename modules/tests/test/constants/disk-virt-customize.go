package constants

const DiskVirtCustomizeClusterTaskName = "disk-virt-customize"

type diskVirtCustomizeParams struct {
	PVCName           string
	CustomizeCommands string
	AdditionalOptions string
	Verbose           string
}

var DiskVirtCustomizeParams = diskVirtCustomizeParams{
	PVCName:           "pvc",
	CustomizeCommands: "customizeCommands",
	AdditionalOptions: "additionalOptions",
	Verbose:           "verbose",
}
