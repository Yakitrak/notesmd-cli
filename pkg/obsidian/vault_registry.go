package obsidian

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Yakitrak/notesmd-cli/pkg/config"
)

// AddVault registers a vault path in the Obsidian config file.
// It creates the config file and directory if they don't exist.
func AddVault(vaultPath string) error {
	absPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("path does not exist: %s", absPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	obsidianConfigFile, vaultsConfig, err := readOrCreateObsidianConfig()
	if err != nil {
		return err
	}

	// Check if vault already registered
	for _, v := range vaultsConfig.Vaults {
		if filepath.Clean(v.Path) == filepath.Clean(absPath) {
			return fmt.Errorf("vault already registered: %s", absPath)
		}
	}

	id, err := generateVaultID()
	if err != nil {
		return fmt.Errorf("failed to generate vault ID: %w", err)
	}

	vaultsConfig.Vaults[id] = struct {
		Path string `json:"path"`
	}{Path: absPath}

	return writeObsidianConfig(obsidianConfigFile, vaultsConfig)
}

// RemoveVault removes a vault from the Obsidian config file by name.
func RemoveVault(name string) error {
	obsidianConfigFile, err := ObsidianConfigFile()
	if err != nil {
		return errors.New(ObsidianConfigReadError)
	}

	content, err := os.ReadFile(obsidianConfigFile)
	if err != nil {
		return errors.New(ObsidianConfigReadError)
	}

	vaultsConfig := ObsidianVaultConfig{}
	if json.Unmarshal(content, &vaultsConfig) != nil {
		return errors.New(ObsidianConfigParseError)
	}

	found := false
	for id, v := range vaultsConfig.Vaults {
		vaultName := filepath.Base(v.Path)
		if vaultName == name || filepath.Clean(v.Path) == filepath.Clean(name) {
			delete(vaultsConfig.Vaults, id)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("vault %q not found", name)
	}

	return writeObsidianConfig(obsidianConfigFile, vaultsConfig)
}

// ClearDefaultIfMatch clears the default vault in CLI config if it matches the given name.
func ClearDefaultIfMatch(name string) error {
	_, cliConfigFile, err := CliConfigPath()
	if err != nil {
		return nil //nolint:nilerr // no config dir means nothing to clear
	}

	content, err := os.ReadFile(cliConfigFile)
	if err != nil {
		return nil //nolint:nilerr // no config file means nothing to clear
	}

	cliConfig := CliConfig{}
	if err := json.Unmarshal(content, &cliConfig); err != nil {
		return nil //nolint:nilerr // unparseable config means nothing to clear
	}

	if cliConfig.DefaultVaultName != name {
		return nil
	}

	cliConfig.DefaultVaultName = ""
	v := &Vault{}
	return v.SetDefaultName("")
}

func readOrCreateObsidianConfig() (string, ObsidianVaultConfig, error) {
	empty := ObsidianVaultConfig{
		Vaults: make(map[string]struct {
			Path string `json:"path"`
		}),
	}

	// Try to find existing config
	obsidianConfigFile, err := ObsidianConfigFile()
	if err == nil {
		content, readErr := os.ReadFile(obsidianConfigFile)
		if readErr == nil {
			vaultsConfig := ObsidianVaultConfig{}
			if json.Unmarshal(content, &vaultsConfig) == nil {
				if vaultsConfig.Vaults == nil {
					vaultsConfig.Vaults = make(map[string]struct {
						Path string `json:"path"`
					})
				}
				return obsidianConfigFile, vaultsConfig, nil
			}
		}
	}

	// Config doesn't exist, create it
	userConfigDir, err := config.UserConfigDirectory()
	if err != nil {
		return "", empty, fmt.Errorf("failed to determine config directory: %w", err)
	}

	configDir := filepath.Join(userConfigDir, config.ObsidianConfigDirectory)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", empty, fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, config.ObsidianConfigFile)
	return configFile, empty, nil
}

func writeObsidianConfig(path string, cfg ObsidianVaultConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func generateVaultID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
