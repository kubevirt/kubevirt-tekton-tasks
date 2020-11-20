#!/usr/bin/env bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

rm -rf dist

visit "${REPO_DIR}"
  visit "templates"
    for TASK_NAME in *; do
      visit "${TASK_NAME}"
        rm -rf dist
      leave
    done
  leave

  visit "modules"
    for MODULE in *; do
      visit "${MODULE}"
        make clean
      leave
    done
  leave
leave

