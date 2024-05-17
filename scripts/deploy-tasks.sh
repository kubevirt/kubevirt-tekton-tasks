#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"
DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(kubectl config view --minify --output 'jsonpath={..namespace}')}"

source "${SCRIPT_DIR}/common.sh"

if [ -z "${RELEASE_VERSION}" ]; then
  echo "RELEASE_VERSION is not defined"
  exit 1
fi

make generate-yaml-tasks

# run only for specified tasks in script arguments
# or default to all if no arguments specified
visit "${REPO_DIR}/release/tasks"
  if [[ $# -eq 0 ]]; then
    TASK_NAMES=(*)
  else
    TASK_NAMES=("$@")
  fi
  for TASK_NAME in ${TASK_NAMES[*]}; do
    if [ ! -d  "${TASK_NAME}" ]; then
      continue
    fi

    # cleanup first
    kubectl delete -f manifests 2> /dev/null || true

    visit "${TASK_NAME}"
      if [[ -z ${TEKTON_TASKS_IMAGE} ]]; then
        kubectl apply -n ${DEPLOY_NAMESPACE} -f "${TASK_NAME}.yaml"
      else
        TASKS_IMAGE="quay.io/kubevirt/tekton-tasks"
        DISK_VIRT_IMAGE="quay.io/kubevirt/tekton-tasks-disk-virt"
        
        sed -i "s!\"${DISK_VIRT_IMAGE}.*!${TEKTON_TASKS_DISK_VIRT_IMAGE}!g" "${TASK_NAME}.yaml" 
        sed -i "s!\"${TASKS_IMAGE}.*!${TEKTON_TASKS_IMAGE}!g" "${TASK_NAME}.yaml"
        kubectl apply -n ${DEPLOY_NAMESPACE} -f "${TASK_NAME}.yaml"
      fi
    leave
  done
leave
