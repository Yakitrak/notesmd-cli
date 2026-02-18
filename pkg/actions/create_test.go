package actions_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

func TestCreateNote(t *testing.T) {
	t.Run("Successful create note", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
		})
		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "note.md"))
	})

	t.Run("Successful create note with content", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
			Content:  "hello world",
		})
		// Assert
		assert.NoError(t, err)
		content, _ := os.ReadFile(filepath.Join(tmpDir, "note.md"))
		assert.Equal(t, "hello world", string(content))
	})

	t.Run("Successful create note with nested path", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "folder/note",
		})
		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "folder", "note.md"))
	})

	t.Run("Existing file is left unchanged without overwrite or append", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, "note.md")
		os.WriteFile(notePath, []byte("original"), 0644)
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
			Content:  "new content",
		})
		// Assert
		assert.NoError(t, err)
		content, _ := os.ReadFile(notePath)
		assert.Equal(t, "original", string(content))
	})

	t.Run("Successful create note with overwrite", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, "note.md")
		os.WriteFile(notePath, []byte("original"), 0644)
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:        "note",
			Content:         "overwritten",
			ShouldOverwrite: true,
		})
		// Assert
		assert.NoError(t, err)
		content, _ := os.ReadFile(notePath)
		assert.Equal(t, "overwritten", string(content))
	})

	t.Run("Successful create note with append", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, "note.md")
		os.WriteFile(notePath, []byte("original"), 0644)
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:     "note",
			Content:      " appended",
			ShouldAppend: true,
		})
		// Assert
		assert.NoError(t, err)
		content, _ := os.ReadFile(notePath)
		assert.Equal(t, "original appended", string(content))
	})

	t.Run("Successful create note with open in Obsidian", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: true,
			UseEditor:  false,
		})
		// Assert
		assert.NoError(t, err)
	})

	t.Run("Successful create note with open in editor", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: true,
			UseEditor:  true,
		})
		// Assert
		assert.NoError(t, err)
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{
			DefaultNameErr: errors.New("Failed to get vault name"),
		}
		// Act
		err := actions.CreateNote(&vault, &mocks.MockUriManager{}, actions.CreateParams{
			NoteName: "note-name",
		})
		// Assert
		assert.Equal(t, vault.DefaultNameErr, err)
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{
			Name:      "myVault",
			PathError: errors.New("Failed to get vault path"),
		}
		// Act
		err := actions.CreateNote(&vault, &mocks.MockUriManager{}, actions.CreateParams{
			NoteName: "note-name",
		})
		// Assert
		assert.Equal(t, vault.PathError, err)
	})

	t.Run("uri.Execute returns an error when opening in Obsidian", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{
			ExecuteErr: errors.New("Failed to execute URI"),
		}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note-name",
			ShouldOpen: true,
			UseEditor:  false,
		})
		// Assert
		assert.Equal(t, uri.ExecuteErr, err)
	})

	t.Run("Open in editor fails when editor command fails", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "false")

		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: true,
			UseEditor:  true,
		})
		// Assert
		assert.Error(t, err)
	})

	t.Run("Uses default folder from Obsidian config", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "folder",
			"newFileFolderPath": "Inbox"
		}`), 0644)

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
			Content:  "hello",
		})
		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "Inbox", "note.md"))
	})

	t.Run("Explicit path ignores default folder config", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "folder",
			"newFileFolderPath": "Inbox"
		}`), 0644)

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "sub/note",
			Content:  "hello",
		})
		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "sub", "note.md"))
		// Verify it's NOT in Inbox/sub/
		assert.NoFileExists(t, filepath.Join(tmpDir, "Inbox", "sub", "note.md"))
	})

	t.Run("UseEditor without open does not use editor", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		// Act — UseEditor is true but ShouldOpen is false
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: false,
			UseEditor:  true,
		})
		// Assert — file is created, editor is not invoked
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "note.md"))
	})
}

func TestNormalizeContent(t *testing.T) {
	t.Run("Replaces escape sequences with actual characters", func(t *testing.T) {
		// Arrange
		input := "Hello\\nWorld\\tTabbed\\rReturn\\\"Quote\\'SingleQuote\\\\Backslash"
		expected := "Hello\nWorld\tTabbed\rReturn\"Quote'SingleQuote\\Backslash"

		// Act
		result := actions.NormalizeContent(input)

		// Assert
		assert.Equal(t, expected, result, "The content should have the escape sequences replaced correctly")
	})

	t.Run("Handles empty input", func(t *testing.T) {
		// Arrange
		input := ""
		expected := ""

		// Act
		result := actions.NormalizeContent(input)

		// Assert
		assert.Equal(t, expected, result, "Empty input should return empty output")
	})

	t.Run("No escape sequences in input", func(t *testing.T) {
		// Arrange
		input := "Plain text with no escapes"
		expected := "Plain text with no escapes"

		// Act
		result := actions.NormalizeContent(input)

		// Assert
		assert.Equal(t, expected, result, "Content without escape sequences should remain unchanged")
	})
}
