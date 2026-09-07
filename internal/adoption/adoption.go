package adoption

import (
	"sort"

	"github.com/aof-framework/aof-cli/internal/model"
)

func Build(profile string, p model.ProjectDefinition, capabilities map[string]bool) model.AdoptionDefinition {
	controls := map[string]model.ControlSelection{}
	add := func(name, applicability, status, reason string) {
		capability := "supported"
		if supported, ok := capabilities[name]; ok && !supported {
			capability = "unsupported"
			if status == model.StatusEnabled || status == model.StatusRequired {
				status = model.StatusDeferred
				reason = "organizational capability not currently available"
			}
		}
		controls[name] = model.ControlSelection{Applicability: applicability, Capability: capability, Status: status, Reason: reason}
	}

	add("governance_root", "required", model.StatusEnabled, "AOF governance requires a Human/Organization governance root")
	add("capability_boundaries", "required", model.StatusEnabled, "Capability must be explicit and separate from Authority")
	add("authority", "required", model.StatusEnabled, "Authority boundaries are foundational")
	add("policy", "required", model.StatusEnabled, "Policy applicability must be explicit")
	add("risk", "required", model.StatusEnabled, "Risk must be assessed within project scope")
	add("trace", "required", model.StatusEnabled, "governance-relevant decisions require traceability")

	if p.AI.Enabled {
		add("agent_governance", "required", model.StatusEnabled, "AI agents are declared in project scope")
	} else {
		add("agent_governance", "not_applicable", model.StatusNotApplicable, "AI agents are not declared")
	}

	if p.MultiAgent {
		add("agent_interaction_contract", "required", model.StatusEnabled, "multi-agent interaction declared")
	} else {
		add("agent_interaction_contract", "not_applicable", model.StatusNotApplicable, "single-agent or unspecified agent architecture")
	}

	consequential := p.AI.ConsequentialExecution || p.AI.StateMutation || p.AI.ExternalAction || p.AI.InfrastructureAction
	if consequential {
		add("action_proposal", "required", model.StatusEnabled, "consequential AI behavior requires Proposal/Decision separation")
		add("authorized_decision", "required", model.StatusEnabled, "consequential AI behavior declared")
		add("execution_contract", "required", model.StatusEnabled, "consequential AI execution declared")
		add("state_validation", "required", model.StatusEnabled, "consequential AI execution declared")
		add("effect_boundary", "required", model.StatusEnabled, "external or persistent effects are in scope")
		add("verification", "required", model.StatusEnabled, "consequential AI execution declared")
		add("evidence", "required", model.StatusEnabled, "verification requires evidence")
		if !p.HumanApproval {
			add("approval", "applicable", model.StatusDeferred, "no human approval mechanism declared; Approval remains distinct from Authority")
		} else {
			add("approval", "applicable", model.StatusEnabled, "human approval declared")
		}
	} else {
		add("action_proposal", "applicable", model.StatusEnabled, "AI output is treated as UntrustedProposal")
		add("authorized_decision", "applicable", model.StatusEnabled, "Decision remains separate from Proposal")
		add("execution_contract", "not_applicable", model.StatusNotApplicable, "no consequential AI execution declared")
		add("state_validation", "not_applicable", model.StatusNotApplicable, "no consequential AI execution declared")
		add("effect_boundary", "not_applicable", model.StatusNotApplicable, "no consequential AI execution declared")
		add("verification", "applicable", model.StatusDeferred, "may be strengthened as assurance maturity increases")
		add("evidence", "applicable", model.StatusDeferred, "may be strengthened as assurance maturity increases")
		add("approval", "optional", model.StatusNotApplicable, "no consequential action requiring approval declared")
	}

	if profile == model.ProfileGoverned || profile == model.ProfileAssured {
		add("escalation", "required", model.StatusEnabled, "selected adoption profile")
		add("conformance_manifest", "required", model.StatusEnabled, "selected adoption profile")
	} else {
		add("escalation", "applicable", model.StatusDeferred, "progressive strengthening for AOF-Core")
		add("conformance_manifest", "applicable", model.StatusDeferred, "progressive strengthening for AOF-Core")
	}

	if profile == model.ProfileAssured {
		add("independent_verification", "required", model.StatusEnabled, "AOF-Assured strengthening")
		add("formal_conformance_evidence", "required", model.StatusEnabled, "AOF-Assured strengthening")
	} else {
		add("independent_verification", "optional", model.StatusDeferred, "not required by selected baseline")
		add("formal_conformance_evidence", "optional", model.StatusDeferred, "not required by selected baseline")
	}

	objects := buildObjects(p, controls)
	skills := []string{"aof-core", "aof-adoption", "aof-conformance"}
	if p.AI.Enabled {
		skills = append(skills, "aof-agent-design", "aof-authority", "aof-policy", "aof-risk")
	}
	if consequential {
		skills = append(skills, "aof-execution", "aof-evidence-verification", "aof-state-trace")
	}
	sort.Strings(skills)

	warnings := []string{}
	for name, c := range controls {
		if c.Status == model.StatusDeferred && c.Applicability == "required" {
			warnings = append(warnings, name+": required control is deferred")
		}
	}
	if consequential && p.AuthorityHolder == "unspecified" {
		warnings = append(warnings, "authority: consequential execution is enabled but the authority holder is unspecified")
	}
	if consequential && p.PolicyMechanism == "prompt" {
		warnings = append(warnings, "policy: prompt instructions do not satisfy PolicyEnforcement")
	}
	sort.Strings(warnings)

	return model.AdoptionDefinition{Profile: profile, Mode: "greenfield", Controls: controls, CanonicalObjects: objects, RequiredSkills: skills, Warnings: warnings}
}

func buildObjects(p model.ProjectDefinition, controls map[string]model.ControlSelection) map[string]model.ObjectSelection {
	objects := map[string]model.ObjectSelection{}
	for _, name := range model.CanonicalObjects {
		objects[name] = model.ObjectSelection{Status: model.StatusNotApplicable, Reason: "not selected by current project applicability rules"}
	}
	required := func(name, reason string) {
		objects[name] = model.ObjectSelection{Status: model.StatusRequired, Reason: reason}
	}
	applicable := func(name, reason string) {
		objects[name] = model.ObjectSelection{Status: model.StatusEnabled, Reason: reason}
	}

	required("Goal", "all governed work originates from declared goals")
	required("Task", "project work is decomposed into tasks")
	required("ContextDescriptor", "context boundaries are needed for governed reasoning")
	required("Resource", "resources accessed by actors must be bounded")
	required("Capability", "Capability is distinct from Authority")
	required("Policy", "Policy is part of execution authorization")
	required("RiskAssessment", "Risk is part of execution authorization")
	required("ActionProposal", "Proposal is distinct from Decision")
	required("Decision", "Decision is distinct from Proposal and Action")
	required("TraceEvent", "governance-relevant events require traceability")
	required("Outcome", "governed actions require outcome representation")
	if p.AI.Enabled {
		required("Agent", "AI agents are declared")
	}
	if p.MultiAgent {
		required("AgentInteractionContract", "multi-agent interaction declared")
	}

	consequential := p.AI.ConsequentialExecution || p.AI.StateMutation || p.AI.ExternalAction || p.AI.InfrastructureAction
	if consequential {
		required("AuthorityGrant", "consequential execution requires explicit Authority")
		required("ExecutionContract", "consequential effects require a bounded execution contract")
		required("Evidence", "Verification requires Evidence")
		required("Verification", "consequential execution requires verification")
		required("StateTransition", "persistent/external effects imply state transitions")
		if p.HumanApproval {
			applicable("Approval", "human approval declared")
		}
	}
	if c, ok := controls["escalation"]; ok && c.Status != model.StatusNotApplicable {
		applicable("EscalationPackage", "escalation control is applicable")
	}
	if c, ok := controls["conformance_manifest"]; ok && c.Status != model.StatusNotApplicable {
		applicable("ConformanceManifest", "selected adoption profile")
	}
	if c, ok := controls["formal_conformance_evidence"]; ok && c.Status == model.StatusEnabled {
		applicable("ConformanceReport", "formal conformance evidence enabled")
	}
	return objects
}
