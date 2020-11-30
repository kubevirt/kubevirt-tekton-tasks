#!/usr/bin/env bash

set -o pipefail

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

RET_CODE=0

ARTIFACT_DIR="${ARTIFACT_DIR:=dist}"
ARTIFACT_DIR="$(readlink -m "${ARTIFACT_DIR}")"

mkdir -p "${ARTIFACT_DIR}"
LINT_OUT="${ARTIFACT_DIR}/lint.out"
rm -f "${LINT_OUT}"

visit "${REPO_DIR}/modules"
  for MODULE_DIR in *; do
    visit "$MODULE_DIR"
      if [ -f go.mod ]; then
        if [ -n "$(gofmt -d $(ls -d */ | grep -v "^vendor/") | tee -a "${LINT_OUT}")" ]; then
          RET_CODE=1
        fi
      fi
    leave
  done
leave

cat "${LINT_OUT}"

if [ "$OPENSHIFT_CI" != "true" ]; then
  rm -f "${LINT_OUT}"
  if [ -z "$(ls -A "${ARTIFACT_DIR}")" ]; then
    rm -rf "${ARTIFACT_DIR}"
  fi
fi

exit $RET_CODE
