# AOF CLI Repository Instructions

AOF CLI is the official repository bootstrap utility for adopting the AI Orchestration Framework (AOF).

Before changing product behavior, read `docs/AOF-CLI-v0.1-PRD-and-Technical-Design.md`.

Canonical AOF upstream: https://github.com/aof-framework/aof

## Product boundaries

- AOF CLI is a bootstrap/adoption tool, not an AOF runtime.
- Do not redefine AOF semantics.
- Do not make network access mandatory for `aof init`.
- Do not clone mutable `main` during initialization.
- Do not silently overwrite user content.
- Do not introduce false conformance claims.
- Do not add telemetry, LLM calls, cloud dependencies, or application runtime dependencies.

## Engineering quality

Run before completing changes:

```bash
go fmt ./...
go vet ./...
go test ./...
```
