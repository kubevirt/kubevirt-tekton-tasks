# Testing

## Unit tests

```bash
# Run all module unit tests
make test

# Run tests for a single module
cd modules/execute-in-vm && make test

# Run tests with JUnit reports
make test-with-reports
```

## E2E tests

Require a running Kubernetes cluster with KubeVirt, Tekton, and CDI deployed.

By default, E2E scripts deploy pre-built images from `quay.io` (`make deploy`). To test with **locally built images** instead, set `DEV_MODE=true` (either export it or edit `automation/e2e-source.sh`). When enabled, the scripts use `make cluster-sync` which builds images locally and push them to the cluster registry.

### DEV_MODE Decision Criteria

The `DEV_MODE` environment variable controls whether E2E tests use locally built images or pre-built quay.io images.

**Use `DEV_MODE=true` (local images) when:**
- Testing code changes you just made
- Debugging image build or runtime issues
- Developing new tasks or features
- Working on a feature branch before pushing

**Use default (quay.io images) when:**
- Testing published releases
- Verifying production-like deployments
- Running tests on main branch without changes
- CI/CD pipelines (handled automatically)

**How it works:**
- `DEV_MODE=true`: Runs `make cluster-sync` (builds locally, pushes to cluster registry)
- Default: Runs `make deploy` (uses images from quay.io)
- Set in `automation/e2e-source.sh` or export before running tests

```bash
# Deploy dependencies (KubeVirt, Tekton, CDI)
./automation/e2e-deploy-resources.sh

# Run E2E tests with locally built images
DEV_MODE=true NUM_NODES=4 make e2e-tests

# Run E2E tests with pre-built quay.io images (default)
NUM_NODES=4 make e2e-tests

# Or deploy first, then test separately
make cluster-sync
NUM_NODES=4 make cluster-test

# Clean up
make cluster-clean
```

## Troubleshooting

- E2E tests fail: check if the Kubernetes cluster has KubeVirt, Tekton, and CDI (`./automation/e2e-deploy-resources.sh`).

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
