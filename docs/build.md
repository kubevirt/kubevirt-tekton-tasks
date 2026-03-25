# Build & Dependencies

## Build system

- Root build: `makefile` (lowercase), not `Makefile`.
- **All Go builds use vendored deps**: `GOFLAGS=-mod=vendor`.
- After any dependency change: **ALWAYS** run `make vendor` (runs `go mod tidy` + `go mod vendor`).
- Go version: defined in `go.mod` (check the `go` directive). Pinned Kubernetes/OpenShift/CDI versions via `replace` directives in `go.mod`.
- **Do NOT bump the Go version without confirmation.** Before upgrading, verify which Go version the midstream is using - the upstream version must stay compatible. On **release branches**, never bump the Go version as it would break midstream builds.

## Build commands

```bash
# Generate task YAML from Ansible templates (only when templates change)
make generate-yaml-tasks

# Generate pipeline YAML
make generate-pipelines

# Build and push images to cluster registry, deploy tasks
make cluster-sync

# Sync a single task
./scripts/cluster-sync.sh execute-in-vm

# Build multi-arch release images (podman, linux/amd64+arm64+s390x)
make build-release-images

# Full release: generate YAML + build + push images
make release
```

## Container images

- Built with **Podman** (not Docker). Multi-arch: `linux/amd64`, `linux/arm64`, `linux/s390x`.
- `build/Containerfile`: CentOS Stream builder + runtime with xorriso, ssh, nbdkit, qemu-img. Check the file for current Go and base image versions.
- `build/Containerfile.DiskVirt`: Builds disk-virt tasks on a libguestfs-tools base. Check the file for the current base image tag.
- Images pushed to `quay.io/kubevirt/tekton-tasks` and `quay.io/kubevirt/tekton-tasks-disk-virt`.

## Supported environments

- **OKD** (OpenShift Kubernetes Distribution)
- **OCP** (OpenShift Container Platform)

## Troubleshooting

- Build fails with module errors: run `make vendor` first.
- Cluster sync issues: verify cluster registry is accessible and permissions are set.

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
