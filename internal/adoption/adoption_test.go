package adoption

import (
	"testing"

	"github.com/aof-framework/aof-cli/internal/model"
)

func TestBuildCoreAdvisory(t *testing.T) {
	p := model.ProjectDefinition{AI: model.AIUsage{Enabled: true}}
	a := Build(model.ProfileCore, p, map[string]bool{"authority": true, "policy": true, "risk": true, "trace": true})
	if a.Controls["execution_contract"].Status != model.StatusNotApplicable {
		t.Fatal("execution contract should be N/A without consequential execution")
	}
	if a.Controls["evidence"].Status != model.StatusDeferred {
		t.Fatal("evidence should be deferred for advisory core")
	}
	if a.CanonicalObjects["Agent"].Status != model.StatusRequired {
		t.Fatal("Agent should be required when AI is enabled")
	}
}

func TestConsequentialRequiresExecutionObjects(t *testing.T) {
	p := model.ProjectDefinition{AI: model.AIUsage{Enabled: true, ConsequentialExecution: true}}
	a := Build(model.ProfileCore, p, map[string]bool{"authority": true, "policy": true, "risk": true, "trace": true, "evidence": true, "verification": true})
	for _, name := range []string{"AuthorityGrant", "ExecutionContract", "Evidence", "Verification", "StateTransition"} {
		if a.CanonicalObjects[name].Status != model.StatusRequired {
			t.Fatalf("%s should be required: %+v", name, a.CanonicalObjects[name])
		}
	}
	if a.Controls["effect_boundary"].Status != model.StatusEnabled {
		t.Fatal("effect boundary should be enabled")
	}
}

func TestUnsupportedRequiredBecomesDeferred(t *testing.T) {
	p := model.ProjectDefinition{AI: model.AIUsage{Enabled: true, ConsequentialExecution: true}}
	a := Build(model.ProfileAssured, p, map[string]bool{"independent_verification": false})
	c := a.Controls["independent_verification"]
	if c.Status != model.StatusDeferred || c.Capability != "unsupported" {
		t.Fatalf("unexpected control: %+v", c)
	}
}
