# CI / GitHub Actions

## Workflows

- **`test-yaml-consistency.yaml`**: PRs to `main` - runs `scripts/test-yaml-consistency.sh`.
- **`validate-no-offensive-lang.yml`**: PRs to `main` - language validation.
- **`release.yaml`**: On release published - builds multi-arch images, pushes to Quay, uploads manifest asset.

## Dependency management

- **Dependabot**: Watches Go modules and GitHub Actions (ignores Ginkgo/Gomega).
- **Renovate**: Runs on a schedule via `.github/workflows/renovate.yml` (GitHub App token; requires
  `RENOVATE_APP_ID`/`RENOVATE_APP_PRIVATE_KEY` repo secrets). Manages `gomod` updates on `main` and
  `release-v0.15`+ branches, grouping patch/minor together and majors separately, and runs
  `make vendor` and `make test` after each update. Excludes packages that typically need code changes
  (`k8s.io`, `kubevirt.io`, `sigs.k8s.io`, `openshift`, `knative.dev`, `tektoncd`, Ginkgo/Gomega).
  Also handles vulnerability/OSV alerts; ignores `vendor/`. Go module bumps may come from either
  Dependabot or Renovate.

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
