#!/usr/bin/env bash

set -ex

if kubectl get namespace tekton-pipelines > /dev/null 2>&1; then
  exit 0
fi

# Deploy Tekton Pipelines
TEKTON_VERSION="$(curl -s https://api.github.com/repos/tektoncd/operator/tags?per_page=10 | jq -r ".[0].name")"
if hash oc 2>/dev/null; then
  oc new-project tekton-pipelines
  oc apply -f "https://github.com/tektoncd/operator/releases/download/${TEKTON_VERSION}/openshift-release.yaml"
else
  kubectl create ns tekton-pipelines
  kubectl apply -f "https://github.com/tektoncd/operator/releases/download/${TEKTON_VERSION}/release.yaml"
fi

# Deploy Kubevirt
kubectl apply -f "https://github.com/kubevirt/kubevirt/releases/latest/download/kubevirt-operator.yaml"

kubectl apply -f "https://github.com/kubevirt/kubevirt/releases/latest/download/kubevirt-cr.yaml"

kubectl patch kubevirt kubevirt -n kubevirt --type merge -p '{"spec":{"configuration":{"developerConfiguration":{"featureGates": ["VMExport", "VMPersistentState"]}}}}'

# Deploy Storage
kubectl apply -f "https://github.com/kubevirt/containerized-data-importer/releases/latest/download/cdi-operator.yaml"

kubectl apply -f "https://github.com/kubevirt/containerized-data-importer/releases/latest/download/cdi-cr.yaml"

# wait for tekton pipelines
if hash oc 2>/dev/null; then
  kubectl rollout status -n openshift-operators deployment/openshift-pipelines-operator --timeout 10m

  # wait until tasks tekton CRD is properly deployed
  timeout 10m bash <<- EOF
  until kubectl get crd tasks.tekton.dev; do
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
else
  kubectl wait -n tekton-operator deployment tekton-operator --for condition=Available --timeout 10m
  kubectl wait -n tekton-operator deployment tekton-operator-webhook --for condition=Available --timeout 10m
fi

# Wait for kubevirt to be available
kubectl rollout status -n cdi deployment/cdi-operator --timeout 10m
kubectl wait -n kubevirt kv kubevirt --for condition=Available --timeout 10m
