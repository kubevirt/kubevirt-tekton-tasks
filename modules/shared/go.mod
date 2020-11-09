module github.com/kubevirt/kubevirt-tekton-tasks/modules/shared

go 1.15

require (
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	k8s.io/apimachinery v0.17.1-beta.0
)

// from https://github.com/kubevirt/client-go/blob/v0.26.5/go.mod
//      https://github.com/kubevirt/containerized-data-importer/blob/master/go.mod
replace k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
