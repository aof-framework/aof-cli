---
name: aof-state-trace
description: Apply AOF authoritative State, controlled StateTransition, concurrency, idempotency, retry/recovery, reconciliation, and Trace reconstruction semantics.
---
# AOF State & Trace

## Required reads
- `aof/execution/state-transitions.yaml`
- `aof/execution/failure-recovery.yaml`
- `aof/assurance/trace-model.md`
- `aof/architecture/planes.yaml`

## Procedure
1. Identify authoritative orchestration State separately from Agent private memory.
2. Bind before_state, after_state, preconditions, postconditions, owner, and version.
3. Define conflict-control, idempotency, replay, TOCTOU/revalidation, and external-state reconciliation.
4. Preserve governance gates across retry, replan, recovery, compensation, and rollback.
5. Trace actor, Proposal, Decision, Authority, Policy, Risk, Action, StateTransition, Evidence, Verification, and Outcome as applicable.
6. Preserve ordering/causality, correction history, integrity, retention, redaction/access control, and failure behavior.
7. Trace MUST NOT require private chain-of-thought and MUST NOT itself grant Authority.
