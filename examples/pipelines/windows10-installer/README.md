# Windows 10 Installer Pipeline

This pipeline installs Windows 10 into a new DataVolume. This DataVolume is suitable to be used as a default boot source
or golden image for Windows 10 VMs.

The pipeline implements this by spinning up a new VM which boots from the Windows installation image (ISO file). The
installation of Windows is automatically executed and controlled by a Windows answer file. Then the pipeline will wait
for the installation to complete and will delete the created VM while keeping the resulting DataVolume with the
installed operating system. The pipeline can be customized to support different installation requirements.

There is a specific version of this pipeline for OKD.
This version is using templates, which are not available on Kubernetes.

## Prerequisites

- KubeVirt `v0.53.1`
- Tekton Pipelines `v0.35.0`

## Links

- [Windows 10 Installer Pipeline for Kubernetes](https://github.com/kubevirt/tekton-tasks-operator/blob/main/data/tekton-pipelines/kubernetes/windows10-installer.yaml)
- [Windows 10 Installer Pipeline for OKD](https://github.com/kubevirt/tekton-tasks-operator/blob/main/data/tekton-pipelines/okd/windows10-installer.yaml)
- For example PipelineRuns see commit message of [2c4daed](https://github.com/kubevirt/tekton-tasks-operator/commit/2c4daed4124654a765f69acdc2b1c7390ee3c2f4)

### Obtain Windows ISO download URL

1. Go to https://www.microsoft.com/en-us/software-download/windows10ISO.
   You can also obtain a server edition for evaluation at https://www.microsoft.com/en-us/evalcenter/evaluate-windows-server-2019. (needs a different answer file!)
2. Fill in the edition and `English` language (other languages need to be updated in windows10-autounattend ConfigMap) and go to the download page.
3. Right-click on the 64-bit download button and copy the download link. The link should be valid for 24 hours.
4. Initialize a WIN_URL variable that will be used to create a DataVolume which will download this ISO into a PVC.

```bash
# Real URL can look differently
WIN_URL="https://software.download.prss.microsoft.com/db/Win10_21H2_English_x64.iso..."
```

### Prepare autounattend.xml ConfigMap

1. Supply, generate or use the default autounattend.xml.
   For information on answer files see [Startup Scripts - KubeVirt User Guide](https://kubevirt.io/user-guide/virtual_machines/startup_scripts/#sysprep).
2. Replace the default example autounattend.xml with your own in the definition of the `windows10-autounattend` ConfigMap in the pipeline YAML.
   Different autounattend.xml can be also passed in a separate ConfigMap with the Pipeline parameter `autounattendConfigMapName` when creating a PipelineRun.

## Pipeline Description (Kubernetes)

```
  create-source-dv --- create-vm-from-manifest --- wait-for-vmi-status --- cleanup-vm
                    |
    create-base-dv --
```

1. `create-source-dv` task downloads a Windows source ISO into a DV called `windows10-source-*`.
2. `create-base-dv` task creates an empty DV for new windows installation called `windows10-base-*`.
3. `create-vm-from-manifest` task creates a VM called `windows10-installer-*`
   from the empty DV and with the `windows10-source-*` DV attached as a CD-ROM.
   A second DV with the virtio-win ISO will also be attached. (Pipeline parameter `virtioContainerDiskName`)
4. `wait-for-vmi-status` task waits until the VM shuts down.
5. `cleanup-vm` deletes the installer VM and ISO DV.
6. The output artifact will be the `windows10-base-*` DV with the basic Windows installation.
   It will boot into the Windows OOBE and needs to be setup further before it can be used.

## Pipeline Description (OKD)

```
  copy-template --- modify-vm-template --- create-vm-from-template --- wait-for-vmi-status --- create-base-dv --- cleanup-vm
```

1. `copy-template` copies the template defined by the pipeline parameters `sourceTemplateName` and `sourceTemplateNamespace`
    to a new template with the name specified by parameter `installerTemplateName` in the same namespace. 
    An already existing template can be overwritten when setting `allowReplaceInstallerTemplate` to `true`.
2. `modify-vm-template` sets the display name of the new Template and the dataVolumeTemplates, Disks and Volumes needed for the installation.
3. `create-vm-from-template` task creates a VM from the newly created Template.
   A DV with the Windows source ISO will be attached as CD-ROM and a second empty DV will be used as installation destination.
   A third DV with the virtio-win ISO will also be attached. (Pipeline parameter `virtioContainerDiskName`)
4. `wait-for-vmi-status` task waits until the VM shuts down.
5. `create-base-dv` task creates an DV with the specified name and namespace (Pipeline parameters `baseDvName` and `baseDvNamespace`).
    Then it clones the second DV of the installation VM into the new DV.
6. `cleanup-vm` deletes the installer VM and all of its DVs.
7. The output artifact will be the `baseDvName`/`baseDvNamespace` DV with the basic Windows installation. 
   It will boot into the Windows OOBE and needs to be setup further before it can be used.

## How to run (Kubernetes)

```bash
WIN_URL="https://software.download.prss.microsoft.com/db/Win10_21H2_English_x64.iso..."
kubectl apply -f windows10-installer-kubernetes.yaml
sed 's!INSERT_WINDOWS_ISO_URL!'"$WIN_URL"'!g' windows10-installer-pipelinerun-kubernetes.yaml | kubectl create -f -
```

## How to run (OKD)

```bash
WIN_URL="https://software.download.prss.microsoft.com/db/Win10_21H2_English_x64.iso..."
oc apply -f windows10-installer-okd.yaml
sed 's!INSERT_WINDOWS_ISO_URL!'"$WIN_URL"'!g' windows10-installer-pipelinerun-okd.yaml | oc create -f -
```

## Possible Optimizations

### Windows Source ISO Caching

Windows source ISO can be downloaded and cached/stored in a DV.
So the subsequent PipelineRuns don't have to download the same ISO multiple times.

You can for example create the following DV first and remove create-source-dv step from the windows10-installer pipeline:

```yaml
 apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: windows10-source
spec:
  pvc:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 7Gi
    volumeMode: Filesystem
  source:
    http:
      url: WIN_IMAGE_DOWNLOAD_URL
```
