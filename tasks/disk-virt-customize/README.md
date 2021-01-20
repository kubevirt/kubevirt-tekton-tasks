# Virt Customize Task

This task uses `virt-customize` to run a customize script on a target pvc.


### Parameters

- **pvc**: PersistentVolumeClaim to run the the virt-customize script in. PVC should be in the same namespace as taskrun/pipelinerun.
- **customizeCommands**: virt-customize commands in `--commands-from-file` format.
- **verbose**: Enable verbose mode and tracing of libguestfs API calls.
- **additionalOptions**: Additional options to pass to virt-customize.


### Usage

Please see [examples](examples)

#### Common Errors

- The input PVC disk should not be accessed by a running VM or other tools like virt-customize task concurrently.
The task will fail with a generic `...guestfs_launch failed...` message.
Verbose parameter can be set to true for more information.

### OS support

- Linux: full; all the customize commands work
- Windows: partial; only some customize commands work
