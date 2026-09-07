---
name: aof-authority
description: Apply AOF Authority semantics to authorization, delegated operational authority, approvals, grants, and authority-sensitive execution.
---

# AOF Authority

## Trigger
Use for approval flows, privileged actions, delegated authority, role checks, or any consequential execution.

## Required reads
- `aof/authority/authority-model.yaml`
- `aof/governance/governance-root.md`
- `aof/execution/consequential-actions.yaml`

## Procedure
1. Identify the actor.
2. Identify the action requiring Authority.
3. Confirm Capability independently.
4. Resolve an explicit AuthorityGrant or equivalent project authority mechanism.
5. Confirm scope, target, environment, and validity of the grant.
6. Treat Approval as a separate fact from AuthorityGrant.
7. Block authority-sensitive execution when no valid grant exists.

## Invariants
- `Capability != Authority`
- `Approval != AuthorityGrant`
- `NoGrant => NoAuthoritySensitiveExecution`
