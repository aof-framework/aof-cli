# AOF CLI

AOF CLI is the repository bootstrap utility for adopting the [AI Orchestration Framework (AOF)](https://github.com/aof-framework/aof).

**Current version:** v0.3.0 — Semantic Fidelity & Adoption Hardening  
**AOF baseline:** v1.0 LTS

AOF CLI configures project-aware adoption context. It is **not** an AOF runtime, does not prove implementation, and does not create an AOF conformance claim.

## Commands

```bash
aof init
aof version
aof help
```

## What changed in v0.3.0

v0.3.0 hardens the semantic projection from AOF v1.0 LTS into generated repositories:

- separates `Applicability`, `NormativeLevel`, `AdoptionState`, `ImplementationState`, and `VerificationState`;
- separates target profile from claimed conformance;
- supports AOF Agent types `LLM`, `Deterministic`, `Human`, `Hybrid`, and `ExternalService`;
- implements literal AOF-Core, AOF-Governed, and AOF-Assured strengthening rules;
- supports `AOF-Secure-SDLC` as a domain profile;
- supports `AOF-High-Assurance` as a strengthening overlay;
- generates requirement-ID traceability;
- generates Four-Plane architecture, Trust Boundaries, Effect Boundaries, and Safety Kernel models;
- expands Authority lifecycle, Policy conflict resolution, dynamic Risk, State/Trace, failure/recovery, Evidence/Verification, Human Governance, and accountability models;
- embeds the pinned canonical AOF v1.0 LTS JSON Schema bundle;
- distinguishes `.aof/bootstrap-manifest.json` from canonical `ConformanceManifest`;
- preserves an existing `AGENTS.md` and generates `AGENTS.aof.md` as an explicit merge sidecar;
- records AOF specification/schema checksums for provenance.

## Build

```bash
go build -o aof ./cmd/aof
```

Quality gates:

```bash
go fmt ./...
go vet ./...
go test ./...
```

## Example

```bash
aof init \
  --name payments-platform \
  --language go \
  --type backend \
  --profile assured \
  --secure-sdlc \
  --high-assurance \
  --agent-types LLM,Human,Deterministic \
  --ai \
  --ai-consequential \
  --effect-types database,external-api \
  --criticality critical \
  --data-sensitivity restricted \
  --non-interactive
```

Profile composition rules:

```text
AOF-Core ⊆ AOF-Governed ⊆ AOF-Assured
AOF-Secure-SDLC requires Governed or Assured base
AOF-High-Assurance requires Assured base
```

## Generated state is not conformance

After `aof init`:

```text
planned != implemented
not_assessed != pass
not_evaluated != pass
target_profile != claimed_profile
BootstrapManifest != ConformanceManifest
```

Generated project context includes `.aof/requirements.yaml`, `aof/control/safety-kernel.yaml`, `aof/architecture/`, governance models, assurance models, and canonical schemas under `aof/schemas/`.

## Canonical AOF

- Repository: https://github.com/aof-framework/aof
- Specification: https://github.com/aof-framework/aof/blob/main/specification/AOF-v1.0-Framework-Specification.md

`aof init` is offline-capable and does not clone mutable upstream `main`.
