package actions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

func TestOpenNote(t *testing.T) {
	t.Run("Successful open note", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.OpenNote(&vault, &uri, actions.OpenParams{
			NoteName: "note.md",
		})
		// Assert
		assert.Equal(t, err, nil)
		assert.Equal(t, "myVault", uri.LastParams["vault"])
		assert.Equal(t, "note.md", uri.LastParams["file"])
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		// Arrange
		vaultDefaultNameErr := errors.New("Failed to get vault name")
		vaultOp := &mocks.MockVaultOperator{
			DefaultNameErr: vaultDefaultNameErr,
		}
		// Act
		err := actions.OpenNote(vaultOp, &mocks.MockUriManager{}, actions.OpenParams{
			NoteName: "note.md",
		})
		// Assert
		assert.Error(t, err, vaultDefaultNameErr)
	})

	t.Run("uri.Execute returns an error", func(t *testing.T) {
		// Arrange
		uri := mocks.MockUriManager{
			ExecuteErr: errors.New("Failed to execute URI"),
		}
		// Act
		err := actions.OpenNote(&mocks.MockVaultOperator{}, &uri, actions.OpenParams{
			NoteName: "note1.md",
			Section:  "Heading One",
		})
		// Assert
		assert.Equal(t, err, uri.ExecuteErr)
		assert.Equal(t, "note1.md#Heading One", uri.LastParams["file"])
	})

	t.Run("Opens note with section", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}

		err := actions.OpenNote(&vault, &uri, actions.OpenParams{
			NoteName: "note.md",
			Section:  "Section Name",
		})

		assert.NoError(t, err)
		assert.Equal(t, "note.md#Section Name", uri.LastParams["file"])
	})

	t.Run("Opens note in editor when UseEditor is true", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		// Create the note file so the editor can open it
		os.WriteFile(tmpDir+"/note.md", []byte("hello"), 0644)

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		err := actions.OpenNote(&vault, &uri, actions.OpenParams{
			NoteName:  "note",
			UseEditor: true,
		})

		assert.NoError(t, err)
	})

	t.Run("Editor open fails when vault.Path returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{
			Name:      "myVault",
			PathError: errors.New("path error"),
		}
		uri := mocks.MockUriManager{}

		err := actions.OpenNote(&vault, &uri, actions.OpenParams{
			NoteName:  "note",
			UseEditor: true,
		})

		assert.Equal(t, vault.PathError, err)
	})
}
