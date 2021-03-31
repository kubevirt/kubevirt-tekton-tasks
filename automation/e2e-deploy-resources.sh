#!/usr/bin/env bash

set -ex

if kubectl get namespace tekton-pipelines > /dev/null 2>&1; then
  exit 0
fi

KUBEVIRT_VERSION="v0.39.0"
# TODO switch once the latest release YAML is fixed: https://github.com/kubevirt/kubevirt/issues/5352
#KUBEVIRT_VERSION=$(curl -s https://github.com/kubevirt/kubevirt/releases/latest | grep -o "v[0-9]\.[0-9]*\.[0-9]*")
CDI_VERSION=$(curl -s https://github.com/kubevirt/containerized-data-importer/releases/latest | grep -o "v[0-9]\.[0-9]*\.[0-9]*")
COMMON_TEMPLATES_VERSION=""
if kubectl get templates > /dev/null 2>&1; then
  # openshift
  COMMON_TEMPLATES_VERSION=$(curl -s https://github.com/kubevirt/common-templates/releases/latest | grep -o "v[0-9]\.[0-9]*\.[0-9]*")

  # Prepare Tekton Pipelines
  oc new-project tekton-pipelines
  oc adm policy add-scc-to-user anyuid -z tekton-pipelines-controller
  oc adm policy add-scc-to-user anyuid -z tekton-pipelines-webhook
fi


# Deploy Tekton Pipelines
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.notags.yaml
kubectl config set-context --current --namespace=default

# Deploy Kubevirt
kubectl create -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-operator.yaml"

kubectl create -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-cr.yaml"

kubectl apply -f - <<EOF
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: kubevirt-config
  namespace: kubevirt
data:
  feature-gates: "DataVolumes"
---
EOF

# Deploy Storage
kubectl create -f "https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-operator.yaml"

kubectl create -f "https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-cr.yaml"

# Deploy Common Templates

if [ -n "${COMMON_TEMPLATES_VERSION}" ]; then
  kubectl create -n openshift -f "https://github.com/kubevirt/common-templates/releases/download/${COMMON_TEMPLATES_VERSION}/common-templates-${COMMON_TEMPLATES_VERSION}.yaml"
fi

# wait for tekton pipelines
kubectl rollout status -n tekton-pipelines deployment/tekton-pipelines-controller --timeout 10m
kubectl rollout status -n tekton-pipelines deployment/tekton-pipelines-webhook --timeout 10m

# Wait for kubevirt to be available
kubectl rollout status -n cdi deployment/cdi-operator --timeout 10m
kubectl wait -n kubevirt kv kubevirt --for condition=Available --timeout 10m
