# Modify OpenShift template Task

This task modifies template.
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Service Account

This task should be run with `modify-vm-template-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **templateName**: Name of an OpenShift template.
- **templateNamespace**: Namespace of an source OpenShift template. (defaults to active namespace)
- **cpuSockets**: Number of CPU sockets
- **cpuCores**: Number of CPU cores
- **cpuThreads**: Number of CPU threads
- **memory**: Number of memory vm can use
- **templateLabels**: Template labels. If template contains same label, it will be replaced. Each param should have KEY:VAL format. Eg [`key:value`, `key:value`].
- **templateAnnotations**: Template Annotations. If template contains same annotation, it will be replaced. Each param should have KEY:VAL format. Eg [`key:value`, `key:value`]
- **vmLabels**: VM labels. If VM contains same label, it will be replaced. Each param should have KEY:VAL format. Eg [`key:value`, `key:value`].
- **vmAnnotations**: VM annotations. If VM contains same annotation, it will be replaced. Each param should have KEY:VAL format. Eg [`key:value`, `key:value`].
- **disks**: VM disks in json format, replace vm disk if same name, otherwise new disk is appended. Eg [{`name`: `test`, `cdrom`: {`bus`: `sata`}}, {`name`: `disk2`}]
- **volumes**: VM volumes in json format, replace vm volume if same name, otherwise new volume is appended. Eg [{`name`: `virtiocontainerdisk`, `containerDisk`: {`image`: `kubevirt/virtio-container-disk`}}]
- **datavolumeTemplates**: Datavolume templates in json format, replace datavolume if same name, otherwise new datavolume is appended. If deleteDatavolumeTemplate is set, first datavolumes are deleted and then datavolumes from this attribute are added. Eg [{`apiVersion`: `cdi.kubevirt.io/v1beta1`, `kind`: `DataVolume`, `metadata`:{`name`: `test1`}, `spec`: {`source`: {`http`: {`url`: `test.somenonexisting`}}}}]
- **deleteDatavolumeTemplate**: Set to `true` or `false` if task should delete datavolume template in template and all associated volumes and disks.

### Results

- **name**: The name of a template that was updated.
- **namespace**: The namespace of a template that was updated.

### Usage

Please see [examples](examples) on how to do a copy template from a template.
