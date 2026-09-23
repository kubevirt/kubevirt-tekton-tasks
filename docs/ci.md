# CI / GitHub Actions

## Workflows

- **`test-yaml-consistency.yaml`**: PRs to `main` - runs `scripts/test-yaml-consistency.sh`.
- **`validate-no-offensive-lang.yml`**: PRs to `main` - language validation.
- **`release.yaml`**: On release published - builds multi-arch images, pushes to Quay, uploads manifest asset.

## Dependency management

- **Dependabot**: Watches Go modules and GitHub Actions (ignores Ginkgo/Gomega). Runs daily on
  `main` only, and opens PRs for all `gomod` updates, routine or security.
- **Renovate**: Runs on a schedule via `.github/workflows/renovate.yml` (GitHub App token; requires
  `RENOVATE_APP_ID`/`RENOVATE_APP_PRIVATE_KEY` repo secrets). Covers `main` and `release-v0.15`+
  branches, but only raises PRs for security fixes (including indirect deps), and runs
  `make vendor` and `make test` after each update. Excludes packages that typically need code
  changes (`k8s.io`, `kubevirt.io`, `sigs.k8s.io`, `openshift`, `knative.dev`, `tektoncd`,
  Ginkgo/Gomega). Also handles vulnerability/OSV alerts; ignores `vendor/`. Because Dependabot
  also covers security fixes on `main`, expect the occasional duplicate PR there; release
  branches only get updates from Renovate.

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
