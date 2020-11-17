#!/usr/bin/env bash

set -e

RET_CODE=0

ARTIFACT_DIR="${ARTIFACT_DIR:=dist}"
ARTIFACT_DIR="$(readlink -m "${ARTIFACT_DIR}")"
TEST_OUT="${ARTIFACT_DIR}/test-yaml-consistency.out"

rm -rf "${TEST_OUT}"
mkdir -p "${ARTIFACT_DIR}"

pushd tasks > /dev/null || exit 1

  for TASK_DIR in *; do
    pushd "$TASK_DIR" > /dev/null || continue
      set +e
      make test-yaml-consistency | tee -a "${TEST_OUT}"
      CURRENT_RET_CODE=${PIPESTATUS[0]}
      if [ "${CURRENT_RET_CODE}" -ne 0 ]; then
        RET_CODE=${CURRENT_RET_CODE}
      fi
      set -e
    popd > /dev/null || continue
  done
popd > /dev/null || exit 1

exit $RET_CODE
