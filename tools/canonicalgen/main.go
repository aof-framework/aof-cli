// Command canonicalgen vendors exact AOF v1.0 LTS source artifacts from a
// local, read-only upstream checkout. It never edits upstream files and never
// adds fields to canonical JSON. aof init does not invoke this tool.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	expectedInvariants   = 162
	expectedRequirements = 332
	expectedSchemas      = 22
)

var sources = []struct{ upstream, destination string }{
	{"release/AOF-v1.0-LTS-Release-Manifest.json", "release-manifest.json"},
	{"specification/AOF-v1.0-Framework-Specification.md", "specification.md"},
	{"conformance/registry/requirement-registry.json", "requirement-registry.json"},
	{"conformance/mapping/lts-a3-traceability-matrix.json", "traceability-matrix.json"},
	{"conformance/profiles/profile-definitions.json", "profile-definitions.json"},
	{"schemas/schema-index.json", "schema-index.json"},
	{"schemas/catalog.json", "schema-catalog.json"},
	{"schemas/SHA256SUMS.txt", "schema-SHA256SUMS.txt"},
}

func main() {
	var upstreamRoot, outDir, schemasOut string
	flag.StringVar(&upstreamRoot, "upstream-root", "", "local read-only AOF v1.0 LTS checkout")
	flag.StringVar(&outDir, "out", "internal/canonical/source", "destination for exact source artifacts")
	flag.StringVar(&schemasOut, "schemas-out", "templates/static/schemas", "destination for exact runtime schema payload")
	flag.Parse()
	if upstreamRoot == "" {
		fatalf("-upstream-root is required")
	}

	artifacts := map[string][]byte{}
	for _, source := range sources {
		path := filepath.Join(upstreamRoot, filepath.FromSlash(source.upstream))
		data, err := os.ReadFile(path)
		if err != nil {
			fatalf("read upstream %s: %v", source.upstream, err)
		}
		artifacts[source.destination] = data
	}
	validate(artifacts)
	schemaArtifacts := map[string][]byte{}
	checksums := parseChecksums(artifacts["schema-SHA256SUMS.txt"])
	for _, relative := range schemaPaths(artifacts) {
		data, err := os.ReadFile(filepath.Join(upstreamRoot, "schemas", filepath.FromSlash(relative)))
		if err != nil {
			fatalf("read upstream schema %s: %v", relative, err)
		}
		if want := checksums[relative]; want != "" && digest(data) != want {
			fatalf("upstream schema checksum mismatch for %s", relative)
		}
		schemaArtifacts[relative] = data
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fatalf("create output directory: %v", err)
	}
	for _, source := range sources {
		if err := os.WriteFile(filepath.Join(outDir, source.destination), artifacts[source.destination], 0o644); err != nil {
			fatalf("write %s: %v", source.destination, err)
		}
	}
	for _, relative := range schemaPaths(artifacts) {
		destination := filepath.Join(schemasOut, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			fatalf("create schema destination: %v", err)
		}
		if err := os.WriteFile(destination, schemaArtifacts[relative], 0o644); err != nil {
			fatalf("write schema %s: %v", relative, err)
		}
	}
	fmt.Printf("vendored %d exact AOF v1.0 LTS artifacts from local read-only source\n", len(sources))
}

func parseChecksums(data []byte) map[string]string {
	result := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			result[strings.ReplaceAll(fields[1], "\\", "/")] = fields[0]
		}
	}
	return result
}

func schemaPaths(files map[string][]byte) []string {
	var index struct {
		Contracts map[string]struct {
			Path string `json:"path"`
		} `json:"contracts"`
	}
	var catalog struct {
		Shared []string `json:"shared_schema_files"`
	}
	mustJSON("schema index", files["schema-index.json"], &index)
	mustJSON("schema catalog", files["schema-catalog.json"], &catalog)
	paths := map[string]bool{
		"catalog.json": true, "schema-index.json": true, "SHA256SUMS.txt": true,
		"provenance/frozen-specification.json": true, "provenance/phase-contract-hashes.json": true,
	}
	for _, contract := range index.Contracts {
		paths[contract.Path] = true
	}
	for _, path := range catalog.Shared {
		paths[path] = true
	}
	result := make([]string, 0, len(paths))
	for path := range paths {
		clean := filepath.Clean(filepath.FromSlash(path))
		if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			fatalf("unsafe upstream schema path %q", path)
		}
		result = append(result, path)
	}
	sort.Strings(result)
	return result
}

func validate(files map[string][]byte) {
	var manifest struct {
		Release                string `json:"release"`
		Status                 string `json:"status"`
		NormativeSpecification struct {
			SHA256 string `json:"sha256"`
		} `json:"normative_specification"`
	}
	mustJSON("release manifest", files["release-manifest.json"], &manifest)
	if manifest.Release != "v1.0 LTS" || manifest.Status != "RELEASED" {
		fatalf("upstream is not the active released AOF v1.0 LTS")
	}
	if got := digest(files["specification.md"]); got != manifest.NormativeSpecification.SHA256 {
		fatalf("active specification checksum mismatch: got %s want %s", got, manifest.NormativeSpecification.SHA256)
	}

	var requirements struct {
		Count        int               `json:"count"`
		Requirements []json.RawMessage `json:"requirements"`
	}
	mustJSON("requirement registry", files["requirement-registry.json"], &requirements)
	if requirements.Count != expectedRequirements || len(requirements.Requirements) != expectedRequirements {
		fatalf("requirements=%d declared=%d want=%d", len(requirements.Requirements), requirements.Count, expectedRequirements)
	}

	var trace struct {
		Counts struct {
			CanonicalInvariants int `json:"canonical_invariants"`
		} `json:"counts"`
	}
	mustJSON("traceability matrix", files["traceability-matrix.json"], &trace)
	if trace.Counts.CanonicalInvariants != expectedInvariants {
		fatalf("invariants=%d want=%d", trace.Counts.CanonicalInvariants, expectedInvariants)
	}

	var profiles struct {
		Profiles []json.RawMessage `json:"profiles"`
	}
	mustJSON("profile definitions", files["profile-definitions.json"], &profiles)
	if len(profiles.Profiles) != 5 {
		fatalf("profiles=%d want=5", len(profiles.Profiles))
	}

	var schemas struct {
		Contracts map[string]json.RawMessage `json:"contracts"`
	}
	mustJSON("schema index", files["schema-index.json"], &schemas)
	if len(schemas.Contracts) != expectedSchemas {
		fatalf("schemas=%d want=%d", len(schemas.Contracts), expectedSchemas)
	}
}

func mustJSON(name string, data []byte, dst any) {
	if err := json.Unmarshal(data, dst); err != nil {
		fatalf("decode %s: %v", name, err)
	}
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
