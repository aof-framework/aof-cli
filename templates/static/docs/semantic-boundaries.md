# AOF Semantic Boundaries

Canonical upstream: https://github.com/aof-framework/aof

This project adopts AOF v1.0 LTS. Preserve these boundaries:

- Reasoning != Decision != Authority != Action
- Capability != Authority
- Proposal != AuthorizedDecision
- AgentOutput = UntrustedProposal
- PolicyPrompt != PolicyEnforcement
- Approval != AuthorityGrant
- RiskAssessment != RiskAcceptance
- Pending != Pass
- Claim != Evidence != Verification
- SchemaValidity != SemanticValidity != AOFConformance

Consequential execution must not be inferred from model output alone.
