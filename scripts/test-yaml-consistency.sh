#!/usr/bin/env bash

set -e

function testConsistency() {
  set +e
    diff -r "${SOURCE_DIR}" "${DESTINATION_DIR}"| tee -a "${TEST_OUT}"
    CURRENT_RET_CODE=${PIPESTATUS[0]}
    if [ "${CURRENT_RET_CODE}" -ne 0 ]; then
      RET_CODE=${CURRENT_RET_CODE}
      echo -e "\n${DESTINATION_DIR} yaml files are not fresh and should be regenerated!" | tee -a "${TEST_OUT}" 1>&2
    fi
  set -e
}

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
  DRY_RUN=true "${SCRIPT_DIR}/generate-yaml-tasks.sh" | tee -a "${TEST_OUT}"
  visit "templates"
    for TASK_NAME in *; do
      visit "${TASK_NAME}"
        visit dist
          for RESOURCE_TYPE in *; do
            DESTINATION_DIR="${REPO_DIR}/tasks/${TASK_NAME}/${RESOURCE_TYPE}"
            SOURCE_DIR="${REPO_DIR}/templates/${TASK_NAME}/dist/${RESOURCE_TYPE}"
            testConsistency
          done
        leave
        rm -rf dist
      leave
    done
  leave

  visit manifests
    for RELEASE_TYPE in kubernetes okd; do
      visit "${RELEASE_TYPE}"
        DESTINATION_DIR="${REPO_DIR}/manifests/${RELEASE_TYPE}"
        SOURCE_DIR="${DESTINATION_DIR}-dist"
        testConsistency
        rm -rf "${SOURCE_DIR}"
      leave
    done
  leave
leave

 if [ "${RET_CODE}" -eq 0 ]; then
   echo -e "\nConsistency: OK!"
 fi

exit ${RET_CODE}
