package obsidian

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ResolveVaultFromInput takes an input string (path, name, or ID) and resolves it to vault info
func ResolveVaultFromInput(input string) (*VaultInfo, error) {
	obsidianConfigFile, err := ObsidianConfigFile()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(obsidianConfigFile)
	if err != nil {
		return nil, errors.New(ObsidianConfigReadError)
	}

	vaultsConfig := ObsidianVaultConfig{}
	err = json.Unmarshal(content, &vaultsConfig)
	if err != nil {
		return nil, errors.New(ObsidianConfigParseError)
	}

	// Clean up input path (resolve symlinks, get absolute path if it's a path)
	cleanInput := input
	if isPath(input) {
		if absPath, err := filepath.Abs(input); err == nil {
			cleanInput = absPath
		}
		// Try to resolve symlinks
		if realPath, err := filepath.EvalSymlinks(cleanInput); err == nil {
			cleanInput = realPath
		}
	}

	// Try to find vault by:
	// 1. Exact ID match
	// 2. Exact path match
	// 3. Vault name (directory name from path)
	
	var nameMatches []*VaultInfo

	for vaultID, vaultData := range vaultsConfig.Vaults {
		vaultPath := vaultData.Path
		vaultName := filepath.Base(vaultPath)

		// Exact ID match
		if vaultID == cleanInput {
			return &VaultInfo{
				ID:   vaultID,
				Name: vaultName,
				Path: vaultPath,
			}, nil
		}

		// Exact path match (handle symlinks)
		realVaultPath := vaultPath
		if realPath, err := filepath.EvalSymlinks(vaultPath); err == nil {
			realVaultPath = realPath
		}
		
		if cleanInput == vaultPath || cleanInput == realVaultPath {
			return &VaultInfo{
				ID:   vaultID,
				Name: vaultName,
				Path: vaultPath,
			}, nil
		}

		// Collect name matches for later
		if vaultName == cleanInput {
			nameMatches = append(nameMatches, &VaultInfo{
				ID:   vaultID,
				Name: vaultName,
				Path: vaultPath,
			})
		}
	}

	// If we found name matches, return the first one (or warn if multiple)
	if len(nameMatches) > 0 {
		if len(nameMatches) > 1 {
			// Multiple vaults with same name - this is ambiguous but we return the first
			// The caller should warn the user
			return nameMatches[0], errors.New("multiple vaults found with name '" + cleanInput + "', using: " + nameMatches[0].Path)
		}
		return nameMatches[0], nil
	}

	return nil, errors.New(ObsidianConfigVaultNotFoundError)
}

// isPath checks if the input looks like a file path
func isPath(input string) bool {
	// Check if it contains path separators or starts with common path prefixes
	return strings.Contains(input, string(filepath.Separator)) ||
		strings.HasPrefix(input, "~") ||
		strings.HasPrefix(input, ".") ||
		strings.HasPrefix(input, "/")
}

// GetVaultInfo retrieves vault info from the vault's ID
func (v *Vault) GetVaultInfo() (*VaultInfo, error) {
	obsidianConfigFile, err := ObsidianConfigFile()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(obsidianConfigFile)
	if err != nil {
		return nil, errors.New(ObsidianConfigReadError)
	}

	vaultsConfig := ObsidianVaultConfig{}
	err = json.Unmarshal(content, &vaultsConfig)
	if err != nil {
		return nil, errors.New(ObsidianConfigParseError)
	}

	// If vault has an ID, use it directly
	if v.ID != "" {
		if vaultData, exists := vaultsConfig.Vaults[v.ID]; exists {
			return &VaultInfo{
				ID:   v.ID,
				Name: filepath.Base(vaultData.Path),
				Path: vaultData.Path,
			}, nil
		}
	}

	// Fallback: try to find by name (for backward compatibility)
	if v.Name != "" {
		for vaultID, vaultData := range vaultsConfig.Vaults {
			vaultName := filepath.Base(vaultData.Path)
			if vaultName == v.Name || strings.HasSuffix(vaultData.Path, v.Name) {
				return &VaultInfo{
					ID:   vaultID,
					Name: vaultName,
					Path: vaultData.Path,
				}, nil
			}
		}
	}

	return nil, errors.New(ObsidianConfigVaultNotFoundError)
}
