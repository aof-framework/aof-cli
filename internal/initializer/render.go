package initializer

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/aof-framework/aof-cli/internal/model"
	bootstraptemplates "github.com/aof-framework/aof-cli/templates"
)

type File struct {
	Path    string
	Content []byte
}

type Plan struct {
	Files []File
}

func BuildPlan(m model.BootstrapModel, cliVersion string) (Plan, error) {
	files := []File{}
	static := []struct{ src, dst string }{
		{"static/docs/semantic-boundaries.md", "docs/aof/semantic-boundaries.md"},
		{"static/docs/governance-model.md", "docs/aof/governance-model.md"},
		{"static/docs/execution-model.md", "docs/aof/execution-model.md"},
		{"static/skills/aof-core/SKILL.md", ".agents/skills/aof-core/SKILL.md"},
		{"static/skills/aof-governance/SKILL.md", ".agents/skills/aof-governance/SKILL.md"},
		{"static/skills/aof-conformance/SKILL.md", ".agents/skills/aof-conformance/SKILL.md"},
	}
	for _, item := range static {
		b, err := fs.ReadFile(bootstraptemplates.FS, item.src)
		if err != nil {
			return Plan{}, fmt.Errorf("read embedded template %s: %w", item.src, err)
		}
		files = append(files, File{Path: item.dst, Content: b})
	}

	dynamic := []struct{ src, dst string }{
		{"dynamic/AGENTS.md.tmpl", "AGENTS.md"},
		{"dynamic/AOF.md.tmpl", "AOF.md"},
		{"dynamic/project-architecture.SKILL.md.tmpl", ".agents/skills/project-architecture/SKILL.md"},
		{"dynamic/project-domain.SKILL.md.tmpl", ".agents/skills/project-domain/SKILL.md"},
		{"dynamic/agent-architect.md.tmpl", ".agents/agents/architect/agent.md"},
		{"dynamic/agent-implementer.md.tmpl", ".agents/agents/implementer/agent.md"},
		{"dynamic/agent-reviewer.md.tmpl", ".agents/agents/reviewer/agent.md"},
		{"dynamic/agent-aof-reviewer.md.tmpl", ".agents/agents/aof-reviewer/agent.md"},
	}
	for _, item := range dynamic {
		b, err := fs.ReadFile(bootstraptemplates.FS, item.src)
		if err != nil {
			return Plan{}, fmt.Errorf("read embedded template %s: %w", item.src, err)
		}
		t, err := template.New(item.src).Parse(string(b))
		if err != nil {
			return Plan{}, fmt.Errorf("parse template %s: %w", item.src, err)
		}
		var out bytes.Buffer
		if err := t.Execute(&out, m); err != nil {
			return Plan{}, fmt.Errorf("render template %s: %w", item.src, err)
		}
		files = append(files, File{Path: item.dst, Content: out.Bytes()})
	}

	projectYAML := renderProjectYAML(m.Project)
	configYAML := renderConfigYAML(m.Adoption)
	files = append(files,
		File{Path: ".aof/project.yaml", Content: []byte(projectYAML)},
		File{Path: ".aof/config.yaml", Content: []byte(configYAML)},
		File{Path: "aof/policies/.gitkeep", Content: []byte{}},
		File{Path: "aof/profiles/.gitkeep", Content: []byte{}},
		File{Path: "aof/fixtures/.gitkeep", Content: []byte{}},
		File{Path: "aof/conformance/.gitkeep", Content: []byte{}},
	)

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	manifest, err := renderManifest(m, cliVersion, files)
	if err != nil {
		return Plan{}, err
	}
	files = append(files, File{Path: ".aof/manifest.json", Content: manifest})
	return Plan{Files: files}, nil
}

func renderProjectYAML(p model.ProjectDefinition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "project:\n  name: %q\n  language: %q\n  type: %q\n  domain: %q\n", p.Name, p.Language, p.Type, p.Domain)
	b.WriteString("features:\n")
	if len(p.Features) == 0 {
		b.WriteString("  []\n")
	} else {
		for _, f := range p.Features {
			fmt.Fprintf(&b, "  - %q\n", f)
		}
	}
	fmt.Fprintf(&b, "ai:\n  enabled: %t\n  analysis: %t\n  classification: %t\n  recommendation: %t\n  consequential_execution: %t\n", p.AI.Enabled, p.AI.Analysis, p.AI.Classification, p.AI.Recommendation, p.AI.ConsequentialExecution)
	return b.String()
}

func renderConfigYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "aof:\n  specification: \"1.0\"\n  release: \"LTS\"\n  profile: %q\n", a.Profile)
	b.WriteString("upstream:\n  repository: \"https://github.com/aof-framework/aof\"\n  specification: \"specification/AOF-v1.0-Framework-Specification.md\"\n")
	fmt.Fprintf(&b, "adoption:\n  mode: %q\ncontrols:\n", a.Mode)
	keys := make([]string, 0, len(a.Controls))
	for k := range a.Controls {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		c := a.Controls[k]
		fmt.Fprintf(&b, "  %s:\n    applicability: %q\n    capability: %q\n    status: %q\n", k, c.Applicability, c.Capability, c.Status)
		if c.Reason != "" {
			fmt.Fprintf(&b, "    reason: %q\n", c.Reason)
		}
	}
	return b.String()
}

func safePath(path string) bool {
	if path == "" {
		return false
	}
	if filepath.IsAbs(path) || filepath.VolumeName(path) != "" || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
		return false
	}
	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || filepath.IsAbs(clean) || filepath.VolumeName(clean) != "" {
		return false
	}
	if strings.HasPrefix(clean, "/") || strings.HasPrefix(clean, "\\") {
		return false
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "..\\") {
		return false
	}
	return true
}
