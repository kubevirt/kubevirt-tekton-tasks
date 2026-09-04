#!/usr/bin/env bash

set -ex

# Cleanup Tekton Pipelines
if kubectl get namespace tekton-pipelines > /dev/null 2>&1; then
  kubectl get crd -o name | grep '\.tekton\.dev$' | xargs --no-run-if-empty kubectl delete --ignore-not-found
  kubectl delete --ignore-not-found ns/tekton-pipelines

  TEKTON_VERSION="$(curl -s https://api.github.com/repos/tektoncd/operator/tags?per_page=10 | jq -r ".[0].name")"
  if hash oc 2>/dev/null; then
    oc delete --ignore-not-found -f "https://github.com/tektoncd/operator/releases/download/${TEKTON_VERSION}/openshift-release.yaml"
  else
    kubectl delete --ignore-not-found -f "https://github.com/tektoncd/operator/releases/download/${TEKTON_VERSION}/release.yaml"
  fi
fi


# Cleanup Storage
if kubectl get namespace cdi > /dev/null 2>&1; then
  kubectl delete --ignore-not-found -f "https://github.com/kubevirt/containerized-data-importer/releases/latest/download/cdi-cr.yaml"
  kubectl wait -n cdi deployment cdi-apiserver cdi-uploadproxy cdi-deployment --for=delete --timeout=10m

  kubectl delete --ignore-not-found -f "https://github.com/kubevirt/containerized-data-importer/releases/latest/download/cdi-operator.yaml"
fi

if kubectl get namespace kubevirt > /dev/null 2>&1; then
  # Cleanup KubeVirt
  kubectl delete --ignore-not-found -f "https://github.com/kubevirt/kubevirt/releases/latest/download/kubevirt-cr.yaml"
  kubectl wait -n kubevirt deployment virt-api virt-controller --for=delete --timeout=10m

  kubectl delete --ignore-not-found -f "https://github.com/kubevirt/kubevirt/releases/latest/download/kubevirt-operator.yaml"
fi
