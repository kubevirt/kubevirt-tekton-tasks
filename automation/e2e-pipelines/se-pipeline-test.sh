#!/bin/bash

set -e

source "./automation/e2e-source.sh"

function wait_until_exists() {
  timeout 10m bash <<- EOF
  until $client get $1; do
    sleep 5
  done
EOF
}

function wait_for_pipelinerun() {
  local sample=60
  local current_time=0
  local timeout=3660  # 1 hours for SE pipeline + buffer to ensure pipeline timeouts first
  while  [ $current_time -lt $timeout ]; do
    sleep $sample

    # Check if pipelinerun exists
    rc=0
    pipeline_run="$($client get pipelinerun -l pipelinerun="$1-run" -o jsonpath='{.items[0]}')"  || rc=1
    if [[ $rc != 0 ]]; then
      echo "Waiting for pipelinerun to be created..."
      (( current_time+=sample ))
      continue
    fi

    # Get both status and reason
    pipelinerun_name=$(echo "$pipeline_run" | jq -r '.metadata.name')
    condition_status=$(echo "$pipeline_run" | jq -r '.status.conditions[]| select(.type=="Succeeded").status')
    condition_reason=$(echo "$pipeline_run" | jq -r '.status.conditions[]| select(.type=="Succeeded").reason')

    # Check for success (status=True and reason=Succeeded or Completed)
    if [ "$condition_status" = "True" ] && { [ "$condition_reason" = "Succeeded" ] || [ "$condition_reason" = "Completed" ]; }; then
      echo "Pipelinerun $1 succeeded"
      break
    fi

    # Check for failure (status=False)
    if [ "$condition_status" = "False" ]; then
      echo "Pipelinerun $1 failed with reason: $condition_reason"
      # Print logs for debugging
      $client get pipelinerun "$pipelinerun_name" -o yaml
      $client get taskrun -l "tekton.dev/pipelineRun=$pipelinerun_name"
      exit 1
    fi

    (( current_time+=sample ))
    if [ $current_time -ge $timeout ]; then
      echo "Pipelinerun $1 timed out after ${timeout}s"
      echo "Last known status: $condition_status, reason: $condition_reason"
      exit 1
    fi

    echo "Waiting on $pipelinerun_name, current status $condition_status and reason $condition_reason"
  done
}

function cleanup_cluster() {
  echo ""
  echo "=== Cleanup ==="
  # Delete tmp directory
  [[ -n "${SSH_KEY_DIR}" ]] && rm -rf "${SSH_KEY_DIR}"

  # Remove test namespace
  $client delete --ignore-not-found ns/${namespace}

  # Cleanup base resources
  echo "Cleaning up base resources"
  ./automation/e2e-cleanup-resources.sh
}

echo "=== Secure Execution Pipeline CI Test ==="

# Record start time
START_TIME=$(date +%s)

# Date
date

# Set KUBECONFIG if not already set
if [ -z "$KUBECONFIG" ]; then
  export KUBECONFIG="$HOME/.kube/config"
fi

cp -L "$KUBECONFIG" /tmp/kubeconfig && export KUBECONFIG=/tmp/kubeconfig
export DEPLOY_NAMESPACE="test-kubevirt-tekton-tasks-secure-execution"

# Ensure we have kubectl or oc, and set client variable
if hash oc 2>/dev/null; then
  client="oc"
elif hash kubectl 2>/dev/null; then
  client="kubectl"
else
  echo "ERROR: Neither kubectl nor oc command found"
  exit 1
fi
export client
echo "Using client: $client"

namespace="${DEPLOY_NAMESPACE}"

# Check architecture
ARCH=$(uname -m)
echo "Running on architecture: $ARCH"

if [[ "$ARCH" != "s390x" ]]; then
  echo "ERROR: SE pipeline requires s390x architecture"
  echo "Current architecture is $ARCH - cannot run SE pipeline tests"
  exit 1
fi

echo "Running SE pipeline integration test on s390x"

# Deploy base resources
echo "Deploying base resources"
./automation/e2e-deploy-resources.sh

# Ensure cleanup is called
trap cleanup_cluster EXIT

# Create test namespace if it doesn't already exist
if ! $client get namespace ${namespace} > /dev/null 2>&1; then
  $client create ns ${namespace}
fi

# Set namespace context
if [ "$client" = "oc" ]; then
  $client project ${namespace}
else
  $client config set-context --current --namespace=${namespace}
fi

# Create secret for container disk puller
accessKeyId="/tmp/secrets/accessKeyId"
secretKey="/tmp/secrets/secretKey"

if test -f "$accessKeyId" && test -f "$secretKey"; then
  echo "Creating container disk puller secret from provided credentials"
  id=$(cat $accessKeyId | tr -d '\n' | base64)
  token=$(cat $secretKey | tr -d '\n' | base64 | tr -d ' \n')

  $client apply -n ${namespace} -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: tekton-tasks-container-disk-puller
  namespace: ${namespace}
type: Opaque
data:
  accessKeyId: "${id}"
  secretKey: "${token}"
EOF
else
  echo "Creating dummy container disk puller secret (using public registry)"
  $client apply -n ${namespace} -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: tekton-tasks-container-disk-puller
  namespace: ${namespace}
type: Opaque
data:
  accessKeyId: ""
  secretKey: ""
EOF
fi

# Deploy tasks and pipelines
# Set namespace context again (in case it changed)
if [ "$client" = "oc" ]; then
  $client project ${namespace}
else
  $client config set-context --current --namespace=${namespace}
fi

if [[ "$DEV_MODE" == "true" ]]; then
  make cluster-sync
else
  make deploy
fi

# Deploy SE pipeline resources
echo "Deploying SE pipeline RBAC"
$client apply -n ${namespace} -f "templates-pipelines/secure-execution-installer/manifests/se-pipeline-rbac.yaml"

echo "Deploying SE pipeline"
$client apply -n ${namespace} -f "release/pipelines/secure-execution-installer/secure-execution-installer.yaml"

wait_until_exists "pipeline secure-execution-installer -n ${namespace}"

# Deploy http-server to serve the configmap
echo "Deploying http-server for serving configmap"
$client apply -n ${namespace} -f "automation/e2e-pipelines/test-files/configmap.yaml"
$client create -n ${namespace} configmap secure-execution-installer --from-file="release/pipelines/secure-execution-installer/configmaps/secure-execution-installer-configmaps.yaml"
$client apply -n ${namespace} -f "automation/e2e-pipelines/test-files/http-server.yaml"
$client patch -n ${namespace} deployment http-server --type='json' -p='[{"op":"replace","path":"/spec/template/spec/volumes/1","value":{"name":"iso-dv","configMap":{"name":"secure-execution-installer"}}}]'
HTTP_SERVER_IP="$($client get -n ${namespace} svc http-server -o jsonpath='{.spec.clusterIP}')"

# Run SE installer pipeline
echo "Generate a temporary SSH keypair"
SSH_KEY_DIR=$(mktemp -d)
ssh-keygen -t ed25519 -N "" -f "${SSH_KEY_DIR}/sec-exec-vm-key" -C "se-pipeline-test"
SSH_PUBKEY=$(cat "${SSH_KEY_DIR}/sec-exec-vm-key.pub")

echo "Fetching hostkey document from cluster"
HOST_DOC="$($client get -n kubevirt-prow-jobs configmaps secex-hostkey -o json | jq -r '.data."secex-hostkey.crt"' | base64 -w 0)"
if [ -z "$HOST_DOC" ]; then
  echo "Missing hostkey document, make sure to create secex-hostkey configmap in kubevirt-prow-jobs namespace"
  exit 1
fi
echo "Running fedora-se-installer pipeline"
PIPELINE_RUN="$(cat "automation/e2e-pipelines/test-files/fedora-se-installer-pipelinerun.yaml")"
PIPELINE_RUN="${PIPELINE_RUN//__HOST_DOC__/$HOST_DOC}"
PIPELINE_RUN="${PIPELINE_RUN//__SSH_PUBKEY__/$SSH_PUBKEY}"
PIPELINE_RUN="${PIPELINE_RUN//__HTTP_SERVER_IP__/$HTTP_SERVER_IP}"
echo "${PIPELINE_RUN}" | $client create -n ${namespace} -f -
wait_until_exists "pipelinerun -n ${namespace} -l pipelinerun=fedora-se-installer-run"

# Wait for pipeline to finish
echo "Waiting for SE pipeline to finish (this may take up to 20 min)"
wait_for_pipelinerun "fedora-se-installer"

# Verify SE is enabled in the VM via virtctl ssh
echo ""
echo "=== Verifying Secure Execution Status ==="
VM_NAME="sec-exec-vm"

VMI_PHASE=$($client get vmi "$VM_NAME" -n ${namespace} -o jsonpath='{.status.phase}')
echo "VMI phase: $VMI_PHASE"
if [ "$VMI_PHASE" != "Running" ]; then
  echo "ERROR: VM $VM_NAME is not running (phase: $VMI_PHASE) — cannot verify SE"
  exit 1
fi

echo "Verifying Secure Execution status of VM"
PROT_VIRT=$(virtctl ssh \
  -i "${SSH_KEY_DIR}/sec-exec-vm-key" \
  -t="-o StrictHostKeyChecking=no" \
  -t="-o ConnectTimeout=30" \
  --username=core \
  -c "cat /sys/firmware/uv/prot_virt_guest" \
  "vm/${VM_NAME}")

if [ "${PROT_VIRT}" != "1" ]; then
  echo "ERROR: Secure Execution is NOT enabled — /sys/firmware/uv/prot_virt_guest returned '${PROT_VIRT}' (expected '1')"
  exit 1
fi

echo "Secure Execution verified: /sys/firmware/uv/prot_virt_guest = 1"

# Calculate and display total time
END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))
MINUTES=$((TOTAL_TIME / 60))
SECONDS=$((TOTAL_TIME % 60))

echo ""
echo "=== SE Pipeline Integration Test Completed Successfully ==="
echo "Total execution time: ${MINUTES}m ${SECONDS}s (${TOTAL_TIME} seconds)"
