package adoption

import (
	"github.com/aof-framework/aof-cli/internal/model"
	"testing"
)

func TestBuildCoreAdvisory(t *testing.T) {
	a := Build(model.ProfileCore, model.AIUsage{Enabled: true}, map[string]bool{"authority": true, "policy": true, "risk": true, "trace": true})
	if a.Controls["execution_contract"].Status != model.StatusNotApplicable {
		t.Fatal("execution contract should be N/A without consequential execution")
	}
	if a.Controls["evidence"].Status != model.StatusDeferred {
		t.Fatal("evidence should be deferred in core")
	}
}

func TestUnsupportedRequiredBecomesDeferred(t *testing.T) {
	a := Build(model.ProfileAssured, model.AIUsage{Enabled: true, ConsequentialExecution: true}, map[string]bool{"independent_verification": false})
	c := a.Controls["independent_verification"]
	if c.Status != model.StatusDeferred || c.Capability != "unsupported" {
		t.Fatalf("unexpected control: %+v", c)
	}
}
