# CI / GitHub Actions

## Workflows

- **`test-yaml-consistency.yaml`**: PRs to `main` - runs `scripts/test-yaml-consistency.sh`.
- **`validate-no-offensive-lang.yml`**: PRs to `main` - language validation.
- **`release.yaml`**: On release published - builds multi-arch images, pushes to Quay, uploads manifest asset.

## Dependency management

- **Dependabot**: Watches Go modules and GitHub Actions (ignores Ginkgo/Gomega). Runs daily on
  `main` only, and opens PRs for all `gomod` updates, routine or security.
- **Renovate**: Runs daily at 06:00 UTC or manually from the Actions tab. The workflow uses a
  GitHub App token and requires the `RENOVATE_APP_ID` and `RENOVATE_APP_PRIVATE_KEY` repository
  secrets. It checks `main` and `release-v0.15`+ for security Go updates, including indirect
  dependencies, and runs `make vendor` after updates. It ignores `vendor/` and
  packages that typically need code changes (`k8s.io`, `kubevirt.io`, `sigs.k8s.io`, `openshift`,
  `knative.dev`, `tektoncd`, Ginkgo/Gomega). Renovate only opens vulnerability PRs, including
  OSV alerts.
  A `pull_request_target` workflow comments on same-repository Renovate PRs when a `go` or
  `toolchain` directive changes. It replaces the comment after each push and removes it if the
  directive change disappears. Review these changes against the midstream build toolchain before
  merging; Go version bumps on release branches break midstream builds. Dependabot can create
  duplicate security PRs on `main`; release branches rely on Renovate.

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
