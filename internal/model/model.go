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

type AIUsage struct {
	Enabled                bool
	Analysis               bool
	Classification         bool
	Recommendation         bool
	ConsequentialExecution bool
}

type ProjectDefinition struct {
	Name     string
	Language string
	Type     string
	Domain   string
	Features []string
	AI       AIUsage
}

type ControlSelection struct {
	Applicability string
	Capability    string
	Status        string
	Reason        string
}

type AdoptionDefinition struct {
	Profile  string
	Mode     string
	Controls map[string]ControlSelection
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
