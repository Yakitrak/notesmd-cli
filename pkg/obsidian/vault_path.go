package obsidian

import (
	"encoding/json"
	"errors"
	"github.com/Yakitrak/notesmd-cli/pkg/config"
	"os"
	"strings"
)

var ObsidianConfigFile = config.ObsidianFile
var RunningInWSL = config.RunningInWSL

func (v *Vault) Path() (string, error) {
	obsidianConfigFile, err := ObsidianConfigFile()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(obsidianConfigFile)
	if err != nil {
		return "", errors.New(ObsidianConfigReadError)
	}

	path, err := getPathForVault(content, v.Name)
	if err != nil {
		return "", err
	}

	if RunningInWSL() {
		return adjustForWslMount(path), nil
	}
	return path, nil
}

func adjustForWslMount(dir string) string {
	// Don't do adjustment if the path is actually linux native (although this only works for C: drive)
	if (!strings.HasPrefix(dir, "C:")) {
		return dir
	}

	mnted := strings.ReplaceAll(dir, "C:", "/mnt/c")
	return strings.ReplaceAll(mnted, "\\", "/")
}

func getPathForVault(content []byte, name string) (string, error) {
	vaultsContent := ObsidianVaultConfig{}
	if json.Unmarshal(content, &vaultsContent) != nil {
		return "", errors.New(ObsidianConfigParseError)
	}

	for _, element := range vaultsContent.Vaults {
		if strings.HasSuffix(element.Path, name) {
			return element.Path, nil
		}
	}

	return "", errors.New(ObsidianConfigVaultNotFoundError)
}