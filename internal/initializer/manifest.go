package initializer

import (
	"encoding/json"

	"github.com/aof-framework/aof-cli/internal/model"
)

type manifest struct {
	AOFSpecification   string   `json:"aof_specification"`
	AOFRelease         string   `json:"aof_release"`
	AOFCLIVersion      string   `json:"aof_cli_version"`
	UpstreamRepository string   `json:"upstream_repository"`
	AdoptionMode       string   `json:"adoption_mode"`
	Profile            string   `json:"profile"`
	GeneratedFiles     []string `json:"generated_files"`
}

func renderManifest(m model.BootstrapModel, cliVersion string, files []File) ([]byte, error) {
	paths := make([]string, 0, len(files)+1)
	for _, f := range files {
		paths = append(paths, f.Path)
	}
	paths = append(paths, ".aof/manifest.json")
	return json.MarshalIndent(manifest{
		AOFSpecification: "1.0", AOFRelease: "LTS", AOFCLIVersion: cliVersion,
		UpstreamRepository: "https://github.com/aof-framework/aof", AdoptionMode: m.Adoption.Mode,
		Profile: m.Adoption.Profile, GeneratedFiles: paths,
	}, "", "  ")
}
