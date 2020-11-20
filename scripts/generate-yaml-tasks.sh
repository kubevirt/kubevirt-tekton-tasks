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
RESOURCE_TYPES=(manifests examples)

visit "${REPO_DIR}/templates"
  if [[ $# -eq 0 ]]; then
    TASK_NAMES=(*)
  else
    TASK_NAMES=("$@")
  fi

  for TASK_NAME in ${TASK_NAMES[*]}; do
    visit "${TASK_NAME}"
      ansible-playbook generate-task.yaml
      for RESOURCE_TYPE in ${RESOURCE_TYPES[*]}; do
        DESTINATION_DIR="${REPO_DIR}/tasks/${TASK_NAME}/${RESOURCE_TYPE}"
        SOURCE_DIR="${REPO_DIR}/templates/${TASK_NAME}/dist/${RESOURCE_TYPE}"

        if [ "${DRY_RUN}" == "false" ] && [ -d "${SOURCE_DIR}" ]; then
          rm -rf "${DESTINATION_DIR}"
          cp -r "${SOURCE_DIR}" "${DESTINATION_DIR}"
        fi
      done
      if [ "${DRY_RUN}" == "false" ]; then
        rm -rf "${REPO_DIR}/templates/${TASK_NAME}/dist"
      fi
    leave
  done
leave
