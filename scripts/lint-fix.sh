#!/usr/bin/env bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

visit "${REPO_DIR}"
  gofmt -w cmd modules test
  status=$?
leave
exit "${status}"
