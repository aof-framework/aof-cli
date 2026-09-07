package model

import "strings"

const (
	ProfileCore          = "AOF-Core"
	ProfileGoverned      = "AOF-Governed"
	ProfileAssured       = "AOF-Assured"
	DomainSecureSDLC     = "AOF-Secure-SDLC"
	OverlayHighAssurance = "AOF-High-Assurance"
)

const (
	ApplicabilityApplicable    = "applicable"
	ApplicabilityConditional   = "conditional"
	ApplicabilityNotApplicable = "not_applicable"
	NormativeMust              = "MUST"
	NormativeShould            = "SHOULD"
	NormativeMay               = "MAY"
	AdoptionPlanned            = "planned"
	AdoptionDeferred           = "deferred"
	AdoptionNotApplicable      = "not_applicable"
	AdoptionUnsupported        = "unsupported"
	ImplementationNotAssessed  = "not_assessed"
	VerificationNotEvaluated   = "not_evaluated"
)

var CanonicalObjects = []string{
	"Goal", "Task", "Agent", "ContextDescriptor", "Resource", "Capability", "AuthorityGrant", "Policy",
	"RiskAssessment", "ActionProposal", "Decision", "ExecutionContract", "Evidence", "Verification", "Approval",
	"StateTransition", "TraceEvent", "AgentInteractionContract", "EscalationPackage", "Outcome", "ConformanceManifest", "ConformanceReport",
}

var AgentTypes = []string{"LLM", "Deterministic", "Human", "Hybrid", "ExternalService"}

type AIUsage struct {
	Enabled                bool
	Analysis               bool
	Classification         bool
	Recommendation         bool
	ToolUse                bool
	StateMutation          bool
	ExternalAction         bool
	InfrastructureAction   bool
	ConsequentialExecution bool
}

type ProjectDefinition struct {
	Name            string
	Language        string
	Type            string
	Domain          string
	Architecture    string
	Deployment      string
	Criticality     string
	DataSensitivity string
	Features        []string
	AI              AIUsage
	AgentTypes      []string
	MultiAgent      bool
	HumanApproval   bool
	AuthorityHolder string
	PolicyMechanism string
	RiskCategories  []string
	EffectTypes     []string
	AdoptionMode    string
}

type ControlSelection struct {
	Applicability       string
	NormativeLevel      string
	Capability          string
	AdoptionState       string
	ImplementationState string
	VerificationState   string
	Reason              string
	RequirementIDs      []string
}

type ObjectSelection struct {
	Applicability  string
	NormativeLevel string
	AdoptionState  string
	Reason         string
	RequirementIDs []string
}

type Requirement struct {
	ID                 string
	Domain             string
	NormativeLevel     string
	AppliesTo          string
	Profiles           []string
	VerificationMethod string
	RequiredEvidence   []string
}

type ProfileSelection struct {
	TargetBaseProfile string
	DomainProfiles    []string
	Overlays          []string
	ClaimedProfile    string
	ClaimStatus       string
}

type AdoptionDefinition struct {
	Profile          ProfileSelection
	Mode             string
	Controls         map[string]ControlSelection
	CanonicalObjects map[string]ObjectSelection
	Requirements     []Requirement
	RequiredSkills   []string
	Warnings         []string
}

type BootstrapModel struct {
	Project  ProjectDefinition
	Adoption AdoptionDefinition
}

func NormalizeProfile(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "core", "aof-core":
		return ProfileCore
	case "governed", "aof-governed":
		return ProfileGoverned
	case "assured", "aof-assured":
		return ProfileAssured
	default:
		return ""
	}
}

func NormalizeAgentType(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "llm", "ai":
		return "LLM"
	case "deterministic", "service", "software":
		return "Deterministic"
	case "human":
		return "Human"
	case "hybrid":
		return "Hybrid"
	case "externalservice", "external-service", "external_service":
		return "ExternalService"
	default:
		return ""
	}
}
