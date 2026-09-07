# AOF CLI

AOF CLI is the repository bootstrap utility for adopting the [AI Orchestration Framework (AOF)](https://github.com/aof-framework/aof).

**Current version:** v0.3.2 — Exact Canonical Source Alignment
**AOF baseline:** v1.0 LTS

AOF CLI configures project-aware adoption context. It is **not** an AOF runtime, does not prove implementation, and does not create an AOF conformance claim.

## Commands

```bash
aof init
aof version
aof help
```

## What changed in v0.3.2

v0.3.2 makes the AOF specification repository the only source of framework data:

- embeds the active release manifest, normative specification, requirement registry, traceability matrix, profile definitions, schema index, schema catalog, and active schema checksums byte-for-byte;
- validates every embedded source artifact by SHA-256 before use;
- preserves all upstream values, including the 17 `Unclassified` architecture requirements, without CLI normalization;
- derives invariant and traceability indexes without modifying canonical source data;
- derives canonical-object names from the upstream schema index rather than a CLI-owned list;
- labels applicability and adoption decisions as non-canonical CLI process metadata;
- represents absent upstream mappings as `source_mapping_absent` and unevaluated requirements as `not_evaluated`;
- synchronizes generated JSON Schemas exactly with the active AOF v1.0 LTS checkout.

## What changed in v0.3.1

v0.3.1 establishes the complete machine-readable AOF v1.0 LTS semantic baseline:

- registers all 162 canonical invariants and 332 canonical requirements;
- hard-gates the exact 494-ID semantic universe and 352 canonical traceability edges;
- derives requirement projection and normative metadata from the canonical registry;
- retains every requirement as `applicable` or explicitly `conditional`;
- records upstream mapping absences without inventing semantic edges;
- pins the verified upstream source commit and separates published/specification freeze checksums.

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

## Canonical semantic coverage

The CLI embeds and validates the complete AOF v1.0 LTS semantic universe before projecting project-specific adoption state:

```text
162 / 162 canonical invariants
332 / 332 canonical requirements
494 / 494 stable semantic IDs
352 canonical traceability edges
```

Build tests fail when any embedded canonical artifact differs from its pinned upstream SHA-256, or on a missing, duplicate, unknown, or malformed semantic ID, broken traceability edge, or incomplete upstream profile/schema index.

`aof init` remains offline. It emits all 332 requirement IDs with their exact canonical statements and normative classifications. An unselected requirement is `not_evaluated`, never silently converted to `Conditional`; missing upstream mappings remain `source_mapping_absent` and are never filled using heuristic similarity.

Exact canonical sources live under `internal/canonical/source/`. Maintainers reproduce them without network access using:

```powershell
go run ./tools/canonicalgen -upstream-root "D:\AI Orchestration Framework\AOF-v1.0-LTS"
```

The importer only reads the upstream checkout. The local absolute path is not embedded in the binary or generated projects.

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
- Pinned source commit: `58ebca64759e2ec66da7a6cab05c3881e39127ca`

`aof init` is offline-capable and does not clone mutable upstream `main`.
