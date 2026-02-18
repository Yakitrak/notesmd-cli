package obsidian

import (
	"encoding/json"
	"errors"
	"github.com/Yakitrak/notesmd-cli/pkg/config"
	"os"
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

	if cliConfig.DefaultVaultName == "" {
		return "", errors.New(ObsidianCLIConfigParseError)
	}

	v.Name = cliConfig.DefaultVaultName
	return cliConfig.DefaultVaultName, nil
}

func (v *Vault) SetDefaultName(name string) error {
	// get cliConfig path
	obsConfigDir, obsConfigFile, err := CliConfigPath()
	if err != nil {
		return err
	}

	// read existing config to preserve other fields
	cliConfig := CliConfig{}
	if content, readErr := os.ReadFile(obsConfigFile); readErr == nil {
		json.Unmarshal(content, &cliConfig) //nolint:errcheck
	}

	cliConfig.DefaultVaultName = name

	// marshal to json
	jsonContent, err := JsonMarshal(cliConfig)
	if err != nil {
		return errors.New(ObsidianCLIConfigGenerateJSONError)
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

	v.Name = name

	return nil
}

func (v *Vault) DefaultOpenType() (string, error) {
	_, cliConfigFile, err := CliConfigPath()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(cliConfigFile)
	if err != nil {
		return "obsidian", nil
	}

	cliConfig := CliConfig{}
	if err := json.Unmarshal(content, &cliConfig); err != nil {
		return "obsidian", nil
	}

	if cliConfig.DefaultOpenType == "" {
		return "obsidian", nil
	}

	return cliConfig.DefaultOpenType, nil
}

func (v *Vault) SetDefaultOpenType(openType string) error {
	obsConfigDir, obsConfigFile, err := CliConfigPath()
	if err != nil {
		return err
	}

	// read existing config to preserve other fields
	cliConfig := CliConfig{}
	if content, readErr := os.ReadFile(obsConfigFile); readErr == nil {
		json.Unmarshal(content, &cliConfig) //nolint:errcheck
	}

	cliConfig.DefaultOpenType = openType

	jsonContent, err := JsonMarshal(cliConfig)
	if err != nil {
		return errors.New(ObsidianCLIConfigGenerateJSONError)
	}

	if err := os.MkdirAll(obsConfigDir, os.ModePerm); err != nil {
		return errors.New(ObsidianCLIConfigDirWriteEror)
	}

	if err := os.WriteFile(obsConfigFile, jsonContent, 0644); err != nil {
		return errors.New(ObsidianCLIConfigWriteError)
	}

	return nil
}
