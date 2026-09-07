# AOF CLI v0.3.0 — v0.2 Audit Gap Resolution

This document maps the v0.2.0 audit findings to v0.3.0 implementation changes.

| # | Audit gap | v0.3.0 resolution |
|---:|---|---|
| 1 | AOF-Governed incomplete | Added Authority lifecycle, Policy conflict resolution, dynamic Risk, bounded delegation, approval/escalation, failure budget rules. |
| 2 | AOF-Assured incomplete | Added Evidence, provenance, Verification, risk-based verifier independence, completion gate, accountability chain. |
| 3 | `Enabled` implied implementation | Replaced with distinct Applicability, NormativeLevel, AdoptionState, ImplementationState, VerificationState. |
| 4 | Agent treated as AI-only | Added LLM, Deterministic, Human, Hybrid, ExternalService Agent types; Core always has identifiable Agent scope. |
| 5 | Applicability mixed with status | Separated applicability from normative/adoption/implementation/verification dimensions. |
| 6 | Profile ambiguous target vs claim | Added target base/domain/overlay composition; claimed profile remains null/none after init. |
| 7 | Bootstrap manifest confused with ConformanceManifest | Renamed to `.aof/bootstrap-manifest.json`, schema type `AOFCLIBootstrapManifest`; canonical schemas remain separate. |
| 8 | AOF-Secure-SDLC absent | Added first-class `--secure-sdlc`, generated profile artifact and Skill; requires Governed/Assured base. |
| 9 | Safety Kernel missing | Added `aof/control/safety-kernel.yaml` with six logical components and fail-controlled semantics. |
| 10 | Four Planes shallow | Added Reasoning/Control/Effect/Assurance plane, Trust Boundary and Effect Boundary artifacts. |
| 11 | Authority model shallow | Added lifecycle, scope, delegation, attenuation/provenance and revalidation placeholders. |
| 12 | Policy conflict resolution absent | Added authoritative source/versioning/conflict-resolution/precedence/override model. |
| 13 | Risk model not dynamic | Added dynamic reassessment, inherent/residual Risk, acceptance authority, treatment, retry/failure budgets. |
| 14 | Evidence/Verification shallow | Added provenance/integrity/freshness/sufficiency and verifier-independence/completion/reverification guidance. |
| 15 | StateTransition metadata only | Added authoritative State, before/after, pre/post, owner/version, conflict, idempotency, replay, TOCTOU and reconciliation model. |
| 16 | Trace model shallow | Added actor/correlation/causality/integrity/retention/redaction/correction/failure requirements in generated guidance/Skill. |
| 17 | Failure/recovery missing | Added containment, reconciliation, retry/replan, compensation, rollback, partial/unknown effect, break-glass and deadlock/livelock model. |
| 18 | Human Governance shallow | Added approval scope/freshness, quorum, separation of duties, RiskAcceptance authority, break-glass and non-delegable responsibility model. |
| 19 | Canonical objects not backed by schemas | Embedded pinned AOF v1.0 LTS canonical JSON Schema bundle under `aof/schemas/`. |
| 20 | Upstream not actually pinned | Added specification and schema-bundle SHA-256 provenance; no unverified commit SHA is asserted. |
| 21 | Requirement Registry missing | Added `.aof/requirements.yaml` and `aof/requirements/registry.yaml` with canonical requirement IDs. |
| 22 | Wizard inputs did not affect rules | Criticality/data sensitivity/consequential scope now strengthen verifier/data controls; profile composition changes control rules. |
| 23 | Typed inputs not validated | Added validation for criticality, data sensitivity, policy mechanism, Agent types and profile composition. |
| 24 | Agent count/multi-agent contradiction | Removed free-form Agent count; model uses Agent types plus explicit multi-Agent interaction scope. |
| 25 | Adoption mode always greenfield | Discovery now distinguishes empty/greenfield from existing/brownfield repositories. |
| 26 | Existing AGENTS hard conflict | Existing `AGENTS.md` is preserved and AOF instructions are generated as `AGENTS.aof.md` for explicit user-reviewed merge; no silent merge occurs. |
| 27 | Dead `aof-governance` Skill | Removed; specific Authority/Policy/Risk/etc. Skills are used instead. |
| 28 | Semantic tests too small | Added profile, composition, non-AI Agent, risk-strengthening, requirement-traceability and E2E tests. |
| 29 | Negative profile tests missing | Added invalid Secure-SDLC/Core and High-Assurance/non-Assured composition tests. |
| 30 | Strong foundations should remain | Preserved offline init, no LLM/runtime dependency, path safety, plan-first conflicts, rollback, semantic boundaries and no false conformance claim. |

All listed v0.2 audit findings are addressed in v0.3.0. Existing-file handling remains conservative: AOF CLI never performs a silent heuristic merge.
