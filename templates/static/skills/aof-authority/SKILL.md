---
name: aof-authority
description: Apply AOF Authority lifecycle, scope, delegation, validity, approval separation, revocation, expiry, consumption, provenance, and revalidation.
---
# AOF Authority

## Trigger
Use for privileged/consequential actions, approval flows, delegated operational authority, role checks, service accounts, release gates, or replacement actors.

## Required reads
- `aof/authority/authority-model.yaml`
- `aof/governance/governance-root.md`
- `aof/execution/consequential-actions.yaml`
- `.aof/requirements.yaml`

## Procedure
1. Identify actor and Authority-sensitive Action.
2. Confirm Capability independently.
3. Resolve an explicit AuthorityGrant or equivalent governed grant.
4. Validate issuer, status, time, environment, resource, operation, quantity, and other scope.
5. Validate delegation depth, attenuation, non-delegability, and provenance where applicable.
6. Revalidate after material replan and at Effect Boundary according to policy.
7. Reject revoked, expired, suspended, or consumed Authority.
8. Keep Approval separate from AuthorityGrant.

## Invariants
- `Capability != Authority`
- `Approval != AuthorityGrant`
- `NoGrant => NoAuthoritySensitiveExecution`
- `AgentType != Authority`
