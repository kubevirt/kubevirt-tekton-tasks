package constants

const DiskVirtCustomizeTaskName = "disk-virt-customize"
const DiskVirtSysprepTaskName = "disk-virt-sysprep"

type diskVirtLibguestfsTasksParams struct {
	PVCName                      string
	VirtCommandsFileNameCommands string
	AdditionalOptions            string
	Verbose                      string
}

var DiskVirtLibguestfsTasksParams = diskVirtLibguestfsTasksParams{
	PVCName:                      "pvc",
	VirtCommandsFileNameCommands: "virtCommands",
	AdditionalOptions:            "additionalVirtOptions",
	Verbose:                      "verbose",
}

type LibguestfsTaskType string

const (
	VirtSysPrepTaskType   LibguestfsTaskType = "virt-sysprep"
	VirtCustomizeTaskType LibguestfsTaskType = "virt-customize"
)
