package config

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	WslInteropFile = "/proc/sys/fs/binfmt_misc/WSLInterop"
)

var ExecCommand = func(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func RunningInWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	_, err := os.Stat(WslInteropFile)
	return err == nil
}

func ObsidianFile() (obsidianConfigFile string, err error) {
	userConfigDir, err := UserConfigDirectory()
	if err != nil {
		return "", errors.New(UserConfigDirectoryNotFoundErrorMessage)
	}

	defaultPath := filepath.Join(userConfigDir, ObsidianConfigDirectory, ObsidianConfigFile)
	// We don't need to do anything more for the default case
	if _, err := os.Stat(defaultPath); !os.IsNotExist(err) {
		return defaultPath, nil
	}

	if runtime.GOOS != "linux" {
		return "", errors.New(UserConfigDirectoryNotFoundErrorMessage)
	}

	if RunningInWSL() {
		return resolveWslCandidates(defaultPath)
	}

	// Otherwise, we should check in case Linux installed to a different location
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", errors.New(UserConfigDirectoryNotFoundErrorMessage)
	}

	var candidatePaths []string
	candidatePaths = append(candidatePaths, defaultPath)
	candidatePaths = append(candidatePaths,
		filepath.Join(homeDir, ".var", "app", "md.obsidian.Obsidian", "config", "obsidian", ObsidianConfigFile))
	candidatePaths = append(candidatePaths,
		filepath.Join(homeDir, "snap", "obsidian", "current", ".config", "obsidian", ObsidianConfigFile))

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

func resolveWslCandidates(defaultPath string) (string, error) {
	// Unfortunately, we cannot just "reuse" UserConfigDirectory or os.UserHomeDir as they are set following linux rules and there is no guarantee
	// the `username` of the wsl user will be the exact same as the `username` for windows.
	out, err := ExecCommand("cmd.exe", "/c", "echo %APPDATA%")
	if err != nil {
		log.Print("Failed to extract user APPDATA location. Assuming non-WSL install.")
		return "", errors.New(UserConfigDirectoryNotFoundErrorMessage)
	}

	// Trim whitespace/newlines from output
	appDataPath := strings.TrimSpace(string(out))
	if len(appDataPath) > 1 && appDataPath[1] == ':' {
		driveLetter := strings.ToLower(string(appDataPath[0]))
		restPath := appDataPath[3:] // Skip "C:\"
		// Replace backslashes with forward slashes
		restPath = strings.ReplaceAll(restPath, "\\", "/")
		wslPath := filepath.Join("/mnt", driveLetter, restPath, ObsidianConfigDirectory, ObsidianConfigFile)
		return wslPath, nil
	}

	return "", errors.New(UserConfigDirectoryNotFoundErrorMessage)
}
