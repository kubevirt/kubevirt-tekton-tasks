# Secure Execution Installer Pipeline

This pipeline automates the creation of IBM Secure Execution (SE) enabled virtual machines on s390x architecture. It performs automated Linux installation, configures LUKS encryption for root and boot volumes, generates SE images, and optionally creates golden images for reuse.

## Overview

The Secure Execution Installer Pipeline provides a complete, automated workflow for creating production-ready SE VMs:

1. **Automated Installation** - Downloads ISO and performs unattended Linux installation
2. **LUKS Encryption** - Encrypts root and boot volumes with auto-generated keys
3. **SE Image Generation** - Creates secure execution images with genprotimg
4. **Security Hardening** - Disables serial console, masks emergency services
5. **Golden Image Creation** - Optionally creates reusable DataSource images
6. **GUI Integration** - Full support for OpenShift Console GUI workflow

## Prerequisites

- **KubeVirt** >= v1.0.0
- **Tekton Pipelines** >= v0.43.0
- **OpenShift Virtualization** (for s390x support)
- **Storage Class** with dynamic provisioning
- **IBM Z or LinuxONE** hardware with Secure Execution support

## Architecture Support

- **Primary**: linux/s390x (IBM Z, LinuxONE)
- **Secure Execution**: Requires IBM Z15 or newer with SE enabled

## Supported Operating Systems

- Fedora Server (s390x) - Tested with Fedora 43
- Red Hat Enterprise Linux (s390x) - RHEL 8.x, 9.x
- Other s390x Linux distributions with kickstart support

## Pipeline Parameters

| Parameter | Type | Description | Default | Required |
|-----------|------|-------------|---------|----------|
| `vmName` | string | Name of the VM to create | - |  
| `memory` | string | Memory allocation (e.g., 8Gi, 16Gi) | 8Gi | 
| `diskSize` | string | Root disk size (e.g., 20Gi, 50Gi) | 20Gi |
| `isoUrl` | string | URL to Linux ISO for s390x | - |
| `bootImage` | string | Container image with kernel/initrd | - |
| `hostDoc` | string | Base64 encoded host key document cert | - |
| `ibmSign` | string | Base64 encoded IBM signing cert | - |
| `sshKey` | string | SSH public key for VM access | - |
| `namespace` | string | Namespace for resources | "" (current) |
| `storageClass` | string | Storage class for DataVolumes | hostpath-sc |
| `createGoldenImage` | string | Create DataSource golden image | false |
| `goldenImageName` | string | Name for golden image DataSource | "" |

## Pipeline Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                    Secure Execution Pipeline                     │
└─────────────────────────────────────────────────────────────────┘

1. create-kickstart-config
   └─> Creates ConfigMap with kickstart configuration
   
2. deploy-kickstart-server
   └─> Deploys HTTP server to serve kickstart file
   
3. create-se-secret
   └─> Creates cloud-init secret with SE automation scripts
   
4. create-iso-datavolume ──┐
                           ├─> Create DataVolumes in parallel
5. create-root-datavolume ─┘
   
6. wait-for-datavolumes
   └─> Waits for both DVs to be ready
   
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

### 1. Deploy the Pipeline

```bash
kubectl apply -f https://github.com/kubevirt/kubevirt-tekton-tasks/releases/download/v0.25.0/kubevirt-tekton-tasks.yaml
```

### 2. Create ConfigMap with SE Templates

```bash
kubectl apply -f configmaps/se-templates-configmaps.yaml
```

### 3. Prepare Certificates

Encode your certificates to base64:

```bash
# Host key document
HOST_DOC=$(base64 -w0 < /path/to/host-key-doc.crt)

# IBM signing certificate
IBM_SIGN=$(base64 -w0 < /path/to/ibm-sign.crt)

# SSH public key
SSH_KEY=$(cat ~/.ssh/id_rsa.pub)
```

### 4. Run the Pipeline

#### Option A: Using OpenShift Console GUI (Recommended)

1. Navigate to **Pipelines** in OpenShift Console
2. Click **Create** → **Pipeline**
3. Search for "Secure Execution Installer"
4. Fill in the form with your parameters
5. Upload certificates (GUI will base64 encode automatically)
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
      value: "https://download.fedoraproject.org/pub/fedora-secondary/releases/43/Server/s390x/iso/Fedora-Server-netinst-s390x-43-1.6.iso"
    - name: bootImage
      value: "quay.io/<repo>/fedora-43:latest"
    - name: hostDoc
      value: "${HOST_DOC}"
    - name: ibmSign
      value: "${IBM_SIGN}"
    - name: sshKey
      value: "${SSH_KEY}"
  workspaces:
    - name: se-templates
      configMap:
        name: se-templates-configmap
  timeout: 2h0m0s
EOF
```

## Creating a Golden Image

To create a reusable golden image:

```yaml
params:
  - name: createGoldenImage
    value: "true"
  - name: goldenImageName
    value: "fedora-43-se-base"
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
          name: fedora-43-se-base
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
podman build -t quay.io/your-org/fedora-43-boot:latest .
podman push quay.io/your-org/fedora-43-boot:latest
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

## Troubleshooting

### Pipeline Fails at create-iso-datavolume

**Issue**: ISO download timeout or storage issues

**Solution**:
- Verify ISO URL is accessible
- Check storage class has sufficient capacity
- Increase timeout in PipelineRun spec

### VM Installation Hangs

**Issue**: Kickstart server not reachable or kickstart config error

**Solution**:
- Check kickstart server pod logs: `kubectl logs deployment/kickstart-server-<vmname>`
- Verify Route is created: `kubectl get route kickstart-server-<vmname>`
- Check VM console for errors: `virtctl console <vmname>`

### SE Image Generation Fails

**Issue**: Invalid certificates or genprotimg errors

**Solution**:
- Verify certificates are valid and base64 encoded correctly
- Check VM logs: `kubectl logs virt-launcher-<vmname>-xxxxx`
- Ensure s390-tools package is installed in the VM

### Cleanup Not Working

**Issue**: Resources remain after pipeline completion

**Solution**:
- Manually delete resources: `kubectl delete vm,dv,deployment,service,route,configmap,secret -l vm-name=<vmname>`
- Check owner references are set correctly

## Security Considerations

1. **Certificate Management**: Store certificates securely, use Secrets for sensitive data
2. **SSH Keys**: Use dedicated keys for SE VMs, rotate regularly
3. **Network Isolation**: Consider network policies for kickstart server
4. **RBAC**: Limit pipeline execution to authorized users
5. **Audit Logging**: Enable audit logs for pipeline runs

## Performance Tuning

- **ISO Download**: Use local mirror or CDN for faster downloads
- **Storage**: Use high-performance storage class for better I/O
- **Memory**: Allocate sufficient memory for installation (minimum 4Gi)
- **Timeout**: Adjust based on ISO size and network speed

## Integration with CI/CD

The pipeline can be integrated into GitOps workflows:

```yaml
apiVersion: triggers.tekton.dev/v1beta1
kind: EventListener
metadata:
  name: se-vm-creator
spec:
  triggers:
    - name: create-se-vm
      bindings:
        - ref: se-vm-binding
      template:
        ref: se-vm-template
```

## Examples

See the `pipelineruns/` directory for complete examples:
- Fedora 43 SE VM
- RHEL 9 SE VM with golden image
- Custom storage class configuration

## Contributing

To contribute improvements to this pipeline:

1. Fork kubevirt-tekton-tasks repository
2. Create feature branch
3. Make changes and test
4. Submit pull request

## Support

- **Documentation**: https://kubevirt.io/
- **IBM Secure Execution**: https://www.ibm.com/docs/en/linux-on-systems?topic=virtualization-secure-execution
- **Issues**: https://github.com/kubevirt/kubevirt-tekton-tasks/issues

## License

Apache License 2.0
