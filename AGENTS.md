# KubeVirt Tekton Tasks - Project Guide

KubeVirt-specific Tekton tasks for CI/CD pipelines on Kubernetes. Provides tasks for creating/managing VMs, executing commands over SSH, manipulating disk images with libguestfs, and uploading disk images to container registries.

**IMPORTANT for AI Agents**: See [AI Workflow Rules](docs/ai-workflow.md) for the required 2-step planning and approval process before modifying any code.

## Jira Task Management

When the user references a Jira task (e.g., "read CNV-82681" or "check issue CNV-12345"):

1. **Use acli to fetch task information**:
   ```bash
   acli jira workitem view <TASK-ID> --fields '*all'
   ```

2. **Examples**:
   ```bash
   # View complete task information
   acli jira workitem view CNV-82681 --fields '*all'

   # View specific fields only
   acli jira workitem view CNV-82681 --fields summary,description,status,assignee
   ```

## Documentation

| Document | Description |
|----------|-------------|
| [AI Workflow Rules](docs/ai-workflow.md) | **Required 2-step process for AI agents** (planning -> approval -> implementation) |
| [Build & Dependencies](docs/build.md) | Build system, Go/container toolchain, supported environments |
| [YAML Generation](docs/yaml-generation.md) | Ansible-driven task YAML generation pipeline (**never hand-edit `release/`**) |
| [Code Conventions](docs/code-conventions.md) | Go module structure, testing patterns, error handling, formatting |
| [Testing](docs/testing.md) | Unit tests, E2E tests, `DEV_MODE` for local images |
| [CI / GitHub Actions](docs/ci.md) | Workflow descriptions, Dependabot, Renovate |
| [Release](docs/release.md) | Release process, version bumping, artifacts |
| [Contributing](docs/contributing.md) | PR requirements, code quality workflow, adding new tasks |
| [Tasks RBAC Permissions](docs/tasks-rbac-permissions.md) | ServiceAccount, ClusterRole, RoleBinding setup |

## Project Structure

```
cmd/                        # One main.go per compiled task binary
  create-vm/
  execute-in-vm/
  generate-ssh-keys/
  modify-data-object/
  wait-for-vmi-status/
  disk-uploader/
  disk-virt-customize/
  disk-virt-sysprep/
modules/                    # Go libraries and per-task packages + tests
  shared/pkg/...            # Shared utilities (env, log, output, zconstants, zerrors, zutils)
  sharedtest/...            # Shared test helpers
  <task>/pkg/...            # Per-task logic
  disk-virt/                # Shared by disk-virt-customize and disk-virt-sysprep
build/                      # Containerfiles and entrypoint scripts
  Containerfile             # Multi-stage: CentOS Stream 10 builder + runtime
  Containerfile.DiskVirt    # Builds disk-virt tasks on libguestfs-tools base
configs/                    # Per-task Ansible variable files
templates/                  # Ansible-driven Tekton task YAML generation
  <task>/manifests/         # Generated task manifests
  <task>/generate-task.yaml # Ansible playbook for YAML generation
templates-pipelines/        # Pipeline YAML generation (Windows pipelines)
release/                    # Published task and pipeline YAML artifacts
  tasks/<task>/
  pipelines/<pipeline>/
scripts/                    # Shell automation for build, test, lint, deploy
  makefile-snippets/        # Included makefile fragments
test/                       # Ginkgo E2E integration tests and framework
automation/                 # E2E orchestration, resource deployment, CI helpers
vendor/                     # Vendored Go modules (GOFLAGS=-mod=vendor)
docs/                       # Developer and user documentation
```

### Key relationships

- `cleanup-vm` is generated from `execute-in-vm` templates (with `is_cleanup: true`), not a separate binary.
- `modify-windows-iso-file` exists as templates and release YAML only (no `cmd/` entry).
- `disk-virt-customize` and `disk-virt-sysprep` share `modules/disk-virt/`.

## Common AI Agent Mistakes & Quick Fixes

| Mistake | Why It Happens | How to Fix |
|---------|---------------|------------|
| Hand-editing `release/` files | Not reading yaml-generation.md | Always edit templates, then run `make generate-yaml-tasks` |
| Missing vendored dependencies | After adding/updating Go modules | Run `make vendor` after any `go.mod` change |
| Skipping the 2-step workflow | Eagerness to implement | Always create `change.md` first, wait for approval |
| Forgetting DCO sign-off | Not following git workflow | Use `git commit --signoff` |
| Using testify instead of Ginkgo | Following common Go patterns | Check code-conventions.md - we use Ginkgo/Gomega |
| Bumping Go version on release branch | Assuming newer is better | Never bump Go on release branches (breaks midstream) |

## Common Commands

| Command | Description |
|---------|-------------|
| `make clean` | Clean build artifacts |
| `make generate-yaml-tasks` | Generate task YAML from templates |
| `make generate-pipelines` | Generate pipeline YAML |
| `make test` | Run unit tests |
| `make lint` / `make lint-fix` | Check / fix formatting |
| `make test-yaml-consistency` | Verify YAML consistency |
| `make cluster-sync` | Build, push, deploy to cluster |
| `make cluster-test` | Run E2E tests on cluster |
| `make cluster-clean` | Clean cluster resources |
| `make e2e-tests` | Full E2E: deploy + test |
| `make vendor` | Vendor Go dependencies |
| `make release` | Full release pipeline |
| `make deploy` / `make undeploy` | Deploy/remove tasks |
