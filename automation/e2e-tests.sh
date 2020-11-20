#!/usr/bin/env bash

set -ex

export SCOPE="${SCOPE:-cluster}"
export DEV_MODE="${DEV_MODE:-false}"
export STORAGE_CLASS="${STORAGE_CLASS:-}"
export DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-e2e-tests-$(shuf -i10000-99999 -n1)}"
export NUM_NODES=${NUM_NODES:-2}

export CREATE_VM_IMAGE="${CREATE_VM_IMAGE:-}"

./automation/e2e-deploy-resources.sh

oc get namespaces -o name | grep -Eq "^namespace/$DEPLOY_NAMESPACE$" || oc new-project "$DEPLOY_NAMESPACE"

oc project $DEPLOY_NAMESPACE

if [[ "$DEV_MODE" == "true" ]]; then
  make cluster-sync
else
  make deploy
fi


# Wait for kubevirt to be available
oc wait -n kubevirt kv kubevirt --for condition=Available --timeout 15m
oc rollout status -n cdi deployment/cdi-operator --timeout 10m

make cluster-test
