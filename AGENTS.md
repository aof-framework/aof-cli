# AOF CLI Repository Instructions

This repository implements **AOF CLI**, the bootstrap/adoption utility for AOF v1.0 LTS.

## Product boundary

```text
AOF Framework != AOF CLI
AOF CLI = Adoption Configurator + Project-Aware Bootstrap Generator
```

AOF CLI MUST NOT become an orchestration runtime, policy/authority/risk runtime engine, LLM client, telemetry system, or application dependency.

## Canonical source

AOF semantics come from:

https://github.com/aof-framework/aof

Before changing semantic generation or profile/applicability behavior, read the relevant AOF v1.0 LTS sections and `docs/AOF-CLI-v0.3-Semantic-Fidelity-and-Adoption-Hardening.md`.

## v0.3 invariants

- `Reasoning != Decision != Authority != Action`
- `Capability != Authority`
- `Proposal != AuthorizedDecision`
- `AgentOutput = UntrustedProposal`
- `PolicyPrompt != PolicyEnforcement`
- `Approval != AuthorityGrant`
- `RiskAssessment != RiskAcceptance`
- `Pending != Pass`
- `Claim != Evidence != Verification`
- `SchemaValidity != SemanticValidity != AOFConformance`
- `Applicability != NormativeLevel != AdoptionState != ImplementationState != VerificationState`
- target profile is not a conformance claim.
- `.aof/bootstrap-manifest.json != ConformanceManifest`.

## Engineering

Prefer deterministic stdlib-heavy Go. `aof init` MUST remain offline-capable, deterministic, safe against path traversal, and non-destructive on conflicts.

Before commit:

```bash
go fmt ./...
go vet ./...
go test ./...
```
