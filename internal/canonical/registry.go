// Package canonical exposes an exact, embedded snapshot of the active AOF
// v1.0 LTS release sources. Canonical source documents are never rewritten or
// enriched by AOF CLI; indexes returned by this package are derived views.
package canonical

import (
	"bufio"
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
	ExpectedInvariantCount   = 162
	ExpectedRequirementCount = 332
	ExpectedSemanticIDCount  = ExpectedInvariantCount + ExpectedRequirementCount
	ExpectedSchemaCount      = 22

	// SourceCheckoutCommit identifies the read-only local AOF checkout used to
	// vendor this snapshot. It is provenance, not normative authority.
	SourceCheckoutCommit = "58ebca64759e2ec66da7a6cab05c3881e39127ca"

	ActiveSpecificationSHA256 = "6197f71416984cca1811d5cd0cdd30327cccb9026acc5b92509f8f2a3a137974"
	OriginalBaselineSHA256    = "57ddbd64671eea615535b20f109064d96fb262e781969ef757a6f4d5efa869d5"
	OriginalBaselineCommit    = "1dd6738fdcd106750194a16666136402917cd2e8"

	releaseManifestSHA256       = "1ca59837b953da8af6e11fb19c91f7410886f7077b81e83d28fa56c872ff2580"
	requirementRegistrySHA256   = "e087d580ce59ccef3f20224d13391fd540781ba7e59b7d251786b8d8c1fe43a2"
	traceabilityMatrixSHA256    = "9ab0d4d74269857253fe4c8a720a134810099a10c2947149c833757dd0c977b6"
	profileDefinitionsSHA256    = "09260ba6a9b7e247d44c411f698a0294b415a4d761e415690d50d4a9a269d63c"
	schemaIndexSHA256           = "c166ba44c4c31d3cb3f4c776a89d551dbc34cbcdc2941c8d72c66253360b77f6"
	schemaCatalogSHA256         = "a0c5740340b72370ef63f40677e8b78db32e22ccb3219b9f0e7d827e425e9a8d"
	ActiveSchemaChecksumsSHA256 = "11ada5bf3e8ba57da2ddea6b169b41670f9aed2f0c34666bc959c1a8889672df"
)

// The files below are byte-for-byte copies from the read-only AOF checkout.
//
//go:embed source/release-manifest.json
var releaseManifestData []byte

//go:embed source/specification.md
var specificationData []byte

//go:embed source/requirement-registry.json
var requirementRegistryData []byte

//go:embed source/traceability-matrix.json
var traceabilityMatrixData []byte

//go:embed source/profile-definitions.json
var profileDefinitionsData []byte

//go:embed source/schema-index.json
var schemaIndexData []byte

//go:embed source/schema-catalog.json
var schemaCatalogData []byte

//go:embed source/schema-SHA256SUMS.txt
var schemaChecksumsData []byte

type RequirementSource struct {
	Specification string `json:"specification"`
	Line          int    `json:"line"`
	SelectionRule string `json:"selection_rule"`
}

// Requirement mirrors the upstream requirement-registry.json record exactly.
type Requirement struct {
	ID                  string            `json:"requirement_id"`
	Domain              string            `json:"domain"`
	Statement           string            `json:"statement"`
	NormativeLevel      string            `json:"normative_level"`
	Source              RequirementSource `json:"source"`
	AppliesTo           []string          `json:"applies_to"`
	Profiles            []string          `json:"profile"`
	Rationale           json.RawMessage   `json:"rationale"`
	VerificationMethods []string          `json:"verification_method"`
	RequiredEvidence    []string          `json:"required_evidence"`
	RelatedInvariants   []string          `json:"related_invariants"`
	RelatedRequirements []string          `json:"related_requirements"`
	Status              string            `json:"status"`
}

// Invariant is an index over the exact invariant record in the embedded
// normative specification. CanonicalRecord is copied verbatim from that file.
type Invariant struct {
	ID              string
	Title           string
	CanonicalRecord string
	SourceLine      int
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
	Name         string
	Schema       string
	ID           string
	IndexSHA256  string
	ActiveSHA256 string
}

// Traceability is a derived lookup over upstream direct_frozen_edges. Empty
// fields mean the upstream source asserts no mapping; AOF CLI invents none.
type Traceability struct {
	RelatedInvariants []string
	ReferenceTests    []string
}

type Registry struct {
	Invariants       []Invariant
	Requirements     []Requirement
	Edges            []Edge
	Profiles         []Profile
	CanonicalObjects []CanonicalObject
	requirements     map[string]Requirement
	traceability     map[string]Traceability
}

var (
	once            sync.Once
	loaded          Registry
	loadedErr       error
	requirementIDRE = regexp.MustCompile(`^AOF-[A-Z]+-[0-9]{3}$`)
	invariantIDRE   = regexp.MustCompile(`^AOF-INV-[0-9]{3}$`)
	referenceTestRE = regexp.MustCompile(`^CT-[A-Z]+-[0-9]{3}$`)
	invariantHeadRE = regexp.MustCompile(`^### (AOF-INV-[0-9]{3}) --- (.+)$`)
)

func Load() (Registry, error) {
	once.Do(func() { loaded, loadedErr = decodeAndValidate() })
	return loaded, loadedErr
}

func Validate() error { _, err := Load(); return err }

func MustLoad() Registry {
	r, err := Load()
	if err != nil {
		panic("invalid embedded AOF canonical sources: " + err.Error())
	}
	return r
}

func (r Registry) Requirement(id string) (Requirement, bool) {
	v, ok := r.requirements[id]
	return v, ok
}

func (r Registry) RequirementTraceability(id string) Traceability {
	return r.traceability[id]
}

func decodeAndValidate() (Registry, error) {
	var problems []error
	for _, file := range []struct {
		name, want string
		data       []byte
	}{
		{"release manifest", releaseManifestSHA256, releaseManifestData},
		{"specification", ActiveSpecificationSHA256, specificationData},
		{"requirement registry", requirementRegistrySHA256, requirementRegistryData},
		{"traceability matrix", traceabilityMatrixSHA256, traceabilityMatrixData},
		{"profile definitions", profileDefinitionsSHA256, profileDefinitionsData},
		{"schema index", schemaIndexSHA256, schemaIndexData},
		{"schema catalog", schemaCatalogSHA256, schemaCatalogData},
		{"schema checksums", ActiveSchemaChecksumsSHA256, schemaChecksumsData},
	} {
		if got := digest(file.data); got != file.want {
			problems = append(problems, fmt.Errorf("%s bytes differ from pinned upstream artifact: got %s want %s", file.name, got, file.want))
		}
	}

	var manifest struct {
		Release                string `json:"release"`
		Status                 string `json:"status"`
		EditorialRevision      string `json:"editorial_revision"`
		NormativeSpecification struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"normative_specification"`
		OriginalBaseline struct {
			SpecificationSHA256 string `json:"specification_sha256"`
			RepositoryCommit    string `json:"repository_commit"`
		} `json:"original_baseline"`
	}
	decode("release manifest", releaseManifestData, &manifest, &problems)
	if manifest.Release != "v1.0 LTS" || manifest.Status != "RELEASED" || manifest.EditorialRevision != "LTS-Editorial-2" || manifest.NormativeSpecification.Path != "specification/AOF-v1.0-Framework-Specification.md" || manifest.NormativeSpecification.SHA256 != ActiveSpecificationSHA256 || manifest.OriginalBaseline.SpecificationSHA256 != OriginalBaselineSHA256 || manifest.OriginalBaseline.RepositoryCommit != OriginalBaselineCommit {
		problems = append(problems, errors.New("active release manifest does not identify the pinned AOF v1.0 LTS editorial source"))
	}

	var reqDoc struct {
		SpecificationSHA256 string        `json:"specification_sha256"`
		Count               int           `json:"count"`
		Requirements        []Requirement `json:"requirements"`
	}
	decode("requirement registry", requirementRegistryData, &reqDoc, &problems)
	if reqDoc.SpecificationSHA256 != OriginalBaselineSHA256 {
		problems = append(problems, errors.New("requirement registry is not bound to the declared original semantic baseline"))
	}
	if reqDoc.Count != ExpectedRequirementCount || len(reqDoc.Requirements) != ExpectedRequirementCount {
		problems = append(problems, fmt.Errorf("requirements=%d declared=%d want=%d", len(reqDoc.Requirements), reqDoc.Count, ExpectedRequirementCount))
	}

	var traceDoc struct {
		Counts struct {
			Requirements        int `json:"requirements"`
			CanonicalInvariants int `json:"canonical_invariants"`
			ExplicitFrozenEdges int `json:"explicit_frozen_edges"`
		} `json:"counts"`
		DirectFrozenEdges []Edge `json:"direct_frozen_edges"`
	}
	decode("traceability matrix", traceabilityMatrixData, &traceDoc, &problems)

	var profileDoc struct {
		Profiles []Profile `json:"profiles"`
	}
	decode("profile definitions", profileDefinitionsData, &profileDoc, &problems)

	var schemaDoc struct {
		Contracts map[string]struct {
			Path   string `json:"path"`
			ID     string `json:"$id"`
			SHA256 string `json:"sha256"`
		} `json:"contracts"`
	}
	decode("schema index", schemaIndexData, &schemaDoc, &problems)

	invariants, err := extractInvariants(string(specificationData))
	if err != nil {
		problems = append(problems, err)
	}
	if len(invariants) != ExpectedInvariantCount {
		problems = append(problems, fmt.Errorf("invariants=%d want=%d", len(invariants), ExpectedInvariantCount))
	}

	r := Registry{
		Invariants: invariants, Requirements: reqDoc.Requirements, Edges: traceDoc.DirectFrozenEdges,
		Profiles: profileDoc.Profiles, requirements: map[string]Requirement{}, traceability: map[string]Traceability{},
	}
	invIDs := map[string]bool{}
	for _, inv := range invariants {
		if !invariantIDRE.MatchString(inv.ID) || invIDs[inv.ID] || inv.Title == "" || inv.CanonicalRecord == "" {
			problems = append(problems, fmt.Errorf("invalid or duplicate invariant %q", inv.ID))
		}
		invIDs[inv.ID] = true
	}
	reqIDs := map[string]bool{}
	allowedLevels := map[string]bool{"MUST": true, "MUST NOT": true, "SHOULD": true, "SHOULD NOT": true, "MAY": true, "Unclassified": true}
	for _, req := range reqDoc.Requirements {
		if !requirementIDRE.MatchString(req.ID) || reqIDs[req.ID] || req.Statement == "" || req.Domain == "" || req.Status != "CanonicalFrozen" {
			problems = append(problems, fmt.Errorf("invalid or duplicate requirement %q", req.ID))
		}
		if !allowedLevels[req.NormativeLevel] {
			problems = append(problems, fmt.Errorf("requirement %s retains unknown upstream classification %q", req.ID, req.NormativeLevel))
		}
		reqIDs[req.ID] = true
		r.requirements[req.ID] = req
	}
	if len(invariants)+len(reqDoc.Requirements) != ExpectedSemanticIDCount {
		problems = append(problems, fmt.Errorf("semantic IDs=%d want=%d", len(invariants)+len(reqDoc.Requirements), ExpectedSemanticIDCount))
	}
	if traceDoc.Counts.Requirements != ExpectedRequirementCount || traceDoc.Counts.CanonicalInvariants != ExpectedInvariantCount || traceDoc.Counts.ExplicitFrozenEdges != len(traceDoc.DirectFrozenEdges) {
		problems = append(problems, errors.New("upstream traceability counts are inconsistent"))
	}

	directInv := map[string][]string{}
	testsByInv := map[string][]string{}
	reqsByTest := map[string][]string{}
	for _, edge := range traceDoc.DirectFrozenEdges {
		switch edge.Type {
		case "InvariantToRequirement":
			if !invIDs[edge.From] || !reqIDs[edge.To] {
				problems = append(problems, fmt.Errorf("broken edge %s -> %s", edge.From, edge.To))
			}
			directInv[edge.To] = append(directInv[edge.To], edge.From)
		case "InvariantToReferenceCT":
			if !invIDs[edge.From] || !referenceTestRE.MatchString(edge.To) {
				problems = append(problems, fmt.Errorf("broken edge %s -> %s", edge.From, edge.To))
			}
			testsByInv[edge.From] = append(testsByInv[edge.From], edge.To)
		case "ReferenceCTToRequirement":
			if !referenceTestRE.MatchString(edge.From) || !reqIDs[edge.To] {
				problems = append(problems, fmt.Errorf("broken edge %s -> %s", edge.From, edge.To))
			}
			reqsByTest[edge.To] = append(reqsByTest[edge.To], edge.From)
		default:
			problems = append(problems, fmt.Errorf("unknown upstream edge type %q", edge.Type))
		}
	}
	for id := range reqIDs {
		invs := append([]string(nil), directInv[id]...)
		tests := append([]string(nil), reqsByTest[id]...)
		for inv, invTests := range testsByInv {
			for _, test := range tests {
				if contains(invTests, test) {
					invs = append(invs, inv)
				}
			}
		}
		r.traceability[id] = Traceability{RelatedInvariants: uniqueSorted(invs), ReferenceTests: uniqueSorted(tests)}
	}
	if len(profileDoc.Profiles) != 5 {
		problems = append(problems, fmt.Errorf("profiles=%d want=5", len(profileDoc.Profiles)))
	}
	activeSchemaHashes := parseChecksums(schemaChecksumsData)
	for name, contract := range schemaDoc.Contracts {
		if name == "" || contract.Path == "" || contract.ID == "" || contract.SHA256 == "" {
			problems = append(problems, fmt.Errorf("incomplete schema contract %q", name))
		}
		activeHash := activeSchemaHashes[contract.Path]
		if activeHash == "" {
			problems = append(problems, fmt.Errorf("schema %q absent from active upstream checksums", contract.Path))
		}
		r.CanonicalObjects = append(r.CanonicalObjects, CanonicalObject{Name: name, Schema: contract.Path, ID: contract.ID, IndexSHA256: contract.SHA256, ActiveSHA256: activeHash})
	}
	sort.Slice(r.CanonicalObjects, func(i, j int) bool { return r.CanonicalObjects[i].Name < r.CanonicalObjects[j].Name })
	if len(r.CanonicalObjects) != ExpectedSchemaCount {
		problems = append(problems, fmt.Errorf("schema contracts=%d want=%d", len(r.CanonicalObjects), ExpectedSchemaCount))
	}
	return r, errors.Join(problems...)
}

func extractInvariants(spec string) ([]Invariant, error) {
	lines := splitLines(spec)
	start, end := -1, len(lines)
	for i, line := range lines {
		if strings.HasPrefix(line, "# Appendix A --- Master Invariant Registry") {
			start = i
		}
		if start >= 0 && strings.HasPrefix(line, "## A.5 ") {
			end = i
			break
		}
	}
	if start < 0 {
		return nil, errors.New("master invariant registry not found in normative specification")
	}
	var out []Invariant
	for i := start; i < end; i++ {
		m := invariantHeadRE.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		next := end
		for j := i + 1; j < end; j++ {
			if invariantHeadRE.MatchString(lines[j]) {
				next = j
				break
			}
		}
		out = append(out, Invariant{ID: m[1], Title: strings.TrimSpace(m[2]), CanonicalRecord: strings.TrimSpace(strings.Join(lines[i+1:next], "\n")), SourceLine: i + 1})
		i = next - 1
	}
	return out, nil
}

func splitLines(text string) []string {
	s := bufio.NewScanner(strings.NewReader(text))
	s.Buffer(make([]byte, 64*1024), 4*1024*1024)
	var lines []string
	for s.Scan() {
		lines = append(lines, strings.TrimSuffix(s.Text(), "\r"))
	}
	return lines
}

func decode(name string, data []byte, dst any, problems *[]error) {
	if err := json.Unmarshal(data, dst); err != nil {
		*problems = append(*problems, fmt.Errorf("decode %s: %w", name, err))
	}
}

func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func parseChecksums(data []byte) map[string]string {
	result := map[string]string{}
	for _, line := range splitLines(string(data)) {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			result[strings.ReplaceAll(fields[1], "\\", "/")] = fields[0]
		}
	}
	return result
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
