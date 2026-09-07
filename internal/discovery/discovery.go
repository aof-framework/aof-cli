package discovery

import (
	"os"
	"path/filepath"
)

type Result struct {
	Language        string
	Type            string
	ExistingProject bool
}

func Detect(root string) Result {
	result := Result{}
	langChecks := []struct {
		path string
		lang string
	}{
		{"go.mod", "Go"},
		{"tsconfig.json", "TypeScript"},
		{"package.json", "Node.js"},
		{"pyproject.toml", "Python"},
		{"requirements.txt", "Python"},
		{"Pipfile", "Python"},
		{"setup.py", "Python"},
		{"Cargo.toml", "Rust"},
		{"pom.xml", "Java"},
		{"build.gradle.kts", "Kotlin"},
		{"build.gradle", "Java"},
		{"composer.json", "PHP"},
		{"Gemfile", "Ruby"},
		{"Package.swift", "Swift"},
		{"CMakeLists.txt", "C/C++"},
	}
	for _, c := range langChecks {
		if _, err := os.Stat(filepath.Join(root, c.path)); err == nil {
			result.Language = c.lang
			result.ExistingProject = true
			break
		}
	}
	if result.Language == "" {
		if matches, _ := filepath.Glob(filepath.Join(root, "*.csproj")); len(matches) > 0 {
			result.Language = "C#"
			result.ExistingProject = true
		} else if matches, _ := filepath.Glob(filepath.Join(root, "*.sln")); len(matches) > 0 {
			result.Language = "C#"
			result.ExistingProject = true
		}
	}
	if result.Language == "Go" || result.Language == "Rust" || result.Language == "Java" || result.Language == "Kotlin" {
		result.Type = "backend"
	} else if result.Language == "TypeScript" {
		result.Type = "fullstack"
	}
	if entries, err := os.ReadDir(root); err == nil {
		for _, e := range entries {
			if e.Name() == ".git" || e.Name() == ".DS_Store" {
				continue
			}
			result.ExistingProject = true
			break
		}
	}
	return result
}
