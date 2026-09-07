package discovery

import (
	"os"
	"path/filepath"
)

type Result struct {
	Language string
	Type     string
}

func Detect(root string) Result {
	res := Result{}

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
			res.Language = c.lang
			break
		}
	}

	if res.Language == "" {
		if matches, _ := filepath.Glob(filepath.Join(root, "*.csproj")); len(matches) > 0 {
			res.Language = "C#"
		} else if matches, _ := filepath.Glob(filepath.Join(root, "*.sln")); len(matches) > 0 {
			res.Language = "C#"
		}
	}

	if res.Language == "Go" || res.Language == "Rust" || res.Language == "Java" || res.Language == "Kotlin" {
		res.Type = "backend"
	} else if res.Language == "TypeScript" {
		res.Type = "fullstack"
	}

	return res
}
