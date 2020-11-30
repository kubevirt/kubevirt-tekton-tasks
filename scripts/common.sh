#!/usr/bin/env bash

function visit() {
  pushd "${1}" > /dev/null
}

function leave() {
  popd > /dev/null
}

export EXCLUDED_NON_IMAGE_MODULES="shared|sharedtest|tests"

declare -A TASK_NAME_TO_ENV_NAME
declare -A TASK_NAME_TO_IMAGE

export CREATE_VM_IMAGE="${CREATE_VM_IMAGE:-}"
TASK_NAME_TO_ENV_NAME["create-vm"]="CREATE_VM_IMAGE"
TASK_NAME_TO_IMAGE["create-vm"]="${CREATE_VM_IMAGE}"
