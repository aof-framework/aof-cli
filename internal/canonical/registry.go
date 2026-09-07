// Package canonical exposes the embedded, pinned AOF v1.0 LTS semantic
// universe. It is source data for deterministic adoption projection; it is
// not a runtime policy or conformance engine.
package canonical

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const (
	ExpectedInvariantCount          = 162
	ExpectedRequirementCount        = 332
	ExpectedSemanticIDCount         = ExpectedInvariantCount + ExpectedRequirementCount
	ExpectedTraceabilityRecordCount = 325
	ExpectedMappedRequirementCount  = 238
	ExpectedMappedInvariantCount    = 106
	SourceCommit                    = "58ebca64759e2ec66da7a6cab05c3881e39127ca"
	PublishedSpecificationSHA256    = "6197f71416984cca1811d5cd0cdd30327cccb9026acc5b92509f8f2a3a137974"
	SemanticFreezeSHA256            = "57ddbd64671eea615535b20f109064d96fb262e781969ef757a6f4d5efa869d5"
	expectedInvariantIDSetSHA256    = "7da193f259dce92ed680ce051e08cb66f9d579e4929e3bc6ea7dac9fc1bc713c"
	expectedRequirementIDSetSHA256  = "b29b1f5d30adc1523e8727dce20dbb302f0fc52d6dc7cd1c44bca3aee6c865c9"
)

//go:embed data/invariants.json
var invariantData []byte

//go:embed data/requirements.json
var requirementData []byte

//go:embed data/traceability.json
var traceabilityData []byte

//go:embed data/profiles.json
var profileData []byte

//go:embed data/canonical_objects.json
var canonicalObjectData []byte

type Source struct {
	Specification string `json:"specification"`
	Commit        string `json:"commit"`
	SHA256        string `json:"sha256"`
	Line          int    `json:"line"`
}

type Invariant struct {
	ID                    string   `json:"invariant_id"`
	Title                 string   `json:"title"`
	DomainAliases         []string `json:"domain_aliases"`
	SemanticDomains       []string `json:"semantic_domains"`
	Enforcement           string   `json:"enforcement"`
	CanonicalRecord       string   `json:"canonical_record"`
	EnforcingRequirements []string `json:"enforcing_requirements"`
	ReferenceTests        []string `json:"reference_tests"`
	TraceabilityStatus    string   `json:"traceability_status"`
	Source                Source   `json:"source"`
}

type RequirementSource struct {
	Specification string `json:"specification"`
	Line          int    `json:"line"`
	SelectionRule string `json:"selection_rule"`
}

type Requirement struct {
	ID                          string            `json:"requirement_id"`
	Domain                      string            `json:"domain"`
	Statement                   string            `json:"statement"`
	NormativeLevel              string            `json:"normative_level"`
	NormativeBasis              string            `json:"normative_basis"`
	Source                      RequirementSource `json:"source"`
	AppliesTo                   []string          `json:"applies_to"`
	Profiles                    []string          `json:"profile"`
	VerificationMethods         []string          `json:"verification_method"`
	RequiredEvidence            []string          `json:"required_evidence"`
	RelatedInvariants           []string          `json:"related_invariants"`
	RelatedRequirements         []string          `json:"related_requirements"`
	RelatedCanonicalObjects     []string          `json:"related_canonical_objects"`
	ApplicabilityRule           string            `json:"applicability_rule"`
	VerificationDisposition     string            `json:"verification_disposition"`
	EvidenceDisposition         string            `json:"evidence_disposition"`
	InvariantMappingDisposition string            `json:"invariant_mapping_disposition"`
	ObjectMappingDisposition    string            `json:"object_mapping_disposition"`
	CanonicalTraceabilityRecord string            `json:"canonical_traceability_record"`
	CanonicalTraceabilityLine   int               `json:"canonical_traceability_source_line"`
	Status                      string            `json:"status"`
}

type Edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`
	Source string `json:"source"`
}

type Profile struct {
	ID       string   `json:"id"`
	Kind     string   `json:"kind"`
	Includes []string `json:"includes"`
	Summary  string   `json:"summary"`
}

type CanonicalObject struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type Registry struct {
	Invariants       []Invariant
	Requirements     []Requirement
	Edges            []Edge
	Profiles         []Profile
	CanonicalObjects []CanonicalObject
	requirements     map[string]Requirement
}

var (
	once            sync.Once
	loaded          Registry
	loadedErr       error
	requirementIDRE = regexp.MustCompile(`^AOF-[A-Z]+-[0-9]{3}$`)
	invariantIDRE   = regexp.MustCompile(`^AOF-INV-[0-9]{3}$`)
	referenceTestRE = regexp.MustCompile(`^CT-[A-Z]+-[0-9]{3}$`)
)

func Load() (Registry, error) {
	once.Do(func() { loaded, loadedErr = decodeAndValidate() })
	return loaded, loadedErr
}

// Validate executes the canonical completeness and traceability hard gates.
func Validate() error {
	_, err := Load()
	return err
}

func MustLoad() Registry {
	registry, err := Load()
	if err != nil {
		panic("invalid embedded AOF canonical registry: " + err.Error())
	}
	return registry
}

func (r Registry) Requirement(id string) (Requirement, bool) {
	requirement, ok := r.requirements[id]
	return requirement, ok
}

func decodeAndValidate() (Registry, error) {
	var invariantDocument struct {
		Count      int         `json:"count"`
		Source     Source      `json:"source"`
		Invariants []Invariant `json:"invariants"`
	}
	var requirementDocument struct {
		Count                        int           `json:"count"`
		ImportedFromCommit           string        `json:"imported_from_commit"`
		PublishedSpecificationSHA256 string        `json:"published_specification_sha256"`
		SemanticFreezeSHA256         string        `json:"semantic_freeze_sha256"`
		Requirements                 []Requirement `json:"requirements"`
	}
	var traceabilityDocument struct {
		Counts struct {
			Requirements        int `json:"requirements"`
			CanonicalInvariants int `json:"canonical_invariants"`
			ExplicitFrozenEdges int `json:"explicit_frozen_edges"`
		} `json:"counts"`
		DirectFrozenEdges []Edge `json:"direct_frozen_edges"`
	}
	var profileDocument struct {
		Profiles []Profile `json:"profiles"`
	}
	var objectDocument struct {
		Count   int               `json:"count"`
		Objects []CanonicalObject `json:"objects"`
	}
	for _, item := range []struct {
		name string
		data []byte
		dst  any
	}{{"invariants", invariantData, &invariantDocument}, {"requirements", requirementData, &requirementDocument}, {"traceability", traceabilityData, &traceabilityDocument}, {"profiles", profileData, &profileDocument}, {"canonical objects", canonicalObjectData, &objectDocument}} {
		if err := json.Unmarshal(item.data, item.dst); err != nil {
			return Registry{}, fmt.Errorf("decode %s registry: %w", item.name, err)
		}
	}
	r := Registry{Invariants: invariantDocument.Invariants, Requirements: requirementDocument.Requirements, Edges: traceabilityDocument.DirectFrozenEdges, Profiles: profileDocument.Profiles, CanonicalObjects: objectDocument.Objects, requirements: map[string]Requirement{}}
	var problems []error
	if invariantDocument.Count != ExpectedInvariantCount || len(r.Invariants) != ExpectedInvariantCount {
		problems = append(problems, fmt.Errorf("registered invariants: declared=%d actual=%d expected=%d", invariantDocument.Count, len(r.Invariants), ExpectedInvariantCount))
	}
	if requirementDocument.Count != ExpectedRequirementCount || len(r.Requirements) != ExpectedRequirementCount {
		problems = append(problems, fmt.Errorf("registered requirements: declared=%d actual=%d expected=%d", requirementDocument.Count, len(r.Requirements), ExpectedRequirementCount))
	}
	if len(r.Invariants)+len(r.Requirements) != ExpectedSemanticIDCount {
		problems = append(problems, fmt.Errorf("stable semantic IDs: actual=%d expected=%d", len(r.Invariants)+len(r.Requirements), ExpectedSemanticIDCount))
	}
	if invariantDocument.Source.Commit != SourceCommit || requirementDocument.ImportedFromCommit != SourceCommit || requirementDocument.PublishedSpecificationSHA256 != PublishedSpecificationSHA256 || requirementDocument.SemanticFreezeSHA256 != SemanticFreezeSHA256 {
		problems = append(problems, errors.New("canonical provenance does not match pinned AOF v1.0 LTS source"))
	}
	invIDs := map[string]bool{}
	mappedInvariantCount := 0
	for _, invariant := range r.Invariants {
		if !invariantIDRE.MatchString(invariant.ID) || invIDs[invariant.ID] {
			problems = append(problems, fmt.Errorf("invalid or duplicate invariant ID %q", invariant.ID))
		}
		invIDs[invariant.ID] = true
		if invariant.Title == "" || invariant.CanonicalRecord == "" || len(invariant.SemanticDomains) == 0 || invariant.Enforcement == "" {
			problems = append(problems, fmt.Errorf("invariant %s has incomplete canonical meaning/domain/enforcement", invariant.ID))
		}
		if len(invariant.EnforcingRequirements) == 0 && invariant.TraceabilityStatus != "NoDirectMappingAssertedByCanonicalSource" {
			problems = append(problems, fmt.Errorf("invariant %s silently lacks requirement traceability", invariant.ID))
		}
		if len(invariant.EnforcingRequirements) > 0 {
			mappedInvariantCount++
		}
	}
	reqIDs := map[string]bool{}
	traceabilityRecordCount := 0
	mappedRequirementCount := 0
	allowedNormative := map[string]bool{"MUST": true, "MUST NOT": true, "SHOULD": true, "SHOULD NOT": true, "MAY": true}
	for _, requirement := range r.Requirements {
		if !requirementIDRE.MatchString(requirement.ID) || reqIDs[requirement.ID] {
			problems = append(problems, fmt.Errorf("invalid or duplicate requirement ID %q", requirement.ID))
		}
		reqIDs[requirement.ID] = true
		r.requirements[requirement.ID] = requirement
		if requirement.Statement == "" || requirement.Domain == "" || requirement.Status != "CanonicalFrozen" {
			problems = append(problems, fmt.Errorf("requirement %s has incomplete canonical identity", requirement.ID))
		}
		if !allowedNormative[requirement.NormativeLevel] || requirement.NormativeBasis == "" {
			problems = append(problems, fmt.Errorf("requirement %s has invalid normative classification %q", requirement.ID, requirement.NormativeLevel))
		}
		if requirement.ApplicabilityRule == "" || requirement.VerificationDisposition == "" || requirement.EvidenceDisposition == "" || requirement.InvariantMappingDisposition == "" || requirement.ObjectMappingDisposition == "" {
			problems = append(problems, fmt.Errorf("requirement %s has a silent mapping/applicability gap", requirement.ID))
		}
		if (requirement.VerificationDisposition == "DeclaredInCanonicalTraceabilityRecord" || requirement.EvidenceDisposition == "DeclaredInCanonicalTraceabilityRecord") && (requirement.CanonicalTraceabilityRecord == "" || requirement.CanonicalTraceabilityLine == 0) {
			problems = append(problems, fmt.Errorf("requirement %s declares Appendix F traceability without its canonical record", requirement.ID))
		}
		if requirement.CanonicalTraceabilityRecord != "" {
			traceabilityRecordCount++
			if len(requirement.VerificationMethods) == 0 {
				problems = append(problems, fmt.Errorf("requirement %s traceability record lacks a verification method", requirement.ID))
			}
		}
		if len(requirement.RelatedInvariants) > 0 {
			mappedRequirementCount++
		}
	}
	if traceabilityRecordCount != ExpectedTraceabilityRecordCount || mappedRequirementCount != ExpectedMappedRequirementCount || mappedInvariantCount != ExpectedMappedInvariantCount {
		problems = append(problems, fmt.Errorf("canonical mapping coverage changed: Appendix-F=%d/%d requirements-mapped=%d/%d invariants-mapped=%d/%d", traceabilityRecordCount, ExpectedTraceabilityRecordCount, mappedRequirementCount, ExpectedMappedRequirementCount, mappedInvariantCount, ExpectedMappedInvariantCount))
	}
	if idSetHash(invIDs) != expectedInvariantIDSetSHA256 {
		problems = append(problems, errors.New("canonical invariant ID set changed"))
	}
	if idSetHash(reqIDs) != expectedRequirementIDSetSHA256 {
		problems = append(problems, errors.New("canonical requirement ID set changed"))
	}
	if traceabilityDocument.Counts.Requirements != ExpectedRequirementCount || traceabilityDocument.Counts.CanonicalInvariants != ExpectedInvariantCount || traceabilityDocument.Counts.ExplicitFrozenEdges != len(r.Edges) || len(r.Edges) != 352 {
		problems = append(problems, errors.New("canonical traceability counts do not match the frozen graph"))
	}
	edgeIDs := map[string]bool{}
	for _, edge := range r.Edges {
		key := edge.Type + "\x00" + edge.From + "\x00" + edge.To
		if edgeIDs[key] {
			problems = append(problems, fmt.Errorf("duplicate traceability edge %s -> %s", edge.From, edge.To))
		}
		edgeIDs[key] = true
		switch edge.Type {
		case "InvariantToRequirement":
			if !invIDs[edge.From] || !reqIDs[edge.To] {
				problems = append(problems, fmt.Errorf("broken invariant-to-requirement edge %s -> %s", edge.From, edge.To))
			}
		case "InvariantToReferenceCT":
			if !invIDs[edge.From] || !referenceTestRE.MatchString(edge.To) {
				problems = append(problems, fmt.Errorf("broken invariant-to-test edge %s -> %s", edge.From, edge.To))
			}
		case "ReferenceCTToRequirement":
			if !referenceTestRE.MatchString(edge.From) || !reqIDs[edge.To] {
				problems = append(problems, fmt.Errorf("broken test-to-requirement edge %s -> %s", edge.From, edge.To))
			}
		default:
			problems = append(problems, fmt.Errorf("unknown traceability edge type %q", edge.Type))
		}
	}
	wantedProfiles := map[string]bool{"AOF-Core": true, "AOF-Governed": true, "AOF-Assured": true, "AOF-Secure-SDLC": true, "AOF-High-Assurance": true}
	for _, profile := range r.Profiles {
		delete(wantedProfiles, profile.ID)
	}
	if len(r.Profiles) != 5 || len(wantedProfiles) != 0 {
		problems = append(problems, errors.New("canonical profile registry is incomplete"))
	}
	objects := map[string]bool{}
	for _, object := range r.CanonicalObjects {
		if object.Name == "" || object.Schema == "" || objects[object.Name] {
			problems = append(problems, fmt.Errorf("invalid or duplicate canonical object %q", object.Name))
		}
		objects[object.Name] = true
	}
	if objectDocument.Count != 22 || len(r.CanonicalObjects) != 22 {
		problems = append(problems, errors.New("canonical object registry is incomplete"))
	}
	return r, errors.Join(problems...)
}

func idSetHash(ids map[string]bool) string {
	values := make([]string, 0, len(ids))
	for id := range ids {
		values = append(values, id)
	}
	sort.Strings(values)
	sum := sha256.Sum256([]byte(strings.Join(values, "\n") + "\n"))
	return hex.EncodeToString(sum[:])
}
