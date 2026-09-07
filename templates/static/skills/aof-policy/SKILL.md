---
name: aof-policy
description: Apply AOF Policy semantics and distinguish enforceable policy controls from prompt-only instructions.
---

# AOF Policy

## Trigger
Use when adding rules, constraints, guardrails, approvals, policy checks, or prompt instructions intended to constrain execution.

## Required reads
- `aof/policies/policy-model.yaml`
- `.aof/config.yaml`

## Procedure
1. Identify the applicable Policy.
2. Identify the enforcement mechanism.
3. Separate prompt guidance from actual PolicyEnforcement.
4. Evaluate mandatory Policy state as Pass, Fail, or Pending.
5. Block execution on mandatory Fail or Pending.

## Blocking finding
A prompt instruction MUST NOT be represented as sufficient PolicyEnforcement for consequential execution.
