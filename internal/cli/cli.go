package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aof-framework/aof-cli/internal/adoption"
	"github.com/aof-framework/aof-cli/internal/discovery"
	"github.com/aof-framework/aof-cli/internal/initializer"
	"github.com/aof-framework/aof-cli/internal/model"
)

const Version = "0.1.0"

func Run(args []string, in io.Reader, out, errOut io.Writer) int {
	if len(args) == 0 {
		printHelp(out)
		return 0
	}
	switch args[0] {
	case "help", "-h", "--help":
		printHelp(out)
		return 0
	case "version", "--version", "-v":
		printVersion(out)
		return 0
	case "init":
		return runInit(args[1:], in, out, errOut)
	default:
		fmt.Fprintf(errOut, "unknown command %q\n\n", args[0])
		printHelp(errOut)
		return 2
	}
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "AOF CLI — Bootstrap utility for AI Orchestration Framework (AOF)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  aof init [flags]     Initialize AOF adoption in current repository")
	fmt.Fprintln(w, "  aof version          Print version and framework specification info")
	fmt.Fprintln(w, "  aof help             Show this help message")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'aof init --help' for details on initialization flags.")
}
func printVersion(w io.Writer) {
	fmt.Fprintf(w, "AOF CLI v%s\n", Version)
	fmt.Fprintln(w, "AOF Specification v1.0 LTS")
	fmt.Fprintln(w, "Canonical Upstream: https://github.com/aof-framework/aof")
}

func runInit(args []string, in io.Reader, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(errOut)
	name := fs.String("name", "", "project name")
	language := fs.String("language", "", "project language")
	projectType := fs.String("type", "", "project type")
	domain := fs.String("domain", "", "project domain")
	profile := fs.String("profile", "core", "AOF profile: core, governed, assured")
	features := fs.String("features", "", "comma-separated project features")
	ai := fs.Bool("ai", false, "declare AI usage")
	analysis := fs.Bool("ai-analysis", false, "AI analysis use case")
	classification := fs.Bool("ai-classification", false, "AI classification use case")
	recommendation := fs.Bool("ai-recommendation", false, "AI recommendation use case")
	consequential := fs.Bool("ai-consequential", false, "AI may perform consequential execution")
	nonInteractive := fs.Bool("non-interactive", false, "disable interactive prompts")
	fs.Usage = func() {
		fmt.Fprintln(out, "Usage of aof init:")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Flags:")
		fs.SetOutput(out)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(errOut, "resolve project root: %v\n", err)
		return 1
	}
	detected := discovery.Detect(root)
	reader := bufio.NewReader(in)
	if *name == "" {
		*name = filepath.Base(root)
	}
	if *language == "" {
		*language = detected.Language
	}
	if *projectType == "" && detected.Type != "" {
		*projectType = detected.Type
	}
	if !*nonInteractive {
		*name = ask(reader, out, "Project name", *name)
		*language = ask(reader, out, "Language", defaultValue(*language, "unspecified"))
		*projectType = ask(reader, out, "Project type", defaultValue(*projectType, "unspecified"))
		*domain = ask(reader, out, "Domain", defaultValue(*domain, "unspecified"))
		*profile = ask(reader, out, "AOF profile (core/governed/assured)", *profile)
		if !*ai {
			*ai = askBool(reader, out, "Will this project use AI agents", false)
		}
		if *ai {
			if !*analysis {
				*analysis = askBool(reader, out, "AI analysis", true)
			}
			if !*classification {
				*classification = askBool(reader, out, "AI classification", false)
			}
			if !*recommendation {
				*recommendation = askBool(reader, out, "AI recommendation", true)
			}
			if !*consequential {
				*consequential = askBool(reader, out, "Consequential AI execution", false)
			}
		}
	}
	if *language == "" {
		*language = "unspecified"
	}
	if *projectType == "" {
		*projectType = "unspecified"
	}
	if *domain == "" {
		*domain = "unspecified"
	}
	p := model.NormalizeProfile(*profile)
	if p == "" {
		fmt.Fprintln(errOut, "invalid profile; use core, governed, or assured")
		return 2
	}
	project := model.ProjectDefinition{Name: *name, Language: *language, Type: *projectType, Domain: *domain, Features: splitList(*features), AI: model.AIUsage{Enabled: *ai || *analysis || *classification || *recommendation || *consequential, Analysis: *analysis, Classification: *classification, Recommendation: *recommendation, ConsequentialExecution: *consequential}}
	capabilities := map[string]bool{"authority": true, "policy": true, "risk": true, "trace": true, "evidence": p != model.ProfileCore, "verification": *consequential, "independent_verification": p == model.ProfileAssured, "formal_conformance_evidence": p == model.ProfileAssured}
	adopt := adoption.Build(p, project.AI, capabilities)
	plan, err := initializer.BuildPlan(model.BootstrapModel{Project: project, Adoption: adopt}, Version)
	if err != nil {
		fmt.Fprintf(errOut, "build initialization plan: %v\n", err)
		return 1
	}
	if err := initializer.Apply(root, plan); err != nil {
		if errors.Is(err, initializer.ErrAlreadyInitialized) {
			fmt.Fprintln(errOut, err.Error())
			return 3
		}
		var ce *initializer.ConflictError
		if errors.As(err, &ce) {
			fmt.Fprintln(errOut, "Initialization cannot continue safely.\n\nConflicts:")
			for _, p := range ce.Paths {
				fmt.Fprintf(errOut, "- %s\n", p)
			}
			fmt.Fprintln(errOut, "\nNo files were changed.")
			return 4
		}
		fmt.Fprintf(errOut, "initialize AOF: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "AOF initialization complete.\n\nProject: %s\nProfile: %s\nGenerated files: %d\nCanonical AOF: https://github.com/aof-framework/aof\n", project.Name, adopt.Profile, len(plan.Files))
	return 0
}

func ask(r *bufio.Reader, w io.Writer, label, def string) string {
	if def != "" {
		fmt.Fprintf(w, "%s [%s]: ", label, def)
	} else {
		fmt.Fprintf(w, "%s: ", label)
	}
	s, _ := r.ReadString('\n')
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	return s
}
func askBool(r *bufio.Reader, w io.Writer, label string, def bool) bool {
	d := "y/N"
	if def {
		d = "Y/n"
	}
	fmt.Fprintf(w, "%s [%s]: ", label, d)
	s, _ := r.ReadString('\n')
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return def
	}
	return s == "y" || s == "yes" || s == "true"
}
func defaultValue(v, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return v
}
func splitList(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}
