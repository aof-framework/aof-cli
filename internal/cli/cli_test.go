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
		t.Fatalf("code %d err %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "AOF CLI v0.2.0") || !strings.Contains(out.String(), "AOF Specification v1.0 LTS") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestHelp(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"help"}, strings.NewReader(""), &out, &errOut); code != 0 {
		t.Fatal(code)
	}
	if !strings.Contains(out.String(), "aof init") {
		t.Fatalf("unexpected help: %s", out.String())
	}
}

func TestInitNonInteractive(t *testing.T) {
	d := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	if err := os.Chdir(d); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := Run([]string{"init", "--name", "itsm-backend", "--language", "go", "--type", "backend", "--domain", "itsm", "--profile", "core", "--ai", "--ai-analysis", "--ai-recommendation", "--non-interactive"}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "AOF initialization complete") || !strings.Contains(out.String(), "Applicable Agent Skills") {
		t.Fatalf("unexpected output: %s", out.String())
	}
	if _, err := os.Stat("AOF.md"); err != nil {
		t.Fatal(err)
	}
}

func TestInitHelp(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"init", "--help"}, strings.NewReader(""), &out, &errOut); code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "Usage of aof init:") || !strings.Contains(out.String(), "-profile") {
		t.Fatalf("unexpected init help: %s", out.String())
	}
}
