#!/bin/bash

source "./automation/e2e-source.sh"

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
  local timeout=2000
  while  [ $current_time -lt $timeout ]; do
    sleep $sample
  
    condition_reason=$(oc get pipelinerun -l pipelinerun="$1"-run -o json | jq -r '.items[0].status.conditions[]| select(.type=="Succeeded").reason')
    if [ "$condition_reason" = "Succeeded" ]; then
      echo "Pipelinerun $1 succeeded"
      break
    fi

    if [ "$condition_reason" = "Failed" ]; then
      echo "Pipelinerun $1 failed"
      exit 1
    fi

    (( current_time+=sample ))
    if [ $current_time -ge $timeout ]; then
      echo "Pipelinerun $1 timed out"
      exit 1
    fi
  done
}

cp -L "$KUBECONFIG" /tmp/kubeconfig && export KUBECONFIG=/tmp/kubeconfig
export DEPLOY_NAMESPACE=kubevirt

if ! hash kubectl 2>/dev/null; then
  oc_path="$(which oc)"
  dir_name="dirname ${oc_path}"
  pushd "$(${dir_name})" || return
  ln -s oc kubectl
  popd || return
fi

# switch to faster storage class for example pipelines tests (slower storage class is causing timeouts due 
# to not able to copy whole windows disk)
if ! oc get storageclass | grep -q 'ssd-csi (default)' > /dev/null; then
  oc annotate storageclass ssd-csi storageclass.kubernetes.io/is-default-class=true --overwrite
  oc annotate storageclass standard-csi storageclass.kubernetes.io/is-default-class- --overwrite
fi

# Deploy resources
echo "Deploying resources"
./automation/e2e-deploy-resources.sh

# SECRET
accessKeyId="/tmp/secrets/accessKeyId"
secretKey="/tmp/secrets/secretKey"
namespace="kubevirt"

oc project ${namespace}

if test -f "$accessKeyId" && test -f "$secretKey"; then
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
fi

echo "Creating datavolume with windows iso"
oc apply -n ${namespace} -f "automation/e2e-pipelines/test-files/${TARGET}-dv.yaml"

echo "Waiting for pvc to be created"
wait_until_exists "pvc -n ${namespace} iso-dv -o jsonpath='{.metadata.annotations.cdi\.kubevirt\.io/storage\.pod\.phase}'"
oc wait -n ${namespace}  pvc iso-dv --timeout=10m --for=jsonpath='{.metadata.annotations.cdi\.kubevirt\.io/storage\.pod\.phase}'='Succeeded'

echo "Create config map for http server"
oc apply -n ${namespace} -f "automation/e2e-pipelines/test-files/configmap.yaml"

echo "Deploying http-server to serve iso file to pipeline"
oc apply -n ${namespace} -f "automation/e2e-pipelines/test-files/http-server.yaml"

wait_until_exists "pods -n ${namespace} -l app=http-server"

echo "Waiting for http server to be ready"
oc wait -n ${namespace}  --for=condition=Ready pod -l app=http-server --timeout=10m

oc project kubevirt

#deploy tasks and pipelines
if [[ "$DEV_MODE" == "true" ]]; then
  make cluster-sync
else
  make deploy
fi

./scripts/deploy-pipelines.sh

wait_until_exists "pipeline windows-efi-installer -n ${namespace}"
wait_until_exists "pipeline windows-customize -n ${namespace}"

# Run windows10/11/2022-installer pipeline
echo "Running ${TARGET}-installer pipeline"
oc create -n ${namespace} -f "automation/e2e-pipelines/test-files/${TARGET}-installer-pipelinerun.yaml"
wait_until_exists "pipelinerun -n ${namespace} -l pipelinerun=${TARGET}-installer-run"

# Wait for pipeline to finish
echo "Waiting for pipeline to finish"
wait_for_pipelinerun "${TARGET}-installer"

# Run windows-customize pipeline
echo "Running windows-customize pipeline"
oc create -n ${namespace} -f "automation/e2e-pipelines/test-files/${TARGET}-customize-pipelinerun.yaml"
wait_until_exists "pipelinerun -n ${namespace} -l pipelinerun=${TARGET}-customize-run"

# Wait for pipeline to finish
echo "Waiting for pipeline to finish"
wait_for_pipelinerun "${TARGET}-customize"
