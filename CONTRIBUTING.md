# Contributing to AOF CLI

Thank you for your interest in contributing to the **AI Orchestration Framework CLI (`aof-cli`)**!

## Code of Conduct

All contributors and maintainers are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Architectural Guidelines & Invariants

AOF CLI is a repository bootstrap utility for adopting AOF, not an AOF execution runtime. When contributing, keep the following invariants in mind:

1. **Do not redefine AOF semantics:** AOF CLI adheres to the canonical [AOF Specification](https://github.com/aof-framework/aof).
2. **Offline-capable:** `aof init` MUST NOT require network connectivity or clone mutable `main` at runtime.
3. **Safe file operations:** Never silently overwrite existing user files without conflict validation.
4. **No unauthorized external dependencies:** Do not add telemetry, LLM API calls, cloud dependencies, or heavy runtime frameworks.

## Development Workflow

### Prerequisites

- [Go](https://go.dev/) 1.23 or newer
- `git`

### Building and Testing

Run the standard Go development checks before submitting any changes:

```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Run all tests
go test -v -count=1 ./...
```

### Submitting a Pull Request

1. Fork the repository and create a descriptive branch name (e.g. `feat/new-language-detection` or `fix/safe-path-windows`).
2. Implement your changes along with appropriate unit tests.
3. Ensure all tests pass and code is formatted cleanly.
4. Open a Pull Request referencing any related issues.
5. Reviewers will examine the PR for adherence to AOF semantic boundaries and engineering quality.
