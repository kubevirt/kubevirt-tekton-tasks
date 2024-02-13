#!/usr/bin/env bash

set -ex

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"
visit "${REPO_DIR}/tasks"
    for TASK_NAME in *; do
        mkdir -p "${REPO_DIR}/release/tasks/${TASK_NAME}/${RELEASE_VERSION:1}"
        cp -r "${TASK_NAME}/." "${REPO_DIR}/release/tasks/${TASK_NAME}/${RELEASE_VERSION:1}/"
    done
leave

#preparation for pipeline releases
#TODO copy pipelines from ssp operator
mkdir -p "${REPO_DIR}/release/pipelines/windows-efi-installer/${RELEASE_VERSION:1}"
mkdir -p "${REPO_DIR}/release/pipelines/windows-bios-installer/${RELEASE_VERSION:1}"
mkdir -p "${REPO_DIR}/release/pipelines/windows-customize/${RELEASE_VERSION:1}"
