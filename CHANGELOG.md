# Changelog

## v0.3.0 — Semantic Fidelity & Adoption Hardening

### Added

- Requirement-driven applicability model with AOF requirement IDs.
- Separate target and claimed profile/conformance state.
- Separate applicability, normative level, adoption, implementation, and verification dimensions.
- Agent types: LLM, Deterministic, Human, Hybrid, ExternalService.
- AOF-Secure-SDLC domain profile support.
- AOF-High-Assurance strengthening overlay support.
- Four-Plane architecture and Trust/Effect Boundary artifacts.
- Safety Kernel and `ExecuteAllowed = C AND H AND P AND S AND R AND V` gate mapping.
- Authority lifecycle/scope/delegation model.
- Policy conflict-resolution model.
- Dynamic Risk/residual-risk/failure-budget model.
- Rich StateTransition and failure/recovery models.
- Evidence provenance, Verification independence/completion, Trace, and accountability guidance.
- Human Governance model.
- Pinned canonical AOF v1.0 LTS JSON Schemas in generated repositories.
- Upstream specification and schema-bundle SHA-256 provenance.
- Secure-SDLC and High-Assurance Agent Skills.
- Broader semantic/profile/negative tests.

### Changed

- `aof init` no longer treats declared/supported controls as implemented.
- AOF Agent is no longer modeled as AI-only.
- AOF-Governed and AOF-Assured now include their canonical strengthening semantics.
- High/Critical consequential scope strengthens verifier-independence requirements.
- Confidential/Restricted data scope activates context/trace data controls.
- Existing repositories are identified as brownfield adoption.
- `.aof/manifest.json` replaced by `.aof/bootstrap-manifest.json` to avoid confusion with canonical ConformanceManifest.
- Initialization summary explicitly reports no conformance claim.

### Removed

- Unused generic `aof-governance` Skill; specific governance Skills are authoritative.

## v0.2.0 — AOF Context & Adoption Hardening

Introduced richer project discovery, canonical-object applicability, substantive Agent Skills, and project governance artifacts.

## v0.1.0 — Complete Bootstrap MVP

Initial safe/offline AOF repository bootstrap with `aof init`, `aof version`, and `aof help`.
