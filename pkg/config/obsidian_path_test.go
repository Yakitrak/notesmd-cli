package config_test

import (
	"errors"
	"github.com/Yakitrak/notesmd-cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConfigObsidianPath(t *testing.T) {
	t.Run("UserConfigDir func successfully returns directory", func(t *testing.T) {
		originalUserConfigDir := config.UserConfigDirectory
		defer func() { config.UserConfigDirectory = originalUserConfigDir }()

		tempDir := t.TempDir()

		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		os.Setenv("HOME", tempDir)

		// Create the config file so os.Stat succeeds on all platforms.
		configDir := filepath.Join(tempDir, "config", "obsidian")
		os.MkdirAll(configDir, 0755)
		os.WriteFile(filepath.Join(configDir, "obsidian.json"), []byte(`{}`), 0644)

		config.UserConfigDirectory = func() (string, error) {
			return filepath.Join(tempDir, "config"), nil
		}
		obsConfigFile, err := config.ObsidianFile()
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(tempDir, "config", "obsidian", "obsidian.json"), obsConfigFile)
	})

	t.Run("UserConfigDir func returns an error", func(t *testing.T) {
		originalUserConfigDir := config.UserConfigDirectory
		defer func() { config.UserConfigDirectory = originalUserConfigDir }()

		config.UserConfigDirectory = func() (string, error) {
			return "", errors.New(config.UserConfigDirectoryNotFoundErrorMessage)
		}
		obsConfigFile, err := config.ObsidianFile()
		assert.Equal(t, config.UserConfigDirectoryNotFoundErrorMessage, err.Error())
		assert.Equal(t, "", obsConfigFile)
	})

	t.Run("Finds Flatpak config on Linux when it exists", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Flatpak test only runs on Linux")
		}

		originalUserConfigDir := config.UserConfigDirectory
		defer func() { config.UserConfigDirectory = originalUserConfigDir }()

		tempDir := t.TempDir()
		flatpakDir := filepath.Join(tempDir, ".var", "app", "md.obsidian.Obsidian", "config", "obsidian")
		err := os.MkdirAll(flatpakDir, 0755)
		assert.NoError(t, err)

		flatpakConfigFile := filepath.Join(flatpakDir, "obsidian.json")
		err = os.WriteFile(flatpakConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		os.Setenv("HOME", tempDir)

		config.UserConfigDirectory = func() (string, error) {
			return filepath.Join(tempDir, ".config"), nil
		}

		obsConfigFile, err := config.ObsidianFile()
		assert.NoError(t, err)
		assert.Equal(t, flatpakConfigFile, obsConfigFile)
	})

	t.Run("Finds Snap config on Linux when Flatpak does not exist", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Snap test only runs on Linux")
		}

		originalUserConfigDir := config.UserConfigDirectory
		defer func() { config.UserConfigDirectory = originalUserConfigDir }()

		tempDir := t.TempDir()
		snapDir := filepath.Join(tempDir, "snap", "obsidian", "current", ".config", "obsidian")
		err := os.MkdirAll(snapDir, 0755)
		assert.NoError(t, err)

		snapConfigFile := filepath.Join(snapDir, "obsidian.json")
		err = os.WriteFile(snapConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		os.Setenv("HOME", tempDir)

		config.UserConfigDirectory = func() (string, error) {
			return filepath.Join(tempDir, ".config"), nil
		}

		obsConfigFile, err := config.ObsidianFile()
		assert.NoError(t, err)
		assert.Equal(t, snapConfigFile, obsConfigFile)
	})

	t.Run("Falls back to native config on Linux when others do not exist", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Native fallback test only runs on Linux")
		}

		originalUserConfigDir := config.UserConfigDirectory
		defer func() { config.UserConfigDirectory = originalUserConfigDir }()

		tempDir := t.TempDir()
		nativeDir := filepath.Join(tempDir, ".config", "obsidian")
		err := os.MkdirAll(nativeDir, 0755)
		assert.NoError(t, err)

		nativeConfigFile := filepath.Join(nativeDir, "obsidian.json")
		err = os.WriteFile(nativeConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		os.Setenv("HOME", tempDir)

		config.UserConfigDirectory = func() (string, error) {
			return filepath.Join(tempDir, ".config"), nil
		}

		obsConfigFile, err := config.ObsidianFile()
		assert.NoError(t, err)
		assert.Equal(t, nativeConfigFile, obsConfigFile)
	})

	t.Run("Returns native path when no config exists on Linux", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Native path fallback test only runs on Linux")
		}

		originalUserConfigDir := config.UserConfigDirectory
		defer func() { config.UserConfigDirectory = originalUserConfigDir }()

		tempDir := t.TempDir()

		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		os.Setenv("HOME", tempDir)

		config.UserConfigDirectory = func() (string, error) {
			return filepath.Join(tempDir, ".config"), nil
		}

		obsConfigFile, err := config.ObsidianFile()
		assert.NoError(t, err)
		expectedPath := filepath.Join(tempDir, ".config", "obsidian", "obsidian.json")
		assert.Equal(t, expectedPath, obsConfigFile)
	})

	t.Run("Prefers native over Flatpak when both exist on Linux", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Precedence test only runs on Linux")
		}

		originalUserConfigDir := config.UserConfigDirectory
		defer func() { config.UserConfigDirectory = originalUserConfigDir }()

		tempDir := t.TempDir()

		nativeDir := filepath.Join(tempDir, ".config", "obsidian")
		err := os.MkdirAll(nativeDir, 0755)
		assert.NoError(t, err)
		nativeConfigFile := filepath.Join(nativeDir, "obsidian.json")
		err = os.WriteFile(nativeConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		flatpakDir := filepath.Join(tempDir, ".var", "app", "md.obsidian.Obsidian", "config", "obsidian")
		err = os.MkdirAll(flatpakDir, 0755)
		assert.NoError(t, err)
		flatpakConfigFile := filepath.Join(flatpakDir, "obsidian.json")
		err = os.WriteFile(flatpakConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		os.Setenv("HOME", tempDir)

		config.UserConfigDirectory = func() (string, error) {
			return filepath.Join(tempDir, ".config"), nil
		}

		obsConfigFile, err := config.ObsidianFile()
		assert.NoError(t, err)
		assert.Equal(t, nativeConfigFile, obsConfigFile)
	})

	t.Run("Finds WSL Install Location", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("WSL test only runs on Linux")
		}

		// Mock variables
		originalUserConfigDir := config.UserConfigDirectory
		originalExecCommand := config.ExecCommand
		originalWslInteropFile := config.WslInteropFile
		defer func() {
			config.UserConfigDirectory = originalUserConfigDir
			config.ExecCommand = originalExecCommand
			config.WslInteropFile = originalWslInteropFile
		}()

		tempDir := t.TempDir()

		// Create WSL interop file to simulate WSL environment
		wslPresenceDir := filepath.Join(tempDir, "proc", "sys", "fs", "binfmt_misc")
		err := os.MkdirAll(wslPresenceDir, 0755)
		assert.NoError(t, err)

		wslPresenceFile := filepath.Join(wslPresenceDir, "WSLInterop")
		config.WslInteropFile = wslPresenceFile
		err = os.WriteFile(wslPresenceFile, []byte{}, 0644)
		assert.NoError(t, err)

		// Mock ExecCommand
		config.ExecCommand = func(name string, arg ...string) ([]byte, error) {
			return []byte("C:\\Users\\user\\AppData\\Roaming\r\n"), nil
		}

		// Setup environment
		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		os.Setenv("HOME", tempDir)

		config.UserConfigDirectory = func() (string, error) {
			return filepath.Join(tempDir, ".config"), nil
		}

		expectedPath := "/mnt/c/Users/user/AppData/Roaming/obsidian/obsidian.json"

		obsConfigFile, err := config.ObsidianFile()
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, obsConfigFile)
	})
}
