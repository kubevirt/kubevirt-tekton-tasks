#!/usr/bin/env bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

visit "${REPO_DIR}/tasks"
  for TASK_NAME in *; do
    visit "${TASK_NAME}"
      kubectl delete -f manifests 2> /dev/null
    leave
  done
leave
