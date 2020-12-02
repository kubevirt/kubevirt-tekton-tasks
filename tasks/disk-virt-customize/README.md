# Virt Customize Task

This task uses `virt-customize` to run a customize script on a target pvc.


### Parameters

- **pvc**: PersistentVolumeClaim to run the the virt-customize script in. PVC should be in the same namespace as taskrun/pipelinerun.
- **script**: virt-customize script in `--commands-from-file` format.


### Usage

Please see [examples](examples)
