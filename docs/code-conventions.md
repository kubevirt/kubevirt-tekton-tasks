# Go Code Conventions

## Module structure

Each task module under `modules/<task>/` contains:
- `pkg/` - Core task logic, organized by domain
- Unit tests colocated with source files (`*_test.go`)
- A `makefile` including shared snippets from `scripts/makefile-snippets/`

The shared module at `modules/shared/pkg/` provides:
- `env/` - Environment variable parsing
- `log/` - Zap logger configuration
- `output/` - Tekton result output utilities
- `zconstants/` - Shared constants
- `zerrors/` - Error utilities
- `zutils/` - General-purpose helpers

## Testing patterns

- Use **Ginkgo v2** (`github.com/onsi/ginkgo/v2`) and **Gomega** (`github.com/onsi/gomega`).
- Use `go.uber.org/mock` for generating mocks.
- **Never use `testify`** - this project standardizes on Ginkgo/Gomega.
- **Prefer `DescribeTable`** for parameterized tests - AI tools often avoid it, but it's the recommended pattern for testing multiple similar cases.

```go
var _ = Describe("MyComponent", func() {
    It("should handle valid input", func() {
        result, err := MyFunc(validInput)
        Expect(err).ToNot(HaveOccurred())
        Expect(result).To(Equal(expected))
    })

    DescribeTable("should handle various inputs", func(input string, expected string) {
        result, err := MyFunc(input)
        Expect(err).ToNot(HaveOccurred())
        Expect(result).To(Equal(expected))
    },
        Entry("case 1", "input1", "expected1"),
        Entry("case 2", "input2", "expected2"),
    )
})
```

## CLI argument parsing

- Use `github.com/alexflint/go-arg` for parsing task parameters from environment variables.
- Task parameters are set as env vars in the Tekton task YAML and parsed in `main.go`.

## Error handling

- Always return and check errors; never silently ignore.
- Use structured logging via `go.uber.org/zap`.
- Wrap errors with context when propagating.

## Formatting & imports

- `gofmt` is the only enforced formatter (no `.golangci.yml` at repo root).
- Run `make lint-fix` to auto-format before committing.
- Import order: standard library, external packages, internal packages.
- Use vendored dependencies (`-mod=vendor`).

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
