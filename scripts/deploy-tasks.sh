#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"
DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(oc config view --minify --output 'jsonpath={..namespace}')}"

source "${SCRIPT_DIR}/common.sh"

if [ -z "${RELEASE_VERSION}" ]; then
  echo "RELEASE_VERSION is not defined"
  exit 1
fi

# run only for specified tasks in script arguments
# or default to all if no arguments specified

visit "${REPO_DIR}/tasks"
  if [[ $# -eq 0 ]]; then
    TASK_NAMES=(*)
  else
    TASK_NAMES=("$@")
  fi
  for TASK_NAME in ${TASK_NAMES[*]}; do
    if [ ! -d  "${TASK_NAME}" ]; then
      continue
    fi
    CONFIG_FILE="${REPO_DIR}/configs/${TASK_NAME}.yaml"
    MAIN_IMAGE="$(sed -n  's/^main_image *: *//p' "${CONFIG_FILE}"):${RELEASE_VERSION}"
    CUSTOM_IMAGE="${TEKTON_TASKS_IMAGE}"

    if [[ "${TASK_NAME}" =~ "disk-virt" ]]; then
      CUSTOM_IMAGE="${TEKTON_TASKS_DISK_VIRT_IMAGE}"
    fi

    # cleanup first
    oc delete -f manifests 2> /dev/null || true

    visit "${TASK_NAME}"

      if [[ -z ${CUSTOM_IMAGE} ]]; then
        oc apply -n ${DEPLOY_NAMESPACE} -f "${TASK_NAME}.yaml"
      else
        sed "s!${MAIN_IMAGE}!${CUSTOM_IMAGE}!g" "${TASK_NAME}.yaml" | oc apply -n ${DEPLOY_NAMESPACE} -f -
      fi
    leave
  done
leave
