#!/usr/bin/env bash

set -ex

if kubectl get namespace tekton-pipelines > /dev/null 2>&1; then
  exit 0
fi

KUBEVIRT_VERSION="v0.59.1"

CDI_VERSION=$(curl -s https://api.github.com/repos/kubevirt/containerized-data-importer/releases | \
            jq '.[] | select(.prerelease==false) | .tag_name' | sort -V | tail -n1 | tr -d '"')

TEKTON_VERSION="v0.64.0"

if kubectl get templates > /dev/null 2>&1; then
  # okd
  COMMON_TEMPLATES_VERSION=$(curl -s https://api.github.com/repos/kubevirt/common-templates/releases | \
            jq '.[] | select(.prerelease==false) | .tag_name' | sort -V | tail -n1 | tr -d '"')
  oc apply -n openshift -f "https://github.com/kubevirt/common-templates/releases/download/${COMMON_TEMPLATES_VERSION}/common-templates.yaml"

  oc new-project tekton-pipelines
fi

# Deploy Tekton Pipelines
oc apply -f "https://github.com/tektoncd/operator/releases/download/${TEKTON_VERSION}/openshift-release.yaml"

# Deploy Kubevirt
kubectl apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-operator.yaml"

kubectl apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-cr.yaml"

kubectl patch kubevirt kubevirt -n kubevirt --type merge -p '{"spec":{"configuration":{"developerConfiguration":{"featureGates": ["DataVolumes"]}}}}'

# Deploy Storage
kubectl apply -f "https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-operator.yaml"

kubectl apply -f "https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-cr.yaml"

# wait for tekton pipelines
kubectl rollout status -n openshift-operators deployment/openshift-pipelines-operator --timeout 10m

# wait until clustertasks tekton CRD is properly deployed
timeout 10m bash <<- EOF
  until kubectl get crd clustertasks.tekton.dev; do
    sleep 5
  done
EOF

# wait until tekton pipelines webhook is created
timeout 10m bash <<- EOF
  until kubectl get deployment tekton-pipelines-webhook -n openshift-pipelines; do
    sleep 5
  done
EOF

# wait until tekton pipelines webhook is online
kubectl wait -n openshift-pipelines deployment tekton-pipelines-webhook --for condition=Available --timeout 10m

# Wait for kubevirt to be available
kubectl rollout status -n cdi deployment/cdi-operator --timeout 10m
kubectl wait -n kubevirt kv kubevirt --for condition=Available --timeout 10m
