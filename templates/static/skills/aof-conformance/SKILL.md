---
name: aof-conformance
description: Review the repository's declared AOF adoption scope without converting partial implementation into false conformance.
---

# AOF Conformance Guidance

Review the project-defined profile, scope, and control status in `.aof/config.yaml`.

Preserve:

- ScopedAdoption -> ScopedConformance
- Deferred != Satisfied
- Unsupported != Satisfied
- NotApplicable must reflect true non-applicability, not convenience
- SchemaValidity != SemanticValidity != AOFConformance
