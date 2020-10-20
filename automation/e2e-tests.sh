#!/usr/bin/env bash

set -ex

export SCOPE="${SCOPE:-cluster}"
export STORAGE_CLASS="${STORAGE_CLASS:-}"
export DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-e2e-tests-$(shuf -i10000-99999 -n1)}"

export IMAGE_REGISTRY_USER="$DEPLOY_NAMESPACE"

if [[ "$SCOPE" == "cluster" ]]; then
    export TEST_NAMESPACE="${TEST_NAMESPACE:-e2e-tests-$(shuf -i10000-99999 -n1)}"
else
    export TEST_NAMESPACE="${TEST_NAMESPACE:-$DEPLOY_NAMESPACE}"
fi
export TARGET_NAMESPACE="$TEST_NAMESPACE"

make lint
make test
make test-generated-tasks-consistency

./automation/e2e-deploy-resources.sh

oc get namespaces -o name | grep -Eq "^namespace/$TEST_NAMESPACE$" || oc new-project "$TEST_NAMESPACE"
oc get namespaces -o name | grep -Eq "^namespace/$DEPLOY_NAMESPACE$" || oc new-project "$DEPLOY_NAMESPACE"

make cluster-sync

# Wait for kubevirt to be available
oc wait -n kubevirt kv kubevirt --for condition=Available --timeout 15m

oc project "$DEPLOY_NAMESPACE"

make e2e-tests-no-deploy
