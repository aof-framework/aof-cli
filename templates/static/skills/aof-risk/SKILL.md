---
name: aof-risk
description: Apply AOF RiskAssessment semantics, project risk categories, and execution gating without conflating assessment with risk acceptance.
---

# AOF Risk

## Trigger
Use for consequential actions, security-sensitive changes, external effects, irreversible state, or changes to project risk controls.

## Required reads
- `aof/risk/risk-model.yaml`
- `.aof/project.yaml`
- `aof/execution/consequential-actions.yaml`

## Procedure
1. Identify affected project risk categories.
2. Evaluate likelihood/impact using project conventions where defined.
3. Record RiskAssessment independently from any RiskAcceptance decision.
4. Escalate when risk exceeds declared authority or organizational capability.
5. Treat unresolved mandatory risk evaluation as Pending, never Pass.

## Invariant
`RiskAssessment != RiskAcceptance`.
