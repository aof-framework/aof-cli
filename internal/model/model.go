package model

import "strings"

const (
	ProfileCore     = "AOF-Core"
	ProfileGoverned = "AOF-Governed"
	ProfileAssured  = "AOF-Assured"
)

const (
	StatusRequired      = "required"
	StatusEnabled       = "enabled"
	StatusDeferred      = "deferred"
	StatusNotApplicable = "not_applicable"
	StatusUnsupported   = "unsupported"
)

var CanonicalObjects = []string{
	"Goal", "Task", "Agent", "ContextDescriptor", "Resource", "Capability", "AuthorityGrant", "Policy",
	"RiskAssessment", "ActionProposal", "Decision", "ExecutionContract", "Evidence", "Verification", "Approval",
	"StateTransition", "TraceEvent", "AgentInteractionContract", "EscalationPackage", "Outcome", "ConformanceManifest", "ConformanceReport",
}

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
	AgentCount      string
	MultiAgent      bool
	HumanApproval   bool
	AuthorityHolder string
	PolicyMechanism string
	RiskCategories  []string
	EffectTypes     []string
}

type ControlSelection struct {
	Applicability string
	Capability    string
	Status        string
	Reason        string
}

type ObjectSelection struct {
	Status string
	Reason string
}

type AdoptionDefinition struct {
	Profile          string
	Mode             string
	Controls         map[string]ControlSelection
	CanonicalObjects map[string]ObjectSelection
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
