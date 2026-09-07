package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	cases := []struct {
		file     string
		wantLang string
		wantType string
	}{
		{"go.mod", "Go", "backend"},
		{"tsconfig.json", "TypeScript", "fullstack"},
		{"package.json", "Node.js", ""},
		{"pyproject.toml", "Python", ""},
		{"requirements.txt", "Python", ""},
		{"Cargo.toml", "Rust", "backend"},
		{"pom.xml", "Java", "backend"},
		{"build.gradle.kts", "Kotlin", "backend"},
		{"composer.json", "PHP", ""},
		{"Gemfile", "Ruby", ""},
		{"Package.swift", "Swift", ""},
		{"app.csproj", "C#", ""},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			d := t.TempDir()
			if err := os.WriteFile(filepath.Join(d, tc.file), []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
			res := Detect(d)
			if res.Language != tc.wantLang {
				t.Fatalf("file %s got lang %q want %q", tc.file, res.Language, tc.wantLang)
			}
			if res.Type != tc.wantType {
				t.Fatalf("file %s got type %q want %q", tc.file, res.Type, tc.wantType)
			}
		})
	}
}
