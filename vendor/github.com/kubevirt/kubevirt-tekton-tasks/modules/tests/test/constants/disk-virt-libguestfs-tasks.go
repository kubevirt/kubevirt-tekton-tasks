package constants

const DiskVirtCustomizeTaskName = "disk-virt-customize"
const DiskVirtSysprepTaskName = "disk-virt-sysprep"

type diskVirtLibguestfsTasksParams struct {
	PVCName           string
	CustomizeCommands string
	SysprepCommands   string
	AdditionalOptions string
	Verbose           string
}

var DiskVirtLibguestfsTasksParams = diskVirtLibguestfsTasksParams{
	PVCName:           "pvc",
	CustomizeCommands: "customizeCommands",
	SysprepCommands:   "sysprepCommands",
	AdditionalOptions: "additionalOptions",
	Verbose:           "verbose",
}

type LibguestfsTaskType string

const (
	VirtSysPrepTaskType   LibguestfsTaskType = "virt-sysprep"
	VirtCustomizeTaskType LibguestfsTaskType = "virt-customize"
)
