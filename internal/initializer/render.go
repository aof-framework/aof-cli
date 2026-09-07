package initializer

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/aof-framework/aof-cli/internal/model"
	bootstraptemplates "github.com/aof-framework/aof-cli/templates"
)

type File struct {
	Path    string
	Content []byte
}
type Plan struct{ Files []File }

const (
	canonicalRepo               = "https://github.com/aof-framework/aof"
	canonicalSpecPath           = "specification/AOF-v1.0-Framework-Specification.md"
	canonicalSpecSHA256         = "57ddbd64671eea615535b20f109064d96fb262e781969ef757a6f4d5efa869d5"
	canonicalSchemaBundleSHA256 = "d9a433c5d5549e61600eb34d84fe8e4d9802b4211542765fb810c8a56baeb975"
)

func BuildPlan(m model.BootstrapModel, cliVersion string) (Plan, error) {
	files := []File{}
	for _, item := range []struct{ src, dst string }{
		{"static/docs/semantic-boundaries.md", "docs/aof/semantic-boundaries.md"},
		{"static/docs/governance-model.md", "docs/aof/governance-model.md"},
		{"static/docs/execution-model.md", "docs/aof/execution-model.md"},
	} {
		b, err := fs.ReadFile(bootstraptemplates.FS, item.src)
		if err != nil {
			return Plan{}, err
		}
		files = append(files, File{Path: item.dst, Content: b})
	}
	for _, skill := range m.Adoption.RequiredSkills {
		src := "static/skills/" + skill + "/SKILL.md"
		b, err := fs.ReadFile(bootstraptemplates.FS, src)
		if err != nil {
			return Plan{}, fmt.Errorf("read embedded skill %s: %w", src, err)
		}
		files = append(files, File{Path: ".agents/skills/" + skill + "/SKILL.md", Content: b})
	}
	for _, item := range []struct{ src, dst string }{
		{"dynamic/AGENTS.md.tmpl", "AGENTS.md"}, {"dynamic/AOF.md.tmpl", "AOF.md"},
		{"dynamic/project-architecture.SKILL.md.tmpl", ".agents/skills/project-architecture/SKILL.md"},
		{"dynamic/project-domain.SKILL.md.tmpl", ".agents/skills/project-domain/SKILL.md"},
		{"dynamic/agent-architect.md.tmpl", ".agents/agents/architect/agent.md"},
		{"dynamic/agent-implementer.md.tmpl", ".agents/agents/implementer/agent.md"},
		{"dynamic/agent-reviewer.md.tmpl", ".agents/agents/reviewer/agent.md"},
		{"dynamic/agent-aof-reviewer.md.tmpl", ".agents/agents/aof-reviewer/agent.md"},
	} {
		b, err := fs.ReadFile(bootstraptemplates.FS, item.src)
		if err != nil {
			return Plan{}, err
		}
		t, err := template.New(item.src).Funcs(template.FuncMap{"join": strings.Join}).Parse(string(b))
		if err != nil {
			return Plan{}, err
		}
		var out bytes.Buffer
		if err := t.Execute(&out, m); err != nil {
			return Plan{}, err
		}
		files = append(files, File{Path: item.dst, Content: out.Bytes()})
	}

	generated := []File{
		{Path: ".aof/project.yaml", Content: []byte(renderProjectYAML(m.Project))},
		{Path: ".aof/config.yaml", Content: []byte(renderConfigYAML(m.Adoption))},
		{Path: ".aof/adoption.yaml", Content: []byte(renderAdoptionYAML(m.Adoption))},
		{Path: ".aof/applicability.yaml", Content: []byte(renderApplicabilityYAML(m.Adoption))},
		{Path: ".aof/requirements.yaml", Content: []byte(renderRequirementsYAML(m.Adoption))},
		{Path: ".aof/capabilities.yaml", Content: []byte(renderCapabilityYAML(m.Project))},
		{Path: ".aof/provenance.json", Content: []byte(renderProvenanceJSON())},
		{Path: "aof/architecture/planes.yaml", Content: []byte(renderPlanes(m))},
		{Path: "aof/architecture/trust-boundaries.yaml", Content: []byte(renderTrustBoundaries(m))},
		{Path: "aof/architecture/effect-boundaries.yaml", Content: []byte(renderEffectBoundaries(m))},
		{Path: "aof/control/safety-kernel.yaml", Content: []byte(renderSafetyKernel(m))},
		{Path: "aof/control/gate-mapping.yaml", Content: []byte(renderGateMapping(m))},
		{Path: "aof/governance/governance-root.md", Content: []byte(renderGovernanceRoot(m))},
		{Path: "aof/governance/governance-envelope.yaml", Content: []byte(renderGovernanceEnvelope(m))},
		{Path: "aof/governance/human-governance.yaml", Content: []byte(renderHumanGovernance(m))},
		{Path: "aof/agents/inventory.yaml", Content: []byte(renderAgentInventory(m))},
		{Path: "aof/agents/capabilities.yaml", Content: []byte(renderAgentCapabilities(m))},
		{Path: "aof/agents/interactions.yaml", Content: []byte(renderAgentInteractions(m))},
		{Path: "aof/authority/authority-model.yaml", Content: []byte(renderAuthorityModel(m))},
		{Path: "aof/authority/README.md", Content: []byte(renderAuthorityReadme(m))},
		{Path: "aof/policies/policy-model.yaml", Content: []byte(renderPolicyModel(m))},
		{Path: "aof/risk/risk-model.yaml", Content: []byte(renderRiskModel(m))},
		{Path: "aof/execution/execution-model.md", Content: []byte(renderExecutionModel(m))},
		{Path: "aof/execution/consequential-actions.yaml", Content: []byte(renderConsequentialActions(m))},
		{Path: "aof/execution/state-transitions.yaml", Content: []byte(renderStateTransitions(m))},
		{Path: "aof/execution/failure-recovery.yaml", Content: []byte(renderFailureRecovery(m))},
		{Path: "aof/assurance/evidence-model.md", Content: []byte(renderEvidenceModel(m))},
		{Path: "aof/assurance/verification-model.md", Content: []byte(renderVerificationModel(m))},
		{Path: "aof/assurance/trace-model.md", Content: []byte(renderTraceModel(m))},
		{Path: "aof/assurance/accountability-chain.yaml", Content: []byte(renderAccountability(m))},
		{Path: "aof/conformance/scope.yaml", Content: []byte(renderConformanceScope(m))},
		{Path: "aof/conformance/README.md", Content: []byte(renderConformanceReadme(m))},
		{Path: "aof/conformance/gaps.md", Content: []byte(renderGaps(m))},
		{Path: "aof/requirements/registry.yaml", Content: []byte(renderRequirementsYAML(m.Adoption))},
	}
	if contains(m.Adoption.Profile.DomainProfiles, model.DomainSecureSDLC) {
		generated = append(generated, File{Path: "aof/secure-sdlc/profile.yaml", Content: []byte(renderSecureSDLC(m))})
	}
	if contains(m.Adoption.Profile.Overlays, model.OverlayHighAssurance) {
		generated = append(generated, File{Path: "aof/high-assurance/overlay.yaml", Content: []byte(renderHighAssurance(m))})
	}
	files = append(files, generated...)

	// Embed the pinned canonical JSON Schema bundle into generated repositories.
	err := fs.WalkDir(bootstraptemplates.FS, "static/schemas", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(bootstraptemplates.FS, path)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, "static/schemas/")
		files = append(files, File{Path: "aof/schemas/" + rel, Content: b})
		return nil
	})
	if err != nil {
		return Plan{}, fmt.Errorf("walk embedded schemas: %w", err)
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	manifest, err := renderManifest(m, cliVersion, files)
	if err != nil {
		return Plan{}, err
	}
	files = append(files, File{Path: ".aof/bootstrap-manifest.json", Content: manifest})
	return Plan{Files: files}, nil
}

func renderProjectYAML(p model.ProjectDefinition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "project:\n  name: %q\n  language: %q\n  type: %q\n  domain: %q\n  architecture: %q\n  deployment: %q\n  criticality: %q\n  data_sensitivity: %q\n  adoption_mode: %q\n", p.Name, p.Language, p.Type, p.Domain, p.Architecture, p.Deployment, p.Criticality, p.DataSensitivity, p.AdoptionMode)
	writeList(&b, "features", p.Features)
	writeList(&b, "agent_types", p.AgentTypes)
	fmt.Fprintf(&b, "multi_agent: %t\nai:\n  enabled: %t\n  analysis: %t\n  classification: %t\n  recommendation: %t\n  tool_use: %t\n  state_mutation: %t\n  external_action: %t\n  infrastructure_action: %t\n  consequential_execution: %t\n", p.MultiAgent, p.AI.Enabled, p.AI.Analysis, p.AI.Classification, p.AI.Recommendation, p.AI.ToolUse, p.AI.StateMutation, p.AI.ExternalAction, p.AI.InfrastructureAction, p.AI.ConsequentialExecution)
	writeList(&b, "effect_types", p.EffectTypes)
	writeList(&b, "risk_categories", p.RiskCategories)
	return b.String()
}

func renderConfigYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "aof:\n  specification: \"1.0\"\n  release: \"LTS\"\n  target_base_profile: %q\n", a.Profile.TargetBaseProfile)
	writeListIndented(&b, "  domain_profiles", a.Profile.DomainProfiles)
	writeListIndented(&b, "  overlays", a.Profile.Overlays)
	b.WriteString("  claimed_profile: null\n  conformance_claim_status: \"none\"\n")
	fmt.Fprintf(&b, "upstream:\n  repository: %q\n  specification: %q\n  specification_sha256: %q\n  schema_bundle_sha256: %q\nadoption:\n  mode: %q\n", canonicalRepo, canonicalSpecPath, canonicalSpecSHA256, canonicalSchemaBundleSHA256, a.Mode)
	return b.String()
}
func renderAdoptionYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "target_base_profile: %q\n", a.Profile.TargetBaseProfile)
	writeList(&b, "domain_profiles", a.Profile.DomainProfiles)
	writeList(&b, "overlays", a.Profile.Overlays)
	b.WriteString("claimed_profile: null\nclaim_status: \"none\"\nsemantics:\n  applicability_is_implementation: false\n  adoption_is_implementation: false\n  implementation_is_verification: false\n  partial_implementation_is_full_conformance: false\n")
	writeList(&b, "required_skills", a.RequiredSkills)
	writeList(&b, "warnings", a.Warnings)
	return b.String()
}
func renderApplicabilityYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	b.WriteString("control_families:\n")
	for _, k := range sortedControlKeys(a.Controls) {
		c := a.Controls[k]
		fmt.Fprintf(&b, "  %s:\n    applicability: %q\n    normative_level: %q\n    capability: %q\n    adoption_state: %q\n    implementation_state: %q\n    verification_state: %q\n    reason: %q\n", k, c.Applicability, c.NormativeLevel, c.Capability, c.AdoptionState, c.ImplementationState, c.VerificationState, c.Reason)
		writeListIndented(&b, "    requirement_ids", c.RequirementIDs)
	}
	b.WriteString("canonical_objects:\n")
	keys := make([]string, 0, len(a.CanonicalObjects))
	for k := range a.CanonicalObjects {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		o := a.CanonicalObjects[k]
		fmt.Fprintf(&b, "  %s:\n    applicability: %q\n    normative_level: %q\n    adoption_state: %q\n    reason: %q\n", k, o.Applicability, o.NormativeLevel, o.AdoptionState, o.Reason)
		writeListIndented(&b, "    requirement_ids", o.RequirementIDs)
	}
	return b.String()
}
func renderRequirementsYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	b.WriteString("requirements:\n")
	for _, r := range a.Requirements {
		fmt.Fprintf(&b, "  - requirement_id: %q\n    domain: %q\n    normative_level: %q\n    applies_to: %q\n    verification_method: %q\n", r.ID, r.Domain, r.NormativeLevel, r.AppliesTo, r.VerificationMethod)
		writeListIndented(&b, "    profiles", r.Profiles)
		writeListIndented(&b, "    required_evidence", r.RequiredEvidence)
	}
	return b.String()
}
func renderCapabilityYAML(p model.ProjectDefinition) string {
	var b strings.Builder
	writeList(&b, "agent_types", p.AgentTypes)
	b.WriteString("declared_ai_capabilities:\n")
	for _, x := range []struct {
		name string
		v    bool
	}{{"analysis", p.AI.Analysis}, {"classification", p.AI.Classification}, {"recommendation", p.AI.Recommendation}, {"tool_use", p.AI.ToolUse}, {"state_mutation", p.AI.StateMutation}, {"external_action", p.AI.ExternalAction}, {"infrastructure_action", p.AI.InfrastructureAction}, {"consequential_execution", p.AI.ConsequentialExecution}} {
		fmt.Fprintf(&b, "  %s: %t\n", x.name, x.v)
	}
	b.WriteString("note: \"Declared capability does not imply Authority or implementation verification.\"\n")
	return b.String()
}
func renderProvenanceJSON() string {
	return fmt.Sprintf("{\n  \"aof_specification\": \"1.0\",\n  \"aof_release\": \"LTS\",\n  \"canonical_repository\": %q,\n  \"canonical_specification_path\": %q,\n  \"canonical_specification_sha256\": %q,\n  \"canonical_schema_bundle_sha256\": %q,\n  \"source_commit\": null,\n  \"note\": \"Checksums pin the embedded LTS artifacts; no unverified Git commit is asserted.\"\n}\n", canonicalRepo, canonicalSpecPath, canonicalSpecSHA256, canonicalSchemaBundleSHA256)
}

func renderPlanes(m model.BootstrapModel) string {
	return "planes:\n  reasoning:\n    responsibility: \"analysis, planning, proposal generation\"\n    implementation_state: \"not_assessed\"\n  control:\n    responsibility: \"Authority, Policy, State, Risk, Verification mediation\"\n    implementation_state: \"not_assessed\"\n  effect:\n    responsibility: \"consequential external or persistent effects\"\n    implementation_state: \"not_assessed\"\n  assurance:\n    responsibility: \"Evidence, Verification, Trace, conformance\"\n    implementation_state: \"not_assessed\"\nseparation_required: true\n"
}
func renderTrustBoundaries(m model.BootstrapModel) string {
	return fmt.Sprintf("trust_boundaries:\n  reasoning_to_control:\n    input_trust: \"UntrustedProposal\"\n  control_to_effect:\n    consequential: %t\n    governed_mediation_required: %t\n  effect_to_assurance:\n    evidence_required: %t\n", m.Project.AI.ConsequentialExecution, m.Project.AI.ConsequentialExecution, m.Adoption.Controls["evidence"].Applicability == model.ApplicabilityApplicable)
}
func renderEffectBoundaries(m model.BootstrapModel) string {
	var b strings.Builder
	fmt.Fprintf(&b, "effect_boundary:\n  consequential_execution: %t\n  revalidation_required: %t\n  implementation_state: \"not_assessed\"\n", m.Project.AI.ConsequentialExecution, m.Project.AI.ConsequentialExecution)
	writeListIndented(&b, "  declared_effect_types", m.Project.EffectTypes)
	return b.String()
}
func renderSafetyKernel(m model.BootstrapModel) string {
	return "safety_kernel:\n  logical_components:\n    AuthorityEvaluator: {implementation_state: \"not_assessed\"}\n    PolicyEvaluator: {implementation_state: \"not_assessed\"}\n    StateValidator: {implementation_state: \"not_assessed\"}\n    RiskGate: {implementation_state: \"not_assessed\"}\n    VerificationGate: {implementation_state: \"not_assessed\"}\n    TraceRecorder: {implementation_state: \"not_assessed\"}\n  consequential_mediation_required: true\n  no_safety_by_prompt_alone: true\n  fail_open_for_mandatory_control: false\n"
}
func renderGateMapping(m model.BootstrapModel) string {
	return "execute_allowed:\n  expression: \"C AND H AND P AND S AND R AND V\"\n  C: \"Capability satisfied\"\n  H: \"Authority satisfied\"\n  P: \"Policy satisfied\"\n  S: \"State valid\"\n  R: \"Risk acceptable\"\n  V: \"Verification satisfied\"\nthree_valued_semantics: [Pass, Fail, Pending]\nmandatory_pending_is_pass: false\n"
}
func renderGovernanceRoot(m model.BootstrapModel) string {
	return fmt.Sprintf("# Governance Root — %s\n\nHuman/Organization remains the ultimate GovernanceRoot. Agent autonomy is bounded by the GovernanceEnvelope.\n\n- Declared operational Authority holder: **%s**\n- Target base profile: **%s**\n- Claimed conformance: **none**\n\nGenerated adoption artifacts are planning context, not proof of implementation.\n", m.Project.Name, m.Project.AuthorityHolder, m.Adoption.Profile.TargetBaseProfile)
}
func renderGovernanceEnvelope(m model.BootstrapModel) string {
	return fmt.Sprintf("governance_root: \"Human/Organization\"\nauthority_holder: %q\nagent_autonomy_bounded: true\nconstraints:\n  capability_implies_authority: false\n  proposal_is_authorized_decision: false\n  approval_is_authority_grant: false\n  risk_assessment_is_risk_acceptance: false\n", m.Project.AuthorityHolder)
}
func renderHumanGovernance(m model.BootstrapModel) string {
	return fmt.Sprintf("human_governance:\n  approval_declared: %t\n  approval_scope: \"not_assessed\"\n  approval_freshness: \"not_assessed\"\n  quorum: \"not_assessed\"\n  separation_of_duties: \"not_assessed\"\n  risk_acceptance_authority: \"not_assessed\"\n  break_glass: \"not_assessed\"\n  non_delegable_responsibilities: []\n", m.Project.HumanApproval)
}
func renderAgentInventory(m model.BootstrapModel) string {
	var b strings.Builder
	b.WriteString("agents:\n")
	for i, t := range m.Project.AgentTypes {
		fmt.Fprintf(&b, "  - id: \"agent-%d\"\n    type: %q\n    role: \"to-be-defined\"\n    authority: \"not_implied_by_type\"\n", i+1, t)
	}
	fmt.Fprintf(&b, "multi_agent: %t\nagent_output_trust: \"UntrustedProposal\"\n", m.Project.MultiAgent)
	return b.String()
}
func renderAgentCapabilities(m model.BootstrapModel) string {
	return renderCapabilityYAML(m.Project) + "rules:\n  capability_is_authority: false\n  self_grant_authority: false\n  technical_reachability_is_authority: false\n"
}
func renderAgentInteractions(m model.BootstrapModel) string {
	return fmt.Sprintf("agent_interactions:\n  multi_agent: %t\n  interaction_contract_applicability: %q\n  bounded_delegation_required: %t\n  disclosure_scope: \"not_assessed\"\n  provenance_preserved: \"not_assessed\"\n", m.Project.MultiAgent, m.Adoption.CanonicalObjects["AgentInteractionContract"].Applicability, m.Project.MultiAgent)
}
func renderAuthorityModel(m model.BootstrapModel) string {
	return fmt.Sprintf("authority:\n  governance_root: \"Human/Organization\"\n  declared_holder: %q\n  implementation_state: \"not_assessed\"\n  lifecycle:\n    issuance: \"not_assessed\"\n    activation: \"not_assessed\"\n    expiry: \"not_assessed\"\n    suspension: \"not_assessed\"\n    revocation: \"not_assessed\"\n    consumption: \"not_assessed\"\n    revalidation: \"not_assessed\"\n  scope:\n    operations: []\n    resources: []\n    environments: []\n    temporal_bounds: \"not_assessed\"\n  delegation:\n    allowed: %t\n    attenuation_required: true\n    provenance_required: true\n  principles:\n    - \"Capability != Authority\"\n    - \"Approval != AuthorityGrant\"\n    - \"NoGrant => NoAuthoritySensitiveExecution\"\n", m.Project.AuthorityHolder, m.Project.MultiAgent)
}
func renderAuthorityReadme(m model.BootstrapModel) string {
	return "# Authority Model\n\nThis directory records the project-specific Authority lifecycle and scope. `aof init` does **not** assert that Authority controls are implemented. Review issuance, validity, delegation, revocation, expiry, scope, provenance, and Effect Boundary revalidation before consequential execution.\n"
}
func renderPolicyModel(m model.BootstrapModel) string {
	return fmt.Sprintf("policy:\n  declared_enforcement_mechanism: %q\n  implementation_state: \"not_assessed\"\n  authoritative_source: \"not_assessed\"\n  versioning: \"not_assessed\"\n  conflict_resolution:\n    mechanism: \"not_assessed\"\n    deterministic_required: true\n    precedence: []\n  overrides:\n    bounded: true\n    traceable: true\n  policy_prompt_is_enforcement: false\n  unresolved_mandatory_policy: \"Pending\"\n  pending_allows_execution: false\n", m.Project.PolicyMechanism)
}
func renderRiskModel(m model.BootstrapModel) string {
	var b strings.Builder
	fmt.Fprintf(&b, "risk:\n  project_criticality: %q\n  data_sensitivity: %q\n  classification_method: \"not_assessed\"\n  dynamic_reassessment: \"not_assessed\"\n  inherent_risk: \"not_assessed\"\n  residual_risk: \"not_assessed\"\n  acceptance_authority: \"not_assessed\"\n  treatment: \"not_assessed\"\n  failure_budget: \"not_assessed\"\n  retry_budget: \"not_assessed\"\n", m.Project.Criticality, m.Project.DataSensitivity)
	writeListIndented(&b, "  categories", m.Project.RiskCategories)
	b.WriteString("  principles:\n    - \"RiskAssessment != RiskAcceptance\"\n    - \"Material risk change triggers reevaluation\"\n    - \"Mandatory unresolved risk != Pass\"\n")
	return b.String()
}
func renderExecutionModel(m model.BootstrapModel) string {
	return fmt.Sprintf("# Execution Model\n\nConsequential execution declared: **%t**.\n\n```text\nAgent output -> UntrustedProposal -> Governance Evaluation -> AuthorizedDecision -> ExecutionContract -> Effect Boundary -> Effect -> Evidence -> Verification -> Trace\n```\n\nExecution uses `ExecuteAllowed = C AND H AND P AND S AND R AND V`. Mandatory `Fail` blocks execution; mandatory `Pending` never becomes implicit Allow. Implementation of each gate remains `not_assessed` after init.\n", m.Project.AI.ConsequentialExecution)
}
func renderConsequentialActions(m model.BootstrapModel) string {
	var b strings.Builder
	fmt.Fprintf(&b, "consequential_execution: %t\n", m.Project.AI.ConsequentialExecution)
	writeList(&b, "effect_types", m.Project.EffectTypes)
	b.WriteString("requirements:\n  proposal_before_decision: true\n  authority_required: true\n  policy_required: true\n  state_validation_required: true\n  risk_required: true\n  verification_required: true\n  trace_required: true\n  effect_boundary_revalidation: true\nimplementation_state: \"not_assessed\"\n")
	return b.String()
}
func renderStateTransitions(m model.BootstrapModel) string {
	return fmt.Sprintf("state_model:\n  consequential_execution: %t\n  authoritative_state_function: \"not_assessed\"\n  transition_schema:\n    before_state: \"required_when_applicable\"\n    after_state: \"required_when_applicable\"\n    preconditions: []\n    postconditions: []\n    owner: \"not_assessed\"\n    versioning: \"not_assessed\"\n    conflict_control: \"not_assessed\"\n    idempotency: \"not_assessed\"\n    replay_behavior: \"not_assessed\"\n    external_state_reconciliation: \"not_assessed\"\n    toctou_revalidation: \"not_assessed\"\n  pending_is_pass: false\n", m.Project.AI.ConsequentialExecution)
}
func renderFailureRecovery(m model.BootstrapModel) string {
	return "failure_recovery:\n  containment: \"not_assessed\"\n  reconciliation: \"not_assessed\"\n  retry_policy: \"not_assessed\"\n  retry_budget: \"not_assessed\"\n  replan_policy: \"not_assessed\"\n  compensation: \"not_assessed\"\n  rollback: \"not_assessed\"\n  partial_effect_handling: \"not_assessed\"\n  unknown_effect_handling: \"not_assessed\"\n  break_glass: \"not_assessed\"\n  deadlock_livelock_controls: \"not_assessed\"\n  mandatory_control_failure_fails_open: false\n"
}
func renderEvidenceModel(m model.BootstrapModel) string {
	c := m.Adoption.Controls["evidence"]
	return fmt.Sprintf("# Evidence Model\n\n- Applicability: `%s`\n- Adoption state: `%s`\n- Implementation state: `%s`\n\nEvidence MUST remain distinct from Claim and Verification. Define provenance, integrity, freshness, relevance, sufficiency, retention, and binding to the relevant effect/outcome before claiming assurance.\n", c.Applicability, c.AdoptionState, c.ImplementationState)
}
func renderVerificationModel(m model.BootstrapModel) string {
	c := m.Adoption.Controls["verification"]
	vi := m.Adoption.Controls["verifier_independence"]
	return fmt.Sprintf("# Verification Model\n\n- Verification applicability: `%s`\n- Adoption state: `%s`\n- Implementation state: `%s`\n- Verifier independence applicability: `%s`\n\nVerification evaluates Evidence against explicit criteria. `Inconclusive`/`Pending` MUST NOT be treated as `Pass`. Define subject binding, result binding, independence, re-verification triggers, completion gate, and conflict handling.\n", c.Applicability, c.AdoptionState, c.ImplementationState, vi.Applicability)
}
func renderTraceModel(m model.BootstrapModel) string {
	return "# Trace Model\n\nTrace must support governance reconstruction without requiring private chain-of-thought. Define actor identity, Proposal/Decision/Action correlation, ordering/causality, Authority/Policy/Risk references, StateTransition, Evidence/Verification references, integrity, retention, redaction/access control, correction semantics, and failure behavior. Implementation state remains `not_assessed`.\n"
}
func renderAccountability(m model.BootstrapModel) string {
	return "accountability_chain:\n  governance_root: \"Human/Organization\"\n  authority_issuer: \"not_assessed\"\n  decision_actor: \"not_assessed\"\n  executor: \"not_assessed\"\n  verifier: \"not_assessed\"\n  risk_acceptor: \"not_assessed\"\n  evidence_provenance: \"not_assessed\"\n  separation_of_duties: \"not_assessed\"\n"
}
func renderConformanceScope(m model.BootstrapModel) string {
	var b strings.Builder
	fmt.Fprintf(&b, "aof_specification: \"1.0 LTS\"\ntarget_base_profile: %q\n", m.Adoption.Profile.TargetBaseProfile)
	writeList(&b, "domain_profiles", m.Adoption.Profile.DomainProfiles)
	writeList(&b, "overlays", m.Adoption.Profile.Overlays)
	b.WriteString("claimed_profile: null\nclaim_status: \"none\"\nscoped_adoption: true\nfull_conformance_claimed: false\nscope:\n  components: []\n  environments: []\n  resources: []\n")
	writeListIndented(&b, "  agent_types", m.Project.AgentTypes)
	b.WriteString("  tool_classes: []\n  governance_boundaries: []\n  excluded_functionality: []\n  security_assumptions: []\nprinciple: \"ScopedAdoption -> ScopedConformance\"\n")
	return b.String()
}
func renderConformanceReadme(m model.BootstrapModel) string {
	return "# Conformance\n\n`aof init` produces an adoption plan, **not** a `ConformanceManifest` or `ConformanceReport` claim. Canonical schemas are available under `aof/schemas/`. A future validation/conformance step must evaluate implementation and Evidence before any AOF profile claim is made.\n\n`BootstrapManifest != ConformanceManifest`.\n"
}
func renderGaps(m model.BootstrapModel) string {
	var b strings.Builder
	b.WriteString("# AOF Adoption Gaps\n\nGenerated applicability and adoption state do not prove implementation. The following controls require implementation/verification assessment.\n\n")
	for _, k := range sortedControlKeys(m.Adoption.Controls) {
		c := m.Adoption.Controls[k]
		if c.Applicability != model.ApplicabilityNotApplicable && (c.ImplementationState == model.ImplementationNotAssessed || c.AdoptionState == model.AdoptionUnsupported) {
			fmt.Fprintf(&b, "- **%s** — applicability `%s`, normative `%s`, adoption `%s`, implementation `%s`, verification `%s`. %s\n", k, c.Applicability, c.NormativeLevel, c.AdoptionState, c.ImplementationState, c.VerificationState, c.Reason)
		}
	}
	if len(m.Adoption.Warnings) > 0 {
		b.WriteString("\n## Warnings\n\n")
		for _, w := range m.Adoption.Warnings {
			fmt.Fprintf(&b, "- %s\n", w)
		}
	}
	return b.String()
}
func renderSecureSDLC(m model.BootstrapModel) string {
	return "profile: \"AOF-Secure-SDLC\"\nclassification: \"domain_profile\"\ncontrols:\n  requirements_and_acceptance_criteria: \"planned\"\n  architecture_security_review: \"planned\"\n  threat_modeling: \"planned\"\n  code_test_security_verification: \"planned\"\n  release_deployment_authority_gates: \"planned\"\n  human_control_gates_by_risk: \"planned\"\n  sast_sca_dast_or_equivalent_evidence: \"conditional\"\nimplementation_state: \"not_assessed\"\n"
}
func renderHighAssurance(m model.BootstrapModel) string {
	return "overlay: \"AOF-High-Assurance\"\nclassification: \"strengthening_overlay\"\ncontrols:\n  strong_separation_of_duties: \"planned\"\n  independent_verification: \"planned\"\n  critical_action_approval: \"planned\"\n  stronger_evidence: \"planned\"\n  strict_authority_scope: \"planned\"\n  effect_boundary_revalidation: \"planned\"\n  trace_protection: \"planned\"\n  bounded_recovery: \"planned\"\n  residual_risk_handling: \"planned\"\nimplementation_state: \"not_assessed\"\n"
}

func writeList(b *strings.Builder, key string, vals []string) {
	fmt.Fprintf(b, "%s:\n", key)
	if len(vals) == 0 {
		b.WriteString("  []\n")
		return
	}
	for _, v := range vals {
		fmt.Fprintf(b, "  - %q\n", v)
	}
}
func writeListIndented(b *strings.Builder, key string, vals []string) {
	fmt.Fprintf(b, "%s:\n", key)
	if len(vals) == 0 {
		b.WriteString("      []\n")
		return
	}
	indent := strings.Repeat(" ", len(key)-len(strings.TrimLeft(key, " "))+2)
	for _, v := range vals {
		fmt.Fprintf(b, "%s- %q\n", indent, v)
	}
}
func sortedControlKeys(m map[string]model.ControlSelection) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func safePath(path string) bool {
	if path == "" {
		return false
	}
	if filepath.IsAbs(path) || filepath.VolumeName(path) != "" || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
		return false
	}
	if len(path) >= 2 && ((path[0] >= 'a' && path[0] <= 'z') || (path[0] >= 'A' && path[0] <= 'Z')) && path[1] == ':' {
		return false
	}
	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || filepath.IsAbs(clean) || filepath.VolumeName(clean) != "" {
		return false
	}
	if strings.HasPrefix(clean, "/") || strings.HasPrefix(clean, "\\") {
		return false
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "..\\") {
		return false
	}
	return true
}
func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
