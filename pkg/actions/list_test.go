package actions_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Yakitrak/notesmd-cli/pkg/actions"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

type vaultStub struct {
	path       string
	defaultErr error
	pathErr    error
}

func (v *vaultStub) DefaultName() (string, error) {
	if v.defaultErr != nil {
		return "", v.defaultErr
	}
	return "example-vault", nil
}

func (v *vaultStub) SetDefaultName(_ string) error {
	return nil
}

func (v *vaultStub) Path() (string, error) {
	if v.pathErr != nil {
		return "", v.pathErr
	}
	return v.path, nil
}

func TestListEntries(t *testing.T) {
	t.Run("List vault root", func(t *testing.T) {
		vaultDir := t.TempDir()

		err := os.Mkdir(filepath.Join(vaultDir, "Project Alpha"), 0755)
		assert.NoError(t, err)

		err = os.WriteFile(filepath.Join(vaultDir, "Ideas.md"), []byte(""), 0644)
		assert.NoError(t, err)

		err = os.WriteFile(filepath.Join(vaultDir, "Meeting Notes.md"), []byte(""), 0644)
		assert.NoError(t, err)

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"Project Alpha/", "Ideas.md", "Meeting Notes.md"}, entries)
	})

	t.Run("List subdirectory", func(t *testing.T) {
		vaultDir := t.TempDir()
		subDir := filepath.Join(vaultDir, "001 Notes")
		assert.NoError(t, os.Mkdir(subDir, 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(subDir, "Daily.md"), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		entries, err := actions.ListEntries(vault, actions.ListParams{Path: "001 Notes"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"Daily.md"}, entries)
	})

	t.Run("Rejects path traversal", func(t *testing.T) {
		vault := &vaultStub{path: t.TempDir()}
		_, err := actions.ListEntries(vault, actions.ListParams{Path: "../"})
		assert.ErrorIs(t, err, obsidian.ErrPathTraversal)
	})

	t.Run("Returns error when path is not directory", func(t *testing.T) {
		vaultDir := t.TempDir()
		fileName := "Ideas.md"
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, fileName), []byte(""), 0644))

		vault := &vaultStub{path: vaultDir}
		_, err := actions.ListEntries(vault, actions.ListParams{Path: fileName})
		assert.EqualError(t, err, obsidian.VaultAccessError)
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		vault := &vaultStub{path: t.TempDir(), defaultErr: errors.New("default error")}
		_, err := actions.ListEntries(vault, actions.ListParams{})
		assert.Equal(t, vault.defaultErr, err)
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		vault := &vaultStub{pathErr: errors.New("path error")}
		_, err := actions.ListEntries(vault, actions.ListParams{})
		assert.Equal(t, vault.pathErr, err)
	})
}
