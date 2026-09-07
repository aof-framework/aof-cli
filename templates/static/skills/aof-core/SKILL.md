---
name: aof-core
description: Apply canonical AOF v1.0 LTS semantic boundaries whenever reasoning about Agents, Capability, Authority, Policy, Risk, Decision, Action, Evidence, Verification, State, Trace, or Conformance.
---

# AOF Core

## Trigger
Use for any change that touches AI behavior, governance, execution, assurance, or AOF artifacts.

## Required reads
1. `AOF.md`
2. `.aof/adoption.yaml`
3. `.aof/applicability.yaml`
4. Relevant AOF project model under `aof/`

## Canonical invariants
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

## Procedure
1. Identify the actor and requested behavior.
2. Identify which canonical AOF objects are applicable in `.aof/applicability.yaml`.
3. Preserve all semantic separations above in code and documentation.
4. Never infer a missing governance state as `Pass`.
5. Report a gap rather than inventing project Authority, Policy, Evidence, or Verification.

## Stop conditions
Stop and surface a governance gap when required Authority, Policy, Risk, State, or Verification cannot be established.
