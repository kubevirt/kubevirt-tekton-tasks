#!/bin/bash

set -e

source "./automation/e2e-source.sh"

# Create oc function if only kubectl is available
if ! hash oc 2>/dev/null && hash kubectl 2>/dev/null; then
  function oc() {
    kubectl "$@"
  }
  export -f oc
fi

function wait_until_exists() {
  timeout 10m bash <<- EOF
  until oc get $1; do
    sleep 5
  done
EOF
}

function wait_for_pipelinerun() {
  local sample=10
  local current_time=0
  local timeout=7200  # 2 hours for SE pipeline
  while  [ $current_time -lt $timeout ]; do
    sleep $sample
  
    # Check if pipelinerun exists
    if ! oc get pipelinerun -l pipelinerun="$1"-run -o json | jq -e '.items[0]' > /dev/null 2>&1; then
      echo "Waiting for pipelinerun to be created..."
      (( current_time+=sample ))
      continue
    fi
    
    # Get both status and reason
    condition_status=$(oc get pipelinerun -l pipelinerun="$1"-run -o json | jq -r '.items[0].status.conditions[]| select(.type=="Succeeded").status')
    condition_reason=$(oc get pipelinerun -l pipelinerun="$1"-run -o json | jq -r '.items[0].status.conditions[]| select(.type=="Succeeded").reason')
    
    # Check for success (status=True and reason=Succeeded or Completed)
    if [ "$condition_status" = "True" ] && { [ "$condition_reason" = "Succeeded" ] || [ "$condition_reason" = "Completed" ]; }; then
      echo "Pipelinerun $1 succeeded"
      break
    fi

    # Check for failure (status=False)
    if [ "$condition_status" = "False" ]; then
      echo "Pipelinerun $1 failed with reason: $condition_reason"
      # Print logs for debugging
      oc get pipelinerun -l pipelinerun="$1"-run -o yaml
      exit 1
    fi

    (( current_time+=sample ))
    if [ $current_time -ge $timeout ]; then
      echo "Pipelinerun $1 timed out after ${timeout}s"
      echo "Last known status: $condition_status, reason: $condition_reason"
      exit 1
    fi
  done
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
export DEPLOY_NAMESPACE=kubevirt

# Ensure we have kubectl or oc
if ! hash kubectl 2>/dev/null && ! hash oc 2>/dev/null; then
  echo "ERROR: Neither kubectl nor oc command found"
  exit 1
fi

# Create kubectl symlink if only oc is available
if ! hash kubectl 2>/dev/null && hash oc 2>/dev/null; then
  oc_path="$(which oc)"
  dir_name=$(dirname "${oc_path}")
  pushd "${dir_name}" || exit 1
  ln -s oc kubectl
  popd || exit 1
fi


namespace="kubevirt"

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

# Set namespace context (use config set-context for kubectl, project for oc)
if hash oc 2>/dev/null && ! hash kubectl 2>/dev/null; then
  oc project ${namespace}
else
  kubectl config set-context --current --namespace=${namespace}
fi

# Create secret for container disk puller
accessKeyId="/tmp/secrets/accessKeyId"
secretKey="/tmp/secrets/secretKey"

if test -f "$accessKeyId" && test -f "$secretKey"; then
  echo "Creating container disk puller secret from provided credentials"
  id=$(cat $accessKeyId | tr -d '\n' | base64)
  token=$(cat $secretKey | tr -d '\n' | base64 | tr -d ' \n')

  oc apply -n ${namespace} -f - <<EOF
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
  oc apply -n ${namespace} -f - <<EOF
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

# Clean up any existing ISO DataVolume from previous runs
echo "Cleaning up any existing ISO DataVolume"
oc delete dv iso-dv -n ${namespace} --ignore-not-found
oc delete pvc iso-dv -n ${namespace} --ignore-not-found

# Create ISO DataVolume
echo "Creating datavolume with Fedora ISO"
oc apply -n ${namespace} -f "automation/e2e-pipelines/test-files/fedora-se-dv.yaml"

echo "Waiting for ISO DV to be ready"
wait_until_exists "pvc -n ${namespace} iso-dv -o jsonpath='{.metadata.annotations.cdi\.kubevirt\.io/storage\.pod\.phase}'"
oc wait -n ${namespace} pvc iso-dv --timeout=15m --for=jsonpath='{.metadata.annotations.cdi\.kubevirt\.io/storage\.pod\.phase}'='Succeeded'

# Deploy HTTP server
echo "Create config map for http server"
oc apply -n ${namespace} -f "automation/e2e-pipelines/test-files/configmap.yaml"

echo "Deploying http-server to serve ISO file"
oc apply -n ${namespace} -f "automation/e2e-pipelines/test-files/http-server.yaml"

wait_until_exists "pods -n ${namespace} -l app=http-server"

echo "Waiting for http server to be ready"
oc wait -n ${namespace} --for=condition=Ready pod -l app=http-server --timeout=10m

# Deploy tasks and pipelines
# Set namespace context again (in case it changed)
if hash oc 2>/dev/null && ! hash kubectl 2>/dev/null; then
  oc project kubevirt
else
  kubectl config set-context --current --namespace=kubevirt
fi

if [[ "$DEV_MODE" == "true" ]]; then
  make cluster-sync
else
  make deploy
fi

# Deploy SE pipeline resources
echo "Deploying SE pipeline RBAC"
oc apply -n ${namespace} -f "templates-pipelines/secure-execution-installer/manifests/se-pipeline-rbac.yaml"

echo "Deploying SE templates ConfigMap"
oc apply -n ${namespace} -f "templates-pipelines/secure-execution-installer/configmaps/se-templates-configmaps.yaml"

echo "Deploying SE pipeline"
oc apply -n ${namespace} -f "templates-pipelines/secure-execution-installer/manifests/secure-execution-installer.yaml"

wait_until_exists "pipeline secure-execution-installer -n ${namespace}"

# Run SE installer pipeline
echo "Running fedora-se-installer pipeline"
oc create -n ${namespace} -f "automation/e2e-pipelines/test-files/fedora-se-installer-pipelinerun.yaml"
wait_until_exists "pipelinerun -n ${namespace} -l pipelinerun=fedora-se-installer-run"

# Wait for pipeline to finish
echo "Waiting for SE pipeline to finish (this may take up to 20 min)"
wait_for_pipelinerun "fedora-se-installer"

# Verify SE is enabled in the VM
echo ""
echo "=== Verifying Secure Execution Status ==="
VM_NAME="sec-exec-vm"

# Check if VM exists and is running
if oc get vmi "$VM_NAME" -n ${namespace} &>/dev/null; then
  VMI_PHASE=$(oc get vmi "$VM_NAME" -n ${namespace} -o jsonpath='{.status.phase}')
  echo "VM Status: $VMI_PHASE"
  
  if [ "$VMI_PHASE" = "Running" ]; then
    echo "VM is running - SE verification task should have checked /sys/firmware/uv/prot_virt_guest"
    echo "Check pipeline logs for SE verification results:"
    echo "  oc logs -n ${namespace} -l tekton.dev/pipelineTask=verify-se-enabled --tail=50"
    
    # Try to get the verification task logs
    echo ""
    echo "=== SE Verification Task Logs ==="
    if oc get pods -n ${namespace} -l tekton.dev/pipelineTask=verify-se-enabled &>/dev/null; then
      oc logs -n ${namespace} -l tekton.dev/pipelineTask=verify-se-enabled --tail=100 || echo "Could not retrieve verification logs"
    else
      echo "Verification task pod not found (may have been cleaned up)"
    fi
  else
    echo "WARNING: VM is not running (phase: $VMI_PHASE)"
  fi
else
  echo "WARNING: VM $VM_NAME not found - may have been deleted by golden image creation"
fi

# Calculate and display total time
END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))
MINUTES=$((TOTAL_TIME / 60))
SECONDS=$((TOTAL_TIME % 60))

echo ""
echo "=== SE Pipeline Integration Test Completed Successfully ==="
echo "Total execution time: ${MINUTES}m ${SECONDS}s (${TOTAL_TIME} seconds)"
