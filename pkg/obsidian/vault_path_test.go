package obsidian_test

import (
	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestVaultPath(t *testing.T) {
	// Temporarily override the ObsidianConfigFile function
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	defer func() { obsidian.ObsidianConfigFile = originalObsidianConfigFile }()

	obsidianConfig := `{
		"vaults": {
			"random1": {
				"path": "/path/to/vault1"
			},
			"random2": {
				"path": "/path/to/vault2"
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

	t.Run("Returns absolute path directly without reading obsidian config", func(t *testing.T) {
		// When the vault name is already an absolute path, Path() should return
		// it without touching ObsidianConfigFile at all.
		obsidian.ObsidianConfigFile = func() (string, error) {
			t.Fatal("ObsidianConfigFile should not be called when Name is an absolute path")
			return "", nil
		}
		vault := obsidian.Vault{Name: "/home/user/Sync/MyVault"}
		vaultPath, err := vault.Path()
		assert.Equal(t, nil, err)
		assert.Equal(t, "/home/user/Sync/MyVault", vaultPath)
		// restore for subsequent tests
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
	})

	t.Run("Gets vault path successfully from vault name without errors", func(t *testing.T) {
		// Arrange
		vault := obsidian.Vault{Name: "vault1"}
		// Act
		vaultPath, err := vault.Path()
		// Assert
		assert.Equal(t, nil, err)
		assert.Equal(t, "/path/to/vault1", vaultPath)
	})

	t.Run("Error in getting obsidian config file ", func(t *testing.T) {
		// Arrange
		obsidian.ObsidianConfigFile = func() (string, error) {
			return "", os.ErrNotExist
		}
		vault := obsidian.Vault{Name: "vault1"}
		// Act
		_, err := vault.Path()
		// Assert
		assert.Equal(t, os.ErrNotExist, err)
	})

	t.Run("Error in reading obsidian config file", func(t *testing.T) {
		// Arrange
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(``), 0000)
		if err != nil {
			t.Fatalf("Failed to create obsidian.json file: %v", err)
		}
		vault := obsidian.Vault{Name: "vault1"}
		// Act
		_, err = vault.Path()
		// Assert
		assert.Equal(t, err.Error(), obsidian.ObsidianConfigReadError)

	})

	t.Run("Error in unmarshalling obsidian config file", func(t *testing.T) {
		// Arrange
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		err := os.WriteFile(mockObsidianConfigFile, []byte(`abc`), 0644)
		if err != nil {
			t.Fatalf("Failed to create obsidian.json file: %v", err)
		}
		vault := obsidian.Vault{Name: "vault1"}
		// Act
		_, err = vault.Path()
		// Assert
		assert.Equal(t, err.Error(), obsidian.ObsidianConfigParseError)

	})

	t.Run("No vault found with given name", func(t *testing.T) {
		// Arrange
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		vault := obsidian.Vault{Name: "vault3"}
		// Act
		_, err = vault.Path()
		// Assert
		assert.Equal(t, err.Error(), obsidian.ObsidianConfigVaultNotFoundError)
	})

	t.Run("Converts windows C: path to WSL path when running in WSL", func(t *testing.T) {
		// Arrange
		originalRunningInWSL := obsidian.RunningInWSL
		obsidian.RunningInWSL = func() bool { return true }
		defer func() { obsidian.RunningInWSL = originalRunningInWSL }()

		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		configContent := `{
			"vaults": {
				"abc123": {
					"path": "C:\\Users\\user\\Documents\\myVault"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create obsidian.json file: %v", err)
		}

		vault := obsidian.Vault{Name: "myVault"}

		// Act
		vaultPath, err := vault.Path()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "/mnt/c/Users/user/Documents/myVault", vaultPath)
	})

	t.Run("Converts windows D: path to WSL path when running in WSL", func(t *testing.T) {
		// Arrange
		originalRunningInWSL := obsidian.RunningInWSL
		obsidian.RunningInWSL = func() bool { return true }
		defer func() { obsidian.RunningInWSL = originalRunningInWSL }()

		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		configContent := `{
			"vaults": {
				"def456": {
					"path": "D:\\Data\\Vaults\\MyVault"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create obsidian.json file: %v", err)
		}

		vault := obsidian.Vault{Name: "MyVault"}

		// Act
		vaultPath, err := vault.Path()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "/mnt/d/Data/Vaults/MyVault", vaultPath)
	})

	t.Run("Does not modify linux-native path when running in WSL", func(t *testing.T) {
		// Arrange
		originalRunningInWSL := obsidian.RunningInWSL
		obsidian.RunningInWSL = func() bool { return true }
		defer func() { obsidian.RunningInWSL = originalRunningInWSL }()

		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		configContent := `{
			"vaults": {
				"ghi789": {
					"path": "/home/user/Documents/Obsidian Vault"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create obsidian.json file: %v", err)
		}

		vault := obsidian.Vault{Name: "Obsidian Vault"}

		// Act
		vaultPath, err := vault.Path()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "/home/user/Documents/Obsidian Vault", vaultPath)
	})
}
