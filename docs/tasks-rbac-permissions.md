# RBAC permissions for running the tasks

Each KubeVirt Tekton Task usually comes with its own `ServiceAccount`, `ClusterRole` and `RoleBinding`.
This allows the task to monitor, create, update and delete relevant resources in the namespace where it was deployed.

## Using the RBAC permissions

When running the `TaskRun`/`Pipelinerun`, correct service accounts should be used for each task.
Please see Tekton documentation on how to specify service accounts in
[TaskRuns](https://github.com/tektoncd/pipeline/blob/master/docs/taskruns.md#configuring-a-taskrun)
and [PipelineRuns](https://github.com/tektoncd/pipeline/blob/master/docs/pipelineruns.md#mapping-serviceaccount-credentials-to-tasks).

This will make sure that the task's pod has right permissions for accessing the resources it needs.
For example `execute-in-vm` task has RBAC permissions for fetching `VirtualMachines` and for executing `start`/`stop` operations on them.


### Example service account mapping

The service account is usually named `${TASK_NAME}-task`.

```yaml
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  name: my-run
spec:
  pipelineRef:
    name: my-pipeline
  serviceAccountNames:
    - serviceAccountName: create-vm-from-manifest-task
      taskName: create-vm-from-manifest
    - serviceAccountName: execute-in-vm-task
      taskName: execute-in-vm

```

## Enhanced RBAC permissions

The default deployment of KubeVirt Tekton Tasks is initialized with permissions (`RoleBinding`) for only a specific namespace.
By default, it is not possible for example to create a pipeline in namespace `task-ns1` that would fetch/create resources (eg `VirtualMachines`) in namespace `vm-ns1`.
To support such use cases additional bindings need to be created.

### Multi Namespace RBAC permissions

Let's assume the tasks and thus their service accounts are deployed in `task-ns1` namespace.
To allow access to `vm-ns1` namespace, the following `RoleBinding` should be created for each task.
Only two tasks are used for simplicity in the following example.

```bash
#!/usr/bin/env bash
for TASK_NAME in generate-ssh-keys create-vm-from-manifest; do
    kubectl apply -f - << EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ${TASK_NAME}-task
  namespace: vm-ns1
roleRef:
  kind: ClusterRole
  name: ${TASK_NAME}-task
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name:  ${TASK_NAME}-task
    namespace: task-ns1
EOF
done
```

This example enables `generate-ssh-keys` and `create-vm-from-manifest` tasks to create secrets/VMs in `vm-ns1` from a pipeline started in `task-ns1`.

Warning: this will allow users with access to `generate-ssh-keys-task` and `create-vm-from-manifest-task` service accounts to run pods which can access `vm-ns1`.


### Cluster RBAC permissions

Let's assume the tasks and thus their service accounts are deployed in `task-ns1` namespace.
To allow access to any other namespace, the following `ClusterRoleBinding` should be created for each task.
Only two tasks are used for simplicity in the following example.

```bash
#!/usr/bin/env bash
for TASK_NAME in generate-ssh-keys create-vm-from-manifest; do
    kubectl apply -f - << EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ${TASK_NAME}-task
roleRef:
  kind: ClusterRole
  name: ${TASK_NAME}-task
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name:  ${TASK_NAME}-task
    namespace: task-ns1
EOF
done
```

This example enables `generate-ssh-keys` and `create-vm-from-manifest` tasks to create secrets/VMs in any namespace in the cluster from a pipeline started in `task-ns1`.

Warning: this will allow users with access to `generate-ssh-keys-task` and `create-vm-from-manifest-task` service accounts to run pods in any namespace.
This will elevate their privileges to cluster admin level.


## Deploying the tasks in additional namespaces

`ServiceAccount` and `RoleBinding` should be created for each new namespace where Tasks and Pipelines will be run. Let's deploy the tasks in a new namespace called `task-ns2`.
Only two tasks are used for simplicity in the following example.

```bash
#!/usr/bin/env bash
for TASK_NAME in generate-ssh-keys create-vm-from-manifest; do
    kubectl apply -f - << EOF
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ${TASK_NAME}-task
  namespace: task-ns2

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ${TASK_NAME}-task
  namespace: task-ns2
roleRef:
  kind: ClusterRole
  name: ${TASK_NAME}-task
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: ${TASK_NAME}-task
EOF
done
```

TIP: it might be easier to `kubectl apply` the original task yaml in a new namespace.
