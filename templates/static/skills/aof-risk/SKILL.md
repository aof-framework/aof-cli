---
name: aof-risk
description: Apply AOF RiskAssessment, dynamic reassessment, inherent/residual risk, acceptance authority, treatment, retry/failure budgets, and risk-based assurance.
---
# AOF Risk

## Required reads
- `aof/risk/risk-model.yaml`
- `.aof/project.yaml`
- `aof/execution/consequential-actions.yaml`
- `.aof/requirements.yaml`

## Procedure
1. Use an explicit classification method and Risk Profile.
2. Evaluate likelihood/impact and inherent Risk where applicable.
3. Evaluate controls and residual Risk separately.
4. Keep RiskAssessment separate from RiskAcceptance and identify acceptance Authority.
5. Reassess after material Risk, State, Policy, Authority, retry, replan, or partial-effect changes.
6. Apply treatment/reject/abort/escalation according to thresholds.
7. Require risk-based verifier independence for High/Critical consequential Risk.
8. Record failure/retry budgets where AOF-Governed applies.

`RiskAssessment != RiskAcceptance`; unresolved mandatory Risk never becomes Pass.
