# Windows Customize Pipeline

This pipeline clones the DataVolume of a basic and generalized Windows 10/11/2k22 installation and runs arbitrary customization
commands through an unattend.xml after startup of the VM. As an example a ConfigMap which installs Microsoft SQL
Server Express and generalizes the VM after is included (`windows-sqlserver`) or a ConfigMap that install VSCode (`windows-vs-code`). 
For basic setup after the first start of a customized VM an example unattend.xml is included in the pipeline's ConfigMap `windows10-unattend`, 
or `windows11-unattend`.

This example pipeline can be used for running windows 10, 11 or 2k22 (or others - not tested!). Always adjust pipeline parameters 
for windows version you are currently using(e.g. different template name, different base image name, ...). It is possible to use 
`windows-sqlserver` config map for win11/2k22 and vice versa (`windows-vs-code` for win10/2k22).

The provided reference ConfigMap (`windows-sqlserver`) boots Windows 10/11/2k22 into Audit mode, applies the customizations as
part of a Powershell script (ran by `SynchronousCommand`) and then generalizes the VM again. The Powershell
script can be adapted as desired to apply other customizations.

There is a specific version of this pipeline for OKD. This version is using templates, which are not available on Kubernetes.
A new golden template is created after a successful customization through this version.

## Prerequisites

- KubeVirt `v1.0.0`
- Tekton Pipelines `v0.44.0`

## Links

- [Windows Customize Pipeline](https://github.com/kubevirt/ssp-operator/blob/main/data/tekton-pipelines/windows-customize-pipeline.yaml)
- [Windows 10 Customize PipelineRun](windows10-customize-pipelinerun.yaml)
- [Windows 11 Customize PipelineRun](windows11-customize-pipelinerun.yaml)
- [Windows 2k22 Customize PipelineRun](windows2k22-customize-pipelinerun.yaml)

### Prepare unattend.xml ConfigMap

1. Supply, generate or use the default unattend.xml.
   For information on answer files see [Startup Scripts - KubeVirt User Guide](https://kubevirt.io/user-guide/virtual_machines/startup_scripts/#sysprep).
2. Create a new ConfigMap with the unattend.xml
3. Pass the name of the new ConfigMap to the PipelineRun with the parameter `customizeConfigMapName`.

## Pipeline Description

```
  copy-vm-root-disk --- create-vm --- wait-for-vmi-status --- cleanup-vm
```

1. `copy-vm-root-disk` task copies PVC defined in `sourceDiskImageName` and `sourceDiskImageNamespace` parameters.
2. `create-vm` task creates a VM called `windows-customize-*`
   from the base DV and with the customize ConfigMap attached as a CD-ROM. (Pipeline parameter `customizeConfigMapName`)
3. `wait-for-vmi-status` task waits until the VM shuts down.
4. `cleanup-vm` deletes the installer VM. (also in case of failure of the previous tasks)
5. The output artifact will be the `win*-customized` DV with the customized Windows installation.
   It will boot into the Windows OOBE and needs to be setup further before it can be used. (depends on the applied customizations)
6. The `windows11-unattend` ConfigMap can be used to boot the VM into the Desktop. (depends on the applied customizations)

## How to run

```bash
oc apply -f windows11-customize.yaml
oc create -f windows11-customize-pipelinerun-okd.yaml
```
## Cancelling/Deleteting pipelineRuns

When running the example pipelines, they create temporary objects (DataVolumes, VMs, templates, ...). Each pipeline has its own clean up system which 
should keep the cluster clean from leftovers. In case user hard deletes or cancels running pipelineRun, the pipelineRun will not clean temporary 
objects and objects will stay in the cluster and then they have to be deleted manually. To prevent this behaviour, cancel the 
[pipelineRun gracefully](https://tekton.dev/docs/pipelines/pipelineruns/#gracefully-cancelling-a-pipelinerun). It triggers special tasks, 
which remove temporary objects and keep only result DataVolume/PVC.
