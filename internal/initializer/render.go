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

type Plan struct{ Files []File }

func BuildPlan(m model.BootstrapModel, cliVersion string) (Plan, error) {
	files := []File{}
	staticDocs := []struct{ src, dst string }{
		{"static/docs/semantic-boundaries.md", "docs/aof/semantic-boundaries.md"},
		{"static/docs/governance-model.md", "docs/aof/governance-model.md"},
		{"static/docs/execution-model.md", "docs/aof/execution-model.md"},
	}
	for _, item := range staticDocs {
		b, err := fs.ReadFile(bootstraptemplates.FS, item.src)
		if err != nil {
			return Plan{}, fmt.Errorf("read embedded template %s: %w", item.src, err)
		}
		files = append(files, File{Path: item.dst, Content: b})
	}

	for _, skill := range m.Adoption.RequiredSkills {
		src := "static/skills/" + skill + "/SKILL.md"
		b, err := fs.ReadFile(bootstraptemplates.FS, src)
		if err != nil {
			return Plan{}, fmt.Errorf("read embedded skill %s: %w", src, err)
		}
		files = append(files, File{Path: ".agents/skills/" + skill + "/SKILL.md", Content: b})
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
		t, err := template.New(item.src).Funcs(template.FuncMap{"join": strings.Join}).Parse(string(b))
		if err != nil {
			return Plan{}, fmt.Errorf("parse template %s: %w", item.src, err)
		}
		var out bytes.Buffer
		if err := t.Execute(&out, m); err != nil {
			return Plan{}, fmt.Errorf("render template %s: %w", item.src, err)
		}
		files = append(files, File{Path: item.dst, Content: out.Bytes()})
	}

	generated := []File{
		{Path: ".aof/project.yaml", Content: []byte(renderProjectYAML(m.Project))},
		{Path: ".aof/config.yaml", Content: []byte(renderConfigYAML(m.Adoption))},
		{Path: ".aof/adoption.yaml", Content: []byte(renderAdoptionYAML(m.Adoption))},
		{Path: ".aof/applicability.yaml", Content: []byte(renderApplicabilityYAML(m.Adoption))},
		{Path: ".aof/capabilities.yaml", Content: []byte(renderCapabilityYAML(m.Project))},
		{Path: "aof/governance/governance-root.md", Content: []byte(renderGovernanceRoot(m))},
		{Path: "aof/governance/governance-envelope.yaml", Content: []byte(renderGovernanceEnvelope(m))},
		{Path: "aof/agents/inventory.yaml", Content: []byte(renderAgentInventory(m))},
		{Path: "aof/agents/capabilities.yaml", Content: []byte(renderAgentCapabilities(m))},
		{Path: "aof/authority/authority-model.yaml", Content: []byte(renderAuthorityModel(m))},
		{Path: "aof/authority/README.md", Content: []byte(renderAuthorityReadme(m))},
		{Path: "aof/policies/policy-model.yaml", Content: []byte(renderPolicyModel(m))},
		{Path: "aof/risk/risk-model.yaml", Content: []byte(renderRiskModel(m))},
		{Path: "aof/execution/execution-model.md", Content: []byte(renderExecutionModel(m))},
		{Path: "aof/execution/consequential-actions.yaml", Content: []byte(renderConsequentialActions(m))},
		{Path: "aof/execution/state-transitions.yaml", Content: []byte(renderStateTransitions(m))},
		{Path: "aof/assurance/evidence-model.md", Content: []byte(renderEvidenceModel(m))},
		{Path: "aof/assurance/verification-model.md", Content: []byte(renderVerificationModel(m))},
		{Path: "aof/assurance/trace-model.md", Content: []byte(renderTraceModel(m))},
		{Path: "aof/conformance/scope.yaml", Content: []byte(renderConformanceScope(m))},
		{Path: "aof/conformance/gaps.md", Content: []byte(renderGaps(m))},
		{Path: "aof/fixtures/.gitkeep", Content: []byte{}},
		{Path: "aof/profiles/.gitkeep", Content: []byte{}},
	}
	files = append(files, generated...)
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
	fmt.Fprintf(&b, "project:\n  name: %q\n  language: %q\n  type: %q\n  domain: %q\n  architecture: %q\n  deployment: %q\n  criticality: %q\n  data_sensitivity: %q\n", p.Name, p.Language, p.Type, p.Domain, p.Architecture, p.Deployment, p.Criticality, p.DataSensitivity)
	writeList(&b, "features", p.Features)
	fmt.Fprintf(&b, "ai:\n  enabled: %t\n  agent_count: %q\n  multi_agent: %t\n  analysis: %t\n  classification: %t\n  recommendation: %t\n  tool_use: %t\n  state_mutation: %t\n  external_action: %t\n  infrastructure_action: %t\n  consequential_execution: %t\n", p.AI.Enabled, p.AgentCount, p.MultiAgent, p.AI.Analysis, p.AI.Classification, p.AI.Recommendation, p.AI.ToolUse, p.AI.StateMutation, p.AI.ExternalAction, p.AI.InfrastructureAction, p.AI.ConsequentialExecution)
	writeList(&b, "effect_types", p.EffectTypes)
	writeList(&b, "risk_categories", p.RiskCategories)
	return b.String()
}

func renderConfigYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "aof:\n  specification: \"1.0\"\n  release: \"LTS\"\n  profile: %q\n", a.Profile)
	b.WriteString("upstream:\n  repository: \"https://github.com/aof-framework/aof\"\n  specification: \"specification/AOF-v1.0-Framework-Specification.md\"\n")
	fmt.Fprintf(&b, "adoption:\n  mode: %q\ncontrols:\n", a.Mode)
	keys := sortedControlKeys(a.Controls)
	for _, k := range keys {
		c := a.Controls[k]
		fmt.Fprintf(&b, "  %s:\n    applicability: %q\n    capability: %q\n    status: %q\n", k, c.Applicability, c.Capability, c.Status)
		if c.Reason != "" {
			fmt.Fprintf(&b, "    reason: %q\n", c.Reason)
		}
	}
	return b.String()
}

func renderAdoptionYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "profile: %q\nmode: %q\nprinciples:\n  - \"AOF semantics are static\"\n  - \"Project adoption is scoped\"\n  - \"Deferred != Satisfied\"\n  - \"NotApplicable != DisabledForConvenience\"\n  - \"PartialImplementation != FullAOFConformance\"\n", a.Profile, a.Mode)
	writeList(&b, "required_skills", a.RequiredSkills)
	writeList(&b, "warnings", a.Warnings)
	return b.String()
}
func renderApplicabilityYAML(a model.AdoptionDefinition) string {
	var b strings.Builder
	b.WriteString("canonical_objects:\n")
	keys := make([]string, 0, len(a.CanonicalObjects))
	for k := range a.CanonicalObjects {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		o := a.CanonicalObjects[k]
		fmt.Fprintf(&b, "  %s:\n    status: %q\n    reason: %q\n", k, o.Status, o.Reason)
	}
	return b.String()
}
func renderCapabilityYAML(p model.ProjectDefinition) string {
	var b strings.Builder
	b.WriteString("project_capabilities:\n")
	pairs := []struct {
		name string
		v    bool
	}{{"analysis", p.AI.Analysis}, {"classification", p.AI.Classification}, {"recommendation", p.AI.Recommendation}, {"tool_use", p.AI.ToolUse}, {"state_mutation", p.AI.StateMutation}, {"external_action", p.AI.ExternalAction}, {"infrastructure_action", p.AI.InfrastructureAction}, {"consequential_execution", p.AI.ConsequentialExecution}}
	for _, x := range pairs {
		fmt.Fprintf(&b, "  %s: %t\n", x.name, x.v)
	}
	return b.String()
}

func renderGovernanceRoot(m model.BootstrapModel) string {
	return fmt.Sprintf("# Governance Root — %s\n\nAOF v1.0 LTS requires the Human/Organization to remain the GovernanceRoot. Agent autonomy is bounded by the Human/Organizational GovernanceEnvelope.\n\n- Declared Authority holder: **%s**\n- AOF profile: **%s**\n- Consequential AI execution: **%t**\n\nAI Agents MUST NOT self-grant Authority or expand their GovernanceEnvelope.\n", m.Project.Name, m.Project.AuthorityHolder, m.Adoption.Profile, m.Project.AI.ConsequentialExecution)
}
func renderGovernanceEnvelope(m model.BootstrapModel) string {
	return fmt.Sprintf("governance_root: %q\nagent_autonomy_bounded: true\nauthority_holder: %q\nprofile: %q\nconstraints:\n  capability_implies_authority: false\n  proposal_is_decision: false\n  approval_is_authority_grant: false\n", "Human/Organization", m.Project.AuthorityHolder, m.Adoption.Profile)
}
func renderAgentInventory(m model.BootstrapModel) string {
	return fmt.Sprintf("ai_enabled: %t\nlogical_agent_count: %q\nmulti_agent: %t\nagent_output_trust: \"UntrustedProposal\"\ninteraction_contract_status: %q\n", m.Project.AI.Enabled, m.Project.AgentCount, m.Project.MultiAgent, m.Adoption.CanonicalObjects["AgentInteractionContract"].Status)
}
func renderAgentCapabilities(m model.BootstrapModel) string {
	return renderCapabilityYAML(m.Project) + "\nrules:\n  capability_is_authority: false\n  self_grant_authority: false\n"
}
func renderAuthorityModel(m model.BootstrapModel) string {
	return fmt.Sprintf("authority:\n  holder: %q\n  governance_root: \"Human/Organization\"\n  human_approval: %t\n  consequential_execution: %t\n  principles:\n    - \"Capability != Authority\"\n    - \"Approval != AuthorityGrant\"\n    - \"NoGrant => NoAuthoritySensitiveExecution\"\n", m.Project.AuthorityHolder, m.Project.HumanApproval, m.Project.AI.ConsequentialExecution)
}
func renderAuthorityReadme(m model.BootstrapModel) string {
	return fmt.Sprintf("# Authority Model\n\nDeclared authority holder: **%s**.\n\nFor consequential execution, implementation MUST establish Authority independently from Capability and Approval. An AI Agent's ability to invoke a tool does not create Authority to use it.\n", m.Project.AuthorityHolder)
}
func renderPolicyModel(m model.BootstrapModel) string {
	return fmt.Sprintf("policy:\n  enforcement_mechanism: %q\n  policy_prompt_is_enforcement: false\n  consequential_execution: %t\n  unresolved_mandatory_policy: \"Pending\"\n  pending_allows_execution: false\n", m.Project.PolicyMechanism, m.Project.AI.ConsequentialExecution)
}
func renderRiskModel(m model.BootstrapModel) string {
	var b strings.Builder
	fmt.Fprintf(&b, "risk:\n  project_criticality: %q\n  data_sensitivity: %q\n", m.Project.Criticality, m.Project.DataSensitivity)
	writeListIndented(&b, "  categories", m.Project.RiskCategories)
	b.WriteString("  principles:\n    - \"RiskAssessment != RiskAcceptance\"\n    - \"mandatory Pending != Pass\"\n")
	return b.String()
}
func renderExecutionModel(m model.BootstrapModel) string {
	return fmt.Sprintf("# Execution Model\n\nConsequential execution declared: **%t**.\n\nCanonical flow:\n\n```text\nLLM/Agent output\n-> UntrustedProposal\n-> Governance Evaluation\n-> AuthorizedDecision\n-> ExecutionContract\n-> Effect Boundary\n-> Effect\n-> Evidence\n-> Verification\n-> Trace\n```\n\nExecution MUST satisfy the AOF control predicate: `ExecuteAllowed = C AND H AND P AND S AND R AND V`. Mandatory `Fail` blocks execution; mandatory `Pending` remains Pending and MUST NOT become implicit Allow.\n", m.Project.AI.ConsequentialExecution)
}
func renderConsequentialActions(m model.BootstrapModel) string {
	var b strings.Builder
	fmt.Fprintf(&b, "consequential_execution: %t\neffect_types:\n", m.Project.AI.ConsequentialExecution)
	if len(m.Project.EffectTypes) == 0 {
		b.WriteString("  []\n")
	} else {
		for _, x := range m.Project.EffectTypes {
			fmt.Fprintf(&b, "  - %q\n", x)
		}
	}
	b.WriteString("requirements:\n  proposal_before_decision: true\n  authority_required: true\n  policy_required: true\n  risk_required: true\n  state_validation_required: true\n  verification_required: true\n  trace_required: true\n")
	return b.String()
}
func renderStateTransitions(m model.BootstrapModel) string {
	return fmt.Sprintf("state_model:\n  consequential_execution: %t\n  explicit_state_validation: %t\n  transition_authorization_required: %t\n  pending_is_pass: false\n", m.Project.AI.ConsequentialExecution, m.Project.AI.ConsequentialExecution, m.Project.AI.ConsequentialExecution)
}
func renderEvidenceModel(m model.BootstrapModel) string {
	c := m.Adoption.Controls["evidence"]
	return fmt.Sprintf("# Evidence Model\n\nStatus: **%s** (`%s`).\n\nEvidence is not a Claim and is not Verification. When Evidence is required, implementations SHOULD capture objective artifacts sufficient for the configured Verification step.\n", c.Status, c.Reason)
}
func renderVerificationModel(m model.BootstrapModel) string {
	c := m.Adoption.Controls["verification"]
	return fmt.Sprintf("# Verification Model\n\nStatus: **%s**.\n\nVerification evaluates Evidence against acceptance conditions. Mandatory unresolved Verification remains `Pending`; it MUST NOT be treated as `Pass`. Independent Verification status is `%s`.\n", c.Status, m.Adoption.Controls["independent_verification"].Status)
}
func renderTraceModel(m model.BootstrapModel) string {
	return "# Trace Model\n\nTrace SHOULD preserve governance-relevant chronology across Proposal, Decision, Authority, execution, Evidence, Verification, StateTransition, and Outcome. Trace is an Assurance-plane artifact and MUST not itself grant Authority.\n"
}
func renderConformanceScope(m model.BootstrapModel) string {
	return fmt.Sprintf("aof_specification: \"1.0 LTS\"\nprofile: %q\nscoped_adoption: true\nfull_conformance_claimed: false\nprinciple: \"ScopedAdoption -> ScopedConformance\"\n", m.Adoption.Profile)
}
func renderGaps(m model.BootstrapModel) string {
	var b strings.Builder
	b.WriteString("# AOF Adoption Gaps\n\nThis file makes incomplete or deferred adoption explicit. It is not a failure by itself and MUST NOT be converted into a full conformance claim.\n\n")
	n := 0
	keys := sortedControlKeys(m.Adoption.Controls)
	for _, k := range keys {
		c := m.Adoption.Controls[k]
		if c.Status == model.StatusDeferred || c.Status == model.StatusUnsupported {
			n++
			fmt.Fprintf(&b, "- **%s** — `%s`; applicability `%s`; capability `%s`. %s\n", k, c.Status, c.Applicability, c.Capability, c.Reason)
		}
	}
	if n == 0 {
		b.WriteString("No Deferred or Unsupported controls were declared during initialization.\n")
	}
	if len(m.Adoption.Warnings) > 0 {
		b.WriteString("\n## Warnings\n\n")
		for _, w := range m.Adoption.Warnings {
			fmt.Fprintf(&b, "- %s\n", w)
		}
	}
	return b.String()
}

func writeList(b *strings.Builder, key string, vals []string) {
	fmt.Fprintf(b, "%s:\n", key)
	if len(vals) == 0 {
		b.WriteString("  []\n")
		return
	}
	for _, v := range vals {
		fmt.Fprintf(b, "  - %q\n", v)
	}
}
func writeListIndented(b *strings.Builder, key string, vals []string) {
	fmt.Fprintf(b, "%s:\n", key)
	if len(vals) == 0 {
		b.WriteString("    []\n")
		return
	}
	for _, v := range vals {
		fmt.Fprintf(b, "    - %q\n", v)
	}
}
func sortedControlKeys(m map[string]model.ControlSelection) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func safePath(path string) bool {
	if path == "" {
		return false
	}
	if filepath.IsAbs(path) || filepath.VolumeName(path) != "" || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
		return false
	}
	if len(path) >= 2 && ((path[0] >= 'a' && path[0] <= 'z') || (path[0] >= 'A' && path[0] <= 'Z')) && path[1] == ':' {
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
