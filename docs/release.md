# Release

Current version is defined in `scripts/makefile-snippets/makefile-release.mk` (see the `RELEASE_VERSION` variable).

## Release process

1. **Bump version** in `scripts/makefile-snippets/makefile-release.mk` (`RELEASE_VERSION`).
2. **Regenerate manifests**: `make generate-yaml-tasks` and `make generate-pipelines`.
3. **Commit** all changed files (updated version, regenerated YAML under `release/`).
4. **Create a PR**, get it reviewed and merged to `main`.
5. **Create a new GitHub release** from `main` with the new version tag. The `release.yaml` workflow will automatically build and push multi-arch container images to `quay.io` and upload the manifest asset.

## Release artifacts

- Container images on `quay.io/kubevirt/`
- `kubevirt-tekton-tasks.yaml` manifest (uploaded as GitHub release asset)
- Per-task YAML under `release/tasks/`
- Pipeline YAML under `release/pipelines/`

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
