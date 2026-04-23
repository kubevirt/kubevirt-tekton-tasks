#!/bin/bash
# Build and push kernel boot container

set -euo pipefail

# Configuration
REGISTRY="${REGISTRY:-quay.io}"
USERNAME="${USERNAME:-your-username}"
IMAGE_NAME="${IMAGE_NAME:-fedora43-s390x-boot}"
TAG="${TAG:-latest}"

FULL_IMAGE="${REGISTRY}/${USERNAME}/${IMAGE_NAME}:${TAG}"

echo "Building kernel boot container..."
echo "Image: ${FULL_IMAGE}"
echo ""

# Check if kernel.img and initrd.img exist
if [ ! -f "kernel.img" ] || [ ! -f "initrd.img" ]; then
    echo "Error: kernel.img and initrd.img must be in the current directory"
    echo ""
    echo "Extract them first using:"
    echo "  ./extract-from-iso.sh <iso-file>"
    echo "  cp output/kernel.img output/initrd.img ."
    exit 1
fi

# Display file sizes
echo "Files to include:"
echo "  kernel.img: $(du -h kernel.img | cut -f1)"
echo "  initrd.img: $(du -h initrd.img | cut -f1)"
echo ""

# Build container
echo "Building container..."
podman build -f Containerfile.manual -t "${FULL_IMAGE}" .

if [ $? -eq 0 ]; then
    echo ""
    echo "Build successful!"
    echo ""
    echo "To push to registry:"
    echo "  podman login ${REGISTRY}"
    echo "  podman push ${FULL_IMAGE}"
    echo ""
    echo "To use in VM:"
    echo "  spec.template.spec.domain.firmware.kernelBoot.container.image: ${FULL_IMAGE}"
else
    echo " Build failed"
    exit 1
fi
