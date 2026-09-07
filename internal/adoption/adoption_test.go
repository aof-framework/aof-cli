package adoption

import (
	"testing"

	"github.com/aof-framework/aof-cli/internal/model"
)

func baseProject() model.ProjectDefinition {
	return model.ProjectDefinition{Name: "x", AgentTypes: []string{"Deterministic"}, Criticality: "moderate", DataSensitivity: "internal"}
}

func TestCoreRequiresIdentifiableAgentEvenWithoutAI(t *testing.T) {
	a := Build(model.ProfileCore, baseProject(), nil)
	if a.Controls["identifiable_agent"].Applicability != model.ApplicabilityApplicable {
		t.Fatalf("identifiable agent not applicable: %+v", a.Controls["identifiable_agent"])
	}
	if a.CanonicalObjects["Agent"].Applicability != model.ApplicabilityApplicable {
		t.Fatalf("Agent should be applicable: %+v", a.CanonicalObjects["Agent"])
	}
	if a.Profile.ClaimStatus != "none" || a.Profile.ClaimedProfile != "" {
		t.Fatal("init must not claim conformance")
	}
}

func TestAdoptionDoesNotAssertImplementation(t *testing.T) {
	a := Build(model.ProfileCore, baseProject(), nil)
	for name, c := range a.Controls {
		if c.Applicability != model.ApplicabilityNotApplicable && c.ImplementationState != model.ImplementationNotAssessed {
			t.Fatalf("%s asserted implementation: %+v", name, c)
		}
		if c.VerificationState != model.VerificationNotEvaluated {
			t.Fatalf("%s asserted verification: %+v", name, c)
		}
	}
}

func TestConsequentialRequiresGovernanceObjects(t *testing.T) {
	p := baseProject()
	p.AI = model.AIUsage{Enabled: true, ConsequentialExecution: true}
	p.AgentTypes = []string{"LLM"}
	p.EffectTypes = []string{"database"}
	a := Build(model.ProfileCore, p, map[string]bool{"evidence": true})
	for _, name := range []string{"AuthorityGrant", "RiskAssessment", "ExecutionContract", "Evidence", "Verification", "StateTransition", "TraceEvent"} {
		if a.CanonicalObjects[name].Applicability != model.ApplicabilityApplicable {
			t.Fatalf("%s should be applicable: %+v", name, a.CanonicalObjects[name])
		}
	}
	if a.Controls["effect_boundary"].NormativeLevel != model.NormativeMust {
		t.Fatal("effect boundary must be mandatory")
	}
}

func TestGovernedProfileLiteralStrengthening(t *testing.T) {
	a := Build(model.ProfileGoverned, baseProject(), nil)
	for _, name := range []string{"authority_lifecycle", "policy_conflict_resolution", "dynamic_risk", "bounded_delegation", "approval", "escalation", "failure_budget"} {
		c := a.Controls[name]
		if c.Applicability != model.ApplicabilityApplicable || c.NormativeLevel != model.NormativeMust {
			t.Fatalf("%s missing Governed semantics: %+v", name, c)
		}
	}
}

func TestAssuredProfileLiteralStrengthening(t *testing.T) {
	a := Build(model.ProfileAssured, baseProject(), nil)
	for _, name := range []string{"evidence", "evidence_provenance", "verification", "verifier_independence", "completion_gate", "accountability_chain"} {
		c := a.Controls[name]
		if c.Applicability != model.ApplicabilityApplicable || c.NormativeLevel != model.NormativeMust {
			t.Fatalf("%s missing Assured semantics: %+v", name, c)
		}
	}
}

func TestSecureSDLCAndHighAssuranceComposition(t *testing.T) {
	a := BuildWithOptions(model.ProfileAssured, baseProject(), nil, Options{DomainProfiles: []string{model.DomainSecureSDLC}, HighAssurance: true})
	for _, name := range []string{"secure_sdlc_requirements", "threat_modeling", "release_authority_gate", "separation_of_duties", "trace_protection", "bounded_recovery", "residual_risk_handling"} {
		if a.Controls[name].Applicability != model.ApplicabilityApplicable {
			t.Fatalf("%s missing composition: %+v", name, a.Controls[name])
		}
	}
	if len(a.Profile.DomainProfiles) != 1 || len(a.Profile.Overlays) != 1 {
		t.Fatalf("bad profile composition: %+v", a.Profile)
	}
}

func TestHighCriticalConsequentialStrengthensVerifierIndependence(t *testing.T) {
	p := baseProject()
	p.Criticality = "critical"
	p.AI = model.AIUsage{ConsequentialExecution: true}
	p.EffectTypes = []string{"infrastructure"}
	a := Build(model.ProfileCore, p, map[string]bool{"verifier_independence": false})
	c := a.Controls["verifier_independence"]
	if c.Applicability != model.ApplicabilityApplicable || c.NormativeLevel != model.NormativeMust || c.AdoptionState != model.AdoptionUnsupported {
		t.Fatalf("unexpected: %+v", c)
	}
}

func TestRequirementRegistryContainsTraceability(t *testing.T) {
	p := baseProject()
	p.AI = model.AIUsage{ConsequentialExecution: true}
	p.EffectTypes = []string{"database"}
	a := Build(model.ProfileGoverned, p, nil)
	seen := map[string]bool{}
	for _, r := range a.Requirements {
		seen[r.ID] = true
	}
	for _, id := range []string{"AOF-ARCH-001", "AOF-AUTH-001", "AOF-POL-004", "AOF-RISK-003", "AOF-PRF-001"} {
		if !seen[id] {
			t.Fatalf("missing %s", id)
		}
	}
}
