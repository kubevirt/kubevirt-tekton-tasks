# Disk Virt Sysprep Task

This task uses [virt-sysprep](https://libguestfs.org/virt-sysprep.1.html) to run a sysprep script on a target pvc.

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Usage

Task run using resolver:
```
{{ task_run_resolver_yaml | to_nice_yaml }}```

#### Common Errors

- The input PVC disk should not be accessed by a running VM or other tools like virt-sysprep task concurrently.
The task will fail with a generic `...guestfs_launch failed...` message.
A verbose parameter can be set to true for more information.

- The task can end with error `Permissions denied`. This error means, the disk-virt-sysprep cannot access the VM image due to wrong permissions set on the file. This issue can be fixed by adding `podTemplate` to the spec of the TaskRun:
```
spec:
  podTemplate:
    securityContext:
      fsGroup: 107
      runAsUser: 107
```

To be able to set the SecurityContext with fsGroup you will probably need custom SCC. Please update variable `${TASK_NAMESPACE}` with namespace where the task is deployed
```
oc apply -f - <<EOF
apiVersion: security.openshift.io/v1
kind: SecurityContextConstraints
metadata:
  name: tekton-tasks-scc
allowPrivilegedContainer: false
allowHostDirVolumePlugin: false
allowHostIPC: false
allowHostNetwork: false
allowHostPID: false
allowHostPorts: false
seccompProfiles:
  - "runtime/default"
readOnlyRootFilesystem: false
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: MustRunAs
fsGroup:
  type: MustRunAs
  ranges:
    - min: 107
      max: 107
users:
- system:serviceaccount:${TASK_NAMESPACE}:pipeline
EOF
```
Then you have to add annotation `"openshift.io/required-scc": "tekton-tasks-scc"` to your TaskRun.

### OS support

- Linux: full; all the sysprep commands work
- Windows: partial; only some sysprep commands work
