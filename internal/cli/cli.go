package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aof-framework/aof-cli/internal/adoption"
	"github.com/aof-framework/aof-cli/internal/discovery"
	"github.com/aof-framework/aof-cli/internal/initializer"
	"github.com/aof-framework/aof-cli/internal/model"
)

const Version = "0.2.0"

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
	fmt.Fprintln(w, "AOF CLI — Context & Adoption Hardening")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  aof init [flags]")
	fmt.Fprintln(w, "  aof version")
	fmt.Fprintln(w, "  aof help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "aof init discovers project context, interviews for applicable AOF governance,")
	fmt.Fprintln(w, "builds a scoped adoption model, and generates project-aware Agent Skills and controls.")
}

func printVersion(w io.Writer) {
	fmt.Fprintf(w, "AOF CLI v%s\n", Version)
	fmt.Fprintln(w, "Release: AOF Context & Adoption Hardening")
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
	architecture := fs.String("architecture", "", "architecture style")
	deployment := fs.String("deployment", "", "deployment model")
	criticality := fs.String("criticality", "", "project criticality: low, moderate, high, critical")
	dataSensitivity := fs.String("data-sensitivity", "", "data sensitivity: public, internal, confidential, restricted")
	profile := fs.String("profile", "core", "AOF profile: core, governed, assured")
	features := fs.String("features", "", "comma-separated project features")
	agentCount := fs.String("agent-count", "", "logical AI agent count or description")
	multiAgent := fs.Bool("multi-agent", false, "declare multi-agent interactions")
	humanApproval := fs.Bool("human-approval", false, "declare human approval for consequential actions")
	authorityHolder := fs.String("authority-holder", "", "ultimate/operational authority holder")
	policyMechanism := fs.String("policy-mechanism", "", "policy enforcement: code, engine, workflow, manual, mixed, prompt, unavailable")
	effectTypes := fs.String("effect-types", "", "comma-separated consequential effect types")
	riskCategories := fs.String("risk-categories", "", "comma-separated risk categories")
	ai := fs.Bool("ai", false, "declare AI usage")
	analysis := fs.Bool("ai-analysis", false, "AI analysis use case")
	classification := fs.Bool("ai-classification", false, "AI classification use case")
	recommendation := fs.Bool("ai-recommendation", false, "AI recommendation use case")
	toolUse := fs.Bool("ai-tool-use", false, "AI may invoke tools")
	stateMutation := fs.Bool("ai-state-mutation", false, "AI may initiate persistent state mutation")
	externalAction := fs.Bool("ai-external-action", false, "AI may initiate external API/system mutation")
	infraAction := fs.Bool("ai-infrastructure-action", false, "AI may initiate infrastructure/deployment actions")
	consequential := fs.Bool("ai-consequential", false, "AI may perform consequential execution")
	supportsEvidence := fs.Bool("supports-evidence", false, "organization supports formal Evidence lifecycle")
	supportsIndependentVerification := fs.Bool("supports-independent-verification", false, "organization supports independent Verification")
	supportsConformance := fs.Bool("supports-conformance", false, "organization supports formal conformance evidence/reporting")
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
		fmt.Fprintln(out, "AOF Context & Adoption Interview")
		fmt.Fprintln(out, "--------------------------------")
		fmt.Fprintln(out, "AOF semantics are fixed; this interview scopes their applicability to this project.")
		fmt.Fprintln(out)
		*name = ask(reader, out, "Project name", *name)
		*language = ask(reader, out, "Language", defaultValue(*language, "unspecified"))
		*projectType = ask(reader, out, "Project type", defaultValue(*projectType, "unspecified"))
		*domain = ask(reader, out, "Domain", defaultValue(*domain, "unspecified"))
		*architecture = ask(reader, out, "Architecture (monolith/modular-monolith/microservices/event-driven/serverless/other)", defaultValue(*architecture, "unspecified"))
		*deployment = ask(reader, out, "Deployment model", defaultValue(*deployment, "unspecified"))
		*criticality = ask(reader, out, "Criticality (low/moderate/high/critical)", defaultValue(*criticality, "moderate"))
		*dataSensitivity = ask(reader, out, "Data sensitivity (public/internal/confidential/restricted)", defaultValue(*dataSensitivity, "internal"))
		*profile = ask(reader, out, "AOF profile (core/governed/assured)", *profile)
		if !*ai {
			*ai = askBool(reader, out, "Will this project use AI agents", false)
		}
		if *ai {
			*agentCount = ask(reader, out, "Logical AI agent count/description", defaultValue(*agentCount, "1"))
			if !*multiAgent {
				*multiAgent = askBool(reader, out, "Will agents interact/delegate to other agents", false)
			}
			if !*analysis {
				*analysis = askBool(reader, out, "AI analysis", true)
			}
			if !*classification {
				*classification = askBool(reader, out, "AI classification", false)
			}
			if !*recommendation {
				*recommendation = askBool(reader, out, "AI recommendation", true)
			}
			if !*toolUse {
				*toolUse = askBool(reader, out, "AI tool invocation", false)
			}
			if !*stateMutation {
				*stateMutation = askBool(reader, out, "AI may initiate persistent state mutation", false)
			}
			if !*externalAction {
				*externalAction = askBool(reader, out, "AI may initiate external-system mutation", false)
			}
			if !*infraAction {
				*infraAction = askBool(reader, out, "AI may initiate infrastructure/deployment action", false)
			}
			if !*consequential {
				*consequential = askBool(reader, out, "Consequential AI execution", *stateMutation || *externalAction || *infraAction)
			}
			if *consequential || *stateMutation || *externalAction || *infraAction {
				*effectTypes = ask(reader, out, "Consequential effect types (comma-separated: database,external-api,infrastructure,account,financial,approval-state,other)", defaultValue(*effectTypes, inferEffects(*stateMutation, *externalAction, *infraAction)))
				*authorityHolder = ask(reader, out, "Who holds/ grants operational Authority", defaultValue(*authorityHolder, "unspecified"))
				if !*humanApproval {
					*humanApproval = askBool(reader, out, "Human approval required for consequential actions", true)
				}
				*policyMechanism = ask(reader, out, "Policy enforcement (code/engine/workflow/manual/mixed/prompt/unavailable)", defaultValue(*policyMechanism, "code"))
				*riskCategories = ask(reader, out, "Risk categories (comma-separated: integrity,confidentiality,financial,availability,security,infrastructure,compliance,irreversible,customer-impact)", defaultValue(*riskCategories, "integrity,availability,security"))
				if !*supportsEvidence {
					*supportsEvidence = askBool(reader, out, "Formal Evidence lifecycle available", true)
				}
				if !*supportsIndependentVerification {
					*supportsIndependentVerification = askBool(reader, out, "Independent Verification available", false)
				}
				if !*supportsConformance {
					*supportsConformance = askBool(reader, out, "Formal conformance evidence/reporting available", false)
				}
			}
		}
	}

	*language = defaultValue(*language, "unspecified")
	*projectType = defaultValue(*projectType, "unspecified")
	*domain = defaultValue(*domain, "unspecified")
	*architecture = defaultValue(*architecture, "unspecified")
	*deployment = defaultValue(*deployment, "unspecified")
	*criticality = defaultValue(*criticality, "moderate")
	*dataSensitivity = defaultValue(*dataSensitivity, "internal")
	*agentCount = defaultValue(*agentCount, func() string {
		if *ai {
			return "1"
		}
		return "0"
	}())
	*authorityHolder = defaultValue(*authorityHolder, "unspecified")
	*policyMechanism = defaultValue(*policyMechanism, "unspecified")

	p := model.NormalizeProfile(*profile)
	if p == "" {
		fmt.Fprintln(errOut, "invalid profile; use core, governed, or assured")
		return 2
	}
	aiUsage := model.AIUsage{
		Enabled:  *ai || *analysis || *classification || *recommendation || *toolUse || *stateMutation || *externalAction || *infraAction || *consequential,
		Analysis: *analysis, Classification: *classification, Recommendation: *recommendation, ToolUse: *toolUse,
		StateMutation: *stateMutation, ExternalAction: *externalAction, InfrastructureAction: *infraAction,
		ConsequentialExecution: *consequential || *stateMutation || *externalAction || *infraAction,
	}
	project := model.ProjectDefinition{
		Name: *name, Language: *language, Type: *projectType, Domain: *domain, Architecture: *architecture, Deployment: *deployment,
		Criticality: *criticality, DataSensitivity: *dataSensitivity, Features: splitList(*features), AI: aiUsage,
		AgentCount: *agentCount, MultiAgent: *multiAgent, HumanApproval: *humanApproval, AuthorityHolder: *authorityHolder,
		PolicyMechanism: *policyMechanism, RiskCategories: splitList(*riskCategories), EffectTypes: splitList(*effectTypes),
	}
	capabilities := map[string]bool{
		"authority": true, "policy": *policyMechanism != "unavailable", "risk": true, "trace": true,
		"evidence": *supportsEvidence || !aiUsage.ConsequentialExecution, "verification": true,
		"independent_verification": *supportsIndependentVerification, "formal_conformance_evidence": *supportsConformance,
	}
	adopt := adoption.Build(p, project, capabilities)
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
	printSummary(out, project, adopt, len(plan.Files))
	return 0
}

func printSummary(w io.Writer, p model.ProjectDefinition, a model.AdoptionDefinition, files int) {
	required, enabled, deferred, na := 0, 0, 0, 0
	for _, c := range a.Controls {
		if c.Applicability == "required" {
			required++
		}
		switch c.Status {
		case model.StatusEnabled:
			enabled++
		case model.StatusDeferred:
			deferred++
		case model.StatusNotApplicable:
			na++
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "AOF initialization complete.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Project")
	fmt.Fprintln(w, "-------")
	fmt.Fprintf(w, "Name: %s\nLanguage: %s\nType: %s\nArchitecture: %s\nProfile: %s\n\n", p.Name, p.Language, p.Type, p.Architecture, a.Profile)
	fmt.Fprintln(w, "AI Governance")
	fmt.Fprintln(w, "-------------")
	fmt.Fprintf(w, "AI agents: %t\nConsequential execution: %t\nMulti-agent: %t\n\n", p.AI.Enabled, p.AI.ConsequentialExecution, p.MultiAgent)
	fmt.Fprintln(w, "Adoption")
	fmt.Fprintln(w, "--------")
	fmt.Fprintf(w, "Required control families: %d\nEnabled: %d\nDeferred: %d\nNot applicable: %d\nApplicable Agent Skills: %d\nGenerated files: %d\n", required, enabled, deferred, na, len(a.RequiredSkills), files)
	if len(a.Warnings) > 0 {
		fmt.Fprintln(w, "\nGovernance warnings")
		fmt.Fprintln(w, "-------------------")
		for _, warning := range a.Warnings {
			fmt.Fprintf(w, "! %s\n", warning)
		}
	}
	fmt.Fprintln(w, "\nNext: review AOF.md and aof/conformance/gaps.md")
	fmt.Fprintln(w, "Canonical AOF: https://github.com/aof-framework/aof")
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
	seen := map[string]bool{}
	out := []string{}
	for _, s := range strings.Split(v, ",") {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
func inferEffects(state, external, infra bool) string {
	parts := []string{}
	if state {
		parts = append(parts, "database")
	}
	if external {
		parts = append(parts, "external-api")
	}
	if infra {
		parts = append(parts, "infrastructure")
	}
	if len(parts) == 0 {
		return "other"
	}
	return strings.Join(parts, ",")
}
