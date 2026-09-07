---
name: aof-evidence-verification
description: Design and review Evidence and Verification flows while preserving Claim, Evidence, and Verification as separate AOF concepts.
---

# AOF Evidence & Verification

## Trigger
Use when proving preconditions/outcomes, verifying execution, collecting artifacts, or implementing assurance controls.

## Required reads
- `aof/assurance/evidence-model.md`
- `aof/assurance/verification-model.md`
- `.aof/config.yaml`

## Procedure
1. Identify the Claim or acceptance condition.
2. Identify objective Evidence that can support/refute it.
3. Define the Verification procedure separately.
4. Record Pass, Fail, or Pending explicitly.
5. Where independent Verification is required, prevent the same actor from silently self-verifying without the declared independent mechanism.

## Invariants
- `Claim != Evidence != Verification`
- `Pending != Pass`
