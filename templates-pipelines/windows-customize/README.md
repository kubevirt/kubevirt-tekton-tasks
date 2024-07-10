# Windows Customize Pipeline

This Pipeline clones the DataVolume of a basic and generalized Windows 10, 11 or Server 2k22 installation and runs arbitrary customization commands through an unattend.xml after startup of the VirtualMachine. As an example a ConfigMap which installs Microsoft SQL Server Express and generalizes the VirtualMachine after (`windows-sqlserver`) and a ConfigMap that installs VSCode (`windows-vs-code`) are included.
For basic setup after the first start of a customized VirtualMachine an example unattend.xml is included in the Pipeline's ConfigMap `windows10-unattend`, or `windows11-unattend`.

This example Pipeline can be used for running Windows 10, 11 or Server 2k22 (or others - not tested!). Always adjust Pipeline parameters for the Windows version you are currently using (e.g. differe name, different base image name, etc.). It is possible to use `windows-sqlserver` ConfigMap for Windows 11 or Server 2k22 and vice versa (`windows-vs-code` for Windows 10 or Server 2k22).

The provided reference ConfigMap (`windows-sqlserver`) boots Windows 10, 11 or Windows Server 2k22 into Audit mode, applies the customizations as part of Powershell script (ran by `SynchronousCommand`) and then generalizes the VirtualMachine again. The Powershell script can be adapted as desired to apply other customizations.

## Prerequisites

- {{ virt_project }} `>= {{ virt_project_version }}`
- {{ tekton_project }} `>= {{ tekton_project_version }}`

### Prepare unattend.xml ConfigMap

1. Supply, generate or use the default unattend.xml. For information on answer files see [Startup Scripts - KubeVirt User Guide](https://kubevirt.io/user-guide/virtual_machines/startup_scripts/#sysprep).
2. Create a new ConfigMap with the unattend.xml
3. Pass the name of the new ConfigMap to the PipelineRun with the parameter `customizeConfigMapName`.

## Pipeline Description

```
  import-unattend-configmaps --- copy-vm-root-disk --- create-vm --- wait-for-vmi-status --- create-datasource-root-disk --- cleanup-vm --- delete-imported-configmaps
```
1. `import-unattend-configmaps` imports ConfigMap with `unattend.xml` needed for automated customization of Windows.
2. `copy-vm-root-disk` Task copies PVC defined in `sourceDiskImageName` and `sourceDiskImageNamespace` parameters.
3. `create-vm` Task creates a VirtualMachine called `windows-customize-*` from the base DataVolume and with the customize ConfigMap attached as a CD-ROM (Pipeline parameter `customizeConfigMapName`). The VirtualMachine has to be created in the same namespace as the source DataVolume.
4. `wait-for-vmi-status` Task waits until the VirtualMachine shuts down.
5. `create-datasource-root-disk` Task creates a DataSource object, which is used by UI for discovering bootable volumes and links PVC created in `copy-vm-root-disk` step.
6. `cleanup-vm` deletes the installer VirtualMachine (also in case of failure of the previous Tasks).
7. The output artifact will be the `win*-customized` DataVolume with the customized Windows installation. It will boot into the Windows OOBE and needs to be setup further before it can be used (depends on the applied customizations).
8. The `windows11-unattend` ConfigMap can be used to boot the VirtualMachine into the Desktop (depends on the applied customizations).
9. `delete-imported-configmaps` deletes imported ConfigMaps.

## How to run

The pipeline uses a ConfigMap containing an `unattend.xml` file for automated customization of Windows. Example ConfigMaps are deployed within the Pipeline. In case you would like to use a different ConfigMap, specify a different URL in the `unattendXMLConfigMapsURL` parameter and adjust `customizeConfigMapName` parameter with correct the `ConfigMap` name. Examples of ConfigMaps can be found [here](https://github.com/kubevirt/kubevirt-tekton-tasks/tree/main/release/pipelines/windows-customize/configmaps).

> [!NOTE]
> By default, the Pipeline requires the ServiceAccount `pipeline` to exist. Tekton does not create this ServiceAccount 
> in namespaces which name starts with `openshift` or `kube`. In case you would like to run this Pipeline in a namespace which 
> starts with `openshift` or `kube`, you have to create the `pipeline` ServiceAccount manually or specify a different ServiceAccount in the PipelineRun.

Pipeline runs with resolvers:
{% for item in pipeline_runs_yaml %}
```yaml
oc create -f - <<EOF
{{ item | to_nice_yaml }}EOF
```
{% endfor %}

## Cancelling/Deleting PipelineRuns

When running the example Pipelines, they create temporary objects (DataVolumes, VirtualMachines, etc.). Each Pipeline has its own clean up system which should keep the cluster clean from leftovers. In case user hard deletes or cancels running PipelineRun, the PipelineRun will not clean temporary objects and objects will stay in the cluster. To prevent this behaviour, cancel the [PipelineRun gracefully](https://tekton.dev/docs/pipelines/pipelineruns/#gracefully-cancelling-a-pipelinerun). It triggers special Tasks, which remove temporary objects and keep only result DataSource/DataVolume/PVC.

Each object created by the Pipeline has OwnerReference to the Pod which created them (result DataVolume does not have it). In the case that the clean up steps are not triggered, by deleting the PipelineRun all leftover objects created by the Pipeline will be deleted.
