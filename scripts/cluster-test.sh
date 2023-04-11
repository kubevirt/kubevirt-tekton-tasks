#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

source "${SCRIPT_DIR}/common.sh"

RET_CODE=0

DEBUG="${DEBUG:-false}"
STORAGE_CLASS="${STORAGE_CLASS:-}"
NUM_NODES="${NUM_NODES:-2}"
DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(kubectl config view --minify --output 'jsonpath={..namespace}')}"
ARTIFACT_DIR="${ARTIFACT_DIR:=dist}"
ARTIFACT_DIR="$(readlink -m "${ARTIFACT_DIR}")"
TEST_OUT="${ARTIFACT_DIR}/test.out"


rm -rf "${TEST_OUT}" "${ARTIFACT_DIR}/"junit*
rm -rf "${ARTIFACT_DIR}"

mkdir -p "${ARTIFACT_DIR}"

export TEST_NAMESPACE="${TEST_NAMESPACE:-e2e-tests-$(shuf -i10000-99999 -n1)}"

kubectl get namespaces -o name | grep -Eq "^namespace/$TEST_NAMESPACE$" || kubectl create namespace "$TEST_NAMESPACE" > /dev/null
kubectl get namespaces -o name | grep -Eq "^namespace/$DEPLOY_NAMESPACE$" || kubectl create namespace "$DEPLOY_NAMESPACE" > /dev/null

kubectl config set-context --current --namespace="$DEPLOY_NAMESPACE"

pushd modules/tests || exit
  rm -rf dist
  mkdir dist

  set +ex
  set -o pipefail

  ginkgo -r -p --randomize-all --randomize-suites --fail-on-pending --trace --race --nodes="${NUM_NODES}" -- \
    --deploy-namespace="${DEPLOY_NAMESPACE}" \
    --test-namespace="${TEST_NAMESPACE}" \
    --kubeconfig-path="${KUBECONFIG}" \
    --is-okd="${IS_OKD}" \
    --ginkgo.junit-report="${ARTIFACT_DIR}/xunit_results.xml" \
    --storage-class="${STORAGE_CLASS}" \
    --debug="${DEBUG}" | tee "${TEST_OUT}"

  RET_CODE="${PIPESTATUS[0]}"
  set -e

popd

exit "${RET_CODE}"
