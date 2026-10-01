# KubeVirt Disk Uploader Task

Automates the extraction of disk and uploads it to a container registry, to be used in multiple Kubernetes clusters.

## Prerequisites

VMExport support must be enabled in the feature gates to be available. The [feature gates](https://kubevirt.io/user-guide/cluster_admin/activating_feature_gates/#how-to-activate-a-feature-gate) field in the KubeVirt Custom Resource (CR) must be expanded by adding the VMExport to it.

# Example Scenario

When user runs [KubeVirt Tekton Tasks](https://github.com/kubevirt/kubevirt-tekton-tasks) example pipelines (windows-installer, windows-customize) to prepare Windows disk images - The newly created disk image is only in a single cluster. If user wants to have it in another cluster, then KubeVirt Disk Uploader can be used to push it out of the cluster.

The uploaded OCI image's architecture comes from the selected export source:

- `vm`: the VM's `spec.template.spec.architecture`.
- `vmsnapshot`: the captured VM's architecture in the snapshot content.
- `pvc`: the PVC's `template.kubevirt.io/architecture` label, when it has a
  supported KubeVirt architecture (`amd64`, `arm64`, or `s390x`).

The architecture stays empty when a VM template has no architecture or a PVC
has no supported architecture label. A snapshot without its content causes an
error. When the architecture is set, the OCI platform OS is `linux`, which
describes the container image and does not indicate the guest OS.

The task's service account needs `get` permission on `virtualmachines` for VM
sources, `virtualmachinesnapshots` and `virtualmachinesnapshotcontents` for
snapshot sources, and `persistentvolumeclaims` for PVC sources, in the source
namespace. The example pipeline's disk-uploader Role includes these permissions.

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Usage

To get the Secret of the task run:
```
{{ secret_task_yaml | to_nice_yaml }}
```

Get `ACCESS_KEY_ID` or `SECRET_KEY` by running: `echo -n "<REGISTRY_USERNAME_OR_PASSWORD>" | base64`.

Task run using resolver:
```
{{ task_run_resolver_yaml | to_nice_yaml }}
```

### Platforms

The Task can be run on linux/amd64 platform.
