# Kernel Boot Container Builder for s390x

This directory contains templates and scripts to create kernel boot containers for use with KubeVirt's direct kernel boot feature on s390x architecture.

## Overview

KubeVirt's `kernelBoot` feature allows VMs to boot directly from a kernel and initrd stored in a container image, bypassing traditional bootloaders. This is particularly useful for:

- **Automated installations** (Anaconda/Kickstart)
- **IBM Secure Execution** VMs
- **Custom boot configurations**
- **Network boot scenarios**

## What You Need

1. **Linux ISO file** for s390x (Fedora, RHEL, etc.)
2. **Podman or Docker** for building containers
3. **Registry access** (quay.io, Docker Hub, or private registry)

## Quick Start

### Method 1: Extract from ISO (Recommended)

```bash
# 1. Download your ISO
Example:
wget https://download.fedoraproject.org/pub/fedora-secondary/releases/43/Server/s390x/iso/Fedora-Server-dvd-s390x-43-1.1.iso

# 2. Extract kernel and initrd
chmod +x extract-from-iso.sh
./extract-from-iso.sh Fedora-Server-dvd-s390x-43-1.1.iso

# 3. Copy extracted files to current directory
cp output/kernel.img output/initrd.img .

# 4. Build and push container
chmod +x build-and-push.sh
export REGISTRY=quay.io
export USERNAME=your-username
export IMAGE_NAME=fedora43-s390x-boot
export TAG=latest
./build-and-push.sh

# 5. Login and push
podman login quay.io
podman push quay.io/your-username/fedora43-s390x-boot:latest
```


### Method 2: Manual Build (Pre-extracted Files)

If you already have `kernel.img` and `initrd.img`:

```bash
# 1. Place files in this directory
ls -lh kernel.img initrd.img

# 2. Build container
podman build -f Containerfile -t quay.io/your-username/custom-boot:latest .

# 3. Push to registry
podman push quay.io/your-username/custom-boot:latest
```

## File Structure

```
kernel-boot-container/
├── Containerfile              # Automated build from ISO
├ files
├── extract-from-iso.sh        # Extract kernel/initrd from ISO
├── build-and-push.sh          # Build and push helper script
├── README.md                  # This file
```

### Secure Execution VM with Kernel Boot

```yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: se-vm
spec:
  running: true
  template:
    spec:
      architecture: s390x
      domain:
        launchSecurity: {}  # Enable Secure Execution
        machine:
          type: s390-ccw-virtio-rhel9.6.0
        resources:
          requests:
            memory: 8Gi
        firmware:
          kernelBoot:
            container:
              image: quay.io/your-username/fedora43-s390x-boot:latest
              kernelPath: /kernel.img
              initrdPath: /initrd.img
            kernelArgs: >-
              console=ttysclp0
              inst.text
              inst.stage2=hd:/dev/vda
              inst.ks=http://kickstart-server/ks.cfg
        devices:
          disks:
            - name: rootdisk
              disk:
                bus: virtio
      volumes:
        - name: rootdisk
          dataVolume:
            name: se-root-dv
```

## Integration with Tekton Pipeline

Use in your Secure Execution pipeline:

```yaml
params:
  - name: bootImage
    type: string
    description: "Kernel boot container image"
    default: "quay.io/your-username/fedora43-s390x-boot:latest"
```

