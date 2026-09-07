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

const Version = "0.3.2"

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
	fmt.Fprintln(w, "AOF CLI — Semantic Fidelity & Adoption Hardening")
	fmt.Fprintln(w, "\nUsage:\n  aof init [flags]\n  aof version\n  aof help")
	fmt.Fprintln(w, "\nKey init flags:")
	fmt.Fprintln(w, "  --profile core|governed|assured")
	fmt.Fprintln(w, "  --secure-sdlc                 add AOF-Secure-SDLC domain profile")
	fmt.Fprintln(w, "  --high-assurance              add AOF-High-Assurance overlay")
	fmt.Fprintln(w, "  --agent-types <csv>           LLM,Deterministic,Human,Hybrid,ExternalService")
	fmt.Fprintln(w, "  --non-interactive")
	fmt.Fprintln(w, "\naof init scopes canonical AOF v1.0 LTS requirements; it does not claim implementation or conformance.")
}

func printVersion(w io.Writer) {
	fmt.Fprintf(w, "AOF CLI v%s\n", Version)
	fmt.Fprintln(w, "Release: Exact Canonical Source Alignment")
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
	criticality := fs.String("criticality", "", "low, moderate, high, critical")
	dataSensitivity := fs.String("data-sensitivity", "", "public, internal, confidential, restricted")
	profile := fs.String("profile", "core", "core, governed, assured")
	secureSDLC := fs.Bool("secure-sdlc", false, "compose AOF-Secure-SDLC domain profile")
	highAssurance := fs.Bool("high-assurance", false, "compose AOF-High-Assurance overlay")
	features := fs.String("features", "", "comma-separated project features")
	agentTypes := fs.String("agent-types", "", "comma-separated Agent types")
	multiAgent := fs.Bool("multi-agent", false, "declare Agent interaction/delegation")
	humanApproval := fs.Bool("human-approval", false, "declare Human approval gate")
	authorityHolder := fs.String("authority-holder", "", "Authority holder/granting authority")
	policyMechanism := fs.String("policy-mechanism", "", "code, engine, workflow, manual, mixed, prompt, unavailable")
	effectTypes := fs.String("effect-types", "", "comma-separated consequential effect types")
	riskCategories := fs.String("risk-categories", "", "comma-separated risk categories")
	ai := fs.Bool("ai", false, "declare LLM/AI Agent usage")
	analysis := fs.Bool("ai-analysis", false, "AI analysis")
	classification := fs.Bool("ai-classification", false, "AI classification")
	recommendation := fs.Bool("ai-recommendation", false, "AI recommendation")
	toolUse := fs.Bool("ai-tool-use", false, "AI tool invocation")
	stateMutation := fs.Bool("ai-state-mutation", false, "AI may initiate persistent state mutation")
	externalAction := fs.Bool("ai-external-action", false, "AI may initiate external mutation")
	infraAction := fs.Bool("ai-infrastructure-action", false, "AI may initiate infrastructure action")
	consequential := fs.Bool("ai-consequential", false, "AI may initiate consequential execution")
	supportsEvidence := fs.Bool("supports-evidence", false, "organization has formal Evidence lifecycle capability")
	supportsIndependentVerification := fs.Bool("supports-independent-verification", false, "organization has independent Verification capability")
	supportsConformance := fs.Bool("supports-conformance", false, "organization has formal conformance reporting capability")
	nonInteractive := fs.Bool("non-interactive", false, "disable interactive prompts")
	if err := fs.Parse(args); err != nil {
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

	if !*nonInteractive {
		fmt.Fprintln(out, "AOF v0.3 Adoption Interview")
		fmt.Fprintln(out, "---------------------------")
		fmt.Fprintln(out, "AOF semantics are fixed. Answers scope applicability; they do NOT prove implementation or conformance.")
		fmt.Fprintln(out)
		*name = ask(reader, out, "Project name", *name)
		*language = ask(reader, out, "Language", defaultValue(*language, "unspecified"))
		*projectType = ask(reader, out, "Project type", defaultValue(*projectType, "unspecified"))
		*domain = ask(reader, out, "Domain", defaultValue(*domain, "unspecified"))
		*architecture = ask(reader, out, "Architecture (monolith/modular-monolith/microservices/event-driven/serverless/other)", defaultValue(*architecture, "unspecified"))
		*deployment = ask(reader, out, "Deployment model", defaultValue(*deployment, "unspecified"))
		*criticality = ask(reader, out, "Criticality (low/moderate/high/critical)", defaultValue(*criticality, "moderate"))
		*dataSensitivity = ask(reader, out, "Data sensitivity (public/internal/confidential/restricted)", defaultValue(*dataSensitivity, "internal"))
		*profile = ask(reader, out, "Base AOF profile (core/governed/assured)", *profile)
		if !*secureSDLC {
			*secureSDLC = askBool(reader, out, "Apply AOF-Secure-SDLC domain profile", isSoftwareProject(*projectType))
		}
		if !*highAssurance {
			*highAssurance = askBool(reader, out, "Apply AOF-High-Assurance strengthening", false)
		}
		if *agentTypes == "" {
			*agentTypes = ask(reader, out, "Operational Agent types (comma-separated: LLM,Deterministic,Human,Hybrid,ExternalService)", "Deterministic")
		}
		if !*ai {
			*ai = containsAgentType(*agentTypes, "LLM") || askBool(reader, out, "Will this project use LLM/AI Agents", false)
		}
		if *ai && !containsAgentType(*agentTypes, "LLM") {
			*agentTypes = appendCSV(*agentTypes, "LLM")
		}
		if !*multiAgent {
			*multiAgent = askBool(reader, out, "Will Agents interact or delegate to other Agents", false)
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
		}
		if *consequential || *stateMutation || *externalAction || *infraAction {
			*effectTypes = ask(reader, out, "Consequential effect types (comma-separated)", defaultValue(*effectTypes, inferEffects(*stateMutation, *externalAction, *infraAction)))
			*authorityHolder = ask(reader, out, "Who holds/grants operational Authority", defaultValue(*authorityHolder, "unspecified"))
			if !*humanApproval {
				*humanApproval = askBool(reader, out, "Human approval required for consequential actions", true)
			}
			*policyMechanism = ask(reader, out, "Policy enforcement (code/engine/workflow/manual/mixed/prompt/unavailable)", defaultValue(*policyMechanism, "code"))
			*riskCategories = ask(reader, out, "Risk categories (comma-separated)", defaultValue(*riskCategories, "integrity,availability,security"))
			if !*supportsEvidence {
				*supportsEvidence = askBool(reader, out, "Formal Evidence lifecycle capability available", true)
			}
			if !*supportsIndependentVerification {
				*supportsIndependentVerification = askBool(reader, out, "Independent Verification capability available", false)
			}
			if !*supportsConformance {
				*supportsConformance = askBool(reader, out, "Formal conformance reporting capability available", false)
			}
		}
	}

	*language = defaultValue(*language, "unspecified")
	*projectType = defaultValue(*projectType, "unspecified")
	*domain = defaultValue(*domain, "unspecified")
	*architecture = defaultValue(*architecture, "unspecified")
	*deployment = defaultValue(*deployment, "unspecified")
	*criticality = strings.ToLower(defaultValue(*criticality, "moderate"))
	*dataSensitivity = strings.ToLower(defaultValue(*dataSensitivity, "internal"))
	*authorityHolder = defaultValue(*authorityHolder, "unspecified")
	*policyMechanism = strings.ToLower(defaultValue(*policyMechanism, "unspecified"))

	p := model.NormalizeProfile(*profile)
	if p == "" {
		fmt.Fprintln(errOut, "invalid profile; use core, governed, or assured")
		return 2
	}
	if *secureSDLC && p == model.ProfileCore {
		fmt.Fprintln(errOut, "AOF-Secure-SDLC requires a governed or assured base profile")
		return 2
	}
	if *highAssurance && p != model.ProfileAssured {
		fmt.Fprintln(errOut, "AOF-High-Assurance requires the assured base profile")
		return 2
	}
	if !oneOf(*criticality, "low", "moderate", "high", "critical") {
		fmt.Fprintln(errOut, "invalid criticality; use low, moderate, high, or critical")
		return 2
	}
	if !oneOf(*dataSensitivity, "public", "internal", "confidential", "restricted") {
		fmt.Fprintln(errOut, "invalid data sensitivity; use public, internal, confidential, or restricted")
		return 2
	}
	if !oneOf(*policyMechanism, "unspecified", "code", "engine", "workflow", "manual", "mixed", "prompt", "unavailable") {
		fmt.Fprintln(errOut, "invalid policy mechanism")
		return 2
	}

	types, err := normalizeAgentTypes(*agentTypes, *ai)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	aiUsage := model.AIUsage{Enabled: *ai || contains(types, "LLM"), Analysis: *analysis, Classification: *classification, Recommendation: *recommendation, ToolUse: *toolUse, StateMutation: *stateMutation, ExternalAction: *externalAction, InfrastructureAction: *infraAction, ConsequentialExecution: *consequential || *stateMutation || *externalAction || *infraAction}
	mode := "greenfield"
	if detected.ExistingProject {
		mode = "brownfield"
	}
	project := model.ProjectDefinition{Name: *name, Language: *language, Type: *projectType, Domain: *domain, Architecture: *architecture, Deployment: *deployment, Criticality: *criticality, DataSensitivity: *dataSensitivity, Features: splitList(*features), AI: aiUsage, AgentTypes: types, MultiAgent: *multiAgent, HumanApproval: *humanApproval, AuthorityHolder: *authorityHolder, PolicyMechanism: *policyMechanism, RiskCategories: splitList(*riskCategories), EffectTypes: splitList(*effectTypes), AdoptionMode: mode}

	capabilities := map[string]bool{}
	if *policyMechanism == "unavailable" {
		capabilities["policy"] = false
	}
	if *supportsEvidence {
		capabilities["evidence"] = true
		capabilities["evidence_provenance"] = true
	}
	if *supportsIndependentVerification {
		capabilities["verifier_independence"] = true
	}
	if *supportsConformance {
		capabilities["conformance_manifest"] = true
		capabilities["conformance_report"] = true
	}
	domainProfiles := []string{}
	if *secureSDLC {
		domainProfiles = append(domainProfiles, model.DomainSecureSDLC)
	}
	adopt := adoption.BuildWithOptions(p, project, capabilities, adoption.Options{DomainProfiles: domainProfiles, HighAssurance: *highAssurance})
	plan, err := initializer.BuildPlan(model.BootstrapModel{Project: project, Adoption: adopt}, Version)
	if err != nil {
		fmt.Fprintf(errOut, "build initialization plan: %v\n", err)
		return 1
	}
	plan, agentsSidecar, err := initializer.PreserveExistingAgents(root, plan)
	if err != nil {
		fmt.Fprintf(errOut, "preserve existing AGENTS.md: %v\n", err)
		return 1
	}
	if agentsSidecar {
		adopt.Warnings = append(adopt.Warnings, "existing AGENTS.md preserved; review and merge generated AGENTS.aof.md explicitly")
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
	applicable, planned, unsupported, conditional, notEvaluated := 0, 0, 0, 0, 0
	for _, c := range a.Controls {
		if c.Applicability == model.ApplicabilityApplicable {
			applicable++
		}
		if c.Applicability == model.ApplicabilityConditional {
			conditional++
		}
		if c.Applicability == model.ApplicabilityNotEvaluated {
			notEvaluated++
		}
		switch c.AdoptionState {
		case model.AdoptionPlanned:
			planned++
		case model.AdoptionUnsupported:
			unsupported++
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "AOF initialization complete.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Project")
	fmt.Fprintln(w, "-------")
	fmt.Fprintf(w, "Name: %s\nLanguage: %s\nType: %s\nArchitecture: %s\nAdoption mode: %s\n\n", p.Name, p.Language, p.Type, p.Architecture, a.Mode)
	fmt.Fprintln(w, "AOF Profile Target")
	fmt.Fprintln(w, "------------------")
	fmt.Fprintf(w, "Base: %s\nDomain: %s\nOverlay: %s\nClaimed conformance: none\n\n", a.Profile.TargetBaseProfile, emptyDash(strings.Join(a.Profile.DomainProfiles, ", ")), emptyDash(strings.Join(a.Profile.Overlays, ", ")))
	fmt.Fprintln(w, "Governance Scope")
	fmt.Fprintln(w, "----------------")
	fmt.Fprintf(w, "Agent types: %s\nConsequential execution: %t\nMulti-Agent: %t\nApplicable controls: %d\nConditional controls: %d\nNot-evaluated controls: %d\nPlanned controls: %d\nUnsupported mandatory/selected capability controls: %d\nRequirements indexed: %d\nGenerated files: %d\n", strings.Join(p.AgentTypes, ", "), p.AI.ConsequentialExecution, p.MultiAgent, applicable, conditional, notEvaluated, planned, unsupported, len(a.Requirements), files)
	if len(a.Warnings) > 0 {
		fmt.Fprintln(w, "\nGovernance warnings\n-------------------")
		for _, warning := range a.Warnings {
			fmt.Fprintf(w, "! %s\n", warning)
		}
	}
	fmt.Fprintln(w, "\nImportant: generated adoption states are plans, NOT implementation or conformance results.")
	fmt.Fprintln(w, "Next: review AOF.md, .aof/applicability.yaml, .aof/requirements.yaml, and aof/conformance/gaps.md")
	fmt.Fprintln(w, "Canonical AOF: https://github.com/aof-framework/aof")
}

func normalizeAgentTypes(v string, ai bool) ([]string, error) {
	raw := splitList(v)
	if len(raw) == 0 {
		if ai {
			raw = []string{"LLM"}
		} else {
			raw = []string{"Deterministic"}
		}
	}
	seen := map[string]bool{}
	out := []string{}
	for _, x := range raw {
		n := model.NormalizeAgentType(x)
		if n == "" {
			return nil, fmt.Errorf("invalid Agent type %q; use LLM, Deterministic, Human, Hybrid, or ExternalService", x)
		}
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out, nil
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
func oneOf(v string, vals ...string) bool {
	for _, x := range vals {
		if v == x {
			return true
		}
	}
	return false
}
func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
func containsAgentType(csv, want string) bool {
	for _, x := range splitList(csv) {
		if model.NormalizeAgentType(x) == want {
			return true
		}
	}
	return false
}
func appendCSV(v, x string) string {
	if strings.TrimSpace(v) == "" {
		return x
	}
	return v + "," + x
}
func isSoftwareProject(t string) bool {
	s := strings.ToLower(t)
	return strings.Contains(s, "backend") || strings.Contains(s, "frontend") || strings.Contains(s, "fullstack") || strings.Contains(s, "software") || strings.Contains(s, "api")
}
func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}
