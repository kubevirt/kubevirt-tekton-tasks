# Modify Windows ISO file

This tasks is modifying windows iso file. It replaces prompt bootloader with non prompt one. This helps with automation of 
windows which requires EFI - the prompt bootloader will not continue with installation until some key is pressed. The non prompt 
bootloader will not require any key pres.

### Parameters

- **pvcName**: PersistentVolumeClaim which contains windows iso.


### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: modify-windows-iso-file-taskrun-resolver-
spec:
    params:
    -   name: pvcName
        value: w11
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: modify-windows-iso-file
        -   name: version
            value: v0.21.0
        resolver: hub
```

### Platforms

The Task can be run on linux/amd64 platform.
