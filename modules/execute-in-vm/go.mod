module github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm

go 1.15

require (
	github.com/alexflint/go-arg v1.3.0
	github.com/kubevirt/kubevirt-tekton-tasks/modules/shared v0.0.0
	github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest v0.0.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	go.uber.org/zap v1.15.0
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v12.0.0+incompatible
	kubevirt.io/client-go v0.35.0
)

// locally referenced modules
replace (
	github.com/kubevirt/kubevirt-tekton-tasks/modules/shared v0.0.0 => ../shared
	github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest v0.0.0 => ../sharedtest
)

// Pinned to kubernetes-1.16.4 (due to kubevirt/client-go v0.35.0)
replace (
	k8s.io/api => k8s.io/api v0.16.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
	k8s.io/apiserver => k8s.io/apiserver v0.16.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.4
	k8s.io/client-go => k8s.io/client-go v0.16.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.4
	k8s.io/code-generator => k8s.io/code-generator v0.16.4
	k8s.io/component-base => k8s.io/component-base v0.16.4
	k8s.io/cri-api => k8s.io/cri-api v0.16.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.4
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.4
	k8s.io/kubectl => k8s.io/kubectl v0.16.4
	k8s.io/kubelet => k8s.io/kubelet v0.16.4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.4
	k8s.io/metrics => k8s.io/metrics v0.16.4
	k8s.io/node-api => k8s.io/node-api v0.16.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.4
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.16.4
	k8s.io/sample-controller => k8s.io/sample-controller v0.16.4
)

// Kubernetes
replace (
	k8s.io/klog => k8s.io/klog v0.4.0
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.2
)

// Kubevirt
replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.3.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
)

// Openshift
// Pinned to openshift-master 20191125 (due to kubernetes-1.16.4)
replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191219222812-2987a591a72c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20191125132246-f6563a70e19a
)

// CDI
// Pinned to v1.21.0 (due to kubernetes-1.16.4)
replace (
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a
	kubevirt.io/containerized-data-importer => kubevirt.io/containerized-data-importer v1.21.0
)
