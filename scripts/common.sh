#!/usr/bin/env bash

function visit() {
  pushd "${1}" > /dev/null || return
}

function leave() {
  popd > /dev/null || return
}

function resolve_deploy_namespace() {
  export DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(kubectl config view --minify --output 'jsonpath={..namespace}')}"
  if [[ -z "${DEPLOY_NAMESPACE}" ]]; then
    echo "DEPLOY_NAMESPACE is not set and no namespace is configured in the current kube context." >&2
    echo "Set it explicitly, e.g. 'DEPLOY_NAMESPACE=kubevirt make <target>'." >&2
    exit 1
  fi
}

export IS_OKD="false"
export IS_MINIKUBE="false"

if kubectl get projects > /dev/null 2>&1 </dev/null; then
  export IS_OKD="true"
elif minikube status 2>&1 | grep -q Running; then
  export IS_MINIKUBE="true"
fi

export EXCLUDED_NON_IMAGE_MODULES="shared|sharedtest|tests"
export DEPENDENCY_MODULES="shared|sharedtest"
