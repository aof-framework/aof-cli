---
name: aof-adoption
description: Maintain the project's scoped and progressive AOF adoption model without requiring every organization to implement the full framework at once.
---

# AOF Adoption

## Trigger
Use when changing AOF profile, project scope, organizational capability, control status, or generated adoption artifacts.

## Required reads
- `.aof/project.yaml`
- `.aof/config.yaml`
- `.aof/adoption.yaml`
- `.aof/applicability.yaml`
- `aof/conformance/gaps.md`

## Decision rules
- `AdoptAOF != ImplementEverything`.
- `ScopedAdoption -> ScopedConformance`.
- `Deferred != Satisfied`.
- `Unsupported != Satisfied`.
- `NotApplicable != DisabledForConvenience`.
- A reduced scope MUST NOT create a full AOF conformance claim.

## Procedure
1. Determine whether a control is applicable to declared project behavior.
2. Separately determine whether the organization can currently support it.
3. Record the result explicitly as Required, Enabled, Deferred, Unsupported, or NotApplicable.
4. Update `aof/conformance/gaps.md` for unresolved applicable controls.
5. Do not weaken canonical semantics to fit organizational capability.
