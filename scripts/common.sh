#!/usr/bin/env bash

function visit() {
  pushd "${1}" > /dev/null
}

function leave() {
  popd > /dev/null
}

export IS_OPENSHIFT="false"
export IS_MINIKUBE="false"

if kubectl get projects > /dev/null 2>&1; then
  export IS_OPENSHIFT="true"
elif minikube status | grep -q Running; then
  export IS_MINIKUBE="true"
fi

export EXCLUDED_NON_IMAGE_MODULES="shared|sharedtest|tests"
export DEPENDENCY_MODULES="shared|sharedtest"

declare -A IMAGE_MODULE_NAME_TO_ENV_NAME
declare -A TASK_NAME_TO_IMAGE

export CREATE_VM_IMAGE="${CREATE_VM_IMAGE:-}"
IMAGE_MODULE_NAME_TO_ENV_NAME["create-vm"]="CREATE_VM_IMAGE"
TASK_NAME_TO_IMAGE["create-vm-from-manifest"]="${CREATE_VM_IMAGE}"
TASK_NAME_TO_IMAGE["create-vm-from-template"]="${CREATE_VM_IMAGE}"

export EXECUTE_IN_VM_IMAGE="${EXECUTE_IN_VM_IMAGE:-}"
IMAGE_MODULE_NAME_TO_ENV_NAME["execute-in-vm"]="EXECUTE_IN_VM_IMAGE"
TASK_NAME_TO_IMAGE["execute-in-vm"]="${EXECUTE_IN_VM_IMAGE}"
TASK_NAME_TO_IMAGE["cleanup-vm"]="${EXECUTE_IN_VM_IMAGE}"

export DISK_VIRT_CUSTOMIZE_IMAGE="${DISK_VIRT_CUSTOMIZE_IMAGE:-}"
IMAGE_MODULE_NAME_TO_ENV_NAME["disk-virt-customize"]="DISK_VIRT_CUSTOMIZE_IMAGE"
TASK_NAME_TO_IMAGE["disk-virt-customize"]="${DISK_VIRT_CUSTOMIZE_IMAGE}"

export GENERATE_SSH_KEYS_IMAGE="${GENERATE_SSH_KEYS_IMAGE:-}"
IMAGE_MODULE_NAME_TO_ENV_NAME["generate-ssh-keys"]="GENERATE_SSH_KEYS_IMAGE"
TASK_NAME_TO_IMAGE["generate-ssh-keys"]="${GENERATE_SSH_KEYS_IMAGE}"
