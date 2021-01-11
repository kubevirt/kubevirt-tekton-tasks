#!/usr/bin/env bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

visit "${REPO_DIR}/modules"
  for MODULE_DIR in shared sharedtest $(ls | grep -vE "^(${DEPENDENCY_MODULES})$"); do
    visit "$MODULE_DIR"
      if [ -f go.mod ]; then
        echo vendoring "$MODULE_DIR"
        go mod tidy
        go mod vendor
      fi
    leave
  done
leave
