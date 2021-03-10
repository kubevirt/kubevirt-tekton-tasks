#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"


# run only for specified tasks in script arguments
# or default to all if no arguments specified

DRY_RUN="${DRY_RUN:=false}"

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"
RESOURCE_TYPES=(manifests examples README.md)
RELEASE_TYPES=(kubernetes openshift)

if [[ $# -eq 0 ]]; then
  visit "${REPO_DIR}/templates"
    TASK_NAMES=(*)
  leave
else
  TASK_NAMES=("$@")
fi

function generateTaskResources() {
  for TASK_NAME in ${TASK_NAMES[*]}; do
    visit "${REPO_DIR}/templates/${TASK_NAME}"
      ansible-playbook generate-task.yaml
      for RESOURCE_TYPE in ${RESOURCE_TYPES[*]}; do
        DESTINATION_PARENT_DIR="${REPO_DIR}/tasks/${TASK_NAME}"
        DESTINATION="${DESTINATION_PARENT_DIR}/${RESOURCE_TYPE}"
        SOURCE="${REPO_DIR}/templates/${TASK_NAME}/dist/${RESOURCE_TYPE}"

        if [ "${DRY_RUN}" == "false" ] && [ -e "${SOURCE}" ]; then
          mkdir -p "${DESTINATION_PARENT_DIR}"
          rm -rf "${DESTINATION}"
          cp -r "${SOURCE}" "${DESTINATION}"
        fi
      done
      if [ "${DRY_RUN}" == "false" ]; then
        rm -rf "${REPO_DIR}/templates/${TASK_NAME}/dist"
      fi
    leave
  done
}

function combineTaskManifestsIntoRelease() {
  for RELEASE_TYPE in ${RELEASE_TYPES[*]}; do
    RESULT_DIR="${REPO_DIR}/manifests/${RELEASE_TYPE}"
    if [ "${DRY_RUN}" == "false" ]; then
      rm -rf "${RESULT_DIR}"
    else
      RESULT_DIR="${RESULT_DIR}-dist"
    fi

    mkdir -p "${RESULT_DIR}"
    RESULT_FILE="${RESULT_DIR}/kubevirt-tekton-tasks.yaml"
    visit "${REPO_DIR}/tasks"
      for TASK_NAME in *; do
        CONFIG_FILE="../configs/${TASK_NAME}.yaml"
        IS_TASK_OPENSHIFT="$(sed -n  's/^is_openshift *: *//p' ${CONFIG_FILE})"
        if [ "${RELEASE_TYPE}" != openshift ] && [ "${IS_TASK_OPENSHIFT}" == true ]; then
          continue
        fi

        cat "${TASK_NAME}"/manifests/* >> "${RESULT_FILE}"
      done
    leave
  done
}

generateTaskResources
combineTaskManifestsIntoRelease
