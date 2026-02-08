package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func ObsidianFile() (obsidianConfigFile string, err error) {
	userConfigDir, err := UserConfigDirectory()
	if err != nil {
		return "", errors.New(UserConfigDirectoryNotFoundErrorMessage)
	}

	defaultPath := filepath.Join(userConfigDir, ObsidianConfigDirectory, ObsidianConfigFile)
	var candidatePaths []string

	switch runtime.GOOS {
	case "linux":
		candidatePaths = append(candidatePaths, defaultPath)

		homeDir, homeErr := os.UserHomeDir()
		if homeErr == nil {
			candidatePaths = append(candidatePaths,
				filepath.Join(homeDir, ".var", "app", "md.obsidian.Obsidian", "config", "obsidian", ObsidianConfigFile))
			candidatePaths = append(candidatePaths,
				filepath.Join(homeDir, "snap", "obsidian", "current", ".config", "obsidian", ObsidianConfigFile))
		}
	default:
		return defaultPath, nil
	}

	var firstNonExistErr error
	for _, path := range candidatePaths {
		if _, statErr := os.Stat(path); statErr == nil {
			return path, nil
		} else if !os.IsNotExist(statErr) && firstNonExistErr == nil {
			firstNonExistErr = statErr
		}
	}

	if firstNonExistErr != nil {
		return "", firstNonExistErr
	}

	return defaultPath, nil
}
