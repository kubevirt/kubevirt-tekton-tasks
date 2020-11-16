#!/usr/bin/env bash

SCOPE="${SCOPE:-cluster}"
DEBUG="${DEBUG:-false}"
STORAGE_CLASS="${STORAGE_CLASS:-}"
NUM_NODES=${NUM_NODES:-2}
DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(oc project --short)}"

if [[ "$SCOPE" == "cluster" ]]; then
    export TEST_NAMESPACE="${TEST_NAMESPACE:-e2e-tests-$(shuf -i10000-99999 -n1)}"
else
    export TEST_NAMESPACE="${TEST_NAMESPACE:-$DEPLOY_NAMESPACE}"
fi
export TARGET_NAMESPACE="$TEST_NAMESPACE"

oc get namespaces -o name | grep -Eq "^namespace/$TEST_NAMESPACE$" || oc new-project "$TEST_NAMESPACE" > /dev/null
oc get namespaces -o name | grep -Eq "^namespace/$DEPLOY_NAMESPACE$" || oc new-project "$DEPLOY_NAMESPACE" > /dev/null

oc project "$DEPLOY_NAMESPACE"

pushd modules/tests || exit
  ginkgo -r -p --randomizeAllSpecs --randomizeSuites --failOnPending --trace --race --nodes="${NUM_NODES}" -- \
    --deploy-namespace="${DEPLOY_NAMESPACE}" \
    --test-namespace="${TEST_NAMESPACE}" \
    --kubeconfig-path="${KUBECONFIG}" \
    --scope="${SCOPE}" \
    --storage-class="${STORAGE_CLASS}" \
    --debug="${DEBUG}"
popd
