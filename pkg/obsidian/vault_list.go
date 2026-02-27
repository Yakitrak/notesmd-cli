package obsidian

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type VaultInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func ListVaults() ([]VaultInfo, error) {
	obsidianConfigFile, err := ObsidianConfigFile()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(obsidianConfigFile)
	if err != nil {
		return nil, errors.New(ObsidianConfigReadError)
	}

	vaultsContent := ObsidianVaultConfig{}
	if json.Unmarshal(content, &vaultsContent) != nil {
		return nil, errors.New(ObsidianConfigParseError)
	}

	vaults := make([]VaultInfo, 0, len(vaultsContent.Vaults))
	for _, element := range vaultsContent.Vaults {
		path := element.Path
		if RunningInWSL() {
			path = adjustForWslMount(path)
		}
		vaults = append(vaults, VaultInfo{
			Name: filepath.Base(path),
			Path: path,
		})
	}

	return vaults, nil
}
