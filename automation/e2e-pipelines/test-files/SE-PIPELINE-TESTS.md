# Secure Execution Pipeline CI Tests

This directory contains CI test files for the Secure Execution (SE) pipeline.

## Test Files

### 1. `fedora-se-installer-pipelinerun.yaml`
PipelineRun definition for testing the SE installer pipeline with Fedora.

**Key Features:**
- Configures for fast CI execution
- Disables golden image creation
- Uses local HTTP server for ISO

### 2. `se-pipeline-test.sh`
Main test script that orchestrates the SE pipeline testing.

**Features:**
- Full integration test on s390x
- Automatic resource deployment
- Pipeline execution and monitoring

## Running Tests

### Prerequisites

```bash
# Required tools
- oc or kubectl CLI
- jq (for JSON parsing)
- Access to Kubernetes/OpenShift cluster
- Tekton Pipelines installed
- IBM Z or LinuxONE hardware (s390x architecture)
- Real IBM Z SE certificates (for actual SE testing)
```

### Preparation

Create a configmap containing the IBM Secure Execution Keys for the test machine.

Rename or copy `HKD-XXXXXXXXXXXX.crt` to `secex-hostkey.crt` and create the configmap:
```bash
kubectl create namespace kubevirt-prow-jobs
kubectl create -n kubevirt-prow-jobs configmap secex-hostkey --from-file=secex-hostkey.crt
```

### Integration Tests (s390x Only)

**IMPORTANT:** This test requires s390x architecture. It will fail on other architectures (x86_64, arm64, etc.).

Run integration tests on s390x:

```bash
# Set environment variables
export KUBECONFIG=/path/to/kubeconfig
export DEV_MODE=true  # Optional: use local images

# Run tests
./automation/e2e-pipelines/se-pipeline-test.sh
```

**What it does:**
- Deploys all required resources
- Creates ISO DataVolume
- Starts HTTP server
- Deploys SE pipeline
- Runs complete pipeline execution
- Waits up to 2 hours for completion

## Test Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    SE Pipeline CI Test                       │
└─────────────────────────────────────────────────────────────┘
                            ↓
                ┌────────────────────────┐
                │ Full Integration Tests │
                ├────────────────────────┤
                │ • Deploy resources     │
                │ • Start HTTP server    │
                │ • Deploy pipeline      │
                │ • Run pipeline         │
                │ • Wait for completion  │
                │ • Verify results       │
                └────────────────────────┘
```

## Customizing Tests

### Change Test VM Configuration

Edit `fedora-se-installer-pipelinerun.yaml`:

```yaml
params:
  - name: memory
    value: "16Gi"  # Increase memory
  - name: diskSize
    value: "50Gi"  # Increase disk
  - name: storageClass
    value: "your-storage-class"  # Use different storage
```

### Use Different ISO

Edit `fedora-se-dv.yaml`:

```yaml
spec:
  source:
    registry:
      url: docker://your-registry/your-iso:tag
```

Or use HTTP source:

```yaml
spec:
  source:
    http:
      url: "https://your-server/fedora-s390x.iso"
```

### Adjust Timeout

Edit `se-pipeline-test.sh`:

```bash
local timeout=3660  # Change from 1 hours to desired value
```

## Troubleshooting

### Test Fails on Non-s390x

**Expected behavior:** Test will fail to start s390x VMs

### Test Fails on s390x

**Common issues:**

1. **ISO Download Fails**
   - Check ISO URL is accessible
   - Verify storage class has enough space
   - Check CDI is working

2. **Pipeline Timeout**
   - SE installation takes 1-2 hours
   - Check VM is actually running
   - Verify kickstart server is accessible

3. **Certificate Errors**
   - Test certificates won't work for real SE
   - Use actual IBM Z certificates for real testing

4. **Storage Issues**
   - Ensure storage class supports dynamic provisioning
   - Check PVC creation succeeds
   - Verify enough storage available

### Debug Commands

```bash
# Check pipeline status
oc get pipelinerun -l pipelinerun=fedora-se-installer-run

# View pipeline logs
tkn pipelinerun logs -f -l pipelinerun=fedora-se-installer-run

# Check VM status
oc get vm test-se-vm

# Check DataVolume status
oc get dv

# Check HTTP server
oc get pods -l app=http-server
oc logs -l app=http-server
```

## Test Results

### Success Criteria

- ✅ All YAML files validate successfully
- ✅ Pipeline deploys without errors
- ✅ Pipeline structure is correct
- ✅ (s390x only) Pipeline completes successfully
- ✅ (s390x only) VM is created and running
- ✅ (s390x only) SE image is generated

When successfull, it will look like this:
```
=== Verifying Secure Execution Status ===
VMI phase: Running
Verifying Secure Execution status of VM
Warning: Permanently added 'vm.sec-exec-vm.test-kubevirt-tekton-tasks-secure-execution' (ED25519) to the list of known hosts.
Warning: Permanently added 'vm.sec-exec-vm.test-kubevirt-tekton-tasks-secure-execution' (ED25519) to the list of known hosts.
Secure Execution verified: /sys/firmware/uv/prot_virt_guest = 1

=== SE Pipeline Integration Test Completed Successfully ===
Total execution time: 22m 54s (1374 seconds)
```

## Contributing

When adding new tests:

1. Follow existing naming conventions
2. Add documentation to this file
3. Ensure tests work on both s390x and non-s390x
4. Use test certificates, not real ones
5. Keep test execution time reasonable

## References

- [SE Pipeline Documentation](../../../templates-pipelines/secure-execution-installer/README.md)
- [SE Pipeline Getting Started](../../../templates-pipelines/secure-execution-installer/GETTING-STARTED.md)
- [IBM Secure Execution Docs](https://www.ibm.com/docs/en/linux-on-systems?topic=virtualization-secure-execution)
