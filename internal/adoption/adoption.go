package adoption

import (
	"sort"
	"strings"

	"github.com/aof-framework/aof-cli/internal/canonical"
	"github.com/aof-framework/aof-cli/internal/model"
)

type Options struct {
	DomainProfiles []string
	HighAssurance  bool
}

func Build(profile string, p model.ProjectDefinition, capabilities map[string]bool) model.AdoptionDefinition {
	return BuildWithOptions(profile, p, capabilities, Options{})
}

func BuildWithOptions(profile string, p model.ProjectDefinition, capabilities map[string]bool, opts Options) model.AdoptionDefinition {
	controls := map[string]model.ControlSelection{}
	add := func(name, applicability, level, reason string, reqs ...string) {
		capability := "not_declared"
		state := model.AdoptionPlanned
		if applicability == model.ApplicabilityNotApplicable {
			state = model.AdoptionNotApplicable
			capability = "not_applicable"
		} else if supported, ok := capabilities[name]; ok {
			if supported {
				capability = "supported"
			} else {
				capability = "unsupported"
				state = model.AdoptionUnsupported
			}
		}
		controls[name] = model.ControlSelection{
			Applicability: applicability, NormativeLevel: level, Capability: capability,
			AdoptionState: state, ImplementationState: model.ImplementationNotAssessed,
			VerificationState: model.VerificationNotEvaluated, Reason: reason, RequirementIDs: reqs,
		}
	}

	// AOF-Core minimum identity and architecture semantics (§21.2 / §8.36).
	add("governance_root", model.ApplicabilityApplicable, model.NormativeMust, "Human/Organization remains GovernanceRoot", "AOF-ARCH-002")
	add("identifiable_agent", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Core requires an identifiable Agent", "AOF-ARCH-001")
	add("task_model", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Core requires Task semantics")
	add("action_separation", model.ApplicabilityApplicable, model.NormativeMust, "Proposal, Decision, and consequential Action remain distinguishable", "AOF-ARCH-001")
	add("capability_boundaries", model.ApplicabilityApplicable, model.NormativeMust, "Capability MUST remain distinct from Authority", "AOF-ARCH-003", "AOF-AUTH-007")
	add("authority", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Core requires explicit Authority semantics", "AOF-AUTH-001", "AOF-AUTH-002")
	add("policy", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Core requires Policy semantics", "AOF-POL-001", "AOF-POL-002")
	add("state", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Core requires authoritative State semantics", "AOF-ARCH-007")
	add("trace", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Core requires reconstruction-capable Trace", "AOF-ARCH-006", "AOF-TRC-001")
	add("controlled_transition", model.ApplicabilityApplicable, model.NormativeMust, "Consequential state mutation MUST use controlled transition", "AOF-ARCH-005")
	add("no_implicit_allow", model.ApplicabilityApplicable, model.NormativeMust, "Mandatory Pending/unknown MUST NOT become implicit allow", "AOF-ARCH-004", "AOF-POL-002")

	consequential := p.AI.ConsequentialExecution || p.AI.StateMutation || p.AI.ExternalAction || p.AI.InfrastructureAction || len(p.EffectTypes) > 0
	if consequential {
		add("action_proposal", model.ApplicabilityApplicable, model.NormativeMust, "consequential execution requires Proposal/Decision separation", "AOF-ARCH-001")
		add("authorized_decision", model.ApplicabilityApplicable, model.NormativeMust, "consequential execution requires governed Decision", "AOF-ARCH-001")
		add("authority_grant", model.ApplicabilityApplicable, model.NormativeMust, "authority-sensitive execution requires applicable grant", "AOF-AUTH-001", "AOF-AUTH-004")
		add("risk_assessment", model.ApplicabilityApplicable, model.NormativeMust, "control predicate requires Risk evaluation for consequential execution", "AOF-RISK-001", "AOF-RISK-006")
		add("execution_contract", model.ApplicabilityApplicable, model.NormativeMust, "consequential effects require explicit execution boundary contract", "AOF-ARCH-002", "AOF-ARCH-009")
		add("effect_boundary", model.ApplicabilityApplicable, model.NormativeMust, "consequential effects MUST be mediated and revalidated", "AOF-ARCH-002", "AOF-ARCH-009")
		add("state_validation", model.ApplicabilityApplicable, model.NormativeMust, "current state MUST be validated before consequential transition", "AOF-ARCH-005")
		add("evidence", model.ApplicabilityApplicable, model.NormativeMust, "effect/result Evidence is required for applicable assurance", "AOF-ARCH-010")
		add("verification", model.ApplicabilityApplicable, model.NormativeMust, "mandatory Verification gates consequential completion", "AOF-VER-001", "AOF-VER-012")
	} else {
		add("action_proposal", model.ApplicabilityConditional, model.NormativeShould, "advisory outputs SHOULD remain explicit proposals")
		add("authorized_decision", model.ApplicabilityConditional, model.NormativeShould, "Decision remains distinct from Proposal")
		add("authority_grant", model.ApplicabilityConditional, model.NormativeMust, "required whenever an authority-sensitive action enters scope", "AOF-AUTH-001")
		add("risk_assessment", model.ApplicabilityConditional, model.NormativeShould, "required when risk-bearing consequential behavior enters scope", "AOF-RISK-001")
		add("execution_contract", model.ApplicabilityNotApplicable, model.NormativeMay, "no consequential execution declared")
		add("effect_boundary", model.ApplicabilityConditional, model.NormativeMust, "required when project crosses a consequential Effect Boundary", "AOF-ARCH-009")
		add("state_validation", model.ApplicabilityConditional, model.NormativeMust, "required for consequential state mutation", "AOF-ARCH-005")
		add("evidence", model.ApplicabilityConditional, model.NormativeShould, "strengthens assurance when outcomes require objective proof")
		add("verification", model.ApplicabilityConditional, model.NormativeShould, "strengthens assurance when acceptance must be evaluated")
	}

	if p.MultiAgent {
		add("agent_interaction_contract", model.ApplicabilityApplicable, model.NormativeMust, "multi-Agent interaction/delegation declared", "AOF-ARCH-012")
	} else {
		add("agent_interaction_contract", model.ApplicabilityConditional, model.NormativeShould, "becomes applicable if multiple Agents interact or delegate")
	}

	if p.HumanApproval || profile == model.ProfileGoverned || profile == model.ProfileAssured {
		add("approval", model.ApplicabilityApplicable, model.NormativeMust, "approval/escalation applies to selected profile or declared human gate", "AOF-AUTH-006")
	} else {
		add("approval", model.ApplicabilityConditional, model.NormativeShould, "applicable when a Human approval gate is introduced")
	}

	if profile == model.ProfileGoverned || profile == model.ProfileAssured {
		add("authority_lifecycle", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Governed requires explicit Authority lifecycle", "AOF-AUTH-004", "AOF-AUTH-008", "AOF-AUTH-016")
		add("policy_conflict_resolution", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Governed requires deterministic Policy conflict resolution", "AOF-POL-004", "AOF-POL-005")
		add("dynamic_risk", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Governed requires dynamic Risk evaluation", "AOF-ARCH-011", "AOF-RISK-003", "AOF-RISK-009")
		add("bounded_delegation", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Governed requires bounded delegation", "AOF-AUTH-009", "AOF-AUTH-011", "AOF-ARCH-012")
		add("escalation", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Governed requires approval/escalation")
		add("failure_budget", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Governed requires failure budgets")
		add("conformance_manifest", model.ApplicabilityApplicable, model.NormativeShould, "profile/scope SHOULD be captured in a canonical ConformanceManifest", "AOF-PRF-003")
	} else {
		for _, n := range []string{"authority_lifecycle", "policy_conflict_resolution", "dynamic_risk", "bounded_delegation", "escalation", "failure_budget", "conformance_manifest"} {
			add(n, model.ApplicabilityConditional, model.NormativeShould, "progressive strengthening beyond AOF-Core")
		}
	}

	if profile == model.ProfileAssured {
		add("evidence", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Assured requires Evidence", "AOF-ARCH-017")
		add("evidence_provenance", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Assured requires Evidence provenance", "AOF-ARCH-017")
		add("verification", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Assured requires Verification", "AOF-VER-001")
		add("verifier_independence", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Assured requires verifier independence based on Risk", "AOF-ARCH-016", "AOF-VER-015")
		add("completion_gate", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Assured requires completion gate", "AOF-VER-012")
		add("accountability_chain", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Assured requires accountability chain")
		add("conformance_report", model.ApplicabilityApplicable, model.NormativeShould, "Assured adoption SHOULD produce formal conformance reporting")
	} else {
		for _, n := range []string{"evidence_provenance", "verifier_independence", "completion_gate", "accountability_chain", "conformance_report"} {
			add(n, model.ApplicabilityConditional, model.NormativeShould, "AOF-Assured strengthening")
		}
	}

	secure := contains(opts.DomainProfiles, model.DomainSecureSDLC)
	if secure {
		add("secure_sdlc_requirements", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Secure-SDLC domain profile")
		add("architecture_security_review", model.ApplicabilityApplicable, model.NormativeShould, "AOF-Secure-SDLC architecture/security review")
		add("threat_modeling", model.ApplicabilityApplicable, model.NormativeShould, "AOF-Secure-SDLC threat modeling")
		add("code_test_security_verification", model.ApplicabilityApplicable, model.NormativeShould, "AOF-Secure-SDLC code/test/security Verification")
		add("release_authority_gate", model.ApplicabilityApplicable, model.NormativeMust, "AOF-Secure-SDLC release/deployment Authority gate")
		add("human_control_gate", model.ApplicabilityApplicable, model.NormativeShould, "AOF-Secure-SDLC Human control gates according to Risk")
		add("security_evidence", model.ApplicabilityConditional, model.NormativeShould, "SAST/SCA/DAST or equivalent Evidence when applicable")
	}

	if opts.HighAssurance {
		add("separation_of_duties", model.ApplicabilityApplicable, model.NormativeMust, "AOF-High-Assurance strengthening")
		add("verifier_independence", model.ApplicabilityApplicable, model.NormativeMust, "independent Verification is mandatory for High-Assurance", "AOF-VER-015", "AOF-VER-016")
		add("critical_action_approval", model.ApplicabilityApplicable, model.NormativeMust, "explicit approval for defined Critical actions")
		add("stronger_evidence", model.ApplicabilityApplicable, model.NormativeMust, "stronger Evidence requirements")
		add("strict_authority_scope", model.ApplicabilityApplicable, model.NormativeMust, "strict Authority scope", "AOF-AUTH-019", "AOF-AUTH-020")
		add("effect_boundary_revalidation", model.ApplicabilityApplicable, model.NormativeMust, "revalidation at Effect Boundary", "AOF-ARCH-009")
		add("trace_protection", model.ApplicabilityApplicable, model.NormativeMust, "tamper-resistant or equivalent Trace protection", "AOF-TRC-013")
		add("bounded_recovery", model.ApplicabilityApplicable, model.NormativeMust, "bounded Recovery")
		add("residual_risk_handling", model.ApplicabilityApplicable, model.NormativeMust, "residual Risk handling", "AOF-RISK-008", "AOF-RISK-017")
	}

	// Project context can strengthen requirements without changing the selected profile.
	if consequential && (p.Criticality == "high" || p.Criticality == "critical") {
		add("verifier_independence", model.ApplicabilityApplicable, model.NormativeMust, "High/Critical consequential Risk requires independent Verification", "AOF-RISK-015", "AOF-RISK-016", "AOF-VER-015", "AOF-VER-016")
	}
	if p.DataSensitivity == "confidential" || p.DataSensitivity == "restricted" {
		add("context_data_controls", model.ApplicabilityApplicable, model.NormativeMust, "sensitive context requires purpose/access/cross-boundary controls")
		add("trace_data_controls", model.ApplicabilityApplicable, model.NormativeMust, "Trace MUST respect data classification/access/retention", "AOF-TRC-008", "AOF-TRC-009")
	}

	objects := buildObjects(p, profile, secure, opts.HighAssurance, consequential)
	requirements := requirementRegistry(controls, secure, opts.HighAssurance)
	skills := []string{"aof-core", "aof-adoption", "aof-agent-design", "aof-authority", "aof-policy", "aof-risk", "aof-state-trace", "aof-conformance"}
	if consequential {
		skills = append(skills, "aof-execution", "aof-evidence-verification")
	}
	if secure {
		skills = append(skills, "aof-secure-sdlc")
	}
	if opts.HighAssurance {
		skills = append(skills, "aof-high-assurance")
	}
	sort.Strings(skills)

	warnings := []string{}
	for name, c := range controls {
		if c.AdoptionState == model.AdoptionUnsupported && c.NormativeLevel == model.NormativeMust {
			warnings = append(warnings, name+": mandatory applicable control is unsupported; no conformance satisfaction is implied")
		}
	}
	if consequential && (strings.TrimSpace(p.AuthorityHolder) == "" || p.AuthorityHolder == "unspecified") {
		warnings = append(warnings, "authority: consequential execution is in scope but Authority holder is unspecified")
	}
	if consequential && p.PolicyMechanism == "prompt" {
		warnings = append(warnings, "policy: prompt instructions do not satisfy PolicyEnforcement")
	}
	if profile == model.ProfileAssured && controls["verifier_independence"].Capability == "unsupported" {
		warnings = append(warnings, "AOF-Assured target requires verifier independence but organizational capability is unavailable")
	}
	if secure && profile == model.ProfileCore {
		warnings = append(warnings, "AOF-Secure-SDLC SHOULD compose with AOF-Governed/AOF-Assured controls according to assurance scope")
	}
	sort.Strings(warnings)

	domains := append([]string(nil), opts.DomainProfiles...)
	overlays := []string{}
	if opts.HighAssurance {
		overlays = append(overlays, model.OverlayHighAssurance)
	}
	mode := p.AdoptionMode
	if mode == "" {
		mode = "greenfield"
	}
	return model.AdoptionDefinition{
		Profile: model.ProfileSelection{TargetBaseProfile: profile, DomainProfiles: domains, Overlays: overlays, ClaimedProfile: "", ClaimStatus: "none"},
		Mode:    mode, Controls: controls, CanonicalObjects: objects, Requirements: requirements,
		RequiredSkills: skills, Warnings: warnings,
	}
}

func buildObjects(p model.ProjectDefinition, profile string, secure, high, consequential bool) map[string]model.ObjectSelection {
	objects := map[string]model.ObjectSelection{}
	set := func(name, applicability, level, state, reason string, reqs ...string) {
		objects[name] = model.ObjectSelection{Applicability: applicability, NormativeLevel: level, AdoptionState: state, Reason: reason, RequirementIDs: reqs}
	}
	for _, object := range canonical.MustLoad().CanonicalObjects {
		set(object.Name, model.ApplicabilityNotEvaluated, "Unclassified", model.AdoptionDeferred, "canonical schema exists; applicability has not been established by an upstream source mapping")
	}
	for _, name := range []string{"Goal", "Task", "Agent", "ContextDescriptor", "Resource", "Capability", "AuthorityGrant", "Policy", "ActionProposal", "Decision", "StateTransition", "TraceEvent", "Outcome"} {
		set(name, model.ApplicabilityApplicable, model.NormativeMust, model.AdoptionPlanned, "AOF-Core project model")
	}
	if consequential {
		for _, name := range []string{"RiskAssessment", "ExecutionContract", "Evidence", "Verification"} {
			set(name, model.ApplicabilityApplicable, model.NormativeMust, model.AdoptionPlanned, "consequential execution scope")
		}
	}
	if p.HumanApproval || profile != model.ProfileCore {
		set("Approval", model.ApplicabilityApplicable, model.NormativeMust, model.AdoptionPlanned, "human/profile governance gate")
	}
	if p.MultiAgent {
		set("AgentInteractionContract", model.ApplicabilityApplicable, model.NormativeMust, model.AdoptionPlanned, "multi-Agent scope")
	}
	if profile == model.ProfileGoverned || profile == model.ProfileAssured || secure || high {
		set("EscalationPackage", model.ApplicabilityApplicable, model.NormativeMust, model.AdoptionPlanned, "governed/escalation scope")
		set("ConformanceManifest", model.ApplicabilityApplicable, model.NormativeShould, model.AdoptionPlanned, "explicit profile/scope declaration")
	}
	if profile == model.ProfileAssured || high {
		set("Evidence", model.ApplicabilityApplicable, model.NormativeMust, model.AdoptionPlanned, "AOF-Assured/High-Assurance")
		set("Verification", model.ApplicabilityApplicable, model.NormativeMust, model.AdoptionPlanned, "AOF-Assured/High-Assurance")
		set("ConformanceReport", model.ApplicabilityApplicable, model.NormativeShould, model.AdoptionPlanned, "formal assurance reporting")
	}
	return objects
}

func requirementRegistry(controls map[string]model.ControlSelection, secure, high bool) []model.Requirement {
	registry := canonical.MustLoad()
	seen := map[string]bool{}
	for _, c := range controls {
		for _, id := range c.RequirementIDs {
			if _, ok := registry.Requirement(id); !ok {
				panic("AOF CLI control references unknown canonical requirement " + id)
			}
			seen[id] = true
		}
	}
	for _, id := range []string{"AOF-PRF-001", "AOF-PRF-002", "AOF-PRF-006"} {
		seen[id] = true
	}
	if secure {
		seen["AOF-PRF-004"] = true
	}
	if high {
		seen["AOF-PRF-005"] = true
	}
	out := make([]model.Requirement, 0, len(registry.Requirements))
	for _, requirement := range registry.Requirements {
		applicability := model.ApplicabilityNotEvaluated
		reason := "canonical requirement retained; upstream applicability mapping is absent and AOF CLI has not evaluated it"
		if seen[requirement.ID] {
			applicability = model.ApplicabilityApplicable
			reason = "selected by the AOF CLI project adoption process; this is non-canonical projection metadata"
		}
		out = append(out, requirementFor(registry, requirement, applicability, reason))
	}
	return out
}

func requirementFor(registry canonical.Registry, requirement canonical.Requirement, applicability, reason string) model.Requirement {
	appliesTo := strings.Join(requirement.AppliesTo, ", ")
	if appliesTo == "" {
		appliesTo = "source_mapping_absent"
	}
	verification := strings.Join(requirement.VerificationMethods, ", ")
	verificationDisposition := "declared_in_requirement_registry"
	if verification == "" {
		verification = "source_mapping_absent"
		verificationDisposition = "source_mapping_absent"
	}
	evidenceDisposition := "declared_in_requirement_registry"
	if len(requirement.RequiredEvidence) == 0 {
		evidenceDisposition = "source_mapping_absent"
	}
	trace := registry.RequirementTraceability(requirement.ID)
	invariantDisposition := "declared_in_traceability_matrix"
	if len(trace.RelatedInvariants) == 0 {
		invariantDisposition = "source_mapping_absent"
	}
	return model.Requirement{
		ID: requirement.ID, Domain: requirement.Domain, Statement: requirement.Statement,
		NormativeLevel: requirement.NormativeLevel, CanonicalSource: requirement.Source.Specification,
		CanonicalSourceLine: requirement.Source.Line, CanonicalStatus: requirement.Status, Applicability: applicability,
		ApplicabilityReason: reason, AppliesTo: appliesTo, Profiles: append([]string(nil), requirement.Profiles...),
		VerificationMethod: verification, RequiredEvidence: append([]string(nil), requirement.RequiredEvidence...),
		RelatedInvariants:       append([]string(nil), trace.RelatedInvariants...),
		RelatedObjects:          nil,
		VerificationDisposition: verificationDisposition, EvidenceDisposition: evidenceDisposition,
		InvariantMappingDisposition: invariantDisposition, ObjectMappingDisposition: "source_mapping_absent",
	}
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
