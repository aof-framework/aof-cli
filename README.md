# AOF CLI

AOF CLI is the official repository bootstrap utility for adopting the [AI Orchestration Framework (AOF)](https://github.com/aof-framework/aof).

It configures project context and developer-facing adoption artifacts. It is **not** an AOF runtime and is not required by applications after initialization.

## MVP commands

```bash
aof init
aof version
aof help
```

## Build

```bash
go build -o aof ./cmd/aof
```

## Example

```bash
mkdir itsm-backend
cd itsm-backend
/path/to/aof init
```

For automation:

```bash
aof init \
  --name itsm-backend \
  --language go \
  --type backend \
  --domain itsm \
  --profile core \
  --ai \
  --non-interactive
```

## Canonical AOF

- Repository: https://github.com/aof-framework/aof
- Specification: https://github.com/aof-framework/aof/blob/main/specification/AOF-v1.0-Framework-Specification.md

`aof init` is offline-capable and does not clone the mutable upstream `main` branch. Static bootstrap guidance is embedded in the binary and records canonical upstream provenance in generated project metadata.
