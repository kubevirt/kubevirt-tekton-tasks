#!/usr/bin/env bash

RET_CODE=0

ARTIFACTS_DIR=${ARTIFACTS_DIR:=dist}
ARTIFACTS_DIR=$(readlink -m "${ARTIFACTS_DIR}")
TEST_OUT="${ARTIFACTS_DIR}/test.out"
COVER_OUT="${ARTIFACTS_DIR}/cover.out"
JUNIT_XML="${ARTIFACTS_DIR}/junit.xml"
COVERAGE_HTML="${ARTIFACTS_DIR}/coverage.html"
FAKE_GOPATH_ROOT="/tmp/goroot-kubevirt-tekton-tasks"
FAKE_KV_GOPATH="${FAKE_GOPATH_ROOT}/src/github.com/kubevirt"
FAKE_REPO_GOPATH="${FAKE_KV_GOPATH}/kubevirt-tekton-tasks"

rm -rf "${TEST_OUT}" "${COVER_OUT}" "${JUNIT_XML}" "${COVERAGE_HTML}" "${FAKE_GOPATH_ROOT}"
mkdir -p "${ARTIFACTS_DIR}"

pushd modules > /dev/null || exit 1
  export DIST_DIR=dist

  for MODULE_DIR in $(ls | grep -vE "^(tests)$"); do
    pushd "$MODULE_DIR" > /dev/null || continue
      if ! make test-verbose; then
        RET_CODE=2
      fi
      cat ${DIST_DIR}/test.out >> "${TEST_OUT}"

      if [ -f "${COVER_OUT}" ]; then
        sed "/^mode.*/d" dist/cover.out >> "${COVER_OUT}" # remove first line with mode
      else
        cp ${DIST_DIR}/cover.out "${COVER_OUT}"
      fi

    popd > /dev/null || continue
  done
popd > /dev/null || exit 1

go-junit-report < "${TEST_OUT}" > "${JUNIT_XML}"

mkdir -p "${FAKE_KV_GOPATH}"
ln -s "$(pwd)" "${FAKE_REPO_GOPATH}"
GOPATH="${FAKE_GOPATH_ROOT}"
export "GOPATH=${GOPATH}"

go tool cover -html "${COVER_OUT}" -o "${COVERAGE_HTML}"

rm -rf "${FAKE_GOPATH_ROOT}"

exit $RET_CODE
