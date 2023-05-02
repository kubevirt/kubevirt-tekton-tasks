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


IMAGE_NAME_AND_TAG="tekton-tasks:${RELEASE_VERSION}"
IMAGE="${REGISTRY}/${REPOSITORY}/${IMAGE_NAME_AND_TAG}"
podman build -f "build/Containerfile" -t "${IMAGE}" .

IMAGE_NAME_AND_TAG="tekton-tasks-disk-virt:${RELEASE_VERSION}"
IMAGE="${REGISTRY}/${REPOSITORY}/${IMAGE_NAME_AND_TAG}"
podman build -f "build/Containerfile.DiskVirt" -t "${IMAGE}" .
