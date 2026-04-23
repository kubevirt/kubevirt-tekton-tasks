# Getting Started with Secure Execution Pipeline

This guide walks you through deploying and using the Secure Execution (SE) pipeline on a fresh OpenShift cluster with KubeVirt.

## Prerequisites

### 1. Cluster Requirements
- OpenShift 4.14+ or Kubernetes 1.28+
- s390x architecture (IBM Z or LinuxONE)
- OpenShift Virtualization 4.21.4+ (or KubeVirt 1.2.0+)
- OpenShift Pipelines (Tekton) installed
- Storage class with ReadWriteOnce support

### 2. Verify Prerequisites

```bash
# Check cluster architecture
oc get nodes -o jsonpath='{.items[0].status.nodeInfo.architecture}'
# Should output: s390x

# Check OpenShift Virtualization
oc get csv -n openshift-cnv | grep kubevirt-hyperconverged
# Should show version 4.21.4 or higher

# Check Tekton Pipelines
oc get csv -n openshift-pipelines | grep pipelines
# Should show OpenShift Pipelines operator

# Check storage classes
oc get storageclass
# Note the name of your RWO storage class
```

### 3. IBM Secure Execution Certificates

You need three certificates for SE:
- **Host Key Document** (`host_key.crt`) - From your IBM Z system
- **IBM Signing Key** (`ibm_sign.crt`) - From IBM
- **CA Certificate** (`ca.crt`) - Certificate Authority certificate

```bash
# Encode certificates to base64
HOST_DOC=$(cat host_key.crt | base64 -w0)
IBM_SIGN=$(cat ibm_sign.crt | base64 -w0)
CA_CERT=$(cat ca.crt | base64 -w0)

# Save for later use
echo "HOST_DOC=$HOST_DOC" > se-certs.env
echo "IBM_SIGN=$IBM_SIGN" >> se-certs.env
echo "CA_CERT=$CA_CERT" >> se-certs.env
```

---

## Quick Start (5 Minutes)

### Step 1: Clone Repository

```bash
git clone https://github.com/kubevirt/kubevirt-tekton-tasks.git
cd kubevirt-tekton-tasks/templates-pipelines/secure-execution-installer
```

### Step 2: Deploy Pipeline and RBAC

```bash
# Deploy RBAC (ServiceAccount, Role, RoleBinding)
oc apply -f manifests/se-pipeline-rbac.yaml

# Deploy the pipeline
oc apply -f manifests/secure-execution-installer.yaml

# Deploy SE templates ConfigMap
oc apply -f configmaps/se-templates-configmaps.yaml

# Verify deployment
oc get pipeline secure-execution-installer
oc get sa se-pipeline-sa
```

### Step 3: Create Kernel Boot Container

The pipeline requires a kernel boot container with kernel and initrd for s390x. You need to build your own:

```bash
# Navigate to kernel boot container utilities
cd ../../utilities/kernel-boot-container

# Extract kernel and initrd from your ISO
./extract-from-iso.sh Fedora-Server-dvd-s390x-43-1.1.iso

# Build and push the container image
./build-and-push.sh

# Set the boot image variable
BOOT_IMAGE="quay.io/your-username/fedora43-s390x-boot:latest"
```

**Note**: There is no pre-built kernel boot container available. You must build your own from your Linux ISO. See `utilities/kernel-boot-container/README.md` for detailed instructions.

### Step 4: Run the Pipeline

```bash
# Load your SE certificates
source se-certs.env

# Create PipelineRun
cat <<EOF | oc apply -f -
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: my-first-se-vm
  namespace: default
spec:
  serviceAccountName: se-pipeline-sa
  pipelineRef:
    name: secure-execution-installer
  params:
    - name: vmName
      value: "my-se-vm"
    - name: memory
      value: "8Gi"
    - name: diskSize
      value: "30Gi"
    - name: isoUrl
      value: "https://download.fedoraproject.org/pub/fedora-secondary/releases/43/Server/s390x/iso/Fedora-Server-dvd-s390x-43-1.1.iso"
    - name: bootImage
      value: "${BOOT_IMAGE}"
    - name: hostDoc
      value: "${HOST_DOC}"
    - name: ibmSign
      value: "${IBM_SIGN}"
    - name: caCert
      value: "${CA_CERT}"
    - name: sshKey
      value: "$(cat ~/.ssh/id_rsa.pub)"
    - name: namespace
      value: "default"
    - name: storageClass
      value: "ocs-storagecluster-ceph-rbd"
    - name: createGoldenImage
      value: "false"
  workspaces:
    - name: se-templates
      configMap:
        name: se-templates
EOF

# Monitor progress
tkn pipelinerun logs my-first-se-vm -f
```

### Step 5: Access Your SE VM

```bash
# Wait for VM to be ready
oc wait --for=condition=Ready vm/my-se-vm --timeout=60m

# Start the VM
virtctl start my-se-vm

# Connect via console
virtctl console my-se-vm

# Or SSH (if configured)
virtctl ssh my-se-vm
```

---

## Detailed Deployment Guide

### 1. Install OpenShift Pipelines (if not installed)

```bash
# Create subscription
cat <<EOF | oc apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: openshift-pipelines-operator
  namespace: openshift-operators
spec:
  channel: latest
  name: openshift-pipelines-operator-rh
  source: redhat-operators
  sourceNamespace: openshift-marketplace
EOF

# Wait for installation
oc wait --for=condition=Ready csv -l operators.coreos.com/openshift-pipelines-operator-rh.openshift-operators --timeout=5m
```

### 2. Install OpenShift Virtualization (if not installed)

```bash
# Create namespace
oc create namespace openshift-cnv

# Create OperatorGroup
cat <<EOF | oc apply -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: kubevirt-hyperconverged-group
  namespace: openshift-cnv
spec:
  targetNamespaces:
    - openshift-cnv
EOF

# Create Subscription
cat <<EOF | oc apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: hco-operatorhub
  namespace: openshift-cnv
spec:
  source: redhat-operators
  sourceNamespace: openshift-marketplace
  name: kubevirt-hyperconverged
  channel: "stable"
  installPlanApproval: Automatic
EOF

# Wait for operator
oc wait --for=condition=Ready csv -n openshift-cnv -l operators.coreos.com/kubevirt-hyperconverged.openshift-cnv --timeout=10m

# Create HyperConverged instance
cat <<EOF | oc apply -f -
apiVersion: hco.kubevirt.io/v1beta1
kind: HyperConverged
metadata:
  name: kubevirt-hyperconverged
  namespace: openshift-cnv
spec:
  featureGates:
    enableCommonBootImageImport: true
EOF

# Wait for deployment
oc wait --for=condition=Available hco/kubevirt-hyperconverged -n openshift-cnv --timeout=20m
```

### 3. Deploy Secure Execution Pipeline

```bash
# Navigate to pipeline directory
cd kubevirt-tekton-tasks/templates-pipelines/secure-execution-installer

# Deploy all components
oc apply -f manifests/se-pipeline-rbac.yaml
oc apply -f manifests/secure-execution-installer.yaml
oc apply -f configmaps/se-templates-configmaps.yaml

# Verify
oc get pipeline,sa,role,rolebinding | grep se-pipeline
```

### 4. Prepare Kernel Boot Container

See `../../utilities/kernel-boot-container/README.md` for detailed instructions.

Quick version:
```bash
cd ../../utilities/kernel-boot-container

# Download ISO
wget https://download.fedoraproject.org/pub/fedora-secondary/releases/43/Server/s390x/iso/Fedora-Server-dvd-s390x-43-1.1.iso

# Extract kernel and initrd
./extract-from-iso.sh Fedora-Server-dvd-s390x-43-1.1.iso
cp output/* .

# Build container
export REGISTRY=quay.io
export USERNAME=your-username
./build-and-push.sh

# Push to registry
podman login quay.io
podman push quay.io/your-username/fedora43-s390x-boot:latest
```

---

## Usage Examples

### Example 1: Basic SE VM

```bash
oc apply -f pipelineruns/pipelineruns.yaml
```

### Example 2: SE VM with Golden Image

```bash
cat <<EOF | oc apply -f -
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: se-vm-with-golden
spec:
  serviceAccountName: se-pipeline-sa
  pipelineRef:
    name: secure-execution-installer
  params:
    - name: vmName
      value: "golden-se-vm"
    - name: createGoldenImage
      value: "true"  # Enable golden image creation
    # ... other params
EOF
```

### Example 3: Create VM from Golden Image

```bash
# First, ensure golden image VM is stopped
oc patch vm golden-se-vm --type merge -p '{"spec":{"runStrategy":"Halted"}}'

# Create DataVolume from golden image
cat <<EOF | oc apply -f -
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: new-se-vm-disk
  annotations:
    cdi.kubevirt.io/storage.usePopulator: "false"
spec:
  source:
    pvc:
      namespace: default
      name: golden-se-vm-dv
  storage:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 30Gi
EOF

# Create VM
cat <<EOF | oc apply -f -
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: new-se-vm
spec:
  running: true
  template:
    spec:
      architecture: s390x
      domain:
        launchSecurity: {}
        machine:
          type: s390-ccw-virtio-rhel9.6.0
        resources:
          requests:
            memory: 8Gi
        devices:
          disks:
            - name: rootdisk
              disk:
                bus: virtio
      volumes:
        - name: rootdisk
          dataVolume:
            name: new-se-vm-disk
EOF
```

---

## Using the OpenShift Console

### 1. Navigate to Pipelines

1. Open OpenShift Console
2. Navigate to **Pipelines** → **Pipelines**
3. Find **secure-execution-installer**
4. Click **Actions** → **Start**

### 2. Fill Parameters

- **vmName**: Name for your VM (e.g., `my-se-vm`)
- **memory**: Memory allocation (e.g., `8Gi`)
- **diskSize**: Disk size (e.g., `30Gi`)
- **isoUrl**: ISO download URL
- **bootImage**: Your kernel boot container image
- **hostDoc**: Base64 encoded host key document
- **ibmSign**: Base64 encoded IBM signing key
- **sshKey**: Your SSH public key
- **storageClass**: Your storage class name
- **createGoldenImage**: `true` or `false`

### 3. Monitor Execution

1. Click on the PipelineRun name
2. View task progress in the graph
3. Click on tasks to see logs
4. Wait for completion (30-60 minutes)

### 4. Access VM

1. Navigate to **Virtualization** → **VirtualMachines**
2. Find your VM
3. Click **Actions** → **Start**
4. Use **Console** tab to access

---

## Troubleshooting

### Pipeline Fails at DataVolume Creation

**Issue**: DataVolumes stuck in Pending

**Solution**: Disable volume populator
```bash
# Already included in pipeline, but verify:
oc get dv <dv-name> -o yaml | grep usePopulator
# Should show: cdi.kubevirt.io/storage.usePopulator: "false"
```

### VM Fails to Boot

**Issue**: Kernel boot fails

**Solution**: Verify kernel boot container
```bash
# Test container
podman run --rm quay.io/your-username/fedora43-s390x-boot:latest ls -lh /
# Should show kernel.img and initrd.img
```

### Cannot Clone from Golden Image

**Issue**: Golden image PVC in use

**Solution**: This is automatically handled by the pipeline
- When `createGoldenImage: "true"` is set, the pipeline automatically:
  1. Stops the VM
  2. Creates the golden image DataSource
  3. Deletes the VM to free the DataVolume
- The original VM is removed so the golden image can be used immediately
- No manual intervention required

**Note**: If you created a golden image manually (outside the pipeline), you need to stop/delete the VM:
```bash
oc patch vm golden-se-vm --type merge -p '{"spec":{"runStrategy":"Halted"}}'
oc wait --for=delete vmi/golden-se-vm --timeout=5m
oc delete vm golden-se-vm
```

---

## Next Steps

1. **Customize Kickstart**: Edit `configmaps/se-templates-configmaps.yaml`
2. **Create Golden Images**: Set `createGoldenImage: "true"`
3. **Automate Deployments**: Use PipelineRuns in CI/CD
4. **Scale**: Create multiple SE VMs from golden images

---

## Additional Resources

- **Pipeline README**: [README.md](README.md)
- **Kernel Boot Containers**: [../../utilities/kernel-boot-container/README.md](../../utilities/kernel-boot-container/README.md)
- **KubeVirt Documentation**: https://kubevirt.io/
- **IBM Secure Execution**: https://www.ibm.com/docs/en/linux-on-systems?topic=virtualization-secure-execution
- **Tekton Documentation**: https://tekton.dev/docs/

---

## Support

For issues and questions:
- **GitHub Issues**: https://github.com/kubevirt/kubevirt-tekton-tasks/issues
- **KubeVirt Slack**: #kubevirt on Kubernetes Slack
- **Mailing List**: kubevirt-dev@googlegroups.com
