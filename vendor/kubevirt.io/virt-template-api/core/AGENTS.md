# AGENTS.md - api/core

Public API types for the `template.kubevirt.io` group (module `kubevirt.io/virt-template-api`).

- `v1alpha1/` or `v1beta1/` - `VirtualMachineTemplate` and `VirtualMachineTemplateRequest` CRD type definitions
- `subresourcesv1alpha1/` or `subresourcesv1beta1/` - subresource types (`ProcessOptions`, `CreateOptions`)

## Rules

- Treat any change to `VirtualMachineTemplate`/`VirtualMachineTemplateRequest` fields as a potentially breaking API change - call it out explicitly in the PR description
- Every exported field requires a kubebuilder validation marker (where applicable) and a godoc comment
- If you modify these type definitions, always run `make generate` (DeepCopy, OpenAPI, typed client) and `make manifests` (CRDs) afterward, and commit the regenerated output
- `VirtualMachineTemplateRequest` spec immutability is enforced by a CEL validation rule (`self == oldSelf`) declared on the API type itself - do not remove this marker
