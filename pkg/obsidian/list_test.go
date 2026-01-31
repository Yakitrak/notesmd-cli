package obsidian

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEntries(t *testing.T) {
	t.Run("Lists directories before files with sorting", func(t *testing.T) {
		vaultDir := t.TempDir()
		assert.NoError(t, os.Mkdir(filepath.Join(vaultDir, "Project Alpha"), 0755))
		assert.NoError(t, os.Mkdir(filepath.Join(vaultDir, "Ideas"), 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "Meeting Notes.md"), []byte(""), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "Ideas.md"), []byte(""), 0644))

		entries, err := ListEntries(vaultDir, "")
		assert.NoError(t, err)
		assert.Equal(t, []string{"Ideas/", "Project Alpha/", "Ideas.md", "Meeting Notes.md"}, entries)
	})

	t.Run("Filters hidden files and folders", func(t *testing.T) {
		vaultDir := t.TempDir()
		assert.NoError(t, os.Mkdir(filepath.Join(vaultDir, ".obsidian"), 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, ".hidden.md"), []byte(""), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(vaultDir, "Visible.md"), []byte(""), 0644))

		entries, err := ListEntries(vaultDir, "")
		assert.NoError(t, err)
		assert.Equal(t, []string{"Visible.md"}, entries)
	})

	t.Run("Empty directory returns empty list", func(t *testing.T) {
		vaultDir := t.TempDir()
		entries, err := ListEntries(vaultDir, "")
		assert.NoError(t, err)
		assert.Empty(t, entries)
	})

	t.Run("Non-directory path returns vault access error", func(t *testing.T) {
		vaultDir := t.TempDir()
		filePath := filepath.Join(vaultDir, "Note.md")
		assert.NoError(t, os.WriteFile(filePath, []byte(""), 0644))

		_, err := ListEntries(vaultDir, "Note.md")
		assert.EqualError(t, err, VaultAccessError)
	})

	t.Run("Path traversal is rejected", func(t *testing.T) {
		vaultDir := t.TempDir()
		_, err := ListEntries(vaultDir, "../")
		assert.ErrorIs(t, err, ErrPathTraversal)
	})
}
