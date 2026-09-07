---
name: aof-agent-design
description: Design or review AOF Agents of type LLM, Deterministic, Human, Hybrid, or ExternalService, including roles, capabilities, delegation, tools, context, and interactions.
---
# AOF Agent Design

## Trigger
Use when adding/changing any operational Agent, Agent role, Capability, delegation, interaction, tool use, context boundary, or Agent replacement.

## Required reads
- `aof/agents/inventory.yaml`
- `aof/agents/capabilities.yaml`
- `aof/agents/interactions.yaml`
- `aof/governance/governance-envelope.yaml`
- `.aof/applicability.yaml`

## Procedure
1. Identify Agent type and role; Agent type does not determine Authority.
2. Declare Capability separately from Authority.
3. Identify Resource and ContextDescriptor boundaries.
4. Treat discretionary/LLM output as `UntrustedProposal` until governed evaluation.
5. Require AgentInteractionContract when multi-Agent interaction/delegation is applicable.
6. Preserve bounded delegation, attenuation, provenance, and GovernanceEnvelope.
7. Prevent any Agent from expanding its own Authority.

## Blocking findings
Self-granted Authority; technical reachability treated as Authority; unbounded delegation; missing interaction contract where required; consequential effect directly connected to discretionary reasoning without governance mediation.
