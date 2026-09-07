# AOF CLI v0.3.x — Semantic Fidelity & Canonical Coverage

## Status

Implementation baseline for v0.3.0.

## Objective

v0.3.0 addresses the semantic gaps found in the v0.2.0 audit against AOF v1.0 LTS. The release changes AOF CLI from a primarily context-rich bootstrap generator into a **requirement-driven adoption model generator** while remaining non-runtime and offline-capable.

## Core state model

v0.3.0 uses distinct dimensions:

```text
Applicability
NormativeLevel
AdoptionState
ImplementationState
VerificationState
```

Therefore:

```text
Applicable != Implemented
Planned != Implemented
Implemented != Verified
TargetProfile != ClaimedProfile
```

Initialization MUST NOT claim AOF conformance.

## Profile model

Base profiles:

```text
AOF-Core
AOF-Governed
AOF-Assured
```

Domain profile:

```text
AOF-Secure-SDLC
```

Strengthening overlay:

```text
AOF-High-Assurance
```

Dependencies enforced by CLI:

```text
AOF-Secure-SDLC -> AOF-Governed or AOF-Assured
AOF-High-Assurance -> AOF-Assured
```

### AOF-Core

Projects model identifiable Agent, Task, Action/Proposal/Decision separation, Authority, Policy, authoritative State, Trace, controlled transition, and no implicit allow.

### AOF-Governed

Adds explicit Authority lifecycle, Policy conflict resolution, dynamic Risk, bounded delegation, approval/escalation, and failure budgets.

### AOF-Assured

Adds Evidence, provenance, Verification, risk-based verifier independence, completion gate, and accountability chain.

### AOF-Secure-SDLC

Adds requirements/acceptance criteria, architecture/security review, threat modeling, code/test/security Verification, release/deployment Authority gates, Human gates by Risk, and applicable security Evidence.

### AOF-High-Assurance

Adds separation of duties, independent Verification, Critical-action approval, stronger Evidence, strict Authority scope, Effect Boundary revalidation, Trace protection, bounded Recovery, and residual Risk handling.

## Agent model

AOF Agent is not AI-only. Supported canonical Agent types:

```text
LLM
Deterministic
Human
Hybrid
ExternalService
```

Agent type never implies Authority.

## Requirement traceability

`.aof/requirements.yaml` and `aof/requirements/registry.yaml` index applicable AOF requirement IDs. The specification remains normative; the generated registry is an index.

The embedded canonical baseline is hard-gated at:

```text
registered_invariants == 162
registered_requirements == 332
stable_semantic_ids == 494
canonical_traceability_edges == 352
```

The gate validates exact ID-set hashes in addition to counts. It rejects duplicate or malformed IDs, invalid normative classifications, unknown edge endpoints, broken profiles/canonical objects, and silent mapping gaps.

Every canonical requirement participates in project projection. Requirements selected by a control/profile rule become `applicable`; all other requirements remain explicitly `conditional` pending profile and project-scope evaluation. They never silently disappear, and initialization does not treat a conditional item as implemented, verified, or conformant.

Canonical upstream intentionally does not assert a direct invariant, verification, Evidence, or canonical-object mapping for every item. The CLI preserves those absences using explicit `CanonicalMappingNotAsserted` dispositions instead of inventing heuristic edges. This satisfies semantic accounting while keeping upstream traceability debt visible.

Traceability direction:

```text
Invariant -> Requirement -> Implementation/Test -> Evidence
```

## Safety Kernel and planes

Generated architecture includes:

```text
Reasoning Plane
Control Plane
Effect Plane
Assurance Plane
```

and logical Safety Kernel components:

```text
AuthorityEvaluator
PolicyEvaluator
StateValidator
RiskGate
VerificationGate
TraceRecorder
```

Consequential effects use:

```text
ExecuteAllowed = C AND H AND P AND S AND R AND V
```

## Canonical schemas and provenance

The CLI embeds the pinned AOF v1.0 LTS canonical JSON Schema bundle and copies it under `aof/schemas/`.

Generated provenance records:

- canonical repository;
- specification path;
- specification SHA-256;
- schema bundle SHA-256.

No unverified Git commit SHA is invented.

The registry import pins upstream commit `58ebca64759e2ec66da7a6cab05c3881e39127ca`. Provenance records the published specification checksum separately from the semantic-freeze checksum declared by upstream, avoiding the previous path/artifact ambiguity.

## Conformance boundary

`.aof/bootstrap-manifest.json` tracks CLI generation. It is explicitly not canonical `ConformanceManifest`.

Generated `aof/conformance/scope.yaml` states:

```text
claimed_profile: null
claim_status: none
full_conformance_claimed: false
```

A future validation/conformance process must assess implementation and Evidence before profile claims are allowed.
