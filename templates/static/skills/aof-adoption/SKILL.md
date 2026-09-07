---
name: aof-adoption
description: Maintain scoped AOF adoption while separating applicability, normative obligation, adoption intent, implementation state, verification state, and conformance claims.
---
# AOF Adoption

## Trigger
Use when changing target profile, scope, Agent types, organizational capability, applicability, control planning, or conformance-facing artifacts.

## Required reads
- `.aof/project.yaml`
- `.aof/config.yaml`
- `.aof/adoption.yaml`
- `.aof/applicability.yaml`
- `.aof/requirements.yaml`
- `aof/conformance/gaps.md`

## Decision rules
- `Applicability != NormativeLevel != AdoptionState != ImplementationState != VerificationState`.
- Target profile is not claimed profile.
- `planned != implemented`.
- `not_assessed != pass`.
- `unsupported != satisfied`.
- `NotApplicable != DisabledForConvenience`.
- `ScopedAdoption -> ScopedConformance`.
- Stronger/profile-composed adoption MUST preserve mandatory base requirements.

## Procedure
1. Determine project scope and Agent/effect boundaries.
2. Determine requirement applicability using requirement IDs and profile composition.
3. Record adoption intent without asserting implementation.
4. Require implementation Evidence before changing implementation state.
5. Require evaluation before changing verification/conformance state.
6. Report unresolved mandatory requirements as gaps rather than weakening applicability.
