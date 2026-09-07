package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGo(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "go.mod"), []byte("module x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := Detect(d)
	if r.Language != "Go" || !r.ExistingProject || r.Type != "backend" {
		t.Fatalf("unexpected: %+v", r)
	}
}

func TestDetectTypeScript(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "tsconfig.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := Detect(d)
	if r.Language != "TypeScript" || !r.ExistingProject || r.Type != "fullstack" {
		t.Fatalf("unexpected: %+v", r)
	}
}

func TestDetectEmpty(t *testing.T) {
	r := Detect(t.TempDir())
	if r.Language != "" || r.ExistingProject {
		t.Fatalf("unexpected: %+v", r)
	}
}
