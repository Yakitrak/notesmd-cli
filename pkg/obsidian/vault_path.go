package obsidian

import (
	"encoding/json"
	"errors"
	"github.com/Yakitrak/notesmd-cli/pkg/config"
	"os"
	"path/filepath"
	"strings"
)

var ObsidianConfigFile = config.ObsidianFile
var RunningInWSL = config.RunningInWSL

func (v *Vault) Path() (string, error) {
	// If the stored name is already an absolute path, return it directly
	// without requiring Obsidian's config file to be present.
	if filepath.IsAbs(v.Name) {
		return v.Name, nil
	}

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
	// Detect any Windows drive letter pattern (e.g. C:, D:, E:)
	if len(dir) >= 2 && dir[1] == ':' && ((dir[0] >= 'A' && dir[0] <= 'Z') || (dir[0] >= 'a' && dir[0] <= 'z')) {
		driveLetter := strings.ToLower(string(dir[0]))
		mnted := "/mnt/" + driveLetter + dir[2:]
		return strings.ReplaceAll(mnted, "\\", "/")
	}

	return dir
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