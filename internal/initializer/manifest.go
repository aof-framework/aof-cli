package initializer

import (
	"encoding/json"

	"github.com/aof-framework/aof-cli/internal/model"
)

type bootstrapManifest struct {
	SchemaType          string   `json:"schema_type"`
	AOFSpecification    string   `json:"aof_specification"`
	AOFRelease          string   `json:"aof_release"`
	AOFCLIVersion       string   `json:"aof_cli_version"`
	UpstreamRepository  string   `json:"upstream_repository"`
	SpecificationSHA256 string   `json:"specification_sha256"`
	SchemaBundleSHA256  string   `json:"schema_bundle_sha256"`
	AdoptionMode        string   `json:"adoption_mode"`
	TargetBaseProfile   string   `json:"target_base_profile"`
	DomainProfiles      []string `json:"domain_profiles"`
	Overlays            []string `json:"overlays"`
	ClaimedConformance  bool     `json:"claimed_conformance"`
	GeneratedFiles      []string `json:"generated_files"`
}

func renderManifest(m model.BootstrapModel, cliVersion string, files []File) ([]byte, error) {
	paths := make([]string, 0, len(files)+1)
	for _, f := range files {
		paths = append(paths, f.Path)
	}
	paths = append(paths, ".aof/bootstrap-manifest.json")
	return json.MarshalIndent(bootstrapManifest{
		SchemaType: "AOFCLIBootstrapManifest", AOFSpecification: "1.0", AOFRelease: "LTS", AOFCLIVersion: cliVersion,
		UpstreamRepository: canonicalRepo, SpecificationSHA256: canonicalSpecSHA256, SchemaBundleSHA256: canonicalSchemaBundleSHA256,
		AdoptionMode: m.Adoption.Mode, TargetBaseProfile: m.Adoption.Profile.TargetBaseProfile,
		DomainProfiles: m.Adoption.Profile.DomainProfiles, Overlays: m.Adoption.Profile.Overlays,
		ClaimedConformance: false, GeneratedFiles: paths,
	}, "", "  ")
}
