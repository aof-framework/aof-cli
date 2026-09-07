# AOF Execution Model

Canonical upstream: https://github.com/aof-framework/aof

AOF preserves separation between Proposal, Decision, Authority, and Action. For consequential effects, project implementations should preserve the logical boundary:

```text
Proposal -> Governance Evaluation -> Decision -> Execution Boundary -> Effect -> Evidence -> Verification -> Trace
```

Mandatory unresolved governance state must not be treated as Allow.
