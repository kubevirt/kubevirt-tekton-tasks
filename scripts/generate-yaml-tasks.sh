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
RESOURCE_TYPES=(manifests examples readme)

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
      DESTINATION_PARENT_DIR="${REPO_DIR}/tasks/${TASK_NAME}"
      rm -rf "${DESTINATION_PARENT_DIR}"
      mkdir -p "${DESTINATION_PARENT_DIR}"
      for RESOURCE_TYPE in ${RESOURCE_TYPES[*]}; do
        DESTINATION="${DESTINATION_PARENT_DIR}"
        SOURCE="${REPO_DIR}/templates/${TASK_NAME}/dist/${RESOURCE_TYPE}/."

        if [ "${DRY_RUN}" == "false" ] && [ -e "${SOURCE}" ]; then
          if [ "${RESOURCE_TYPE}" == "examples" ]; then
            DESTINATION="${DESTINATION}/${RESOURCE_TYPE}"
          fi

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
    RESULT_DIR="${REPO_DIR}/manifests/"
    if [ "${DRY_RUN}" == "false" ]; then
      rm -rf "${RESULT_DIR}"
    else
      RESULT_DIR="${RESULT_DIR}-dist"
    fi

    mkdir -p "${RESULT_DIR}"
    RESULT_FILE="${RESULT_DIR}/kubevirt-tekton-tasks.yaml"
    visit "${REPO_DIR}/tasks"
      for TASK_NAME in *; do
        cat "${TASK_NAME}"/*.yaml >> "${RESULT_FILE}"
      done
    leave
}

generateTaskResources
combineTaskManifestsIntoRelease
