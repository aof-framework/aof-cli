package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"version"}, strings.NewReader(""), &out, &errOut); code != 0 {
		t.Fatal(code)
	}
	if !strings.Contains(out.String(), "AOF CLI v0.3.1") {
		t.Fatalf("unexpected: %s", out.String())
	}
	if !strings.Contains(out.String(), "Release: Canonical Semantic Coverage") {
		t.Fatalf("unexpected release label: %s", out.String())
	}
}
func TestHelp(t *testing.T) {
	var out, errOut bytes.Buffer
	if Run([]string{"help"}, strings.NewReader(""), &out, &errOut) != 0 {
		t.Fatal("help")
	}
	if !strings.Contains(out.String(), "--secure-sdlc") {
		t.Fatal(out.String())
	}
}

func TestInitNonInteractiveAssuredSecureSDLC(t *testing.T) {
	d := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	if err := os.Chdir(d); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	args := []string{"init", "--name", "itsm-backend", "--language", "go", "--type", "backend", "--domain", "itsm", "--profile", "assured", "--secure-sdlc", "--high-assurance", "--agent-types", "LLM,Human", "--ai", "--ai-consequential", "--effect-types", "database,external-api", "--criticality", "critical", "--data-sensitivity", "restricted", "--non-interactive"}
	if code := Run(args, strings.NewReader(""), &out, &errOut); code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	for _, p := range []string{"AOF.md", ".aof/requirements.yaml", "aof/control/safety-kernel.yaml", "aof/secure-sdlc/profile.yaml", "aof/high-assurance/overlay.yaml", "aof/schemas/core/agent.schema.json"} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing %s: %v", p, err)
		}
	}
	if !strings.Contains(out.String(), "Claimed conformance: none") || !strings.Contains(out.String(), "Requirements indexed") {
		t.Fatal(out.String())
	}
}

func TestInvalidTypedInput(t *testing.T) {
	d := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(d)
	var out, errOut bytes.Buffer
	if code := Run([]string{"init", "--criticality", "banana", "--non-interactive"}, strings.NewReader(""), &out, &errOut); code != 2 {
		t.Fatalf("got %d", code)
	}
}

func TestNoAIDefaultsToDeterministicAgent(t *testing.T) {
	d := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(d)
	var out, errOut bytes.Buffer
	if code := Run([]string{"init", "--name", "svc", "--profile", "core", "--non-interactive"}, strings.NewReader(""), &out, &errOut); code != 0 {
		t.Fatalf("%d %s", code, errOut.String())
	}
	b, err := os.ReadFile("aof/agents/inventory.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "Deterministic") {
		t.Fatal(string(b))
	}
}

func TestInvalidProfileComposition(t *testing.T) {
	cases := [][]string{
		{"init", "--profile", "core", "--secure-sdlc", "--non-interactive"},
		{"init", "--profile", "governed", "--high-assurance", "--non-interactive"},
	}
	for _, args := range cases {
		d := t.TempDir()
		old, _ := os.Getwd()
		_ = os.Chdir(d)
		var out, errOut bytes.Buffer
		if code := Run(args, strings.NewReader(""), &out, &errOut); code != 2 {
			t.Fatalf("args=%v code=%d stderr=%s", args, code, errOut.String())
		}
		_ = os.Chdir(old)
	}
}
