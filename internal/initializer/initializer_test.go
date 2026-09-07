package initializer

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aof-framework/aof-cli/internal/adoption"
	"github.com/aof-framework/aof-cli/internal/model"
)

func testModel() model.BootstrapModel {
	p := model.ProjectDefinition{Name: "itsm-backend", Language: "Go", Type: "backend", Domain: "itsm", AI: model.AIUsage{Enabled: true, Analysis: true, Recommendation: true}}
	return model.BootstrapModel{Project: p, Adoption: adoption.Build(model.ProfileCore, p.AI, map[string]bool{"authority": true, "policy": true, "risk": true, "trace": true})}
}

func TestBuildAndApply(t *testing.T) {
	d := t.TempDir()
	p, err := BuildPlan(testModel(), "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if err := Apply(d, p); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"AGENTS.md", "AOF.md", ".aof/project.yaml", ".aof/config.yaml", ".aof/manifest.json", ".agents/skills/aof-core/SKILL.md"} {
		if _, err := os.Stat(filepath.Join(d, path)); err != nil {
			t.Fatalf("missing %s: %v", path, err)
		}
	}
	b, err := os.ReadFile(filepath.Join(d, ".aof/manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	if v["upstream_repository"] != "https://github.com/aof-framework/aof" {
		t.Fatalf("unexpected upstream: %v", v["upstream_repository"])
	}
}

func TestConflictAbortsBeforeWrites(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "AGENTS.md"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, _ := BuildPlan(testModel(), "0.1.0")
	err := Apply(d, p)
	var ce *ConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(d, "AOF.md")); !os.IsNotExist(err) {
		t.Fatal("AOF.md should not be written on conflict")
	}
}

func TestSecondInit(t *testing.T) {
	d := t.TempDir()
	p, _ := BuildPlan(testModel(), "0.1.0")
	if err := Apply(d, p); err != nil {
		t.Fatal(err)
	}
	if err := Apply(d, p); !errors.Is(err, ErrAlreadyInitialized) {
		t.Fatalf("expected already initialized, got %v", err)
	}
}

func TestSafePath(t *testing.T) {
	for _, bad := range []string{"", "../escape", "/absolute", "\\absolute", "C:\\escape", "C:foo", "a/../../escape"} {
		if safePath(bad) {
			t.Fatalf("expected unsafe: %q", bad)
		}
	}
	for _, good := range []string{"AOF.md", ".agents/skills/aof-core/SKILL.md"} {
		if !safePath(good) {
			t.Fatalf("expected safe: %q", good)
		}
	}
}
