module github.com/kubevirt/kubevirt-tekton-tasks/modules/tests

go 1.20

require (
	github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template v0.0.0
	github.com/kubevirt/kubevirt-tekton-tasks/modules/shared v0.0.0
	github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest v0.0.0
	github.com/onsi/ginkgo/v2 v2.1.6
	github.com/onsi/gomega v1.20.1
	github.com/openshift/api v3.9.0+incompatible
	github.com/openshift/client-go v3.9.0+incompatible
	github.com/tektoncd/pipeline v0.52.1
	k8s.io/api v0.27.1
	k8s.io/apimachinery v0.27.3
	k8s.io/client-go v12.0.0+incompatible
	knative.dev/pkg v0.0.0-20230718152110-aef227e72ead
	kubevirt.io/api v0.59.0
	kubevirt.io/client-go v0.59.0
	kubevirt.io/containerized-data-importer v1.55.0
	kubevirt.io/containerized-data-importer-api v1.55.0
	kubevirt.io/qe-tools v0.1.8
	sigs.k8s.io/yaml v1.3.0
)

require (
	contrib.go.opencensus.io/exporter/ocagent 7399e0f8ee5e // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudevents/sdk-go/v2 v2.16.2 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.18.2 // indirect
	github.com/coreos/prometheus-operator v0.94.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/cli v24.0.0+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v28.0.0+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.9 // indirect
	github.com/emicklei/go-restful/v3 v3.13.0 // indirect
	github.com/evanphx/json-patch v4.13.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.1 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-openapi/jsonpointer v0.24.0 // indirect
	github.com/go-openapi/jsonreference v0.21.6 // indirect
	github.com/go-openapi/spec v0.22.11 // indirect
	github.com/go-openapi/swag v0.29.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.5 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic v0.7.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-containerregistry v0.22.1 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.30.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20191119172530-79f836b90111 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0 // indirect
	github.com/letsencrypt/boulder v0.20260908.0 // indirect
	github.com/mailru/easyjson v0.9.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.71.0 // indirect
	github.com/prometheus/procfs v0.22.0 // indirect
	github.com/prometheus/statsd_exporter v0.31.0 // indirect
	github.com/sigstore/sigstore v1.7.3 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/theupdateframework/go-tuf v0.7.0 // indirect
	github.com/titanous/rocacheck v0.0.0-20171023193734-afe73141d399 // indirect
	github.com/vbatts/tar-split v0.12.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20260908205506-85c1c2202aba // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/oauth2 v0.13.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.46.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.16.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	google.golang.org/api v0.298.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20260918162117-cecb64721679 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260918162117-cecb64721679 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260918162117-cecb64721679 // indirect
	google.golang.org/grpc v1.58.3 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.26.5 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230515203736-54b630e78af5 // indirect
	k8s.io/utils v0.0.0-20230505201702-9f6742963106 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

// Kubernetes
replace (
	k8s.io/api => k8s.io/api v0.24.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.24.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.24.6
	k8s.io/client-go => k8s.io/client-go v0.24.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.24.6
)

// CDI
replace (
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.46.0
	kubevirt.io/containerized-data-importer => kubevirt.io/containerized-data-importer v1.55.0
	kubevirt.io/containerized-data-importer-api => kubevirt.io/containerized-data-importer-api v1.55.0
	kubevirt.io/controller-lifecycle-operator-sdk/api => kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90
)

// locally referenced modules
replace (
	github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template v0.0.0 => ../copy-template
	github.com/kubevirt/kubevirt-tekton-tasks/modules/shared v0.0.0 => ../shared
	github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest v0.0.0 => ../sharedtest
)

// Openshift
replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20220325173635-8107b7a38e53
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20220316161609-20d926360175
)
