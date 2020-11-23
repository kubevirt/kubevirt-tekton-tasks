#!/usr/bin/env bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

RET_CODE=0

visit "${REPO_DIR}/modules"
  for MODULE_DIR in $(ls | grep -vE "^(tests)$"); do
    visit "$MODULE_DIR"
      if [ -f go.mod ]; then
        make test
        CURRENT_RET_CODE=$?
        if [ "${CURRENT_RET_CODE}" -ne 0 ]; then
          RET_CODE=${CURRENT_RET_CODE}
        fi
      fi
    leave
  done
leave

exit $RET_CODE
