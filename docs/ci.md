# CI / GitHub Actions

## Workflows

- **`test-yaml-consistency.yaml`**: PRs to `main` - runs `scripts/test-yaml-consistency.sh`.
- **`validate-no-offensive-lang.yml`**: PRs to `main` - language validation.
- **`release.yaml`**: On release published - builds multi-arch images, pushes to Quay, uploads manifest asset.

## Dependency management

- **Dependabot**: Watches Go modules and GitHub Actions (ignores Ginkgo/Gomega).
- **Renovate**: Vulnerability/OSV alerts enabled; ignores `vendor/`.

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
