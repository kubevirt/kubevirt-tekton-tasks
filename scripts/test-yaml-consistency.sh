#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"


RET_CODE=0

ARTIFACT_DIR="${ARTIFACT_DIR:=dist}"
ARTIFACT_DIR="$(readlink -m "${ARTIFACT_DIR}")"
TEST_OUT="${ARTIFACT_DIR}/test-yaml-consistency.out"

rm -rf "${TEST_OUT}"
mkdir -p "${ARTIFACT_DIR}"

visit "${REPO_DIR}"
  DRY_RUN=true make generate-yaml-tasks | tee -a "${TEST_OUT}"
  visit "templates"
    for TASK_NAME in *; do
      visit "${TASK_NAME}"
        visit dist
          for RESOURCE_TYPE in *; do
            DESTINATION_DIR="${REPO_DIR}/tasks/${TASK_NAME}/${RESOURCE_TYPE}"
            SOURCE_DIR="${REPO_DIR}/templates/${TASK_NAME}/dist/${RESOURCE_TYPE}"

            set +e
              diff -r "${SOURCE_DIR}" "${DESTINATION_DIR}"| tee -a "${TEST_OUT}"
              CURRENT_RET_CODE=${PIPESTATUS[0]}
              if [ "${CURRENT_RET_CODE}" -ne 0 ]; then
                RET_CODE=${CURRENT_RET_CODE}
                echo -e "\n${TASK_NAME} yaml files are not fresh and should be regenerated!" | tee -a "${TEST_OUT}" 1>&2
              fi
            set -e
          done
        leave
        rm -rf dist
      leave
    done
  leave
leave

exit $RET_CODE
