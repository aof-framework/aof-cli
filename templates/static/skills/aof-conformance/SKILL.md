---
name: aof-conformance
description: Evaluate AOF adoption and conformance without converting target profiles, schema validity, maturity, planned controls, or partial implementation into AOF conformance.
---
# AOF Conformance

## Required reads
- `.aof/applicability.yaml`
- `.aof/requirements.yaml`
- `.aof/config.yaml`
- `aof/conformance/scope.yaml`
- `aof/conformance/gaps.md`
- canonical schemas under `aof/schemas/`

## Procedure
1. Distinguish target profile from claimed profile.
2. Establish explicit scope, exclusions, environments, Agent types, tool classes, and governance boundaries.
3. Evaluate each applicable requirement independently.
4. Treat `not_assessed`, `not_evaluated`, unsupported, and unresolved mandatory requirements as non-satisfaction.
5. Validate any NotApplicable claim structurally and semantically.
6. Preserve profile dependencies/composition without weakening mandatory requirements.
7. Keep `SchemaValidity != SemanticValidity != AOFConformance` and `Conformance != Maturity`.
8. Do not treat `.aof/bootstrap-manifest.json` as canonical `ConformanceManifest`.
