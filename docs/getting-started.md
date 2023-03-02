# Getting started

## Supported Environments

Cluster oriented dev scripts are optimized for these environments:

- Minikube
    - example setup:
      ```bash
      minikube config set driver kvm2
      minikube config set memory 10240
      minikube start
      minikube addons enable registry
      ```
- OKD

## Build and Deployment

### Deploying Dependencies

To deploy `KubeVirt`, `Tekton Pipelines`, `Containerized Data Importer (CDI)` and `Common Templates` (OKD only):

```bash
./automation/e2e-deploy-resources.sh
```

### Building and Deploying KubeVirt Tekton Tasks
```bash
# generate yaml tasks from yaml templates (required only when changing the templates)
make generate-yaml-tasks
# builds all images (only podman supported ATM) and pushes them to the cluster registry.
# Then it deploys all tasks (which include the registry images) with required RBAC.
make cluster-sync
# to sync only one task
./scripts/cluster-sync.sh execute-in-vm
```

## Testing

### Running Unit Tests
```bash
# test all modules
make test
# test one module
cd modules/execute-in-vm
make test
```

### Running E2E Tests
```bash
NUM_NODES=4 make cluster-test

# clean all used resources, namespaces and images
make cluster-clean
```

### Deploying and Running E2E Tests
Alternatively deployment of resources and running the tests can be run simply with one command:

```bash
NUM_NODES=4 make e2e-tests
```

### Onboarding a New Task

Please see [Creating and Onboarding a New Task](onboarding-new-task.md) on how to create and onboard a new task.
