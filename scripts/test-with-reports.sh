#!/usr/bin/env bash

set -o pipefail

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

RET_CODE=0

ARTIFACT_DIR="${ARTIFACT_DIR:=dist}"
ARTIFACT_DIR="$(readlink -m "${ARTIFACT_DIR}")"
TEST_OUT="${ARTIFACT_DIR}/test.out"
COVER_OUT="${ARTIFACT_DIR}/coverage.out"
JUNIT_XML="${ARTIFACT_DIR}/junit.xml"
COVERAGE_HTML="${ARTIFACT_DIR}/coverage.html"
FAKE_GOPATH_ROOT="/tmp/goroot-kubevirt-tekton-tasks"
FAKE_KV_GOPATH="${FAKE_GOPATH_ROOT}/src/github.com/kubevirt"
FAKE_REPO_GOPATH="${FAKE_KV_GOPATH}/kubevirt-tekton-tasks"

rm -rf "${TEST_OUT}" "${COVER_OUT}" "${JUNIT_XML}" "${COVERAGE_HTML}" "${FAKE_GOPATH_ROOT}"
mkdir -p "${ARTIFACT_DIR}"

visit "${REPO_DIR}/modules"
  for MODULE_DIR in */ ; do
    visit "$MODULE_DIR"
      DIST_DIR=dist
      mkdir -p ${DIST_DIR}
      go test -v -coverprofile=${DIST_DIR}/coverage.out -covermode=atomic \
        $(go list ./... | grep -v utilstest) | tee ${DIST_DIR}/test.out
      CURRENT_RET_CODE=$?
      if [ "${CURRENT_RET_CODE}" -ne 0 ]; then
        RET_CODE=${CURRENT_RET_CODE}
      fi
      cat ${DIST_DIR}/test.out >> "${TEST_OUT}"

      if [ -f "${COVER_OUT}" ]; then
        sed "/^mode.*/d" dist/coverage.out >> "${COVER_OUT}" # remove first line with mode
      else
        cp ${DIST_DIR}/coverage.out "${COVER_OUT}"
      fi
    leave
  done
leave

if type go-junit-report > /dev/null; then
  sed 's/\x1B\[[0-9;]\{1,\}[A-Za-z]//g' "${TEST_OUT}" | go-junit-report > "${JUNIT_XML}"
fi

mkdir -p "${FAKE_KV_GOPATH}"
ln -s "$(pwd)" "${FAKE_REPO_GOPATH}"
GOPATH="${FAKE_GOPATH_ROOT}"
export "GOPATH=${GOPATH}"

go tool cover -html "${COVER_OUT}" -o "${COVERAGE_HTML}"

rm -rf "${FAKE_GOPATH_ROOT}"

exit $RET_CODE
