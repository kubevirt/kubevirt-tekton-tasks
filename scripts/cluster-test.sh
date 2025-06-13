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

kubectl get namespaces -o name | grep -Eq "^namespace/$DEPLOY_NAMESPACE$" || kubectl create namespace "$DEPLOY_NAMESPACE" > /dev/null

# for purposes of disk-virt tests, we need to create SCC, which will allow pod to 
# set FSgroup to uid 107
oc apply -f - <<EOF
apiVersion: security.openshift.io/v1
kind: SecurityContextConstraints
metadata:
  name: tekton-tasks-scc
allowPrivilegedContainer: false
allowHostDirVolumePlugin: false
allowHostIPC: false
allowHostNetwork: false
allowHostPID: false
allowHostPorts: false
seccompProfiles:
  - "runtime/default"
readOnlyRootFilesystem: false
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: MustRunAs
fsGroup:
  type: MustRunAs
  ranges:
    - min: 107
      max: 107
users:
- system:serviceaccount:${DEPLOY_NAMESPACE}:pipeline
EOF
pushd test || exit
  rm -rf dist
  mkdir dist

  set +ex
  set -o pipefail

  ginkgo -r -p --randomize-all --randomize-suites --fail-on-pending --trace --race --nodes="${NUM_NODES}" -- \
    --deploy-namespace="${DEPLOY_NAMESPACE}" \
    --kubeconfig-path="${KUBECONFIG}" \
    --ginkgo.junit-report="${ARTIFACT_DIR}/xunit_results.xml" \
    --storage-class="${STORAGE_CLASS}" \
    --debug="${DEBUG}" | tee "${TEST_OUT}"

  RET_CODE="${PIPESTATUS[0]}"
  set -e

popd

oc delete scc tekton-tasks-scc
exit "${RET_CODE}"
