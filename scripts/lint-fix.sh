#!/usr/bin/env bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

visit "${REPO_DIR}/modules"
  for MODULE_DIR in $(ls | grep -vE "^(tests)$"); do
    visit "$MODULE_DIR"
      if [ -f go.mod ]; then
        gofmt -w $(ls -d */ | grep -v "^vendor/")
      fi
    leave
  done
leave
