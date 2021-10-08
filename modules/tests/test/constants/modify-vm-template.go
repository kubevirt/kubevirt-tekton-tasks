package constants

const (
	ModifyTemplateClusterTaskName    = "modify-vm-template"
	ModifyTemplateServiceAccountName = "modify-vm-template-task"
	ModifyTemplateTaskRunName        = "taskrun-modify-vm-template"

	TemplateNameOptionName        = "templateName"
	TemplateNamespaceOptionName   = "templateNamespace"
	CPUCoresOptionName            = "cpuCores"
	CPUSocketsOptionName          = "cpuSockets"
	CPUThreadsOptionName          = "cpuThreads"
	MemoryOptionName              = "memory"
	TemplateLabelsOptionName      = "templateLabels"
	TemplateAnnotationsOptionName = "templateAnnotations"
	VMLabelsOptionName            = "vmLabels"
	VMAnnotationsOptionName       = "vmAnnotations"

	CPUSocketsTopologyNumber    uint32 = 1
	CPUCoresTopologyNumber      uint32 = 2
	CPUThreadsTopologyNumber    uint32 = 3
	CPUSocketsTopologyNumberStr        = "1"
	CPUCoresTopologyNumberStr          = "2"
	CPUThreadsTopologyNumberStr        = "3"

	MemoryValue = "180M"
)

var (
	MockArray = []string{"newKey: value", "test: true"}

	LabelsAnnotationsMap = map[string]string{
		"newKey": "value",
		"test":   "true",
	}
)
