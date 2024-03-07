#!/usr/bin/env bash

function visit() {
  pushd "${1}" > /dev/null || return
}

function leave() {
  popd > /dev/null || return
}

export IS_OKD="false"
export IS_MINIKUBE="false"

if kubectl get projects > /dev/null 2>&1; then
  export IS_OKD="true"
elif minikube status 2>&1 | grep -q Running; then
  export IS_MINIKUBE="true"
fi

export EXCLUDED_NON_IMAGE_MODULES="shared|sharedtest|tests"
export DEPENDENCY_MODULES="shared|sharedtest"
