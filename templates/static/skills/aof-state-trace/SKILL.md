---
name: aof-state-trace
description: Apply AOF StateTransition and Trace semantics to consequential actions, retries, recovery, and governance chronology.
---

# AOF State & Trace

## Trigger
Use for persistent state mutation, workflow transitions, retries, recovery, execution history, or audit/trace implementation.

## Required reads
- `aof/execution/state-transitions.yaml`
- `aof/assurance/trace-model.md`

## Procedure
1. Identify current state and intended transition.
2. Validate transition preconditions before effect.
3. Associate Proposal, Decision, Authority, Policy, Risk, execution, Evidence, Verification, Outcome, and transition in Trace where applicable.
4. Ensure retry/recovery does not bypass governance gates.
5. Never treat Trace itself as Authority.
