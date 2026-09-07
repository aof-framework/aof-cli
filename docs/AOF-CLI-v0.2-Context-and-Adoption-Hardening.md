# AOF CLI v0.2.0 — Context & Adoption Hardening Specification

**Status:** Implementation Baseline  
**AOF baseline:** AOF v1.0 LTS  
**Canonical upstream:** https://github.com/aof-framework/aof

## 1. Objective

v0.2.0 corrects a weakness in v0.1.0: successful initialization produced an AOF-ready skeleton but projected too little of the framework's semantic depth into the software project.

v0.2.0 therefore changes the product model from:

```text
Project metadata -> generic templates
```

to:

```text
Project discovery
+ Adoption interview
+ AOF applicability analysis
+ Governance model generation
+ Agent context generation
```

AOF semantics remain static. Project adoption remains dynamic.

## 2. Design rule

The goal is **not** to copy the ~25k-line framework specification into every repository.

```text
Richness != Size
```

The goal is to project the applicable subset into operational project guidance:

```text
AOF Specification
-> Applicable Semantics
-> Project-Specific Governance Model
-> Agent Skills / Responsibilities
```

## 3. Context domains

`aof init` models eight project/adoption domains:

1. Project Context
2. Agent Model
3. Capability Model
4. Authority Model
5. Risk Model
6. Policy Model
7. Execution Model
8. Assurance / Conformance Model

Questions are conditional. Consequential execution materially expands the interview and generated context.

## 4. Canonical object applicability

The CLI maps all 22 AOF canonical object families:

Goal, Task, Agent, ContextDescriptor, Resource, Capability, AuthorityGrant, Policy, RiskAssessment, ActionProposal, Decision, ExecutionContract, Evidence, Verification, Approval, StateTransition, TraceEvent, AgentInteractionContract, EscalationPackage, Outcome, ConformanceManifest, ConformanceReport.

The mapping is written to `.aof/applicability.yaml`.

Object applicability is not an implementation or conformance claim.

## 5. Consequential execution rule

When project inputs indicate persistent state mutation, external-system mutation, infrastructure action, or consequential execution, the following control families become applicable/required:

- ActionProposal;
- AuthorizedDecision;
- Authority;
- Policy;
- Risk;
- State validation;
- ExecutionContract;
- Effect Boundary;
- Evidence;
- Verification;
- Trace.

Canonical flow:

```text
Agent output
-> UntrustedProposal
-> Governance Evaluation
-> AuthorizedDecision
-> ExecutionContract
-> Effect Boundary
-> Effect
-> Evidence
-> Verification
-> Trace
```

Execution predicate:

```text
ExecuteAllowed = C AND H AND P AND S AND R AND V
```

Mandatory `Fail` blocks. Mandatory `Pending` remains Pending.

## 6. Generated project knowledge

v0.2 generates structured context under:

```text
.aof/
aof/governance/
aof/agents/
aof/authority/
aof/policies/
aof/risk/
aof/execution/
aof/assurance/
aof/conformance/
```

`AGENTS.md` remains a router. It MUST not become a copy of the full AOF specification.

## 7. Agent Skills

Static AOF skills are procedural and selectively installed:

- aof-core
- aof-adoption
- aof-agent-design
- aof-authority
- aof-policy
- aof-risk
- aof-execution
- aof-evidence-verification
- aof-state-trace
- aof-conformance

Project-specific skills remain dynamic:

- project-architecture
- project-domain

Each skill defines trigger, required reads, procedure, invariants/checks, and stop/blocking conditions where applicable.

## 8. Agent roles

Generated Agent definitions are responsibility-driven:

- architect;
- implementer;
- reviewer;
- aof-reviewer.

Agent roles consume Skills rather than duplicating AOF knowledge.

## 9. Progressive adoption

v0.2 preserves:

```text
AOF-Core ⊆ AOF-Governed ⊆ AOF-Assured
```

and:

```text
AdoptAOF != ImplementEverything
ScopedAdoption -> ScopedConformance
Deferred != Satisfied
Unsupported != Satisfied
NotApplicable != DisabledForConvenience
```

Organizations may adopt progressively without weakening frozen semantics.

## 10. Definition of done

v0.2 is complete when:

- CLI version reports 0.2.0;
- richer interview is available;
- 22 canonical objects are mapped;
- consequential execution activates the required governance projection;
- project governance models are generated;
- applicable Agent Skills are substantive and selectively installed;
- Agent roles are responsibility-driven;
- explicit adoption gaps are generated;
- initialization remains deterministic/offline;
- filesystem safety/idempotency are preserved;
- `go test ./...` and `go vet ./...` pass.
