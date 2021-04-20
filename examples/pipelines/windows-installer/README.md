# Windows Installer Pipeline

Downloads a Windows Source ISO into a PVC and installs Windows into a new base PVC.
It will then spin up an installation VM and use Windows Answer Files to automatically install the VM.
Then the pipeline will wait for the installation to complete and will delete the installation VM while keeping the artifact PVC (backed by DataVolume) with the installed operating system.
The pipeline can be customized to support different installation requirements.

## Prerequisites

- KubeVirt `v0.39.0`
- Tekton Pipelines `v0.19.0`

### Obtain Windows ISO download URL

1. Go to https://www.microsoft.com/en-us/software-download/windows10ISO. You can also obtain a server edition for evaluation at https://www.microsoft.com/en-us/evalcenter/evaluate-windows-server-2019.
2. Fill in the edition and `English` language (other languages need to be updated in windows-10-autounattend ConfigMap) and go to the download page.
3. Right-click on the 64-bit download button and copy the download link. The link should be valid for 24 hours.
4. Initialize a WIN_URL variable that will be used to create a DataVolume which will download this ISO into a PVC.

```bash
# Real URL can look differently
WIN_URL="https://software-download.microsoft.com/db/Win10_20H2_v2_English_x64.iso..."
```

### Prepare autounattend.xml ConfigMap

1. Supply, generate or use the default autounattend.xml.
   The configuration file can be generated with [Windows SIM](https://docs.microsoft.com/en-us/windows-hardware/customize/desktop/wsim/windows-system-image-manager-overview-topics)
   or it can be specified manually according to [Answer File Reference](https://docs.microsoft.com/en-us/windows-hardware/customize/desktop/wsim/answer-files-overview)
   and [Answer File Components Reference](https://docs.microsoft.com/en-us/windows-hardware/customize/desktop/unattend/components-b-unattend).
2. Replace the default example autounattend.xml with your own in `windows-installer-pipeline.yaml`.
   You can also store the config inside a secret (VM definition requires changes in that case).
   Different autounattend.xml can be also passed in the Pipeline parameters when creating a PipelineRun.

## Pipeline Description

```
  create-source-dv --- create-vm-from-manifest --- wait-for-vmi-status --- cleanup-vm
                    |
    create-base-dv --
```

1. `create-source-dv` task downloads a Windows source ISO into a PVC called `windows-10-source-*`.
2. `create-base-dv` task creates empty PVC for new windows installation called `windows-10-base-*`.
3. `create-vm-from-manifest` task creates a VM called `windows-installer-*`
   from the empty PVC and with `windows-10-source-*` PVC attached as a CD-ROM.
4. ` wait-for-vmi-status` task waits until the VM shutdowns.
5. `cleanup-vm` deletes the installer VM and ISO PVC.
6.  The output artifact will be the `windows-10-base-*` PVC with the Windows installation. 
    It includes an `Administrator` user with `changepassword` password.

## How to run

```bash
WIN_URL="https://software-download.microsoft.com/db/Win10_20H2_v2_English_x64.iso..."
kubectl apply -f windows-installer-pipeline.yaml
sed 's!DOWNLOAD_URL!'"$WIN_URL"'!g' windows-installer-pipelinerun.yaml | kubectl create -f -
```

## Possible Optimizations

### Windows Source ISO Caching

Windows source ISO can be downloaded and cached/stored in a PVC.
So the subsequent PipelineRuns don't have to download the same ISO multiple times.

You can create the following DV first and remove create-source-dv step from the windows-installer-pipeline:

```yaml
 apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: windows-10-source
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
