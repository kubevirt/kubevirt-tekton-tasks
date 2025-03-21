#!/usr/bin/env bash

set -e

if [ -z "${RELEASE_VERSION}" ]; then
  echo "RELEASE_VERSION is not defined"
  exit 1
fi

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

source "${SCRIPT_DIR}/release-var.sh"
source "${SCRIPT_DIR}/common.sh"

# add qemu-user-static
sudo podman run --rm --privileged docker.io/multiarch/qemu-user-static --reset -p yes

IMAGE_NAME_AND_TAG="tekton-tasks:${RELEASE_VERSION}"
IMAGE="${REGISTRY}/${REPOSITORY}/${IMAGE_NAME_AND_TAG}"
# Remove any existing manifest and image
podman manifest rm "${IMAGE}" || true
podman image exists "$IMAGE" && podman rmi "${IMAGE}" || true
podman build --platform=linux/amd64,linux/s390x,linux/arm64 --manifest "${IMAGE}" -f "build/Containerfile" .

IMAGE_NAME_AND_TAG="tekton-tasks-disk-virt:${RELEASE_VERSION}"
IMAGE="${REGISTRY}/${REPOSITORY}/${IMAGE_NAME_AND_TAG}"
# Remove any existing manifest and image
podman manifest rm "${IMAGE}" || true
podman image exists "$IMAGE" && podman rmi "${IMAGE}" || true
podman build --platform=linux/amd64,linux/arm64 --manifest "${IMAGE}" -f "build/Containerfile.DiskVirt" .
