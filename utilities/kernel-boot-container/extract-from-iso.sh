#!/bin/bash
# Script to extract kernel and initrd from ISO for s390x

set -euo pipefail

# Configuration
ISO_FILE="${1:-}"
OUTPUT_DIR="${2:-./output}"

if [ -z "$ISO_FILE" ]; then
    echo "Usage: $0 <iso-file> [output-dir]"
    echo ""
    echo "Example:"
    echo "  $0 Fedora-Server-dvd-s390x-43-1.1.iso ./extracted"
    exit 1
fi

if [ ! -f "$ISO_FILE" ]; then
    echo "Error: ISO file not found: $ISO_FILE"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Create mount point
MOUNT_POINT=$(mktemp -d)
trap "umount $MOUNT_POINT 2>/dev/null || true; rmdir $MOUNT_POINT" EXIT

echo "Mounting ISO: $ISO_FILE"
sudo mount -o loop,ro "$ISO_FILE" "$MOUNT_POINT"

# Common locations for kernel and initrd in s390x ISOs
KERNEL_PATHS=(
    "images/kernel.img"
    "boot/kernel.img"
    "images/pxeboot/kernel.img"
    "isolinux/kernel.img"
)

INITRD_PATHS=(
    "images/initrd.img"
    "boot/initrd.img"
    "images/pxeboot/initrd.img"
    "isolinux/initrd.img"
)

# Find and copy kernel
KERNEL_FOUND=false
for path in "${KERNEL_PATHS[@]}"; do
    if [ -f "$MOUNT_POINT/$path" ]; then
        echo "Found kernel at: $path"
        cp "$MOUNT_POINT/$path" "$OUTPUT_DIR/kernel.img"
        KERNEL_FOUND=true
        break
    fi
done

if [ "$KERNEL_FOUND" = false ]; then
    echo "Error: Could not find kernel in ISO"
    echo "Searched paths:"
    printf '  %s\n' "${KERNEL_PATHS[@]}"
    exit 1
fi

# Find and copy initrd
INITRD_FOUND=false
for path in "${INITRD_PATHS[@]}"; do
    if [ -f "$MOUNT_POINT/$path" ]; then
        echo "Found initrd at: $path"
        cp "$MOUNT_POINT/$path" "$OUTPUT_DIR/initrd.img"
        INITRD_FOUND=true
        break
    fi
done

if [ "$INITRD_FOUND" = false ]; then
    echo "Error: Could not find initrd in ISO"
    echo "Searched paths:"
    printf '  %s\n' "${INITRD_PATHS[@]}"
    exit 1
fi

# Unmount
sudo umount "$MOUNT_POINT"

# Display results
echo ""
echo "✅ Extraction complete!"
echo "Kernel: $OUTPUT_DIR/kernel.img ($(du -h "$OUTPUT_DIR/kernel.img" | cut -f1))"
echo "Initrd: $OUTPUT_DIR/initrd.img ($(du -h "$OUTPUT_DIR/initrd.img" | cut -f1))"
echo ""
echo "Next steps:"
echo "1. Build container: podman build -f Containerfile.manual -t quay.io/your-username/fedora43-s390x-boot:latest ."
echo "2. Push to registry: podman push quay.io/your-username/fedora43-s390x-boot:latest"
