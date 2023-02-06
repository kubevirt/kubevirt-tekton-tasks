#!/usr/bin/env bash

set -e

if [ -z "${RELEASE_VERSION}" ]; then
  echo "RELEASE_VERSION is not defined"
  exit 1
fi

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/release-var.sh"
source "${SCRIPT_DIR}/common.sh"


visit "${REPO_DIR}"
  visit modules
    for TASK_NAME in *; do
      if echo "${TASK_NAME}" | grep -vqE "^(${EXCLUDED_NON_IMAGE_MODULES})$"; then
        if [ ! -d  "${TASK_NAME}" ]; then
          continue
        fi
        visit "${TASK_NAME}"
          IMAGE_NAME_AND_TAG="tekton-task-${TASK_NAME}:${RELEASE_VERSION}"
          IMAGE="${REGISTRY}/${REPOSITORY}/${IMAGE_NAME_AND_TAG}"
          podman build -f "build/${TASK_NAME}/Dockerfile" --platform=linux/amd64,linux/arm64 --manifest "${IMAGE}" .
        leave
      fi
    done
  leave
leave
