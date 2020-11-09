module github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm

go 1.15

require (
	github.com/alexflint/go-arg v1.3.0
	github.com/kubevirt/kubevirt-tekton-tasks/modules/shared v0.0.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	go.uber.org/zap v1.15.0
	k8s.io/apimachinery v0.17.1-beta.0
	k8s.io/client-go v12.0.0+incompatible
	kubevirt.io/client-go v0.26.5
)

// remove for versioning the shared module separately
replace github.com/kubevirt/kubevirt-tekton-tasks/modules/shared v0.0.0 => ../shared

replace k8s.io/client-go => k8s.io/client-go v0.16.4
