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

- KubeVirt `v0.59.1`
- Tekton Pipelines `v0.44.0`

## Links

- [Windows Customize Pipeline for Kubernetes](https://github.com/kubevirt/tekton-tasks-operator/blob/main/data/tekton-pipelines/kubernetes/windows-customize-pipeline.yaml)
- [Windows 10 Customize PipelineRun for Kubernetes](windows10-customize-pipelinerun-kubernetes.yaml)
- [Windows Customize Pipeline for OKD](https://github.com/kubevirt/tekton-tasks-operator/blob/main/data/tekton-pipelines/okd/windows-customize-pipeline.yaml)
- [Windows 10 Customize PipelineRun for OKD](windows10-customize-pipelinerun-okd.yaml)
- [Windows 11 Customize PipelineRun for OKD](windows11-customize-pipelinerun-okd.yaml)
- [Windows 2k22 Customize PipelineRun for OKD](windows2k22-customize-pipelinerun-okd.yaml)

### Prepare unattend.xml ConfigMap

1. Supply, generate or use the default unattend.xml.
   For information on answer files see [Startup Scripts - KubeVirt User Guide](https://kubevirt.io/user-guide/virtual_machines/startup_scripts/#sysprep).
2. Create a new ConfigMap with the unattend.xml
3. Pass the name of the new ConfigMap to the PipelineRun with the parameter `customizeConfigMapName`.

## Pipeline Description (Kubernetes)

```
  create-base-dv --- create-vm-from-manifest --- wait-for-vmi-status --- cleanup-vm
```

1. `create-base-dv` task creates an empty DV for the customized windows installation called `windows10-base-*`.
2. `create-vm-from-manifest` task creates a VM called `windows-installer-*`
   from the base DV and with the customize ConfigMap attached as a CD-ROM. (Pipeline parameter `customizeConfigMapName`)
3. `wait-for-vmi-status` task waits until the VM shuts down.
4. `cleanup-vm` deletes the installer VM. (also in case of failure of the previous tasks)
5. The output artifact will be the `windows-base-*` DV with the customized Windows installation.
   It will boot into the Windows OOBE and needs to be setup further before it can be used. (depends on the applied customizations)
6. The `windows10-unattend` ConfigMap can be used to boot the VM into the Desktop. (depends on the applied customizations)

## Pipeline Description (OKD)

```
  copy-template-customize --- modify-vm-template-customize --- create-vm-from-template --- wait-for-vmi-status --- create-base-dv --- copy-template-golden --- modify-vm-template-golden --- cleanup-vm
```

1. `copy-template-customize` copies the template defined by the pipeline parameters `sourceTemplateName` and `sourceTemplateNamespace`
    to a new template with the name specified by parameter `customizeTemplateName` in the same namespace.
    An already existing template can be overwritten when setting `allowReplaceCustomizationTemplate` to `true`.
2. `modify-vm-template-customize` sets the display name of the new Template and the dataVolumeTemplates, Disks and Volumes needed for the customization.
3. `create-vm-from-template` task creates a VM from the newly created Template.
   A DV with the customize ConfigMap will be attached as CD-ROM. (Pipeline parameter `customizeConfigMapName`)
4. `wait-for-vmi-status` task waits until the VM shuts down.
5. `create-base-dv` task creates an DV called `windows-base-*`, then it clones the DV of the customize VM into the new DV.
6. `copy-template-golden` copies the template defined by the pipeline parameters `sourceTemplateName` and `sourceTemplateNamespace`
   to a new template with the name specified by parameter `goldenTemplateName` in the same namespace.
   An already existing template can be overwritten when setting `allowReplaceGoldenTemplate` to `true`.
7. `modify-vm-template-golden` sets the display name of the new Template and the dataVolumeTemplates, Disks and Volumes needed to create customized VMs.
8. `cleanup-vm` deletes the customize VM and all of its DVs. (also in case of failure of the previous tasks)
9. The output artifact will be the `goldenTemplateName` Template with the customized Windows installation.
   From this template the user can create VMs with customizations applied.
   With the windows10-sqlserver ConfigMap VMs will boot into the Windows OOBE and need to be setup further before they can be used.

## How to run (Kubernetes)

```bash
SOURCE_DV_NAME=example-dvname
SOURCE_DV_NAMESPACE=example-dvnamespace
kubectl apply -f windows10-customize.yaml
sed 's!INSERT_NAME_OF_SOURCE_DATAVOLUME!'"$SOURCE_DV_NAME"'!g' windows10-customize-pipelinerun-kubernetes.yaml | \
sed 's!INSERT_NAMESPACE_OF_SOURCE_DATAVOLUME!'"$SOURCE_DV_NAMESPACE"'!g' | \
kubectl create -f -
```

## How to run (OKD)

```bash
oc apply -f windows10-customize.yaml
oc create -f windows10-customize-pipelinerun-okd.yaml
```
