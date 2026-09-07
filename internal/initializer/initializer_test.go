package initializer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aof-framework/aof-cli/internal/adoption"
	"github.com/aof-framework/aof-cli/internal/canonical"
	"github.com/aof-framework/aof-cli/internal/model"
)

func testModel() model.BootstrapModel {
	p := model.ProjectDefinition{Name: "itsm-backend", Language: "Go", Type: "backend", Domain: "itsm", AgentTypes: []string{"LLM"}, AI: model.AIUsage{Enabled: true, Analysis: true, Recommendation: true}, Criticality: "moderate", DataSensitivity: "internal", AdoptionMode: "greenfield"}
	return model.BootstrapModel{Project: p, Adoption: adoption.Build(model.ProfileCore, p, nil)}
}

func TestGeneratedCanonicalSchemasMatchUpstreamIndex(t *testing.T) {
	plan, err := BuildPlan(testModel(), "0.3.2")
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{}
	for _, file := range plan.Files {
		files[file.Path] = file.Content
	}
	for _, object := range canonical.MustLoad().CanonicalObjects {
		path := "aof/schemas/" + object.Schema
		data, ok := files[path]
		if !ok {
			t.Fatalf("generated schema missing: %s", path)
		}
		sum := sha256.Sum256(data)
		if got := hex.EncodeToString(sum[:]); got != object.ActiveSHA256 {
			t.Fatalf("generated schema %s differs from active upstream file: got %s want %s", path, got, object.ActiveSHA256)
		}
	}
}
func TestBuildAndApply(t *testing.T) {
	d := t.TempDir()
	p, err := BuildPlan(testModel(), "0.3.2")
	if err != nil {
		t.Fatal(err)
	}
	if err := Apply(d, p); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"AGENTS.md", "AOF.md", ".aof/project.yaml", ".aof/config.yaml", ".aof/requirements.yaml", ".aof/bootstrap-manifest.json", "aof/control/safety-kernel.yaml", "aof/architecture/planes.yaml", "aof/authority/authority-model.yaml", "aof/execution/failure-recovery.yaml", "aof/schemas/core/agent.schema.json"} {
		if _, err := os.Stat(filepath.Join(d, path)); err != nil {
			t.Fatalf("missing %s: %v", path, err)
		}
	}
	b, err := os.ReadFile(filepath.Join(d, ".aof/bootstrap-manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	if v["schema_type"] != "AOFCLIBootstrapManifest" || v["claimed_conformance"] != false {
		t.Fatalf("bad manifest: %v", v)
	}
	if v["registered_invariants"] != float64(162) || v["registered_requirements"] != float64(332) || v["stable_semantic_ids"] != float64(494) {
		t.Fatalf("manifest lacks canonical coverage gate: %v", v)
	}
	for _, key := range []string{"source_checkout_commit", "active_specification_sha256", "original_semantic_baseline_commit", "original_semantic_baseline_sha256", "active_schema_checksums_sha256"} {
		if v[key] == nil || v[key] == "" {
			t.Fatalf("manifest lacks explicit upstream provenance %s: %v", key, v)
		}
	}
	requirements, err := os.ReadFile(filepath.Join(d, ".aof/requirements.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"requirement_count: 332", "stable_semantic_id_count: 494", `requirement_id: "AOF-AGT-001"`, `statement: "Agent MUST operate as bounded actor within applicable Governance Envelope."`, `requirement_id: "AOF-ARCH-001"`, `normative_level: "Unclassified"`, `requirement_id: "AOF-VER-018"`, `unselected_requirement_disposition: "not_evaluated"`} {
		if !strings.Contains(string(requirements), expected) {
			t.Fatalf("requirements projection missing %q", expected)
		}
	}
	applicability, err := os.ReadFile(filepath.Join(d, ".aof/applicability.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(applicability), `metadata_classification: "non_canonical_aof_cli_process_metadata"`) || strings.Contains(string(applicability), "normative_level:") {
		t.Fatalf("CLI planning metadata is presented as canonical semantics:\n%s", applicability)
	}
}
func TestConflictAbortsBeforeWrites(t *testing.T) {
	d := t.TempDir()
	_ = os.WriteFile(filepath.Join(d, "AGENTS.md"), []byte("existing"), 0o644)
	p, _ := BuildPlan(testModel(), "0.3.2")
	err := Apply(d, p)
	var ce *ConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected conflict,got %v", err)
	}
	if _, err := os.Stat(filepath.Join(d, "AOF.md")); !os.IsNotExist(err) {
		t.Fatal("partial write")
	}
}
func TestSecondInit(t *testing.T) {
	d := t.TempDir()
	p, _ := BuildPlan(testModel(), "0.3.2")
	if err := Apply(d, p); err != nil {
		t.Fatal(err)
	}
	if err := Apply(d, p); !errors.Is(err, ErrAlreadyInitialized) {
		t.Fatalf("expected already initialized,got %v", err)
	}
}
func TestSafePath(t *testing.T) {
	for _, bad := range []string{"", "../escape", "/absolute", "\\absolute", "C:\\escape", "C:foo", "a/../../escape"} {
		if safePath(bad) {
			t.Fatalf("unsafe accepted %q", bad)
		}
	}
}

func TestPreserveExistingAgentsUsesSidecar(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "AGENTS.md"), []byte("existing project instructions"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := BuildPlan(testModel(), "0.3.2")
	if err != nil {
		t.Fatal(err)
	}
	p, sidecar, err := PreserveExistingAgents(d, p)
	if err != nil {
		t.Fatal(err)
	}
	if !sidecar {
		t.Fatal("expected sidecar mode")
	}
	if err := Apply(d, p); err != nil {
		t.Fatal(err)
	}
	original, _ := os.ReadFile(filepath.Join(d, "AGENTS.md"))
	if string(original) != "existing project instructions" {
		t.Fatal("existing AGENTS.md changed")
	}
	if _, err := os.Stat(filepath.Join(d, "AGENTS.aof.md")); err != nil {
		t.Fatalf("missing sidecar: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(d, ".aof/bootstrap-manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "AGENTS.aof.md") || strings.Contains(string(b), `"AGENTS.md"`) {
		t.Fatal("manifest did not record sidecar")
	}
}
