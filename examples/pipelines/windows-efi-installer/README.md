# Windows EFI Installer Pipeline

This pipeline installs Windows 11/2k22 into a new DataVolume. This DataVolume is suitable to be used as a default boot source
or golden image for Windows 11/2k22 VMs.

This example pipeline is suitable only for windows 11/2k22 (or other windows versions which require EFI - not tested!). When using 
this example pipeline always adjust pipeline parameters for windows version you are currently using (e.g. different template name, 
different autoattend config map, different base image name, ...). Each windows version requires change in autounattendConfigMapName 
parameter (e.g. using `windows2k22-autounattend` config map will not work with Windows 11 and vice versa - e.g. due to different storage 
drivers path).

The pipeline implements this by modifying the windows iso - extracts iso files from iso, replaces prompt bootloader with no-prompt bootloader. 
This helps with automated installation of Windows in EFI boot mode. By default Windows in EFI boot mode uses a prompt bootloader, which will not 
continue with the boot process until a key is pressed. By replacing it with the non-prompt bootloader no key press is required to boot into the 
Windows installer. Then task packs updated packages to new iso, converts it with qemu-img and replaces original iso file in PVC.

After the iso is modified it creates a new VM which boots from the modified Windows installation image (ISO file). The installation of Windows is 
automatically executed and controlled by a Windows answer file. Then the pipeline will wait for the installation to complete and will delete the 
created VM while keeping the resulting DataVolume with the installed operating system. The pipeline can be customized to support different 
installation requirements.

## Prerequisites

- KubeVirt `v1.0.0`
- Tekton Pipelines `v0.44.0`

## Links

- [Windows EFI Pipeline](https://github.com/kubevirt/ssp-operator/blob/main/data/tekton-pipelines/windows-efi-installer-pipeline.yaml)
- [Windows 11 EFI Installer PipelineRun](windows11-installer-pipelinerun.yaml)
- [Windows 2k22 EFI Installer PipelineRun](windows2k22-installer-pipelinerun.yaml)


### Obtain Windows ISO 11 download URL

1. Go to https://www.microsoft.com/en-us/software-download/windows11.
2. Fill in the edition and `English` language (other languages need to be updated in `windows11-autounattend` ConfigMap) and go to the download page.
3. Right-click on the 64-bit download button and copy the download link. The link should be valid for 24 hours.
4. Initialize a WIN_URL variable that will be used to create a DataVolume which will download this ISO into a PVC.
   Make sure to escape `&` with `\&` for the example commands below to work.

### Obtain Windows ISO server 2022 download URL

1. Go to https://www.microsoft.com/en-us/evalcenter/evaluate-windows-server-2022.
2. Right-click Download the ISO button.
3. Fill in all required informations and click on Download now button
4. Select English (United States) (other languages need to be updated in `windows2k22-autounattend` ConfigMap) - 64-bit edition iso download
5. Replace the `<INSERT_WINDOWS_ISO_URL>` tag with the new link in [windows2k22-installer-pipelinerun-okd](windows2k22-installer-pipelinerun-okd.yaml) file

```bash
# Real URL can look differently
WIN_URL="https://software.download.prss.microsoft.com/dbazure/Win11_22H2_English_x64v1.iso..."
```

#### Obtaining a download URL in an automated way

The script [`getisourl.py`](getisourl.py) can be used to automatically obtain a Windows 11 ISO download URL.

The prerequisites are:

- python3-selenium
- chromedriver
- chromium

Run it as follows to initialize a WIN_URL variable.

```bash
# Real URL can look differently
WIN_URL=$(./getisourl.py | sed 's/&/\\&/g')
```

### Prepare autounattend.xml ConfigMap

1. Supply, generate or use the default autounattend.xml.
   For information on answer files see [Startup Scripts - KubeVirt User Guide](https://kubevirt.io/user-guide/virtual_machines/startup_scripts/#sysprep).
2. Replace the default example autounattend.xml with your own in the definition of the `windows11-autounattend` ConfigMap in the pipeline YAML.
   Different autounattend.xml can be also passed in a separate ConfigMap with the Pipeline parameter `autounattendConfigMapName` when creating a PipelineRun.

## Pipeline Description

```
  import-win-iso --- modify-windows-iso-file --- create-vm --- wait-for-vmi-status --- create-base-dv --- cleanup-vm
                                              |
                       create-vm-root-disk --- 
```

1. `create-vm-root-disk` creates empty DV which is used for windows installation.
2. `import-win-iso` creates new datavolume with windows iso file with name defined in `isoDVName` parameter.
3. `modify-windows-iso-file` extracts imported iso, replaces prompt bootloader (which is used as a default one when EFI is used) 
   with no-prompt bootloader, pack the updated files back to new iso, convert the iso and replaces original iso with updated one.
   Replacement of bootloader is needed to be able to automate installation of windows versions which require EFI.
4. `create-vm` task creates a VM.
   A DV with the Windows source ISO will be attached as CD-ROM and a second empty DV will be used as installation destination.
   A third DV with the virtio-win ISO will also be attached. (Pipeline parameter `virtioContainerDiskName`)
5. `wait-for-vmi-status` task waits until the VM shuts down.
6. `create-base-dv` task creates an DV with the specified name and namespace (Pipeline parameters `baseDvName` and 
   `baseDvNamespace`).
   Then it clones the second DV of the installation VM into the new DV.
7. `cleanup-vm` deletes the installer VM and all of its DVs.
8. The output artifact will be the `baseDvName`/`baseDvNamespace` DV with the basic Windows installation.
   It will boot into the Windows OOBE and needs to be setup further before it can be used.

## How to run

```bash
WIN_URL="https://software.download.prss.microsoft.com/dbazure/Win11_22H2_English_x64v1.iso..."
oc apply -f windows11-installer-pipelinerun.yaml
sed 's!INSERT_WINDOWS_ISO_URL!'"$WIN_URL"'!g' windows11-installer-pipelinerun.yaml | oc create -f -
```

## Possible Optimizations

### Windows Source ISO Caching

Windows source ISO can be downloaded and cached/stored in a DV.
So the subsequent PipelineRuns don't have to download the same ISO multiple times.

You can for example create the following DV first and remove import-win-iso and modify-windows-iso-file steps from the windows-efi-installer pipeline:

```yaml
 apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: windows11-source
spec:
  pvc:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 9Gi
    volumeMode: Filesystem
  source:
    http:
      url: WIN_IMAGE_DOWNLOAD_URL
```

## Cancelling/Deleteting pipelineRuns

When running the example pipelines, they create temporary objects (DataVolumes, VMs, templates, ...). Each pipeline has its own clean up system which 
should keep the cluster clean from leftovers. In case user hard deletes or cancels running pipelineRun, the pipelineRun will not clean temporary 
objects and objects will stay in the cluster and then they have to be deleted manually. To prevent this behaviour, cancel the 
[pipelineRun gracefully](https://tekton.dev/docs/pipelines/pipelineruns/#gracefully-cancelling-a-pipelinerun). It triggers special tasks, 
which remove temporary objects and keep only result DataVolume/PVC.

windows-efi-installer pipeline generates for each pipelineRun new source datavolume which contains imported iso. This DV has generated name and is 
deleted after pipeline succeeds. However, the created PVC will stay in cluster, but it will have terminating state. It will wait, until pipelinRun is 
deleted. This behaviour is caused by a fact, that PVC is mounted into modify-windows-iso taskRun pod and pvc can be deleted only when the pod does not 
exist.
