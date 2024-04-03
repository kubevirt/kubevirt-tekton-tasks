# Windows BIOS Installer Pipeline

This Pipeline installs Windows 10 into a new DataVolume. This DataVolume is suitable to be used as a default boot source or golden image for Windows 10 VirtualMachines.

The Pipeline implements this by spinning up a new VirtualMachine which boots from the Windows installation image (ISO file). The installation of Windows is automatically executed and controlled by a Windows answer file. Then the Pipeline will wait for the installation to complete and will delete the created VirtualMachine while keeping the resulting DataVolume with the installed operating system. The Pipeline can be customized to support different installation requirements.

## Prerequisites

- KubeVirt `>=v1.0.0`
- Tekton Pipelines `>=v0.43.0`

### Obtain Windows ISO download URL

1. Go to https://www.microsoft.com/en-us/software-download/windows10ISO.
   You can also obtain a server edition for evaluation at https://www.microsoft.com/en-us/evalcenter/evaluate-windows-server-2019 (needs a different answer file!).
2. Fill in the edition and `English` language (other languages need to be updated in windows10-autounattend ConfigMap) and go to the download page.
3. Right-click on the 64-bit download button and copy the download link. The link should be valid for 24 hours.

### Prepare autounattend.xml ConfigMap

1. Supply, generate or use the default autounattend.xml.
   For information on answer files see [Startup Scripts - KubeVirt User Guide](https://kubevirt.io/user-guide/virtual_machines/startup_scripts/#sysprep).
2. Replace the default example autounattend.xml with your own in the definition of the `windows10-autounattend` ConfigMap in the Pipeline YAML.
   Different autounattend.xml can be also passed in a separate ConfigMap with the Pipeline parameter `autounattendConfigMapName` when creating a PipelineRun.

## Pipeline Description

```
  create-vm-root-disk --- create-vm --- wait-for-vmi-status --- cleanup-vm
```

1. `create-vm-root-disk` Task creates an empty DataVolume.
2. `create-vm` Task creates a VirtualMachine called `windows-bios-installer-*`
   from the empty DataVolume and with the `windows-bios-installer-cd-rom` DataVolume attached as a CD-ROM.
   A second DataVolume with the virtio-win ISO will also be attached (Pipeline parameter `virtioContainerDiskName`). The VirtualMachine has to be created in the same namespace as the empty DataVolume.
3. `wait-for-vmi-status` Task waits until the VirtualMachine shuts down.
4. `cleanup-vm` deletes the installer VirtualMachine and ISO DataVolume (also in case of failure of the previous Tasks).
5. The output artifact will be the `win10` DataVolume with the basic Windows installation.
   It will boot into the Windows OOBE and needs to be setup further before it can be used.

## How to run

Before you create PipelineRuns, you must create ConfigMaps with an autounattend.xml in the same namespace in which the VirtualMachine will be created.
Examples of ConfigMaps can be found [here](https://github.com/kubevirt/kubevirt-tekton-tasks/tree/main/release/pipelines/windows-bios-installer/configmaps).

Pipeline run with resolver:
{% for item in pipeline_runs_yaml %}
```yaml
export WIN_IMAGE_DOWNLOAD_URL=$(./getisourl.py) # see paragraph Obtaining a download URL in an automated way

oc create -f - <<EOF
{{ item | to_nice_yaml }}EOF
```
{% endfor %}

## Possible Optimizations

#### Obtaining a download URL in an automated way

The script [`getisourl.py`](https://github.com/kubevirt/kubevirt-tekton-tasks/blob/main/release/pipelines/windows-bios-installer/getisourl.py) can be used to automatically obtain a Windows 10 ISO download URL.

The prerequisites are:

- python3-selenium
- chromedriver
- chromium

Run it as follows to initialize a WIN_URL variable.

```bash
# Real URL can look differently
WIN_IMAGE_DOWNLOAD_URL=$(./getisourl.py)
```

### Windows Source ISO Caching

Windows source ISO can be downloaded and cached/stored in a DataVolume. So the subsequent PipelineRuns don't have to download the same ISO multiple times.

You can for example create the following DataVolume first and remove create-source-dv step from the windows10-installer Pipeline:

```yaml
WIN_IMAGE_DOWNLOAD_URL=$(./getisourl.py)

oc create -f - <<EOF
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  annotations:
    "cdi.kubevirt.io/storage.bind.immediate.requested": "true"
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
      url: ${WIN_IMAGE_DOWNLOAD_URL}
EOF
```

## Cancelling/Deleting PipelineRun

When running the example Pipelines, they create temporary objects (DataVolumes, VirtualMachines, etc.). Each Pipeline has its own clean up system which should keep the cluster clean from leftovers. In case user hard deletes or cancels running PipelineRun, the PipelineRun will not clean temporary objects and objects will stay in the cluster and then they have to be deleted manually. To prevent this behaviour, cancel the [PipelineRun gracefully](https://tekton.dev/docs/pipelines/pipelineruns/#gracefully-cancelling-a-pipelinerun). It triggers special Tasks, which remove temporary objects and keep only result DataVolume/PVC.
