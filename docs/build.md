# Build & Dependencies

## Build system

- Root build: `makefile` (lowercase), not `Makefile`.
- **All Go builds use vendored deps**: `GOFLAGS=-mod=vendor`.
- After any dependency change: **ALWAYS** run `make vendor` (runs `go mod tidy` + `go mod vendor`).
- Go version: defined in `go.mod` (check the `go` directive). Pinned Kubernetes/OpenShift/CDI versions via `replace` directives in `go.mod`.
- **Do NOT bump the Go version without confirmation.** Before upgrading, verify which Go version the midstream is using - the upstream version must stay compatible. On **release branches**, never bump the Go version as it would break midstream builds.

## Python and Ansible requirements

- **Python 3.9 or newer** is the compatibility baseline for repository Python scripts. Python 2 and Python 3.0-3.8 are unsupported. For example, `scripts/check-go-version-bump.py` uses [`str.removeprefix()`](https://docs.python.org/3.9/library/stdtypes.html#str.removeprefix), introduced in Python 3.9. Keep Python changes compatible with this baseline; do not introduce newer syntax or standard-library APIs without an agreed baseline update.
- Invoke helpers explicitly with `python3`, rather than relying on `python` or a script's shebang to select the interpreter. For example:

  ```bash
  python3 scripts/check-go-version-bump.py BASE_SHA HEAD_SHA
  python3 release/pipelines/windows-efi-installer/getisourl.py
  ```

- `getisourl.py` additionally requires Selenium installed for the selected Python interpreter, Chromium, and a compatible ChromeDriver. See the [Windows EFI installer prerequisites](../release/pipelines/windows-efi-installer/README.md#obtaining-a-download-url-in-an-automated-way).
- There is no Python version test matrix. The baseline describes script compatibility, not a guarantee that every Python/dependency combination works. Third-party dependencies such as Selenium and Ansible can require a newer Python version; use a Python version supported by the dependencies you install.
- **Ansible is required for YAML generation** (`make generate-yaml-tasks` and `make generate-pipelines`). The playbooks require **ansible-core 2.13 or newer** because they use the [`ansible.builtin.env` lookup's `default` option](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/env_lookup.html#keyword-parameters), added in 2.13. Older versions are unsupported for generation. This is a feature minimum, not a recommendation to install an obsolete release: choose a maintained release and a controller Python version from its [support matrix](https://docs.ansible.com/ansible/latest/reference_appendices/release_and_maintenance.html#ansible-core-support-matrix). Check the installed pair with `ansible-playbook --version`.
- The [generation CI image](../automation/generatetasks/build/Dockerfile) installs Fedora's `ansible` package without pinning an exact version. The repository does not define a broader Ansible version test matrix.

### Contributor hosts

The existing build and generation workflows use **Linux with Bash**. Native Windows and macOS execution of these workflows, including Python helpers such as `getisourl.py`, is not validated by repository CI. Use a Linux environment to follow the documented workflow. This host expectation is separate from the OKD/OCP target clusters below and from Windows guests managed by the tasks.

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
