# Windows BIOS Installer Pipeline

This pipeline installs Windows 10 into a new DataVolume. This DataVolume is suitable to be used as a default boot source
or golden image for Windows 10 VMs.

The pipeline implements this by spinning up a new VM which boots from the Windows installation image (ISO file). The
installation of Windows is automatically executed and controlled by a Windows answer file. Then the pipeline will wait
for the installation to complete and will delete the created VM while keeping the resulting DataVolume with the
installed operating system. The pipeline can be customized to support different installation requirements.

## Prerequisites

- KubeVirt `v1.0.0`
- Tekton Pipelines `v0.44.0`

## Links

- [Windows BIOS Installer Pipeline](https://github.com/kubevirt/ssp-operator/blob/main/data/tekton-pipelines/windows-bios-installer-pipeline.yaml)
- [Windows BIOS Installer PipelineRun](windows10-installer-pipelinerun.yaml)

### Obtain Windows ISO download URL

1. Go to https://www.microsoft.com/en-us/software-download/windows10ISO.
   You can also obtain a server edition for evaluation at https://www.microsoft.com/en-us/evalcenter/evaluate-windows-server-2019 (needs a different answer file!).
2. Fill in the edition and `English` language (other languages need to be updated in windows10-autounattend ConfigMap) and go to the download page.
3. Right-click on the 64-bit download button and copy the download link. The link should be valid for 24 hours.
4. Initialize a WIN_URL variable that will be used to create a DataVolume which will download this ISO into a PVC.
   Make sure to escape `&` with `\&` for the example commands below to work.

```bash
# Real URL can look differently
WIN_URL="https://software.download.prss.microsoft.com/db/Win10_21H2_English_x64.iso..."
```

#### Obtaining a download URL in an automated way

The script [`getisourl.py`](getisourl.py) can be used to automatically obtain a Windows 10 ISO download URL.

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
2. Replace the default example autounattend.xml with your own in the definition of the `windows10-autounattend` ConfigMap in the pipeline YAML.
   Different autounattend.xml can be also passed in a separate ConfigMap with the Pipeline parameter `autounattendConfigMapName` when creating a PipelineRun.

## Pipeline Description

```
  create-vm-root-disk --- create-vm --- wait-for-vmi-status --- cleanup-vm
```

1. `create-vm-root-disk` task creates an empty DV.
2. `create-vm` task creates a VM called `windows-bios-installer-*`
   from the empty DV and with the `windows-bios-installer-cd-rom` DV attached as a CD-ROM.
   A second DV with the virtio-win ISO will also be attached. (Pipeline parameter `virtioContainerDiskName`)
3. `wait-for-vmi-status` task waits until the VM shuts down.
4. `cleanup-vm` deletes the installer VM and ISO DV. (also in case of failure of the previous tasks)
5. The output artifact will be the `win10` DV with the basic Windows installation.
   It will boot into the Windows OOBE and needs to be setup further before it can be used.

## How to run

```bash
WIN_URL="https://software.download.prss.microsoft.com/db/Win10_21H2_English_x64.iso..."
kubectl apply -f windows10-installer.yaml
sed 's!INSERT_WINDOWS_ISO_URL!'"$WIN_URL"'!g' windows10-installer-pipelinerun.yaml | kubectl create -f -
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
## Cancelling/Deleteting pipelineRuns

When running the example pipelines, they create temporary objects (DataVolumes, VMs, templates, ...). Each pipeline has its own clean up system which 
should keep the cluster clean from leftovers. In case user hard deletes or cancels running pipelineRun, the pipelineRun will not clean temporary 
objects and objects will stay in the cluster and then they have to be deleted manually. To prevent this behaviour, cancel the 
[pipelineRun gracefully](https://tekton.dev/docs/pipelines/pipelineruns/#gracefully-cancelling-a-pipelinerun). It triggers special tasks, 
which remove temporary objects and keep only result DataVolume/PVC.
