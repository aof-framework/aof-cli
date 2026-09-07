# AOF CLI v0.1 — Product Requirements & Technical Design

**Document Status:** Draft Baseline for Implementation  
**Product:** AOF CLI  
**Target Release:** v0.1.0 — Complete Bootstrap MVP  
**Implementation Language:** Go  
**AOF Baseline:** AOF v1.0 LTS  
**Canonical AOF Repository:** https://github.com/aof-framework/aof  
**Canonical AOF Specification:** https://github.com/aof-framework/aof/blob/main/specification/AOF-v1.0-Framework-Specification.md

---

## 1. Executive Summary

AOF CLI adalah developer bootstrap utility untuk mengadopsi **AI Orchestration Framework (AOF) v1.0 LTS** ke dalam software repository secara cepat, aman, reproducible, dan project-aware.

AOF CLI tidak menggantikan AOF Framework dan tidak menjadi runtime dependency aplikasi. Tanggung jawab utamanya adalah mengubah repository kosong maupun repository yang sudah berjalan menjadi repository yang memiliki konteks AOF yang dapat dipahami oleh manusia dan AI coding agents.

Target user experience utama:

```bash
mkdir itsm-backend
cd itsm-backend
aof init
```

Hasil konseptual:

```text
Empty / Existing Repository
          |
          v
       aof init
          |
          v
    AOF-ready Repository
```

AOF CLI dirancang berdasarkan dua fakta utama:

1. **AOF Framework bersifat statis.** Semantics AOF v1.0 LTS, canonical terminology, governance boundaries, canonical contracts, dan conformance semantics tidak berubah berdasarkan project.
2. **AOF adoption context bersifat dinamis.** Setiap project mempunyai domain, architecture, features, AI use cases, authority boundaries, risk profile, engineering practices, dan organizational capability yang berbeda.

Karena itu:

```text
GeneratedContext
=
StaticAOFSubset
+
DynamicProjectAdoptionModel
```

AOF CLI bukan sekadar file copier. AOF CLI adalah:

```text
Project-aware Adoption Configurator
+
Bootstrap Generator
```

---

# 2. Product Context

## 2.1 Apa itu AOF

AOF adalah **AI Orchestration Framework** yang mendefinisikan governance, bounded agency, execution control, evidence, verification, state, traceability, dan conformance untuk sistem yang menggunakan AI Agents maupun automation.

AOF v1.0 LTS merupakan baseline specification yang frozen.

Canonical public repository:

```text
https://github.com/aof-framework/aof
```

Canonical framework specification:

```text
https://github.com/aof-framework/aof/blob/main/specification/AOF-v1.0-Framework-Specification.md
```

AOF CLI MUST memperlakukan repository tersebut sebagai canonical upstream untuk provenance AOF.

AOF CLI MUST NOT menciptakan semantics AOF baru.

Beberapa semantic boundaries fundamental yang MUST dipertahankan:

```text
Reasoning != Decision != Authority != Action
Capability != Authority
Proposal != AuthorizedDecision
AgentOutput = UntrustedProposal
PolicyPrompt != PolicyEnforcement
Approval != AuthorityGrant
RiskAssessment != RiskAcceptance
Pending != Pass
Claim != Evidence != Verification
SchemaValidity != SemanticValidity != AOFConformance
```

Canonical execution predicate:

```text
ExecuteAllowed = C AND H AND P AND S AND R AND V
```

Mandatory unresolved governance state tidak boleh berubah menjadi implicit Allow.

---

# 3. Problem Statement

Tanpa tooling, developer atau organisasi yang ingin mengadopsi AOF harus melakukan banyak pekerjaan manual sebelum project dapat menggunakan AOF secara efektif, misalnya:

- membaca specification yang besar;
- menentukan AOF adoption scope;
- memilih profile;
- membuat `AGENTS.md`;
- membuat `AOF.md`;
- membuat Agent Skills;
- membuat Agent roles;
- menentukan project-specific governance context;
- membuat policy placeholders;
- membuat authority model;
- membuat risk model;
- membuat conformance structure;
- mengintegrasikan project knowledge untuk AI coding tools;
- menjaga konsistensi antar tool-specific configuration.

Hal ini menghasilkan:

```text
High Adoption Friction
+
High Probability of Semantic Drift
+
High Context Duplication
+
Inconsistent Agent Instructions
```

Target AOF CLI adalah mengubah:

```text
Read Framework
-> Configure Manually
-> Build Context
-> Start Development
```

menjadi:

```text
aof init
-> Review Generated Adoption Model
-> Start Development
```

---

# 4. Product Goal

## 4.1 Primary Goal

Menghasilkan repository yang **AOF-ready** melalui satu command dengan tetap menjaga:

- AOF semantic integrity;
- project-specific flexibility;
- progressive adoption;
- safe filesystem behavior;
- tool neutrality;
- offline-capable initialization;
- deterministic output;
- no runtime dependency.

Primary contract:

```text
EmptyOrExistingRepo
    -- aof init -->
AOFReadyRepo
```

## 4.2 Secondary Goals

AOF CLI SHOULD:

- membuat AOF mudah ditemukan oleh AI coding agents;
- mengurangi kebutuhan developer membaca seluruh specification sebelum memulai;
- memetakan AOF ke kebutuhan aktual project;
- mendukung organisasi dengan maturity S-SDLC yang berbeda;
- menghasilkan configuration yang eksplisit dan dapat direview;
- mencegah false conformance claims;
- menjaga project tetap portable antar AI coding tools.

---

# 5. Non-Goals

AOF CLI v0.1 MUST NOT menjadi:

- AOF runtime;
- Agent orchestration engine;
- AI Agent runtime;
- LLM client;
- LLM provider abstraction;
- Policy enforcement engine;
- Authority engine;
- Risk engine;
- Safety Kernel runtime;
- Evidence database;
- Verification runtime;
- tracing backend;
- HTTP server;
- daemon;
- SDK;
- application framework;
- cloud service;
- plugin marketplace;
- telemetry collector;
- replacement untuk AOF Reference Implementation.

AOF CLI MUST NOT menjadi runtime dependency software yang di-bootstrap.

Setelah `aof init` selesai, software project MUST tetap usable apabila binary `aof` dihapus dari komputer developer.

```text
AOFFramework != AOFCLI
AOFCLI -> ProjectBootstrap
ApplicationRuntime -/-> AOFCLI
```

---

# 6. Product Principles

1. **Static Semantics, Dynamic Adoption**  
   `AOF Semantics = Static`, `Project Adoption = Dynamic`.

2. **Progressive Adoption**  
   `AOF Adoption != Implement Everything`.

3. **No False Conformance**  
   `Partial Implementation != Full AOF Conformance`.

4. **Repository Context Is Authoritative**  
   `PretrainedKnowledge(AOF) = Optional`; `RepositoryContext(AOF) = Authoritative`.

5. **Tool Neutrality**  
   AOF project knowledge MUST tidak dimiliki oleh configuration satu AI coding tool.

6. **No Silent Overwrite**  
   Existing project content MUST tidak ditimpa secara silent.

7. **Deterministic Bootstrap**  
   Input yang sama terhadap AOF CLI version dan AOF baseline yang sama SHOULD menghasilkan bootstrap artifacts yang semantically equivalent.

8. **Reproducible Upstream**  
   `aof init` MUST NOT mengambil mutable branch `main` sebagai runtime source of truth.

---

# 7. Canonical AOF Upstream Strategy

## 7.1 Canonical Repository

Canonical upstream AOF:

```text
https://github.com/aof-framework/aof
```

Canonical specification:

```text
https://github.com/aof-framework/aof/blob/main/specification/AOF-v1.0-Framework-Specification.md
```

Repository tersebut adalah provenance source AOF Framework.

## 7.2 Design Decision: `aof init` MUST NOT Clone `main`

Walaupun AOF public repository menjadi canonical source, AOF CLI v0.1 SHOULD NOT menjalankan:

```bash
git clone https://github.com/aof-framework/aof
```

sebagai bagian normal dari `aof init`.

Alasan:

1. `main` adalah mutable branch;
2. hasil bootstrap dapat berubah tanpa perubahan AOF CLI version;
3. network menjadi mandatory dependency;
4. initialization menjadi lebih lambat;
5. Git executable menjadi dependency;
6. upstream outage dapat menghentikan bootstrap;
7. sulit menjamin exact AOF baseline;
8. reproducible tests menjadi lebih sulit;
9. project dapat menerima artifact yang belum sengaja dirilis.

Karena itu:

```text
Canonical Upstream != Runtime Clone Source
```

## 7.3 MVP Distribution Model

AOF CLI v0.1 SHOULD menggunakan **embedded bootstrap bundle** yang berasal dari pinned AOF v1.0 LTS baseline.

```text
AOF canonical upstream
        |
        | pinned release / verified source
        v
Bootstrap Asset Build
        |
        v
Embedded Assets
        |
        v
AOF CLI Binary
```

Dengan Go:

```go
//go:embed templates/**
var templates embed.FS
```

Hasil:

```text
Download one binary
-> aof init
-> AOF-ready repository
```

tanpa network dependency saat initialization.

## 7.4 Future Optional Online Source Mode

Future version MAY mendukung online source mode seperti:

```bash
aof init --source github
```

atau `aof sync`, tetapi MUST mengambil immutable Git tag, GitHub Release asset, atau exact commit SHA dan SHOULD melakukan integrity verification.

Future online mode MUST NOT default ke mutable `main`.

---

# 8. Static vs Dynamic Artifact Model

## 8.1 Static AOF Layer

Static artifacts berasal dari AOF v1.0 LTS dan tidak berubah berdasarkan software project:

- semantic boundaries;
- canonical terminology;
- governance model;
- bounded agency principles;
- Safety Kernel semantics;
- execution semantics;
- Evidence / Verification distinction;
- conformance concepts;
- canonical schema references;
- generic AOF Agent Skills.

Static AOF content MUST tidak dimodifikasi untuk membuat project terlihat conformant.

## 8.2 Dynamic Project Layer

Dynamic artifacts dihasilkan berdasarkan project-specific information:

- project name;
- programming language;
- project type;
- domain;
- architecture;
- features;
- deployment characteristics;
- AI use cases;
- Agent responsibilities;
- consequential actions;
- authority boundaries;
- risk classifications;
- Policy applicability;
- organizational capability;
- verification capability;
- S-SDLC maturity;
- selected AOF profile;
- selected controls;
- deferred controls.

## 8.3 Artifact Classification Matrix

| Artifact | Classification | Notes |
|---|---|---|
| AOF semantic summary | Static | berasal dari AOF v1.0 LTS |
| canonical terminology | Static | tidak diubah per project |
| governance model summary | Static | canonical |
| execution model summary | Static | canonical |
| AOF core Agent Skill | Static | reusable |
| AOF governance Agent Skill | Static | reusable |
| AOF conformance guidance | Static / Parameterized | profile/scope aware |
| `AOF.md` | Dynamic | project adoption entry point |
| root `AGENTS.md` | Dynamic | project-aware AI instructions |
| project architecture skill | Dynamic | project-specific |
| project domain skill | Dynamic | project-specific |
| project security skill | Dynamic | project-specific |
| `.agents/agents/*` | Dynamic | responsibility-driven |
| `.aof/project.yaml` | Dynamic | project model |
| `.aof/config.yaml` | Dynamic | AOF adoption config |
| `.aof/manifest.json` | Dynamic | generated artifact metadata |
| policies | Dynamic | project-specific |
| authority mappings | Dynamic | project-specific |
| risk model | Dynamic | project-specific |
| verification expectations | Dynamic | project-specific |

---

# 9. Progressive Adoption Model

Tidak semua organisasi memiliki kemampuan S-SDLC, governance, assurance, atau operational maturity yang sama.

Karena itu:

```text
AdoptAOF != ImplementEverything
```

AOF CLI MUST mendukung progressive governance integration.

## 9.1 Base Profiles

Base profile:

```text
AOF-Core
AOF-Governed
AOF-Assured
```

Relationship:

```text
AOF-Core ⊆ AOF-Governed ⊆ AOF-Assured
```

### AOF-Core

Target:

- small engineering teams;
- early adoption;
- advisory AI;
- low/medium consequential automation;
- organization dengan assurance maturity terbatas.

Typical baseline:

- Governance Root;
- Agent responsibility boundaries;
- Capability boundaries;
- Authority awareness;
- basic Policy;
- Risk awareness;
- Proposal / Decision separation;
- basic consequential action control;
- basic Trace.

### AOF-Governed

Typical additions:

- explicit AuthorityGrant;
- enforceable Policy;
- formal RiskAssessment;
- ExecutionContract;
- State validation;
- stronger Trace;
- Evidence;
- Verification;
- Escalation.

### AOF-Assured

Typical strengthening:

- independent Verification;
- formal assurance;
- stronger evidence lifecycle;
- high-risk execution controls;
- stronger traceability;
- formal conformance evidence;
- high-assurance controls.

---

# 10. Adoption Scope

Profile saja tidak cukup. Dua project dengan `AOF-Governed` dapat memiliki risk surface yang sangat berbeda.

Karena itu:

```text
EffectiveAdoption = Profile + Scope
```

Contoh:

```yaml
profile: AOF-Core

scope:
  ai_classification: true
  ai_recommendation: true
  autonomous_execution: false
  human_approval: true
  external_effects: false
```

---

# 11. Control Status Model

AOF CLI SHOULD tidak hanya memakai boolean `true/false`.

Control state SHOULD mendukung minimal:

```text
Required
Enabled
Deferred
NotApplicable
Unsupported
```

### Required

Control diwajibkan oleh selected profile/scope atau applicable AOF requirement.

### Enabled

Control akan diadopsi/diimplementasikan dalam project scope.

### Deferred

Control applicable tetapi belum dapat diimplementasikan saat ini.

`Deferred` MUST tidak dianggap satisfied.

### NotApplicable

Control benar-benar tidak applicable terhadap declared project scope.

```text
NotApplicable != DisabledForConvenience
```

### Unsupported

Organisasi/tooling saat ini tidak memiliki capability yang diperlukan.

`Unsupported` MUST menjadi visibility signal, bukan silent omission.

---

# 12. Applicability vs Organizational Capability

AOF CLI MUST membedakan:

```text
Applicability
```

dengan:

```text
OrganizationalCapability
```

Contoh:

| Control | Applicability | Capability | Result |
|---|---|---|---|
| Authority | Required | Supported | Enabled |
| Trace | Required | Supported | Enabled |
| Evidence | Applicable | Limited | Deferred |
| Independent Verification | Optional | Unsupported | Deferred |
| AgentInteractionContract | N/A | Supported | NotApplicable |

Ini mencegah:

```text
"We cannot implement it"
```

disalahartikan sebagai:

```text
"It is not applicable."
```

---

# 13. Scoped Conformance Rule

AOF CLI MUST menjaga:

```text
ScopedAdoption -> ScopedConformance
```

dan MUST mencegah implication:

```text
PartialImplementation -> FullAOFConformance
```

Generated documentation SHOULD secara eksplisit menyatakan declared adoption profile dan scope.

---

# 14. Project Discovery

`aof init` MUST mempunyai project-aware initialization.

Primary phases:

```text
DISCOVER
    |
    v
MODEL
    |
    v
GENERATE
```

## 14.1 Deterministic Discovery

CLI MAY mendeteksi language/tooling dari file seperti:

```text
go.mod          -> Go
package.json    -> Node.js / JavaScript / TypeScript
pyproject.toml  -> Python
Cargo.toml      -> Rust
pom.xml         -> Java / Maven
build.gradle    -> Java / Gradle
```

Discovery MUST bersifat conservative.

CLI MUST tidak menginvent domain hanya berdasarkan nama directory apabila confidence rendah.

## 14.2 User Input

Project kosong hampir tidak memiliki discoverable context.

Interactive mode MUST dapat meminta minimal:

- Project name;
- Project type;
- Language;
- Domain;
- selected AOF profile;
- AI usage;
- consequential automation usage;
- organizational control capability.

---

# 15. CLI Modes

## 15.1 Interactive

Default:

```bash
aof init
```

Example:

```text
AOF CLI

Project name:
> itsm-backend

Project type:
> Backend API

Language:
> Go

Domain:
> IT Service Management

Select AOF adoption profile:
> AOF-Core
  AOF-Governed
  AOF-Assured

Will this project use AI agents?
> Yes

Will AI be allowed to perform consequential external actions?
> No
```

## 15.2 Non-Interactive

CLI SHOULD menyediakan non-interactive input untuk automation.

Conceptual syntax:

```bash
aof init \
  --name itsm-backend \
  --language go \
  --type backend \
  --domain itsm \
  --profile core \
  --non-interactive
```

Exact flag set MAY disederhanakan pada implementation selama memenuhi product behavior.

---

# 16. AI Usage Discovery

AOF CLI SHOULD mengetahui bagaimana AI digunakan.

Example selection:

```text
AI use cases:

[x] classification
[x] analysis
[x] recommendation
[ ] automatic state mutation
[ ] external action
[ ] production infrastructure action
```

Ini mempengaruhi:

- project instructions;
- applicable controls;
- Agent Skills;
- risk guidance;
- generated roles;
- execution guidance.

---

# 17. Generated Project Model

AOF CLI SHOULD menghasilkan:

```text
.aof/project.yaml
```

Example:

```yaml
project:
  name: itsm-backend
  language: go
  type: backend
  domain: itsm

architecture:
  style: unspecified

features: []

ai:
  enabled: true
  use_cases:
    - classification
    - recommendation
  consequential_execution: false
```

Project model merupakan dynamic source untuk generated context.

---

# 18. AOF Adoption Configuration

Generated:

```text
.aof/config.yaml
```

Example:

```yaml
aof:
  specification: "1.0"
  release: "LTS"
  profile: "AOF-Core"

upstream:
  repository: "https://github.com/aof-framework/aof"
  specification: "specification/AOF-v1.0-Framework-Specification.md"

adoption:
  mode: "greenfield"

controls:
  authority:
    status: enabled

  policy:
    status: enabled

  trace:
    status: enabled

  evidence:
    status: deferred

  verification:
    status: deferred

  agent_interaction_contract:
    status: not_applicable
```

---

# 19. Manifest

Generated:

```text
.aof/manifest.json
```

Purpose:

- detect previous initialization;
- record CLI version;
- record AOF baseline;
- record generated files;
- support idempotency;
- provide future migration provenance.

Example:

```json
{
  "aof_specification": "1.0",
  "aof_release": "LTS",
  "aof_cli_version": "0.1.0",
  "upstream_repository": "https://github.com/aof-framework/aof",
  "adoption_mode": "greenfield",
  "generated_files": [
    "AOF.md",
    ".aof/project.yaml",
    ".agents/skills/aof-core/SKILL.md"
  ]
}
```

Sensitive information MUST NOT disimpan.

---

# 20. Repository Context Architecture

Generated project context SHOULD mengikuti:

```text
                 CANONICAL AOF
                      |
                      v
                  AOF.md
                      |
             +--------+--------+
             |                 |
             v                 v
         AGENTS.md       .agents/skills/
                               |
                               v
                        .agents/agents/
                               |
                               v
                       AI Coding Tools
```

AOF canonical project knowledge MUST berada di layer netral, bukan tool-specific layer.

---

# 21. `AGENTS.md`

`AGENTS.md` adalah dynamic bootstrap instruction untuk AI coding agents.

It SHOULD berisi:

- project description;
- architecture summary;
- AOF adoption declaration;
- selected profile;
- project scope;
- location of `AOF.md`;
- location of AOF skills;
- critical semantic boundaries;
- project-specific constraints;
- tests/commands bila discoverable.

It MUST NOT copy full AOF specification.

Conceptual template:

```markdown
# Repository Instructions

This repository implements <project>.

## Project

Language: <language>
Type: <type>
Domain: <domain>

## AOF

This project adopts AOF v1.0 LTS using profile <profile>.

Read `AOF.md` before modifying:
- AI automation
- authority
- approval
- risk
- consequential execution
- verification

Preserve:
- Capability != Authority
- Proposal != AuthorizedDecision
- Pending != Pass
- Claim != Evidence != Verification
```

---

# 22. `AOF.md`

`AOF.md` adalah **dynamic project adoption entry point**.

It SHOULD menjawab:

> How does this repository adopt AOF?

Suggested sections:

```text
AOF baseline
Adoption profile
Adoption scope
AI usage
Governed areas
Consequential actions
Authority posture
Risk posture
Enabled controls
Deferred controls
NotApplicable controls
Canonical references
```

Example:

```markdown
# AOF Adoption — itsm-backend

This repository adopts AOF v1.0 LTS.

Canonical upstream:
https://github.com/aof-framework/aof

Profile:
AOF-Core

AI agents are advisory only.

They may:
- analyze
- classify
- recommend
- produce proposals

They may not:
- grant themselves authority
- approve their own consequential actions
- directly execute production changes
```

---

# 23. `.agents/skills/`

Agent Skills dibagi menjadi static dan dynamic.

Example:

```text
.agents/skills/
├── aof-core/
│   └── SKILL.md
├── aof-governance/
│   └── SKILL.md
├── aof-conformance/
│   └── SKILL.md
├── project-architecture/
│   └── SKILL.md
└── project-domain/
    └── SKILL.md
```

Classification:

```text
aof-core             -> static
aof-governance       -> static
aof-conformance      -> parameterized
project-architecture -> dynamic
project-domain       -> dynamic
```

Not every project MUST receive every skill.

Install only applicable skills.

---

# 24. `.agents/agents/`

Agent roles SHOULD bersifat responsibility-based.

Default archetypes MAY include:

```text
architect
implementer
reviewer
aof-reviewer
```

Project-specific agents MAY ditambahkan berdasarkan project model.

Agent definitions MUST tidak menduplikasi seluruh AOF knowledge.

Preferred relation:

```text
Agent -> uses Skill
Agent != Skill
```

---

# 25. Tool-Specific Configuration

MVP v0.1 SHOULD prioritize portable artifacts:

```text
AGENTS.md
AOF.md
.agents/
```

Tool-specific adapters seperti `.opencode/` SHOULD NOT menjadi mandatory v0.1 baseline.

Future adapters MUST bergantung pada portable project knowledge.

```text
Project Knowledge
      |
      v
Agent Skills
      |
      v
Tool Adapter
```

Deleting a tool adapter MUST NOT menghapus AOF knowledge project.

---

# 26. Expected Generated Structure

Baseline conceptual output:

```text
project/
├── AGENTS.md
├── AOF.md
├── .aof/
│   ├── project.yaml
│   ├── config.yaml
│   └── manifest.json
├── .agents/
│   ├── agents/
│   │   ├── architect/
│   │   │   └── agent.md
│   │   ├── implementer/
│   │   │   └── agent.md
│   │   ├── reviewer/
│   │   │   └── agent.md
│   │   └── aof-reviewer/
│   │       └── agent.md
│   └── skills/
│       ├── aof-core/
│       │   └── SKILL.md
│       ├── aof-governance/
│       │   └── SKILL.md
│       └── aof-conformance/
│           └── SKILL.md
├── docs/
│   └── aof/
│       ├── README.md
│       ├── semantic-boundaries.md
│       ├── governance-model.md
│       └── execution-model.md
└── aof/
    ├── policies/
    ├── profiles/
    ├── fixtures/
    └── conformance/
```

Empty directories MAY memerlukan `.gitkeep` apabila harus di-version-control.

Exact output SHOULD mempertimbangkan hanya artifacts yang applicable.

---

# 27. CLI Surface — v0.1

MVP command surface:

```text
aof init
aof version
aof help
```

No additional mandatory commands.

---

# 28. `aof init`

Primary command.

Lifecycle:

```text
Resolve Project Root
        |
        v
Inspect Existing Files
        |
        v
Discover Project
        |
        v
Collect User Input
        |
        v
Build Project Model
        |
        v
Build Adoption Model
        |
        v
Determine Applicable AOF Subset
        |
        v
Build Write Plan
        |
        v
Detect Conflicts
        |
        v
Generate Artifacts
        |
        v
Write Manifest
        |
        v
Report Result
```

---

# 29. `aof version`

Expected:

```text
AOF CLI v0.1.0
AOF Specification v1.0 LTS
Canonical Upstream: https://github.com/aof-framework/aof
```

AOF CLI version MUST berbeda lifecycle dari AOF specification version.

```text
AOF Specification Version != AOF CLI Version
```

---

# 30. `aof help`

Minimal:

```text
AOF CLI

Usage:
  aof init
  aof version
  aof help
```

---

# 31. Filesystem Safety

Because AOF CLI modifies developer repositories, filesystem safety adalah hard requirement.

AOF CLI MUST:

- restrict writes to resolved project root;
- normalize paths;
- reject template path traversal;
- never write using untrusted absolute template paths;
- inspect conflicts before mutation;
- avoid silent overwrite;
- not execute generated scripts;
- not execute repository code;
- not run package-manager install commands;
- not transmit project contents;
- not collect telemetry.

---

# 32. Conflict Handling

Preferred MVP model:

```text
INSPECT
  |
BUILD PLAN
  |
DETECT CONFLICT
  |
  +-- conflict -> ABORT BEFORE WRITES
  |
  +-- no conflict -> WRITE
```

This is preferable to clever merging in v0.1.

Examples of conflict:

```text
AGENTS.md already exists
AOF.md already exists
.aof/manifest.json exists unexpectedly
target generated skill file exists
```

CLI SHOULD report all detected conflicts in one run when possible.

Example:

```text
Initialization cannot continue safely.

Conflicts:
- AGENTS.md already exists
- .agents/skills/aof-core/SKILL.md already exists

No files were changed.
```

Future version MAY implement explicit merge/import behavior.

---

# 33. Existing AOF Initialization

If:

```text
.aof/manifest.json
```

indicates project already initialized, running `aof init` again SHOULD return a clear non-destructive result.

Example:

```text
AOF is already initialized in this project.
AOF Specification: v1.0 LTS
AOF CLI Bootstrap: v0.1.0
```

---

# 34. Idempotency

Practical target:

```text
Init(Init(Project)) = Init(Project)
```

A second initialization MUST NOT:

- duplicate generated sections;
- regenerate arbitrary different content;
- overwrite user modifications;
- create multiple conflicting manifests.

---

# 35. Embedded Template Architecture

Suggested source repository:

```text
aof-cli/
├── AGENTS.md
├── README.md
├── go.mod
├── cmd/
│   └── aof/
│       └── main.go
├── internal/
│   ├── cli/
│   ├── discovery/
│   ├── model/
│   ├── adoption/
│   ├── initializer/
│   └── template/
├── templates/
│   ├── static/
│   │   ├── docs/
│   │   └── skills/
│   └── dynamic/
│       ├── AGENTS.md.tmpl
│       ├── AOF.md.tmpl
│       ├── project.yaml.tmpl
│       └── agents/
└── tests/
```

Avoid premature abstraction.

Packages SHOULD exist only where responsibilities justify them.

---

# 36. Internal Domain Model

Conceptual Go structures:

```go
type ProjectDefinition struct {
    Name       string
    Language   string
    Type       string
    Domain     string
    Features   []string
    AI         AIUsage
}

type AdoptionDefinition struct {
    Profile  string
    Mode     string
    Controls map[string]ControlSelection
}

type ControlSelection struct {
    Applicability string
    Capability    string
    Status        string
    Reason        string
}
```

Exact structs MAY differ.

Key requirement: configuration MUST represent distinction between applicability, capability, and selected status.

---

# 37. No LLM Dependency

Dynamic generation MUST NOT imply LLM dependency.

MVP MUST be deterministic.

```text
User Input
+
Project Discovery
+
Rules
+
Templates
=
Generated Context
```

No:

- API key;
- external model call;
- cloud inference;
- prompt generation service.

AOF CLI itself SHOULD remain usable in offline environments.

---

# 38. Generic MVP, Not Domain Pack MVP

AOF CLI v0.1 SHOULD remain generic.

Do not include mandatory:

```text
Go pack
ITSM pack
E-commerce pack
Healthcare pack
Finance pack
```

The CLI MAY collect domain metadata and generate generic domain placeholders.

Domain packs are a future ecosystem feature.

---

# 39. Example: Empty Go ITSM Project

Input:

```bash
mkdir itsm-backend
cd itsm-backend
aof init
```

Interactive answers:

```text
Project name: itsm-backend
Project type: Backend API
Language: Go
Domain: IT Service Management
Profile: AOF-Core
AI agents: Yes
Consequential AI execution: No
Independent verification capability: Deferred
```

Output `AOF.md` could describe:

```text
Project: itsm-backend
Domain: IT Service Management
Profile: AOF-Core

AI agents may:
- analyze
- classify
- recommend
- produce proposals

AI agents may not:
- approve their own consequential actions
- self-grant authority
- directly execute production changes
```

This content is dynamic project adoption context, not canonical AOF semantics.

---

# 40. Example: Minimal Organization

Organization capability:

```text
Basic Trace                    Supported
Explicit Authority            Supported
Formal Evidence Lifecycle     Unsupported
Independent Verification      Unsupported
Formal Conformance Reporting  Unsupported
```

AOF CLI SHOULD not reject adoption solely because high-assurance capability tidak tersedia.

Instead it may generate:

```yaml
controls:
  authority:
    status: enabled

  trace:
    status: enabled

  evidence:
    status: deferred
    reason: organizational capability not yet available

  independent_verification:
    status: deferred
    reason: organizational capability not yet available
```

Generated AOF documentation MUST tidak mengklaim controls tersebut sebagai satisfied.

---

# 41. UX Requirements

AOF CLI SHOULD:

- use concise language;
- explain why a question matters;
- provide safe defaults;
- allow user to review before write;
- show final summary;
- clearly differentiate detected values from user selections;
- not overwhelm beginner users with all canonical schema names.

Advanced details MAY be represented in generated docs rather than wizard questions.

---

# 42. Proposed Interactive Wizard

Suggested sequence:

```text
Step 1 — Project
Step 2 — AOF Profile
Step 3 — AI Usage
Step 4 — Consequential Effects
Step 5 — Organizational Capability
Step 6 — Review
Step 7 — Generate
```

Example:

```text
AOF CLI v0.1

Project
-------
Name: itsm-backend
Language: Go
Type: Backend API
Domain: IT Service Management

AOF Adoption
------------
Profile:
> AOF-Core
  AOF-Governed
  AOF-Assured

AI Usage
--------
[x] analysis
[x] classification
[x] recommendation
[ ] consequential automation

Organizational Capability
-------------------------
[x] authority boundaries
[x] policy rules
[x] basic trace
[ ] independent verification
[ ] formal evidence lifecycle

Review
------
Profile: AOF-Core
Applicable controls: 8
Enabled: 6
Deferred: 2

Generate AOF project context?
> Yes
```

---

# 43. Generated Content Must Be Explainable

AOF CLI SHOULD enable developer to understand:

```text
Why was this file created?
Why was this control selected?
Why was this skill installed?
```

Generated files SHOULD contain concise provenance comments or references where appropriate.

---

# 44. Source Provenance

Generated `AOF.md`, config, or manifest SHOULD reference:

```text
https://github.com/aof-framework/aof
```

and the AOF baseline used.

Example:

```text
AOF Baseline: v1.0 LTS
Canonical Repository: https://github.com/aof-framework/aof
```

Future builds SHOULD additionally record exact tag/commit/checksum used to produce embedded static assets.

---

# 45. Security Requirements

MVP MUST:

- perform no telemetry;
- perform no repository upload;
- require no network during `aof init`;
- execute no downloaded code;
- execute no project code;
- execute no generated script;
- enforce project-root write boundary;
- reject path traversal;
- preserve existing content;
- not store secrets;
- not infer or persist credentials.

---

# 46. Privacy Requirements

AOF CLI v0.1:

```text
Data Collection: None
Telemetry: None
Remote Upload: None
LLM Calls: None
```

Project discovery runs locally.

---

# 47. Performance Requirements

For normal repositories, `aof init` SHOULD feel effectively immediate.

Target behavior:

- no Git clone;
- no package download;
- no remote API call;
- small embedded bootstrap asset set;
- filesystem-only generation.

---

# 48. Test Strategy

## 48.1 CLI Tests

- `aof help`;
- `aof version`;
- unknown command behavior;
- exit codes.

## 48.2 Initialization Tests

- empty directory;
- existing unrelated files;
- existing `AGENTS.md`;
- existing `AOF.md`;
- existing `.agents/`;
- existing `.aof/`;
- initialized project;
- non-writable target where testable;
- partial conflict set.

## 48.3 Discovery Tests

- Go project;
- Node project;
- Python project;
- unknown project;
- empty project.

## 48.4 Model Tests

- AOF-Core;
- AOF-Governed;
- AOF-Assured;
- Enabled;
- Deferred;
- NotApplicable;
- Unsupported;
- Required.

## 48.5 Template Tests

- all embedded templates accessible;
- no escaping project root;
- expected variables rendered;
- deterministic render.

## 48.6 Manifest Tests

- valid JSON;
- expected version;
- expected upstream;
- generated paths correct.

## 48.7 Safety Tests

- path traversal rejected;
- existing file not overwritten;
- conflict abort produces no partial writes;
- templates cannot write outside project root.

---

# 49. Engineering Quality Gate

Before v0.1 completion:

```bash
go fmt ./...
go vet ./...
go test ./...
```

MUST pass.

The project SHOULD additionally avoid unnecessary third-party dependencies.

---

# 50. Acceptance Criteria

AOF CLI v0.1 is complete when all of the following are true:

```text
[ ] Single Go binary can be built
[ ] aof init works in empty directory
[ ] aof init works safely in existing repository
[ ] project discovery works for supported basic signals
[ ] interactive project model can be created
[ ] AOF profile can be selected
[ ] adoption scope can be declared
[ ] organizational capability can be declared
[ ] Required/Enabled/Deferred/NotApplicable/Unsupported are represented
[ ] static AOF assets are embedded
[ ] dynamic project artifacts are generated
[ ] AGENTS.md is project-aware
[ ] AOF.md is project-aware
[ ] .agents/ skills are selectively generated
[ ] .agents/ agents are responsibility-based
[ ] .aof/project.yaml is generated
[ ] .aof/config.yaml is generated
[ ] .aof/manifest.json is generated
[ ] canonical upstream repository is recorded
[ ] init requires no network
[ ] init performs no Git clone
[ ] existing user files are never silently overwritten
[ ] second init is safe
[ ] path traversal protection exists
[ ] no telemetry exists
[ ] no application runtime dependency exists
[ ] go fmt passes
[ ] go vet passes
[ ] go test ./... passes
```

---

# 51. Definition of Done

The release SHOULD be considered `v0.1.0 Complete Bootstrap MVP` when:

```text
One binary
+
Safe initialization
+
Project discovery
+
Progressive adoption configuration
+
Static/dynamic generation
+
Portable agent context
+
Canonical AOF provenance
+
Tests
```

are complete.

---

# 52. Initial Release Scope

Commands:

```text
aof init
aof version
aof help
```

No additional command is required for v0.1.

---

# 53. Explicitly Deferred Features

Potential v0.2+:

```text
aof doctor
aof validate
aof upgrade
aof sync
aof add
tool-specific adapters
language packs
domain packs
policy packs
conformance execution integration
migration between bootstrap versions
online immutable release source
signature/checksum verification
```

These MUST NOT expand v0.1 scope unless required to satisfy a v0.1 acceptance criterion.

---

# 54. Future `aof doctor`

Potential role:

```text
AOF Project Health

[PASS] AOF.md
[PASS] project model
[PASS] canonical baseline
[PASS] Agent Skills
[WARN] Verification deferred
```

Not part of v0.1.

---

# 55. Future Domain Packs

Future architecture MAY support:

```text
AOF Core
+
Language Pack
+
Domain Pack
+
Tool Adapter
```

Example:

```text
AOF Core
+
Go
+
ITSM
+
OpenCode adapter
```

But:

```text
DomainPack != AOFCore
ToolAdapter != AOFProjectKnowledge
```

---

# 56. Future Upgrade Model

Future `aof upgrade` MUST distinguish:

```text
AOF Specification Version
AOF CLI Version
Bootstrap Template Version
Project Adoption Model Version
```

No upgrade may silently change frozen AOF v1.0 semantics.

---

# 57. Recommended GitHub Repository Separation

Canonical framework:

```text
aof-framework/aof
```

CLI:

```text
aof-framework/aof-cli
```

Relationship:

```text
aof-framework/aof
        |
        | canonical specification
        | schemas
        | conformance
        | release artifacts
        v
    AOF v1.0 LTS

aof-framework/aof-cli
        |
        | adoption bootstrap tooling
        v
       aof
        |
        +--> aof init
```

AOF CLI repository MUST not become an alternate AOF specification.

---

# 58. Recommended README Positioning

Suggested short description:

> **AOF CLI is the official repository bootstrap utility for adopting the AI Orchestration Framework (AOF).**

Suggested clarification:

> AOF CLI configures project context and developer-facing adoption artifacts. It is not an AOF runtime and is not required by applications after initialization.

---

# 59. Recommended `AGENTS.md` for the `aof-cli` Repository

The `aof-cli` source repository SHOULD itself contain an `AGENTS.md` stating:

```text
AOF CLI is a bootstrap tool.

Before changing product behavior, read:
docs/AOF-CLI-v0.1-PRD-and-Technical-Design.md

Canonical AOF:
https://github.com/aof-framework/aof

Do not:
- redefine AOF semantics
- add runtime responsibilities
- make network mandatory for init
- clone mutable main during init
- silently overwrite project content
- add false conformance claims

Primary product:
aof init
```

---

# 60. Open Implementation Decisions

The following details MAY be decided during coding while preserving this document:

- exact CLI parser implementation;
- whether to use standard `flag` package or minimal custom parser;
- exact YAML library choice;
- exact text UI formatting;
- exact internal package boundaries;
- exact input question wording;
- whether `.gitkeep` is used for empty generated directories;
- exact template engine (`text/template` is preferred unless insufficient).

Changes to these do not require product redesign.

---

# 61. Frozen Product Decisions for v0.1

The following SHOULD be treated as baseline decisions:

1. `aof-cli` is implemented in Go.
2. It is distributed as a single executable.
3. v0.1 primary command is `aof init`.
4. `aof version` and `aof help` are included.
5. AOF v1.0 LTS is the framework baseline.
6. `https://github.com/aof-framework/aof` is canonical upstream.
7. `aof init` does **not** clone mutable `main`.
8. Static AOF bootstrap content is embedded in the CLI.
9. Project adoption context is dynamic.
10. `AGENTS.md` is dynamic.
11. `AOF.md` is dynamic.
12. `.agents/` contains static and dynamic components selectively.
13. AOF adoption is progressive and scoped.
14. Not all AOF capabilities need to be implemented by every organization.
15. Reduced scope MUST not create false conformance claims.
16. Applicability and organizational capability are distinct.
17. `Deferred` does not mean satisfied.
18. `NotApplicable` is not a convenience bypass.
19. CLI initialization is local and offline-capable.
20. CLI performs no telemetry.
21. CLI is not a runtime dependency.
22. Existing user content is never silently overwritten.
23. v0.1 favors safe abort over automatic merge.
24. No LLM dependency exists in v0.1.
25. Tool-specific AI integrations are adapters, not canonical project knowledge.

---

# 62. Product Success Definition

The product succeeds when a developer who has never manually configured AOF can run:

```bash
aof init
```

and receive a repository that:

- identifies its AOF baseline;
- declares an appropriate adoption profile and scope;
- accurately reflects organizational capability;
- provides relevant AOF knowledge to humans and AI coding agents;
- avoids irrelevant AOF context;
- preserves canonical AOF semantics;
- prevents misleading conformance implications;
- remains usable without the AOF CLI;
- is portable across coding-agent ecosystems.

In compact form:

```text
AOFCLI
=
AdoptionConfigurator
+
ProjectAwareBootstrapGenerator
```

and:

```text
AOFReadyProject
=
StaticApplicableAOF
+
DynamicProjectContext
+
ExplicitAdoptionScope
```

---

# 63. Final Architecture

```text
                     AOF v1.0 LTS
                          |
                   Canonical Upstream
                          |
           https://github.com/aof-framework/aof
                          |
                   pinned at build time
                          |
                          v
                   Static AOF Assets
                          |
                          |
Existing Repository -----+------ User Choices
          |                       |
          v                       v
     Discovery              Adoption Input
          |                       |
          +-----------+-----------+
                      |
                      v
                Project Model
                      |
                      v
               Adoption Model
                      |
                      v
             Applicability Engine
                      |
                      v
             Generation Plan
                      |
              conflict check
                      |
                      v
               AOF-ready Repo
                      |
       +--------------+---------------+
       |              |               |
       v              v               v
    AOF.md        AGENTS.md       .agents/
       |
       v
  .aof/config
```

---

# 64. Implementation Starting Point

Implementation SHOULD begin in this order:

```text
1. Repository bootstrap
2. Version/help command
3. Embedded static assets
4. Project discovery
5. Project model
6. Adoption profile model
7. Control applicability/capability model
8. Interactive init
9. Generation plan
10. Conflict detection
11. Static artifact generation
12. Dynamic artifact generation
13. Manifest/config generation
14. Idempotency
15. Safety tests
16. End-to-end init tests
17. README and release documentation
```

Do not expand scope until this sequence is complete.

---

# 65. Canonical References

## AOF Framework Repository

https://github.com/aof-framework/aof

## AOF v1.0 Framework Specification

https://github.com/aof-framework/aof/blob/main/specification/AOF-v1.0-Framework-Specification.md

AOF CLI implementation MUST treat the AOF Framework repository and its versioned release artifacts as canonical upstream sources.

For v0.1, canonical content used by the CLI SHOULD be pinned and embedded at build time rather than cloned dynamically during `aof init`.

---

# 66. Closing Statement

AOF CLI exists to reduce the distance between:

```text
"I want to adopt AOF"
```

and:

```text
"This repository is ready to develop under an explicit AOF adoption model."
```

The CLI MUST preserve a clean boundary:

```text
AOF Framework
     |
     | defines semantics
     v
AOF CLI
     |
     | configures adoption
     v
Software Project
```

AOF CLI SHOULD make adoption easy.

It MUST NOT make AOF semantics ambiguous.

It SHOULD make partial, realistic adoption possible.

It MUST NOT convert partial adoption into false assurance.

The central v0.1 product contract is:

```text
aof init
->
Safe, Scoped, ProjectAware, AOFReadyRepository
```
