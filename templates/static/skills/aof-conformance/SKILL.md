---
name: aof-conformance
description: Review declared AOF adoption and conformance without converting schema validity, maturity, or partial implementation into AOF conformance.
---

# AOF Conformance

## Trigger
Use for conformance claims, release review, adoption status, schema validation, or changes to `.aof/` and `aof/conformance/`.

## Required reads
- `.aof/applicability.yaml`
- `.aof/config.yaml`
- `aof/conformance/scope.yaml`
- `aof/conformance/gaps.md`

## Procedure
1. Establish declared scope and profile.
2. Review every applicable/required control status.
3. Treat Deferred and Unsupported as unresolved gaps.
4. Verify NotApplicable has a defensible scope reason.
5. Keep `SchemaValidity != SemanticValidity != AOFConformance`.
6. Keep `Conformance != Maturity`.
7. Reject full-conformance language unless supported by the applicable conformance process/evidence.
