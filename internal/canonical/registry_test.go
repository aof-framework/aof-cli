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

func TestCanonicalRequirementClassificationIsNeverRewritten(t *testing.T) {
	registry := MustLoad()
	unclassified := 0
	for _, requirement := range registry.Requirements {
		if requirement.NormativeLevel == "Unclassified" {
			unclassified++
		}
	}
	if unclassified != 17 {
		t.Fatalf("upstream Unclassified requirements=%d want=17", unclassified)
	}
	for _, id := range []string{"AOF-ARCH-001", "AOF-ARCH-011", "AOF-ARCH-017"} {
		requirement, ok := registry.Requirement(id)
		if !ok || requirement.NormativeLevel != "Unclassified" {
			t.Fatalf("%s was reclassified: %+v", id, requirement)
		}
	}
}

func TestCanonicalObjectsComeFromUpstreamSchemaIndex(t *testing.T) {
	registry := MustLoad()
	if got := len(registry.CanonicalObjects); got != ExpectedSchemaCount {
		t.Fatalf("schema contracts=%d want=%d", got, ExpectedSchemaCount)
	}
	for _, object := range registry.CanonicalObjects {
		if object.Schema == "" || object.ID == "" || len(object.IndexSHA256) != 64 || len(object.ActiveSHA256) != 64 {
			t.Fatalf("incomplete upstream schema index record: %+v", object)
		}
	}
}
