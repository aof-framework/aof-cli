// Command canonicalgen imports the machine-readable AOF v1.0 LTS registries
// and extracts the Master Invariant Registry from the published specification.
// It is a maintainer tool; aof init never contacts the network.
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	expectedInvariants   = 162
	expectedRequirements = 332
)

type source struct {
	Specification string `json:"specification"`
	Commit        string `json:"commit"`
	SHA256        string `json:"sha256"`
	Line          int    `json:"line"`
}

type invariant struct {
	ID                    string   `json:"invariant_id"`
	Title                 string   `json:"title"`
	DomainAliases         []string `json:"domain_aliases"`
	SemanticDomains       []string `json:"semantic_domains"`
	Enforcement           string   `json:"enforcement"`
	CanonicalRecord       string   `json:"canonical_record"`
	EnforcingRequirements []string `json:"enforcing_requirements"`
	ReferenceTests        []string `json:"reference_tests"`
	TraceabilityStatus    string   `json:"traceability_status"`
	Source                source   `json:"source"`
}

type invariantRegistry struct {
	RegistryVersion string      `json:"registry_version"`
	Source          source      `json:"source"`
	Count           int         `json:"count"`
	Invariants      []invariant `json:"invariants"`
}

type edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`
	Source string `json:"source"`
}

type traceInput struct {
	DirectFrozenEdges []edge `json:"direct_frozen_edges"`
}

var (
	headingRE = regexp.MustCompile(`^### (AOF-INV-[0-9]{3}) --- (.+)$`)
	aliasRE   = regexp.MustCompile("`([A-Z]+-INV-[0-9]+)`")
)

func main() {
	var specPath, requirementsPath, tracePath, profilesPath, outDir, commit string
	flag.StringVar(&specPath, "spec", "", "published AOF specification")
	flag.StringVar(&requirementsPath, "requirements", "", "canonical requirement registry")
	flag.StringVar(&tracePath, "traceability", "", "canonical traceability matrix")
	flag.StringVar(&profilesPath, "profiles", "", "canonical profile definitions")
	flag.StringVar(&outDir, "out", "internal/canonical/data", "output directory")
	flag.StringVar(&commit, "commit", "", "upstream source commit")
	flag.Parse()
	if specPath == "" || requirementsPath == "" || tracePath == "" || profilesPath == "" || commit == "" {
		fatalf("all input flags and -commit are required")
	}

	specBytes := mustRead(specPath)
	specHash := sha256.Sum256(specBytes)
	specSHA := hex.EncodeToString(specHash[:])
	var trace traceInput
	mustJSON(mustRead(tracePath), &trace)
	invariants := extractInvariants(string(specBytes), trace.DirectFrozenEdges, commit, specSHA)
	if len(invariants) != expectedInvariants {
		fatalf("extracted %d invariants, want %d", len(invariants), expectedInvariants)
	}

	requirementBytes := mustRead(requirementsPath)
	var requirements map[string]any
	mustJSON(requirementBytes, &requirements)
	rows, ok := requirements["requirements"].([]any)
	if !ok || int(requirements["count"].(float64)) != expectedRequirements || len(rows) != expectedRequirements {
		fatalf("requirement registry does not contain %d rows", expectedRequirements)
	}
	enrichRequirements(requirements, trace.DirectFrozenEdges, string(specBytes), commit, specSHA)

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fatalf("create output directory: %v", err)
	}
	writeJSON(filepath.Join(outDir, "invariants.json"), invariantRegistry{
		RegistryVersion: "1.0-lts",
		Source:          source{Specification: filepath.Base(specPath), Commit: commit, SHA256: specSHA},
		Count:           len(invariants), Invariants: invariants,
	})
	writeJSON(filepath.Join(outDir, "requirements.json"), requirements)
	copyValidatedJSON(tracePath, filepath.Join(outDir, "traceability.json"))
	copyValidatedJSON(profilesPath, filepath.Join(outDir, "profiles.json"))
}

type requirementTraceRecord struct {
	Record  string
	Methods []string
	Line    int
}

func enrichRequirements(document map[string]any, edges []edge, spec, commit, specSHA string) {
	traceRecords := extractRequirementTraceRecords(spec)
	direct := map[string][]string{}
	invariantsByTest := map[string][]string{}
	requirementsByTest := map[string][]string{}
	for _, e := range edges {
		switch e.Type {
		case "InvariantToRequirement":
			direct[e.To] = append(direct[e.To], e.From)
		case "InvariantToReferenceCT":
			invariantsByTest[e.To] = append(invariantsByTest[e.To], e.From)
		case "ReferenceCTToRequirement":
			requirementsByTest[e.To] = append(requirementsByTest[e.To], e.From)
		}
	}
	document["imported_from_commit"] = commit
	document["published_specification_sha256"] = specSHA
	document["semantic_freeze_sha256"] = document["specification_sha256"]
	document["projection_policy"] = "Every canonical requirement is retained and evaluated by profile and project scope; absence of a canonical mapping is explicit and is never inferred."
	for _, raw := range document["requirements"].([]any) {
		requirement := raw.(map[string]any)
		id := requirement["requirement_id"].(string)
		if strings.HasPrefix(id, "AOF-ARCH-") {
			var number int
			_, _ = fmt.Sscanf(id, "AOF-ARCH-%03d", &number)
			if number >= 1 && number <= 10 {
				requirement["normative_level"] = "MUST"
				requirement["normative_basis"] = "AOF-Core conformance context in section 8.36"
			} else if number >= 11 && number <= 17 {
				requirement["normative_level"] = "SHOULD"
				requirement["normative_basis"] = "AOF-Governed/AOF-Assured strengthening context in section 8.36"
			}
		} else {
			requirement["normative_basis"] = "explicit normative keyword in canonical statement"
		}
		requirement["applicability_rule"] = "EvaluateByProfileAndProjectScope"
		traceRecord, hasTraceRecord := traceRecords[id]
		if hasTraceRecord {
			methods := make([]any, len(traceRecord.Methods))
			for i, method := range traceRecord.Methods {
				methods[i] = method
			}
			requirement["verification_method"] = methods
			requirement["canonical_traceability_record"] = traceRecord.Record
			requirement["canonical_traceability_source_line"] = traceRecord.Line
		}
		methods, _ := requirement["verification_method"].([]any)
		if len(methods) == 0 {
			requirement["verification_disposition"] = "CanonicalMappingNotAsserted"
		} else {
			requirement["verification_disposition"] = "DeclaredInCanonicalTraceabilityRecord"
		}
		evidence, _ := requirement["required_evidence"].([]any)
		if len(evidence) == 0 && !hasTraceRecord {
			requirement["evidence_disposition"] = "CanonicalMappingNotAsserted"
		} else if len(evidence) == 0 {
			requirement["evidence_disposition"] = "DeclaredInCanonicalTraceabilityRecord"
		} else {
			requirement["evidence_disposition"] = "CanonicalEvidenceDeclared"
		}
		invariants := append([]string(nil), direct[id]...)
		for _, testID := range requirementsByTest[id] {
			invariants = append(invariants, invariantsByTest[testID]...)
		}
		requirement["related_invariants"] = uniqueSorted(invariants)
		if len(invariants) == 0 {
			requirement["invariant_mapping_disposition"] = "CanonicalMappingNotAsserted"
		} else {
			requirement["invariant_mapping_disposition"] = "CanonicalFrozenTrace"
		}
		requirement["related_canonical_objects"] = []string{}
		requirement["object_mapping_disposition"] = "CanonicalMappingNotAsserted"
	}
}

func extractRequirementTraceRecords(spec string) map[string]requirementTraceRecord {
	lines := splitLines(spec)
	start, end := -1, len(lines)
	for i, line := range lines {
		if strings.HasPrefix(line, "## F.5 Requirement-Test-Evidence Matrix") {
			start = i
		}
		if start >= 0 && strings.HasPrefix(line, "## F.6 ") {
			end = i
			break
		}
	}
	if start < 0 {
		fatalf("Appendix F.5 traceability matrix not found")
	}
	rowRE := regexp.MustCompile("^(?:  |\\| )`((?:AOF-)[A-Z]+-[0-9]{3})`")
	methodRE := regexp.MustCompile(`\b(?:AT|NT|DI|CI|TI|HR)(?:/(?:AT|NT|DI|CI|TI|HR))*\b`)
	records := map[string]requirementTraceRecord{}
	for i := start; i < end; i++ {
		match := rowRE.FindStringSubmatch(lines[i])
		if match == nil || strings.HasPrefix(match[1], "AOF-INV-") {
			continue
		}
		next := end
		for j := i + 1; j < end; j++ {
			nextMatch := rowRE.FindStringSubmatch(lines[j])
			if nextMatch != nil && !strings.HasPrefix(nextMatch[1], "AOF-INV-") {
				next = j
				break
			}
		}
		record := strings.TrimSpace(strings.Join(lines[i:next], "\n"))
		methods := uniqueSorted(methodRE.FindAllString(record, -1))
		records[match[1]] = requirementTraceRecord{Record: record, Methods: methods, Line: i + 1}
		i = next - 1
	}
	return records
}

func extractInvariants(spec string, edges []edge, commit, specSHA string) []invariant {
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
		fatalf("Master Invariant Registry not found")
	}

	requirementsByInvariant := map[string][]string{}
	testsByInvariant := map[string][]string{}
	requirementsByTest := map[string][]string{}
	for _, e := range edges {
		switch e.Type {
		case "InvariantToRequirement":
			requirementsByInvariant[e.From] = append(requirementsByInvariant[e.From], e.To)
		case "InvariantToReferenceCT":
			testsByInvariant[e.From] = append(testsByInvariant[e.From], e.To)
		case "ReferenceCTToRequirement":
			requirementsByTest[e.From] = append(requirementsByTest[e.From], e.To)
		}
	}

	var out []invariant
	for i := start; i < end; i++ {
		match := headingRE.FindStringSubmatch(lines[i])
		if match == nil {
			continue
		}
		next := end
		for j := i + 1; j < end; j++ {
			if headingRE.MatchString(lines[j]) {
				next = j
				break
			}
		}
		body := strings.TrimSpace(strings.Join(lines[i+1:next], "\n"))
		aliases := uniqueMatches(aliasRE, body)
		domains := make([]string, 0, len(aliases))
		for _, alias := range aliases {
			if p := strings.Index(alias, "-INV-"); p > 0 {
				domains = append(domains, alias[:p])
			}
		}
		domains = uniqueSorted(domains)
		enforcement := extractField(body, "**Enforcement:**")
		reqs := append([]string(nil), requirementsByInvariant[match[1]]...)
		tests := uniqueSorted(testsByInvariant[match[1]])
		for _, testID := range tests {
			reqs = append(reqs, requirementsByTest[testID]...)
		}
		reqs = uniqueSorted(reqs)
		status := "MappedToCanonicalRequirement"
		if len(reqs) == 0 {
			status = "NoDirectMappingAssertedByCanonicalSource"
		} else if len(requirementsByInvariant[match[1]]) == 0 {
			status = "MappedViaCanonicalReferenceTest"
		}
		out = append(out, invariant{
			ID: match[1], Title: strings.TrimSpace(match[2]), DomainAliases: aliases,
			SemanticDomains: domains, Enforcement: enforcement, CanonicalRecord: body,
			EnforcingRequirements: reqs, ReferenceTests: tests, TraceabilityStatus: status,
			Source: source{Specification: "AOF-v1.0-Framework-Specification.md", Commit: commit, SHA256: specSHA, Line: i + 1},
		})
		i = next - 1
	}
	return out
}

func extractField(body, marker string) string {
	lines := splitLines(body)
	for i, line := range lines {
		if !strings.HasPrefix(line, marker) {
			continue
		}
		parts := []string{strings.TrimSpace(strings.TrimPrefix(line, marker))}
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "" || strings.HasPrefix(lines[j], "**") {
				break
			}
			parts = append(parts, strings.TrimSpace(lines[j]))
		}
		return strings.Join(parts, " ")
	}
	return ""
}

func uniqueMatches(re *regexp.Regexp, text string) []string {
	var values []string
	for _, match := range re.FindAllStringSubmatch(text, -1) {
		values = append(values, match[1])
	}
	return uniqueSorted(values)
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

func splitLines(text string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSuffix(scanner.Text(), "\r"))
	}
	if err := scanner.Err(); err != nil {
		fatalf("scan specification: %v", err)
	}
	return lines
}

func mustRead(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		fatalf("read %s: %v", path, err)
	}
	return b
}

func mustJSON(data []byte, dst any) {
	if err := json.Unmarshal(data, dst); err != nil {
		fatalf("decode JSON: %v", err)
	}
}

func copyValidatedJSON(src, dst string) {
	var value any
	mustJSON(mustRead(src), &value)
	writeJSON(dst, value)
}

func writeJSON(path string, value any) {
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fatalf("encode %s: %v", path, err)
	}
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0o644); err != nil {
		fatalf("write %s: %v", path, err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
