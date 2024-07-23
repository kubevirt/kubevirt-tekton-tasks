# Windows EFI Installer Pipeline

This Pipeline installs Windows 10, 11 or Server 2k22 into a new DataVolume. This DataVolume is suitable to be used as a default boot source
or golden image for Windows 10, 11 or Server 2k22 VirtualMachines.

This example Pipeline is suitable only for Windows 10, 11 or Server 2k22 (or other Windows versions which require EFI - not tested!). When using this example Pipeline always adjust Pipeline parameters for Windows version you are currently using (e.g. different name, different autoattend config map, different base image name, etc.). Each Windows version requires change in autounattendConfigMapName parameter (e.g. using `windows2k22-autounattend` config map will not work with Windows 11 and vice versa - e.g. due to different storage drivers path).

The Pipeline implements this by modifying the supplied Windows ISO. It extracts all files from the ISO, replaces the prompt bootloader with the no-prompt bootloader and creates a new bootable ISO.
This helps with automated installation of Windows in EFI boot mode. By default Windows in EFI boot mode uses a prompt bootloader, which will not continue with the boot process until a key is pressed. By replacing it with the no-prompt bootloader no key press is required to boot into the Windows installer. Then Task packs updated packages to new ISO, converts it with qemu-img and replaces original ISO file in PVC.

> [!IMPORTANT]
> If the ISO file does not contain files `efisys_noprompt.bin` and `cdboot_noprompt.efi` located at `efi/microsoft/boot/` 
> inside ISO file, the `modify-windows-iso-file` task will exit and whole Pipeline will fail. In case your ISO file has 
> different file structure, update `modify-windows-iso-file` task accordingly.

After the ISO is modified it creates a new VirtualMachine which boots from the modified Windows installation image (ISO file). The installation of Windows is automatically executed and controlled by a Windows answer file. Then the Pipeline will wait for the installation to complete and will delete the created VirtualMachine while keeping the resulting DataSource/DataVolume with the installed operating system. The Pipeline can be customized to support different installation requirements.

## Prerequisites

- {{ virt_project }} `>= {{ virt_project_version }}`
- {{ tekton_project }} `>= {{ tekton_project_version }}`
- Apply ConfigMaps with Windows sysprep (or select one which you will need) - examples can be found here: https://github.com/kubevirt/kubevirt-tekton-tasks/tree/main/release/pipelines/windows-efi-installer/configmaps

### Obtain Windows 10 ISO download URL

1. Go to https://www.microsoft.com/en-us/software-download/windows10ISO.
2. Fill in the edition and `English` language (other languages need to be updated in `windows10-autounattend` ConfigMap) and go to the download page.
3. Right-click on the 64-bit download button and copy the download link. The link should be valid for 24 hours.

### Obtain Windows 11 ISO download URL

1. Go to https://www.microsoft.com/en-us/software-download/windows11.
2. Fill in the edition and `English` language (other languages need to be updated in `windows11-autounattend` ConfigMap) and go to the download page.
3. Right-click on the 64-bit download button and copy the download link. The link should be valid for 24 hours.

### Obtain Windows Server 2022 ISO download URL

1. Go to https://www.microsoft.com/en-us/evalcenter/evaluate-windows-server-2022.
2. Right-click Download the ISO button.
3. Fill in all required informations and click on Download now button.
4. Select English (United States) (other languages need to be updated in `windows2k22-autounattend` ConfigMap) - 64-bit edition ISO download.

### Prepare autounattend.xml ConfigMap

1. Supply, generate or use the default autounattend.xml.
   For information on answer files see [Startup Scripts - KubeVirt User Guide](https://kubevirt.io/user-guide/virtual_machines/startup_scripts/#sysprep).
2. Replace the default example autounattend.xml with your own in the definition of the `windows11-autounattend` ConfigMap in the Pipeline YAML.
   Different autounattend.xml can be also passed in a separate ConfigMap with the Pipeline parameter `autounattendConfigMapName` when creating a PipelineRun.

## Pipeline Description

```
  import-autounattend-configmaps --- import-win-iso --- modify-windows-iso-file --- create-vm --- wait-for-vmi-status --- create-datasource-root-disk --- cleanup-vm --- delete-imported-configmaps
                                                     |
                              create-vm-root-disk ---
```
1. `import-autounattend-configmaps` imports ConfigMap with `autounattend.xml` needed for automated installation of Windows.
2. `create-vm-root-disk` creates empty DataVolume which is used for Windows installation.
3. `import-win-iso` creates new DataVolume with Windows ISO file with name defined in `isoDVName` parameter. The DataVolume needs to be in the same namespace as the PipelineRun (because the PVC is mounted into the PipelineRun pod).
4. `modify-windows-iso-file` extracts imported ISO file, replaces prompt bootloader (which is used as a default one when EFI is used) with no-prompt bootloader, pack the updated files back to new ISO, convert the ISO and replaces original ISO with updated one.Replacement of bootloader is needed to be able to automate installation of Windows versions which require EFI.
5. `create-vm` Task creates a VirtualMachine. A DataVolume with the Windows source ISO will be attached as CD-ROM and a second empty DataVolume will be used as installation destination. A third DataVolume with the virtio-win ISO will also be attached (Pipeline parameter `virtioContainerDiskName`). The VirtualMachine has to be created in the same namespace as the DataVolume with the ISO file. In case you would like to run the VirtualMachine in a different namespace, both Datavolumes have to be copied to the same namespace as the VirtualMachine.
6. `wait-for-vmi-status` Task waits until the VirtualMachine shuts down.
7. `create-datasource-root-disk` Task creates a DataSource object, which is used by UI for discovering bootable volumes and links PVC created in `create-vm-root-disk` step.
8. `cleanup-vm` deletes the installer VirtualMachine and all of its DataVolumes.
9. The output artifact will be the `baseDvName`/`baseDvNamespace` DataVolume with the basic Windows installation. It will boot into the Windows OOBE and needs to be setup further before it can be used.
10. `delete-imported-configmaps` deletes imported ConfigMaps.

## How to run

Pipeline uses ConfigMaps with `autounattend.xml` file for automated installation of Windows from ISO file. Example ConfigMaps are deployed within the Pipeline. In case you would like to use a different ConfigMap, specify a different URL in the `autounattendXMLConfigMapsURL` parameter and adjust `autounattendConfigMapName` parameter with the correct ConfigMap name. Examples of ConfigMaps can be found [here](https://github.com/kubevirt/kubevirt-tekton-tasks/tree/main/release/pipelines/windows-efi-installer/configmaps). Pipeline automatically removes ConfigMaps at the end of the PipelineRun. If the content available in `autounattendXMLConfigMapsURL` changes during the Pipeline run, Pipeline may not remove all the objects it created.

> [!IMPORTANT]
> Example PipelineRuns have special parameter acceptEula. By setting this parameter, you are agreeing to the applicable 
> Microsoft end user license agreement(s) for each deployment or installation for the Microsoft product(s). In case you 
> set it to false, the Pipeline will exit in first task.

> [!NOTE]
> By default, the Pipeline requires the ServiceAccount `pipeline` to exist. Tekton does not create this ServiceAccount 
> in namespaces which name starts with `openshift` or `kube`. In case you would like to run this Pipeline in a namespace which 
> starts with `openshift` or `kube`, you have to create the `pipeline` ServiceAccount manually or specify a different ServiceAccount in the PipelineRun.

Pipeline runs with resolvers:
{% for item in pipeline_runs_yaml %}
```yaml
{% if item.metadata.generateName == "windows11-installer-run-" %}
export WIN_IMAGE_DOWNLOAD_URL=$(./getisourl.py) # see paragraph Obtaining a download URL in an automated way

{% endif %}
oc create -f - <<EOF
{{ item | to_nice_yaml }}EOF
```
{% endfor %}

## Cancelling/Deleting PipelineRuns

When running the example Pipelines, they create temporary objects (DataVolumes, VirtualMachines, etc.). Each Pipeline has its own clean up system which should keep the cluster clean from leftovers. In case user hard deletes or cancels running PipelineRun, the PipelineRun will not clean temporary objects and objects will stay in the cluster. To prevent this behaviour, cancel the [PipelineRun gracefully](https://tekton.dev/docs/pipelines/pipelineruns/#gracefully-cancelling-a-pipelinerun). It triggers special Tasks, which remove temporary objects and keep only result DataSource/DataVolume/PVC.
Each object created by the Pipeline has OwnerReference to the Pod which created them (result DataSource and DataVolume do not have it). In the case that the clean up steps are not triggered, by deleting the PipelineRun all leftover objects created by the Pipeline will be deleted.

windows-efi-installer Pipeline generates for each PipelineRun new source DataVolume which contains imported ISO file. This DataVolume has generated name and is deleted after Pipeline succeeds. However, the created PVC will stay in cluster, but it will have terminating state. It will wait, until pipelinRun is deleted. This behaviour is caused by a fact, that PVC is mounted into modify-windows-iso TaskRun pod and PVC can be deleted only when the pod does not exist.

#### Obtaining a download URL in an automated way

The script [`getisourl.py`](https://github.com/kubevirt/kubevirt-tekton-tasks/blob/main/release/pipelines/windows-efi-installer/getisourl.py) can be used to automatically obtain a Windows 11 ISO download URL.

The prerequisites are:

- python3-selenium
- chromedriver
- chromium

Run it as follows to initialize a WIN_URL variable.

```bash
# Real URL can look differently
WIN_IMAGE_DOWNLOAD_URL=$(./getisourl.py)
```

### Common Errors
- The `modify-windows-iso-file` task can end with error `guestfish: access: /tmp/target-pvc/disk.img: Permissions denied`. This error means, the guestfish cannot access the Windows ISO file due to wrong permissions set on the file. This issue can be fixed by adding `taskRunSpecs` to the spec of the PipelineRun:
```
spec:
  taskRunSpecs:
    - pipelineTaskName: modify-windows-iso-file
      podTemplate:
        securityContext:
          fsGroup: 107
          runAsUser: 107
```
- Windows preferences in older KubeVirt versions might still use Bios mode. In that case, set the `useBiosMode` parameter 
to `true`. This will skip the `modify-windows-iso-file` task. In case the the Windows preference uses Bios and the 
`useBiosMode` parameter is not set to true, the Windows VM will not work.
