package obsidian

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// ResolveVaultName validates user input against registered Obsidian vaults.
// It accepts a vault name or a path and resolves it to the correct vault name.
func ResolveVaultName(input string) (string, error) {
	vaults, err := ListVaults()
	if err != nil {
		return "", err
	}

	if len(vaults) == 0 {
		return "", errors.New("no vaults registered in Obsidian. Please create a vault in Obsidian first")
	}

	// Collect all name matches
	var nameMatches []VaultInfo
	for _, v := range vaults {
		if v.Name == input {
			nameMatches = append(nameMatches, v)
		}
	}
	if len(nameMatches) == 1 {
		return nameMatches[0].Name, nil
	}
	if len(nameMatches) > 1 {
		var paths []string
		for _, m := range nameMatches {
			paths = append(paths, fmt.Sprintf("  %s", m.Path))
		}
		return "", fmt.Errorf(
			"multiple vaults named %q found. Use the full path to disambiguate:\n%s",
			input, strings.Join(paths, "\n"),
		)
	}

	// Exact path match (user passed a full path)
	cleanInput := filepath.Clean(input)
	for _, v := range vaults {
		if filepath.Clean(v.Path) == cleanInput {
			return v.Name, nil
		}
	}

	// Build available vault list for the error message
	var available []string
	for _, v := range vaults {
		available = append(available, fmt.Sprintf("  %s\t(%s)", v.Name, v.Path))
	}

	return "", fmt.Errorf("vault %q not found in Obsidian.\nAvailable vaults:\n%s", input, strings.Join(available, "\n"))
}
