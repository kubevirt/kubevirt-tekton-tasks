# Secure Execution Installer Pipeline

This pipeline automates the creation of IBM Secure Execution (SE) enabled virtual machines on s390x architecture. It performs automated Linux installation, configures LUKS encryption for root and boot volumes, generates SE images, and optionally creates golden images for reuse.

## Overview

The Secure Execution Installer Pipeline provides a complete, automated workflow for creating production-ready SE VMs:

1. **Automated Installation** - Downloads ISO and performs unattended Linux installation
2. **LUKS Encryption** - Encrypts root and boot volumes with auto-generated keys
3. **SE Image Generation** - Creates secure execution images with genprotimg
4. **Security Hardening** - Disables serial console, masks emergency services
5. **Golden Image Creation** - Optionally creates reusable DataSource images
6. **GUI Integration** - Full support for Kubernetes dashboard workflow

## Prerequisites

### Required Components

- **KubeVirt** >= v1.9.0
- **Tekton Pipelines** >= v0.43.0
- **Kubernetes** cluster with s390x node support
- **Storage Class** with dynamic provisioning
- **IBM Z or LinuxONE** hardware with Secure Execution support

### Deploy Required Resources

Before running the pipeline, deploy the RBAC and ConfigMap resources:

#### 1. Deploy RBAC Resources

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/release/pipelines/secure-execution-installer/secure-execution-installer-rbac.yaml
```

This creates the ServiceAccount, Role, and RoleBinding required for the pipeline.

#### 2. Deploy SE Templates ConfigMap

The pipeline automatically deploys the SE templates ConfigMap from the URL specified in the `seTemplatesConfigMapURL` parameter. However, you can also deploy it manually:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/release/pipelines/secure-execution-installer/configmaps/secure-execution-installer-configmaps.yaml
```

This ConfigMap contains the cloud-init scripts for SE automation.

### Obtain SE Certificates

You need real IBM Z Secure Execution certificates:

- **Host key document** (`hostDoc`) - Base64 encoded
- **IBM signing key** (`ibmSign`) - Base64 encoded
- **CA certificate** (`caCert`) - Base64 encoded

See [IBM Secure Execution Documentation](https://www.ibm.com/docs/en/linux-on-systems?topic=virtualization-secure-execution) for details on obtaining these certificates from your IBM Z system.

> **Note**: Test certificates will NOT work for actual Secure Execution. You must use real certificates from your IBM Z hardware.

## Architecture Support

- **Primary**: linux/s390x (IBM Z, LinuxONE)
- **Secure Execution**: Requires IBM Z15 or newer with SE enabled

## Supported Operating Systems

- Fedora Server (s390x) - Tested with Fedora 44
- Red Hat Enterprise Linux (s390x) - RHEL 8.x, 9.x
- Other s390x Linux distributions with kickstart support

## Pipeline Parameters

| Parameter | Type | Description | Default | Required |
|-----------|------|-------------|---------|----------|
| `vmName` | string | Name of the VM to create | - | ✅ |
| `memory` | string | Memory allocation (e.g., 8Gi, 16Gi) | 8Gi | |
| `diskSize` | string | Root disk size (e.g., 20Gi, 50Gi) | 20Gi | |
| `isoUrl` | string | URL to Linux ISO for s390x | - | ✅ |
| `bootImage` | string | Container image with kernel/initrd | - | ✅ |
| `hostDoc` | string | Base64 encoded host key document cert | - | ✅ |
| `ibmSign` | string | Base64 encoded IBM signing cert | - | ✅ |
| `caCert` | string | Base64 encoded CA cert | - | ✅ |
| `sshKey` | string | SSH public key for VM access | - | ✅ |
| `namespace` | string | Namespace for resources | "" (current) | |
| `storageClass` | string | Storage class for DataVolumes | hostpath-sc | |
| `createGoldenImage` | string | Create DataSource golden image | false | |
| `goldenImageName` | string | Name for golden image DataSource | "" | |
| `osPreference` | string | OS preference label for golden image (fedora or rhel) | fedora | |
| `machineType` | string | VM machine type for s390x architecture | s390-ccw-virtio | |
| `seTemplatesConfigMapURL` | string | URL to SE templates ConfigMap YAML | GitHub URL | |
| `cliImage` | string | kubectl/oc CLI image | quay.io/kata-containers/kubectl:v1.37.1 | |
| `kickstartServerImage` | string | HTTP server image for kickstart file | quay.io/nginx/nginx-unprivileged:1.31.6 | |

### Image Parameter Notes

- **`cliImage`**: Use public image for vanilla Kubernetes, Red Hat image for OpenShift
- **`kickstartServerImage`**: Supports both apache and nginx images
- **Kubernetes**: Use `quay.io/nginx/nginx-unprivileged:1.31.6` (default, no auth required)
- **OpenShift**: Use `registry.redhat.io/rhel8/httpd-24:latest` (requires Red Hat registry auth)
- **Alternatives**: `docker.io/nginx:alpine`, `docker.io/bitnami/apache:latest`

## Pipeline Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                    Secure Execution Pipeline                     │
└─────────────────────────────────────────────────────────────────┘

0. deploy-se-templates
   └─> Fetches and deploys SE templates ConfigMap from GitHub (automatic)

1. create-kickstart-config
   └─> Creates ConfigMap with kickstart configuration

2. deploy-kickstart-server
   └─> Deploys HTTP server to serve kickstart file

3. create-se-secret
   └─> Creates cloud-init secret with SE automation scripts (reads from deployed ConfigMap)

4. create-iso-datavolume ──┐
                           ├─> Create DataVolumes in parallel
5. create-root-datavolume ─┘

6. get-kickstart-url
   └─> Returns url to the kickstart server

7. create-vm
   └─> Creates VM with kernel boot and kickstart

8. wait-for-installation
   └─> Waits for automated installation to complete

9. finalize-vm
   └─> Removes ISO, updates VM config, cleans up

10. create-golden-image (optional)
    └─> Stops VM, creates DataSource, deletes VM to free DataVolume

Finally: cleanup-on-failure
    └─> Cleans up resources if pipeline fails
```

## Quick Start

### 1. Deploy RBAC

First, deploy the required RBAC resources:

```bash
kubectl apply -f release/pipelines/secure-execution-installer/secure-execution-installer-rbac.yaml
```

### 2. Prepare Certificates

Encode your certificates to base64:

```bash
# Host key document
HOST_DOC=$(base64 -w0 < /path/to/host-key-doc.crt)

# IBM signing certificate
IBM_SIGN=$(base64 -w0 < /path/to/ibm-sign.crt)

# IBM ca certificate
CA_CERT=$(base64 -w0 < /path/to/ca-cert.crt)

# SSH public key
SSH_KEY=$(cat ~/.ssh/id_rsa.pub)
```

### 3. Run the Pipeline

#### Option A: Using OpenShift web console (Recommended)

1. Navigate to **Pipelines** in the OpenShift web console
2. Click **Create** → **Pipeline**
3. Search for "Secure Execution Installer"
4. Fill in the form with your parameters
5. Upload certificates
6. Click **Start**
7. Monitor progress in **Pipeline Runs** tab

#### Option B: Using CLI

```bash
kubectl create -f - <<EOF
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  generateName: fedora-se-installer-
spec:
  pipelineRef:
    name: secure-execution-installer
  params:
    - name: vmName
      value: "my-fedora-se-vm"
    - name: memory
      value: "8Gi"
    - name: diskSize
      value: "20Gi"
    - name: isoUrl
      value: "https://download.fedoraproject.org/pub/fedora-secondary/releases/44/Server/s390x/iso/Fedora-Server-dvd-s390x-44-1.7.iso"
    - name: bootImage
      value: "quay.io/<repo>/fedora-44:latest"
    - name: hostDoc
      value: "${HOST_DOC}"
    - name: ibmSign
      value: "${IBM_SIGN}"
    - name: caCert
      value: "${CA_CERT}"
    - name: sshKey
      value: "${SSH_KEY}"
    - name: seTemplatesConfigMapURL
      value: "https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/release/pipelines/secure-execution-installer/configmaps/secure-execution-installer-configmaps.yaml"
  timeouts:
    pipeline: 1h0m0s
EOF
```

## Creating a Golden Image

To create a reusable golden image:

```yaml
params:
  - name: createGoldenImage
    value: "true"
  - name: goldenImageName
    value: "fedora-44-se-base"
```

The golden image can then be used to quickly create new VMs:

```yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: new-vm-from-golden
spec:
  dataVolumeTemplates:
    - metadata:
        name: new-vm-disk
      spec:
        sourceRef:
          kind: DataSource
          name: fedora-44-se-base
        storage:
          resources:
            requests:
              storage: 20Gi
```

## Boot Image Requirements

The boot image must contain:
- Kernel image (`/kernel.img`)
- Initrd image (`/initrd.img`)
- Compatible with s390x architecture

Example Dockerfile:

```dockerfile
FROM scratch
COPY kernel.img /kernel.img
COPY initrd.img /initrd.img
```

Build and push:

```bash
podman build -t quay.io/your-org/fedora-44-boot:latest .
podman push quay.io/your-org/fedora-44-boot:latest
```

## Monitoring Pipeline Execution

### Via GUI

1. Go to **Pipelines** → **Pipeline Runs**
2. Click on your pipeline run
3. View task progress and logs in real-time
4. Check results and outputs

### Via CLI

```bash
# List pipeline runs
kubectl get pipelineruns

# Watch pipeline run
tkn pipelinerun logs <pipelinerun-name> -f

# Describe pipeline run
kubectl describe pipelinerun <pipelinerun-name>
```

## Examples

See the `pipelineruns/` directory for complete examples:
- Fedora 44 SE VM
- RHEL 9 SE VM with golden image
- Custom storage class configuration

## Contributing

To contribute improvements to this pipeline:

1. Fork kubevirt-tekton-tasks repository
2. Create feature branch
3. Make changes and test
4. Submit pull request

## Additional Resources
- [Kernel Boot Container Guide](../../utilities/kernel-boot-container/README.md)
- [KubeVirt Documentation](https://kubevirt.io/)
- [IBM Secure Execution Documentation](https://www.ibm.com/docs/en/linux-on-systems?topic=virtualization-secure-execution)

## License

Apache License 2.0
