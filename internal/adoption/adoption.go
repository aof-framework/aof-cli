package adoption

import "github.com/aof-framework/aof-cli/internal/model"

func Build(profile string, ai model.AIUsage, capabilities map[string]bool) model.AdoptionDefinition {
	controls := map[string]model.ControlSelection{}
	add := func(name, applicability, status, reason string) {
		capability := "supported"
		if supported, ok := capabilities[name]; ok && !supported {
			capability = "unsupported"
			if status == model.StatusEnabled {
				status = model.StatusDeferred
				reason = "organizational capability not currently available"
			}
		}
		controls[name] = model.ControlSelection{Applicability: applicability, Capability: capability, Status: status, Reason: reason}
	}

	add("governance_root", "required", model.StatusEnabled, "")
	add("capability_boundaries", "required", model.StatusEnabled, "")
	add("authority", "required", model.StatusEnabled, "")
	add("policy", "required", model.StatusEnabled, "")
	add("risk", "required", model.StatusEnabled, "")
	add("trace", "required", model.StatusEnabled, "")

	if ai.ConsequentialExecution {
		add("execution_contract", "required", model.StatusEnabled, "consequential AI execution declared")
		add("state_validation", "required", model.StatusEnabled, "consequential AI execution declared")
		add("verification", "required", model.StatusEnabled, "consequential AI execution declared")
	} else {
		add("execution_contract", "not_applicable", model.StatusNotApplicable, "no consequential AI execution declared")
		add("state_validation", "not_applicable", model.StatusNotApplicable, "no consequential AI execution declared")
		add("verification", "applicable", model.StatusDeferred, "may be strengthened as assurance maturity increases")
	}

	if profile == model.ProfileGoverned || profile == model.ProfileAssured {
		add("evidence", "required", model.StatusEnabled, "selected profile")
		add("escalation", "required", model.StatusEnabled, "selected profile")
	} else {
		add("evidence", "applicable", model.StatusDeferred, "optional strengthening for AOF-Core")
		add("escalation", "applicable", model.StatusDeferred, "optional strengthening for AOF-Core")
	}

	if profile == model.ProfileAssured {
		add("independent_verification", "required", model.StatusEnabled, "AOF-Assured strengthening")
		add("formal_conformance_evidence", "required", model.StatusEnabled, "AOF-Assured strengthening")
	} else {
		add("independent_verification", "optional", model.StatusDeferred, "not required by selected baseline")
		add("formal_conformance_evidence", "optional", model.StatusDeferred, "not required by selected baseline")
	}

	return model.AdoptionDefinition{Profile: profile, Mode: "greenfield", Controls: controls}
}
