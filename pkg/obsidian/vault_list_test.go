package obsidian_test

import (
	"os"
	"testing"

	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

func TestListVaults(t *testing.T) {
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	originalRunningInWSL := obsidian.RunningInWSL
	defer func() {
		obsidian.ObsidianConfigFile = originalObsidianConfigFile
		obsidian.RunningInWSL = originalRunningInWSL
	}()

	// Default: not running in WSL
	obsidian.RunningInWSL = func() bool { return false }

	t.Run("Lists all vaults with derived names", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		obsidianConfig := `{
			"vaults": {
				"abc123": {
					"path": "/Users/user/Documents/Personal"
				},
				"def456": {
					"path": "/Users/user/Documents/Work"
				},
				"ghi789": {
					"path": "/Users/user/Documents/Projects/Notes"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(obsidianConfig), 0644)
		assert.NoError(t, err)

		vaults, err := obsidian.ListVaults()

		assert.NoError(t, err)
		assert.Len(t, vaults, 3)

		names := make(map[string]string)
		for _, v := range vaults {
			names[v.Name] = v.Path
		}
		assert.Equal(t, "/Users/user/Documents/Personal", names["Personal"])
		assert.Equal(t, "/Users/user/Documents/Work", names["Work"])
		assert.Equal(t, "/Users/user/Documents/Projects/Notes", names["Notes"])
	})

	t.Run("Empty vaults map returns empty slice", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		vaults, err := obsidian.ListVaults()

		assert.NoError(t, err)
		assert.Empty(t, vaults)
	})

	t.Run("Config file locator error is propagated", func(t *testing.T) {
		obsidian.ObsidianConfigFile = func() (string, error) {
			return "", os.ErrNotExist
		}

		_, err := obsidian.ListVaults()

		assert.Equal(t, os.ErrNotExist, err)
	})

	t.Run("Config file unreadable returns read error", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(``), 0000)
		assert.NoError(t, err)

		_, err = obsidian.ListVaults()

		assert.Equal(t, obsidian.ObsidianConfigReadError, err.Error())
	})

	t.Run("Invalid JSON returns parse error", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(`not valid json`), 0644)
		assert.NoError(t, err)

		_, err = obsidian.ListVaults()

		assert.Equal(t, obsidian.ObsidianConfigParseError, err.Error())
	})

	t.Run("WSL path adjustment converts Windows path and derives name", func(t *testing.T) {
		obsidian.RunningInWSL = func() bool { return true }
		defer func() { obsidian.RunningInWSL = func() bool { return false } }()

		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		configContent := `{
			"vaults": {
				"abc123": {
					"path": "C:\\Users\\user\\Documents\\MyVault"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		vaults, err := obsidian.ListVaults()

		assert.NoError(t, err)
		assert.Len(t, vaults, 1)
		assert.Equal(t, "MyVault", vaults[0].Name)
		assert.Equal(t, "/mnt/c/Users/user/Documents/MyVault", vaults[0].Path)
	})

	t.Run("Single vault returns one entry", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}

		configContent := `{
			"vaults": {
				"abc123": {
					"path": "/home/user/Notes"
				}
			}
		}`
		err := os.WriteFile(mockObsidianConfigFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		vaults, err := obsidian.ListVaults()

		assert.NoError(t, err)
		assert.Len(t, vaults, 1)
		assert.Equal(t, "Notes", vaults[0].Name)
		assert.Equal(t, "/home/user/Notes", vaults[0].Path)
	})
}

func TestResolveVaultName(t *testing.T) {
	originalObsidianConfigFile := obsidian.ObsidianConfigFile
	originalRunningInWSL := obsidian.RunningInWSL
	defer func() {
		obsidian.ObsidianConfigFile = originalObsidianConfigFile
		obsidian.RunningInWSL = originalRunningInWSL
	}()

	obsidian.RunningInWSL = func() bool { return false }

	obsidianConfig := `{
		"vaults": {
			"abc123": {
				"path": "/Users/user/Documents/Personal"
			},
			"def456": {
				"path": "/Users/user/Documents/Work"
			}
		}
	}`

	setupConfig := func(t *testing.T) {
		t.Helper()
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(obsidianConfig), 0644)
		assert.NoError(t, err)
	}

	t.Run("Resolves exact vault name", func(t *testing.T) {
		setupConfig(t)

		name, err := obsidian.ResolveVaultName("Personal")

		assert.NoError(t, err)
		assert.Equal(t, "Personal", name)
	})

	t.Run("Resolves full path to vault name", func(t *testing.T) {
		setupConfig(t)

		name, err := obsidian.ResolveVaultName("/Users/user/Documents/Personal")

		assert.NoError(t, err)
		assert.Equal(t, "Personal", name)
	})

	t.Run("Resolves path with trailing slash to vault name", func(t *testing.T) {
		setupConfig(t)

		name, err := obsidian.ResolveVaultName("/Users/user/Documents/Work/")

		assert.NoError(t, err)
		assert.Equal(t, "Work", name)
	})

	t.Run("Returns error for unregistered vault", func(t *testing.T) {
		setupConfig(t)

		_, err := obsidian.ResolveVaultName("NonExistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found in Obsidian")
		assert.Contains(t, err.Error(), "Personal")
		assert.Contains(t, err.Error(), "Work")
	})

	t.Run("Returns error for unregistered path", func(t *testing.T) {
		setupConfig(t)

		_, err := obsidian.ResolveVaultName("/var/workspace/obsidian")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found in Obsidian")
	})

	t.Run("Returns error when no vaults registered", func(t *testing.T) {
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(`{"vaults":{}}`), 0644)
		assert.NoError(t, err)

		_, err = obsidian.ResolveVaultName("anything")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no vaults registered")
	})

	t.Run("Returns error for ambiguous vault name", func(t *testing.T) {
		ambiguousConfig := `{
			"vaults": {
				"abc123": {"path": "/home/user/work/Notes"},
				"def456": {"path": "/home/user/personal/Notes"}
			}
		}`
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(ambiguousConfig), 0644)
		assert.NoError(t, err)

		_, err = obsidian.ResolveVaultName("Notes")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "multiple vaults named")
		assert.Contains(t, err.Error(), "/home/user/work/Notes")
		assert.Contains(t, err.Error(), "/home/user/personal/Notes")
	})

	t.Run("Resolves ambiguous name via full path", func(t *testing.T) {
		ambiguousConfig := `{
			"vaults": {
				"abc123": {"path": "/home/user/work/Notes"},
				"def456": {"path": "/home/user/personal/Notes"}
			}
		}`
		mockObsidianConfigFile := mocks.CreateMockObsidianConfigFile(t)
		obsidian.ObsidianConfigFile = func() (string, error) {
			return mockObsidianConfigFile, nil
		}
		err := os.WriteFile(mockObsidianConfigFile, []byte(ambiguousConfig), 0644)
		assert.NoError(t, err)

		name, err := obsidian.ResolveVaultName("/home/user/work/Notes")

		assert.NoError(t, err)
		assert.Equal(t, "Notes", name)
	})

	t.Run("Propagates config file errors", func(t *testing.T) {
		obsidian.ObsidianConfigFile = func() (string, error) {
			return "", os.ErrNotExist
		}

		_, err := obsidian.ResolveVaultName("Personal")

		assert.Error(t, err)
	})
}
