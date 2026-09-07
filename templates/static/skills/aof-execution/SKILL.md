---
name: aof-execution
description: Govern consequential AI execution through Proposal, Decision, ExecutionContract, Effect Boundary, Evidence, Verification, State, and Trace.
---

# AOF Execution

## Trigger
Use whenever an AI-controlled or AI-initiated action may alter persistent state, external systems, accounts, infrastructure, financial state, approvals, or other consequential effects.

## Required reads
- `AOF.md`
- `aof/authority/authority-model.yaml`
- `aof/policies/policy-model.yaml`
- `aof/risk/risk-model.yaml`
- `aof/execution/execution-model.md`
- `aof/execution/consequential-actions.yaml`
- `aof/assurance/verification-model.md`

## Procedure
1. Identify actor and proposed Action.
2. Represent Agent output as ActionProposal / UntrustedProposal.
3. Confirm Capability.
4. Establish Authority.
5. Evaluate Policy.
6. Validate current State.
7. evaluate Risk.
8. Resolve required Verification conditions.
9. Produce an AuthorizedDecision.
10. Construct/bind the ExecutionContract.
11. Cross the Effect Boundary only after all mandatory controls pass.
12. Capture Evidence, Verification result, StateTransition, Outcome, and Trace.

## Execution predicate
`ExecuteAllowed = C AND H AND P AND S AND R AND V`

Mandatory Fail blocks execution. Mandatory Pending remains Pending and MUST NOT become implicit Allow.
