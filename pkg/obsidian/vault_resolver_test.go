package obsidian_test

import (
	"github.com/Yakitrak/obsidian-cli/mocks"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveVaultFromInput(t *testing.T) {
	// Temporarily override the ObsidianConfigFile function
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	defer func() { obsidian.ObsidianConfigFile = originalObsidianConfigFile }()

	obsidianConfig := `{
		"vaults": {
			"abc123": {
				"path": "/path/to/my-vault"
			},
			"def456": {
				"path": "/path/to/another-vault"
			},
			"ghi789": {
				"path": "/path/to/duplicate-name"
			},
			"jkl012": {
				"path": "/other/path/duplicate-name"
			}
		}
	}`
	mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
	obsidian.ObsidianConfigFile = func() (string, error) {
		return mockObsidianConfigFile, nil
	}
	err := os.WriteFile(mockObsidianConfigFile, []byte(obsidianConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create obsidian.json file: %v", err)
	}

	t.Run("Resolve vault by ID", func(t *testing.T) {
		// Act
		vaultInfo, err := obsidian.ResolveVaultFromInput("abc123")
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "abc123", vaultInfo.ID)
		assert.Equal(t, "my-vault", vaultInfo.Name)
		assert.Equal(t, "/path/to/my-vault", vaultInfo.Path)
	})

	t.Run("Resolve vault by name", func(t *testing.T) {
		// Act
		vaultInfo, err := obsidian.ResolveVaultFromInput("my-vault")
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "abc123", vaultInfo.ID)
		assert.Equal(t, "my-vault", vaultInfo.Name)
		assert.Equal(t, "/path/to/my-vault", vaultInfo.Path)
	})

	t.Run("Resolve vault by path", func(t *testing.T) {
		// Act
		vaultInfo, err := obsidian.ResolveVaultFromInput("/path/to/another-vault")
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "def456", vaultInfo.ID)
		assert.Equal(t, "another-vault", vaultInfo.Name)
		assert.Equal(t, "/path/to/another-vault", vaultInfo.Path)
	})

	t.Run("Warn when multiple vaults have the same name", func(t *testing.T) {
		// Act
		vaultInfo, err := obsidian.ResolveVaultFromInput("duplicate-name")
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "multiple vaults found")
		assert.NotNil(t, vaultInfo)
		// Should return one of the duplicates
		assert.Equal(t, "duplicate-name", vaultInfo.Name)
	})

	t.Run("Error when vault not found", func(t *testing.T) {
		// Act
		vaultInfo, err := obsidian.ResolveVaultFromInput("non-existent-vault")
		// Assert
		assert.Error(t, err)
		assert.Equal(t, obsidian.ObsidianConfigVaultNotFoundError, err.Error())
		assert.Nil(t, vaultInfo)
	})

	t.Run("Error when vault not found by path", func(t *testing.T) {
		// Act
		vaultInfo, err := obsidian.ResolveVaultFromInput("/non/existent/path")
		// Assert
		assert.Error(t, err)
		assert.Equal(t, obsidian.ObsidianConfigVaultNotFoundError, err.Error())
		assert.Nil(t, vaultInfo)
	})
}

func TestGetVaultInfo(t *testing.T) {
	// Temporarily override the ObsidianConfigFile function
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	defer func() { obsidian.ObsidianConfigFile = originalObsidianConfigFile }()

	obsidianConfig := `{
		"vaults": {
			"abc123": {
				"path": "/path/to/my-vault"
			},
			"def456": {
				"path": "/path/to/another-vault"
			}
		}
	}`
	mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
	obsidian.ObsidianConfigFile = func() (string, error) {
		return mockObsidianConfigFile, nil
	}
	err := os.WriteFile(mockObsidianConfigFile, []byte(obsidianConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create obsidian.json file: %v", err)
	}

	t.Run("Get vault info by ID", func(t *testing.T) {
		// Arrange
		vault := obsidian.Vault{ID: "abc123"}
		// Act
		vaultInfo, err := vault.GetVaultInfo()
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "abc123", vaultInfo.ID)
		assert.Equal(t, "my-vault", vaultInfo.Name)
		assert.Equal(t, "/path/to/my-vault", vaultInfo.Path)
	})

	t.Run("Get vault info by name (fallback)", func(t *testing.T) {
		// Arrange
		vault := obsidian.Vault{Name: "another-vault"}
		// Act
		vaultInfo, err := vault.GetVaultInfo()
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "def456", vaultInfo.ID)
		assert.Equal(t, "another-vault", vaultInfo.Name)
		assert.Equal(t, "/path/to/another-vault", vaultInfo.Path)
	})

	t.Run("Error when vault ID not found", func(t *testing.T) {
		// Arrange
		vault := obsidian.Vault{ID: "nonexistent"}
		// Act
		_, err := vault.GetVaultInfo()
		// Assert
		assert.Error(t, err)
		assert.Equal(t, obsidian.ObsidianConfigVaultNotFoundError, err.Error())
	})
}

func TestSetDefaultNameWithPath(t *testing.T) {
	// Temporarily override the CliConfigPath and ObsidianConfigFile functions
	originalCliConfigPath := obsidian.CliConfigPath
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	defer func() {
		obsidian.CliConfigPath = originalCliConfigPath
		obsidian.ObsidianConfigFile = originalObsidianConfigFile
	}()

	t.Run("Set default with path input", func(t *testing.T) {
		// Arrange
		mockCliConfigDir, mockCliConfigFile := mocks.CreateMockCliConfigDirectories(t)
		obsidian.CliConfigPath = func() (string, string, error) {
			return mockCliConfigDir, mockCliConfigFile, nil
		}

		// Create temp directory to simulate a real vault path
		tmpDir := t.TempDir()
		vaultPath := filepath.Join(tmpDir, "my-test-vault")
		err := os.MkdirAll(vaultPath, 0755)
		assert.NoError(t, err)

		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		obsidianConfig := `{
			"vaults": {
				"xyz789": {
					"path": "` + vaultPath + `"
				}
			}
		}`
		err = os.WriteFile(mockObsidianConfigFile, []byte(obsidianConfig), 0644)
		assert.NoError(t, err)

		vault := obsidian.Vault{}
		// Act - pass a path instead of a name
		err = vault.SetDefaultName(vaultPath)

		// Assert
		assert.NoError(t, err)
		content, err := os.ReadFile(mockCliConfigFile)
		assert.NoError(t, err)
		// Should have resolved the path to the vault ID
		assert.Contains(t, string(content), `"default_vault_id":"xyz789"`)
		assert.Contains(t, string(content), `"default_vault_name":"my-test-vault"`)
	})
}
