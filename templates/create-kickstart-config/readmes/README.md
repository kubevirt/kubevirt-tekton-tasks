# create-kickstart-config task

This task creates a ConfigMap containing kickstart configuration for automated Linux installation on s390x architecture. It's specifically designed for Secure Execution VM provisioning workflows.

## Parameters

| Name | Description | Default |
|------|-------------|---------|
| vmName | Name of the VM for which to create the kickstart configuration | - |
| namespace | Namespace where to create the ConfigMap (defaults to active namespace) | "" |
| kickstartTemplate | Kickstart configuration template content | "" (uses default template) |
| setOwnerReference | Set owner reference to the new object created by the task run pod | "false" |

## Results

| Name | Description |
|------|-------------|
| name | The name of the created ConfigMap |
| namespace | The namespace of the created ConfigMap |

## Usage

### Basic Usage

Create a kickstart ConfigMap with default template:

```yaml
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  generateName: create-kickstart-config-
spec:
  taskRef:
    kind: Task
    name: create-kickstart-config
  params:
    - name: vmName
      value: "my-fedora-vm"
```

### Custom Kickstart Template

Provide your own kickstart configuration:

```yaml
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  generateName: create-kickstart-config-custom-
spec:
  taskRef:
    kind: Task
    name: create-kickstart-config
  params:
    - name: vmName
      value: "my-rhel-vm"
    - name: kickstartTemplate
      value: |
        text
        lang en_US.UTF-8
        keyboard us
        network --bootproto=dhcp
        rootpw --plaintext mypassword
        # ... rest of your kickstart config
```

### With Owner Reference

Set owner reference for automatic cleanup:

```yaml
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  generateName: create-kickstart-config-owner-
spec:
  taskRef:
    kind: Task
    name: create-kickstart-config
  params:
    - name: vmName
      value: "my-vm"
    - name: setOwnerReference
      value: "true"
```

## Default Kickstart Template

The default template includes:
- Automated partitioning for Secure Execution (boot, se, newroot, root partitions)
- Essential packages (s390-tools, cryptsetup, cloud-init, etc.)
- GPT partition table
- XFS and ext4 filesystems
- Minimal Fedora/RHEL installation

## Integration with Secure Execution Pipeline

This task is typically used as the first step in the Secure Execution VM provisioning pipeline:

1. **create-kickstart-config** - Creates kickstart ConfigMap
2. deploy-kickstart-server - Serves the kickstart file via HTTP
3. create-vm - Creates VM that uses the kickstart for automated installation

## Notes

- The ConfigMap is labeled with `vm-name` for easy identification and cleanup
- The default template is optimized for s390x Secure Execution workflows
- Partition labels (boot, se, NEWROOT, root) are critical for SE automation scripts
- The kickstart file is stored in the ConfigMap under the key `ks.cfg`

## RBAC Requirements

This task requires permissions to:
- Create and manage ConfigMaps
- Get and list namespaces (for namespace resolution)

See the generated RBAC files for complete permission requirements.