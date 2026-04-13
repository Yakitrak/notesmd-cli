package obsidian_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

func TestAddVault(t *testing.T) {
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	originalRunningInWSL := obsidian.RunningInWSL
	defer func() {
		obsidian.ObsidianConfigFile = originalObsidianConfigFile
		obsidian.RunningInWSL = originalRunningInWSL
	}()

	obsidian.RunningInWSL = func() bool { return false }

	t.Run("Adds vault to existing config", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		vaultDir := t.TempDir()
		err = obsidian.AddVault(vaultDir)
		assert.NoError(t, err)

		// Verify vault was added
		content, err := os.ReadFile(mockObsidianConfigFile)
		assert.NoError(t, err)

		var cfg obsidian.ObsidianVaultConfig
		err = json.Unmarshal(content, &cfg)
		assert.NoError(t, err)
		assert.Len(t, cfg.Vaults, 1)

		for _, v := range cfg.Vaults {
			assert.Equal(t, vaultDir, v.Path)
		}
	})

	t.Run("Rejects non-existent path", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		err = obsidian.AddVault("/nonexistent/path/to/vault")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path does not exist")
	})

	t.Run("Rejects file path", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		tmpFile := filepath.Join(t.TempDir(), "file.txt")
		err = os.WriteFile(tmpFile, []byte("test"), 0644)
		assert.NoError(t, err)

		err = obsidian.AddVault(tmpFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})

	t.Run("Rejects duplicate vault", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		vaultDir := t.TempDir()
		configContent := `{
			"vaults": {
				"abc123": {
					"path": "` + vaultDir + `"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		err = obsidian.AddVault(vaultDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already registered")
	})

	t.Run("Adds multiple vaults", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		vault1 := t.TempDir()
		vault2 := t.TempDir()

		err = obsidian.AddVault(vault1)
		assert.NoError(t, err)

		err = obsidian.AddVault(vault2)
		assert.NoError(t, err)

		content, err := os.ReadFile(mockObsidianConfigFile)
		assert.NoError(t, err)

		var cfg obsidian.ObsidianVaultConfig
		err = json.Unmarshal(content, &cfg)
		assert.NoError(t, err)
		assert.Len(t, cfg.Vaults, 2)
	})
}

func TestRemoveVault(t *testing.T) {
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	originalRunningInWSL := obsidian.RunningInWSL
	defer func() {
		obsidian.ObsidianConfigFile = originalObsidianConfigFile
		obsidian.RunningInWSL = originalRunningInWSL
	}()

	obsidian.RunningInWSL = func() bool { return false }

	t.Run("Removes vault by name", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		configContent := `{
			"vaults": {
				"abc123": {
					"path": "/Users/user/Documents/Personal"
				},
				"def456": {
					"path": "/Users/user/Documents/Work"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		err = obsidian.RemoveVault("Personal")
		assert.NoError(t, err)

		content, err := os.ReadFile(mockObsidianConfigFile)
		assert.NoError(t, err)

		var cfg obsidian.ObsidianVaultConfig
		err = json.Unmarshal(content, &cfg)
		assert.NoError(t, err)
		assert.Len(t, cfg.Vaults, 1)

		for _, v := range cfg.Vaults {
			assert.Equal(t, "/Users/user/Documents/Work", v.Path)
		}
	})

	t.Run("Removes vault by path", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		configContent := `{
			"vaults": {
				"abc123": {
					"path": "/Users/user/Documents/Personal"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		err = obsidian.RemoveVault("/Users/user/Documents/Personal")
		assert.NoError(t, err)

		content, err := os.ReadFile(mockObsidianConfigFile)
		assert.NoError(t, err)

		var cfg obsidian.ObsidianVaultConfig
		err = json.Unmarshal(content, &cfg)
		assert.NoError(t, err)
		assert.Len(t, cfg.Vaults, 0)
	})

	t.Run("Returns error for non-existent vault", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		err = obsidian.RemoveVault("NonExistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
