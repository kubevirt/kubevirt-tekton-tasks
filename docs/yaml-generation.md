# YAML Generation Pipeline (CRITICAL)

Task YAML is **generated**, not hand-edited:

1. Define task variables in `configs/<task-name>.yaml`
2. Create/update Ansible templates in `templates/<task>/`
3. Run `make generate-yaml-tasks` to regenerate `release/tasks/<task>/`
4. Run `make generate-pipelines` for pipeline YAML in `release/pipelines/`

**NEVER hand-edit files under `release/tasks/` or `release/pipelines/` directly.** Always modify the source templates and regenerate. CI runs `make test-yaml-consistency` on every PR to enforce this.

## Task template structure

Each task under `templates/<task>/` typically contains:
- `generate-task.yaml` - Ansible playbook for generation
- `manifests/` - Jinja2 templates for the Tekton Task YAML
- `readmes/` - README templates
- `examples/` - Example TaskRun YAML

Task YAML references container images via variables set in `scripts/common.sh`. In CI, these are overridden to point to CI-built images.

## Troubleshooting

- YAML consistency fails: run `make generate-yaml-tasks` and/or `make generate-pipelines`.

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
