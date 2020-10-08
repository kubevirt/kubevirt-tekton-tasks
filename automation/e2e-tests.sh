#!/usr/bin/env bash

set -ex

export SCOPE="${SCOPE:-cluster}"
export NAMESPACE="${NAMESPACE:-e2e-tests-$(shuf -i10000-99999 -n1)}"
export TARGET_NAMESPACE="$NAMESPACE"
export IMAGE_REGISTRY_USER="$NAMESPACE"

make lint
make test
make test-generated-tasks-consistency

./automation/e2e-deploy-resources.sh

oc new-project "$NAMESPACE"

make cluster-sync

# Wait for kubevirt to be available
oc wait -n kubevirt kv kubevirt --for condition=Available --timeout 15m
