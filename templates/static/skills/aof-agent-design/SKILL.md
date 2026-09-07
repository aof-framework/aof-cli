---
name: aof-agent-design
description: Design or review AI Agents, their responsibilities, capabilities, delegation, tool use, and interactions under AOF bounded-agency rules.
---

# AOF Agent Design

## Trigger
Use when adding/changing an AI Agent, tool capability, delegation, multi-agent interaction, or agent responsibility.

## Required reads
- `aof/agents/inventory.yaml`
- `aof/agents/capabilities.yaml`
- `aof/governance/governance-envelope.yaml`
- `.aof/applicability.yaml`

## Procedure
1. State the Agent responsibility in operational terms.
2. Declare required Capability separately from Authority.
3. Identify Resources and ContextDescriptor boundaries.
4. Treat Agent output as `UntrustedProposal` unless a later governance step authorizes action.
5. For multi-agent systems, require an applicable AgentInteractionContract.
6. Prevent Agents from expanding their own GovernanceEnvelope.

## Blocking findings
- Self-granted Authority.
- Tool access treated as Authority.
- Delegation without bounded authority/interaction contract where applicable.
- Agent output directly connected to consequential effect without governance evaluation.
