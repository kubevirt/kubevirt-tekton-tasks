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

	CPUTopologyNumber    uint32 = 3
	CPUTopologyNumberStr        = "3"
	MemoryValue                 = "180M"
)
