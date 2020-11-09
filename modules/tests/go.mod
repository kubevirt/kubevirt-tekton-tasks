module github.com/kubevirt/kubevirt-tekton-tasks/modules/tests

go 1.15

require (
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/openshift/api v0.0.0
	github.com/openshift/client-go v0.0.0
	github.com/tektoncd/pipeline v0.17.0
	k8s.io/api v0.18.7-rc.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v12.0.0+incompatible
	kubevirt.io/client-go v0.26.5
	kubevirt.io/containerized-data-importer v1.20.1
	sigs.k8s.io/yaml v1.2.0
)

// from https://github.com/kubevirt/client-go/blob/v0.26.5/go.mod
//      https://github.com/kubevirt/containerized-data-importer/blob/master/go.mod
replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191219222812-2987a591a72c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20191125132246-f6563a70e19a
	k8s.io/api => k8s.io/api v0.16.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
	k8s.io/apiserver => k8s.io/apiserver v0.16.4
	k8s.io/client-go => k8s.io/client-go v0.16.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.4
	k8s.io/code-generator => k8s.io/code-generator v0.16.4
	k8s.io/component-base => k8s.io/component-base v0.16.4
	k8s.io/klog => k8s.io/klog v0.4.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.4
	kubevirt.io/containerized-data-importer => kubevirt.io/containerized-data-importer v1.20.1

	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190302045857-e85c7b244fd2
)
