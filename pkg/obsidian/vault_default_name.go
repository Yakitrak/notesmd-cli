package obsidian

import (
	"encoding/json"
	"errors"
	"github.com/Yakitrak/obsidian-cli/pkg/config"
	"os"
	"strings"
)

var CliConfigPath = config.CliPath
var JsonMarshal = json.Marshal

func (v *Vault) DefaultName() (string, error) {
	if v.Name != "" {
		return v.Name, nil
	}

	// get cliConfig path
	_, cliConfigFile, err := CliConfigPath()
	if err != nil {
		return "", err
	}

	// read file
	content, err := os.ReadFile(cliConfigFile)
	if err != nil {
		return "", errors.New(ObsidianCLIConfigReadError)
	}

	// unmarshal json
	cliConfig := CliConfig{}
	err = json.Unmarshal(content, &cliConfig)

	if err != nil {
		return "", errors.New(ObsidianCLIConfigParseError)
	}

	// Prefer vault ID if available, otherwise use name for backward compatibility
	if cliConfig.DefaultVaultID != "" {
		v.ID = cliConfig.DefaultVaultID
		// Get the name from vault info
		info, err := v.GetVaultInfo()
		if err == nil {
			v.Name = info.Name
			return info.Name, nil
		}
		// If we can't get info, fall back to using name if available
	}

	if cliConfig.DefaultVaultName == "" {
		return "", errors.New(ObsidianCLIConfigParseError)
	}

	v.Name = cliConfig.DefaultVaultName
	return cliConfig.DefaultVaultName, nil
}

func (v *Vault) SetDefaultName(input string) error {
	// Resolve the input to a vault (could be path, name, or ID)
	vaultInfo, err := ResolveVaultFromInput(input)
	if err != nil {
		// Check if it's a warning about multiple vaults
		if strings.Contains(err.Error(), "multiple vaults found") {
			// Still proceed, but the error message contains the warning
			// The caller can check for this
		} else {
			return err
		}
	}

	// Store both name and ID for backward compatibility and future use
	cliConfig := CliConfig{
		DefaultVaultName: vaultInfo.Name,
		DefaultVaultID:   vaultInfo.ID,
	}
	
	jsonContent, err := JsonMarshal(cliConfig)
	if err != nil {
		return errors.New(ObsidianCLIConfigGenerateJSONError)
	}

	// get cliConfig path
	obsConfigDir, obsConfigFile, err := CliConfigPath()
	if err != nil {
		return err
	}

	// create directory
	err = os.MkdirAll(obsConfigDir, os.ModePerm)
	if err != nil {
		return errors.New(ObsidianCLIConfigDirWriteEror)
	}

	// create and write file
	err = os.WriteFile(obsConfigFile, jsonContent, 0644)
	if err != nil {
		return errors.New(ObsidianCLIConfigWriteError)
	}

	v.Name = vaultInfo.Name
	v.ID = vaultInfo.ID

	return nil
}
