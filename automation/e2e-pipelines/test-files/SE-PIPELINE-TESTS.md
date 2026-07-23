# Secure Execution Pipeline CI Tests

This directory contains CI test files for the Secure Execution (SE) pipeline.

## Test Files

### 1. `fedora-se-installer-pipelinerun.yaml`
PipelineRun definition for testing the SE installer pipeline with Fedora.

**Key Features:**
- Uses test certificates (not real IBM Z certificates)
- Configures for fast CI execution
- Disables golden image creation
- Uses local HTTP server for ISO

### 2. `fedora-se-dv.yaml`
DataVolume definition for creating a test ISO volume.

**Default ISO:** Fedora 44 Server netinst for s390x (downloaded from Fedora mirrors)

**Supported Operating Systems:**
- Fedora Server (s390x) - Fedora 44 or later
- Red Hat Enterprise Linux (s390x) - RHEL 8.x, 9.x

**Note:** The ISO download may take several minutes depending on network speed (~1GB file).

### 3. `se-pipeline-test.sh`
Main test script that orchestrates the SE pipeline testing.

**Features:**
- Architecture detection (requires s390x)
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
                    Architecture Check
                            ↓
        ┌───────────────────┴───────────────────┐
        ↓                                       ↓
   Non-s390x                                s390x
        ↓                                       ↓
┌───────────────────┐              ┌────────────────────────┐
│   Exit with Error │              │ Full Integration Tests │
│                   │              ├────────────────────────┤
│ SE requires s390x │              │ • Deploy resources     │
└───────────────────┘              │ • Create ISO DV        │
                                   │ • Start HTTP server    │
                                   │ • Deploy pipeline      │
                                   │ • Run pipeline         │
                                   │ • Wait for completion  │
                                   │ • Verify results       │
                                   └────────────────────────┘
```

## Test Certificates

The test uses **dummy certificates** for CI purposes. These are NOT real IBM Z certificates and will NOT work for actual Secure Execution.

**Test Certificate Format:**
```yaml
hostDoc: "LS0tLS1CRUdJTi..."  # Base64 encoded test cert
ibmSign: "LS0tLS1CRUdJTi..."  # Base64 encoded test cert
caCert: "LS0tLS1CRUdJTi..."   # Base64 encoded test cert
```

**For Real Testing:**
Replace these with actual IBM Z SE certificates obtained from your IBM Z system.

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
local timeout=7200  # Change from 2 hours to desired value
```

## Troubleshooting

### Test Fails on Non-s390x

**Expected behavior:** Test will exit with error message stating SE requires s390x architecture.

This is correct behavior - SE pipeline can only run on s390x systems.

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

## Integration with CI/CD

### GitHub Actions (Limited)

GitHub Actions doesn't support s390x runners, so only validation tests can run on x86_64 runners:

```yaml
name: SE Pipeline Validation
on: [pull_request]
jobs:
  validate:
    runs-on: ubuntu-latest  # x86_64 runner - validation only
    steps:
      - uses: actions/checkout@v4
      - name: Run validation tests
        run: ./automation/e2e-pipelines/se-pipeline-test.sh
```

**Note**: This only runs YAML validation. Full integration tests require s390x hardware with Fedora or RHEL.

### IBM Cloud CI (Recommended for Full Testing)

For full integration tests, use IBM Cloud CI with s390x nodes running Fedora or RHEL:

```yaml
# Jenkins/Tekton pipeline on IBM Cloud
stages:
  - name: SE Pipeline Test
    agent: s390x
    steps:
      - checkout
      - sh './automation/e2e-pipelines/se-pipeline-test.sh'
```

## Test Results

### Success Criteria

- ✅ All YAML files validate successfully
- ✅ Pipeline deploys without errors
- ✅ Pipeline structure is correct
- ✅ (s390x only) Pipeline completes successfully
- ✅ (s390x only) VM is created and running
- ✅ (s390x only) SE image is generated

### Expected Output

```
=== Secure Execution Pipeline CI Test ===
Running on architecture: x86_64
WARNING: SE pipeline requires s390x architecture
Current architecture is x86_64 - running validation tests only
Validating pipeline YAML files
All required YAML files exist
Deploying SE pipeline resources for validation
...
Pipeline structure validation passed
=== SE Pipeline Validation Tests Completed Successfully ===
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