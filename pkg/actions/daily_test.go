package actions_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/actions"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

func TestDailyNote(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	t.Run("Creates daily note in vault root with defaults", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		err := actions.DailyNote(&vault, &uri, actions.DailyParams{})
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, today+".md"))
	})

	t.Run("Creates daily note in configured folder", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "daily-notes.json"), []byte(`{
			"folder": "Daily"
		}`), 0644)

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		err := actions.DailyNote(&vault, &uri, actions.DailyParams{})
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "Daily", today+".md"))
	})

	t.Run("Uses template when configured", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "daily-notes.json"), []byte(`{
			"template": "Templates/Daily"
		}`), 0644)

		// Create template file
		os.MkdirAll(filepath.Join(tmpDir, "Templates"), 0755)
		os.WriteFile(filepath.Join(tmpDir, "Templates", "Daily.md"), []byte("# Daily Note\n- [ ] Task"), 0644)

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		err := actions.DailyNote(&vault, &uri, actions.DailyParams{})
		assert.NoError(t, err)

		content, _ := os.ReadFile(filepath.Join(tmpDir, today+".md"))
		assert.Equal(t, "# Daily Note\n- [ ] Task", string(content))
	})

	t.Run("Does not overwrite existing daily note", func(t *testing.T) {
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, today+".md")
		os.WriteFile(notePath, []byte("existing content"), 0644)

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		err := actions.DailyNote(&vault, &uri, actions.DailyParams{})
		assert.NoError(t, err)

		content, _ := os.ReadFile(notePath)
		assert.Equal(t, "existing content", string(content))
	})

	t.Run("Opens in Obsidian by default", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		err := actions.DailyNote(&vault, &uri, actions.DailyParams{})
		assert.NoError(t, err)
		assert.Equal(t, "myVault", uri.LastParams["vault"])
		assert.Equal(t, today, uri.LastParams["file"])
	})

	t.Run("Opens in editor when UseEditor is true", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		err := actions.DailyNote(&vault, &uri, actions.DailyParams{UseEditor: true})
		assert.NoError(t, err)
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		vaultDefaultNameErr := errors.New("Failed to get vault name")
		vaultOp := &mocks.MockVaultOperator{
			DefaultNameErr: vaultDefaultNameErr,
		}
		err := actions.DailyNote(vaultOp, &mocks.MockUriManager{}, actions.DailyParams{})
		assert.Error(t, err, vaultDefaultNameErr)
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{
			Name:      "myVault",
			PathError: errors.New("path error"),
		}
		err := actions.DailyNote(&vault, &mocks.MockUriManager{}, actions.DailyParams{})
		assert.Equal(t, vault.PathError, err)
	})

	t.Run("uri.Execute returns an error", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{
			ExecuteErr: errors.New("Failed to execute URI"),
		}
		err := actions.DailyNote(&vault, &uri, actions.DailyParams{})
		assert.Equal(t, uri.ExecuteErr, err)
	})

	t.Run("Creates daily note with custom format", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "daily-notes.json"), []byte(`{
			"format": "DD-MM-YYYY"
		}`), 0644)

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		err := actions.DailyNote(&vault, &uri, actions.DailyParams{})
		assert.NoError(t, err)

		expectedName := time.Now().Format(obsidian.MomentToGoFormat("DD-MM-YYYY"))
		assert.FileExists(t, filepath.Join(tmpDir, expectedName+".md"))
	})
}
