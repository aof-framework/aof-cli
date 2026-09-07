package canonical

import "testing"

func TestCanonicalSemanticUniverseIsComplete(t *testing.T) {
	if err := Validate(); err != nil {
		t.Fatal(err)
	}
	registry, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := len(registry.Invariants); got != ExpectedInvariantCount {
		t.Fatalf("registered invariants=%d want=%d", got, ExpectedInvariantCount)
	}
	if got := len(registry.Requirements); got != ExpectedRequirementCount {
		t.Fatalf("registered requirements=%d want=%d", got, ExpectedRequirementCount)
	}
	if got := len(registry.Invariants) + len(registry.Requirements); got != ExpectedSemanticIDCount {
		t.Fatalf("stable semantic IDs=%d want=%d", got, ExpectedSemanticIDCount)
	}
}

func TestEveryCanonicalRequirementIsResolvable(t *testing.T) {
	registry := MustLoad()
	for _, requirement := range registry.Requirements {
		resolved, ok := registry.Requirement(requirement.ID)
		if !ok || resolved.Statement == "" {
			t.Fatalf("canonical requirement %s is not resolvable", requirement.ID)
		}
	}
}
