module github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest

go 1.20

require (
	github.com/onsi/ginkgo/v2 v2.4.0
	github.com/openshift/api v0.0.0-20230503133300-8bbcb7ca7183
	k8s.io/api v0.27.1
	k8s.io/apimachinery v0.27.1
	kubevirt.io/api v1.1.0
	kubevirt.io/containerized-data-importer-api v1.58.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.26.3 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/utils v0.0.0-20230505201702-9f6742963106 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

// Kubernetes
replace (
	k8s.io/api => k8s.io/api v0.26.11
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.11
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.11
	k8s.io/client-go => k8s.io/client-go v0.26.11
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.26.11
)

// Openshift
replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20231118005202-0f638a8a4705
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20230120202327-72f107311084
)
