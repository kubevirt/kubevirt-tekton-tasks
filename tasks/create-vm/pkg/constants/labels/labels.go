package labels

const (
	// TemplateOsLabel is a label that specifies the OS id of the template
	TemplateOsLabel = "os.template.kubevirt.io"

	// TemplateWorkloadLabel is a label that specifies the workload of the template
	TemplateWorkloadLabel = "workload.template.kubevirt.io"

	// TemplateFlavorLabel is a label that specifies the flavor of the template
	TemplateFlavorLabel = "flavor.template.kubevirt.io"

	// TemplateNameOsAnnotation is an annotation that specifies human readable os name
	TemplateNameOsAnnotation = "name.os.template.kubevirt.io"

	// TemplateNameLabel defines a label of the template name which was used to created the VM
	TemplateNameLabel = "vm.kubevirt.io/template"

	// TemplateNamespace defines a label of the template namespace which was used to create the VM
	TemplateNamespace = "vm.kubevirt.io/template.namespace"

	// VMNameLabel defines a label of virtual machine name which was used to create the VM
	VMNameLabel = "vm.kubevirt.io/name"
)
