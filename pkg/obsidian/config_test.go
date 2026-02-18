package obsidian_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

func TestDefaultNoteFolder(t *testing.T) {
	t.Run("Returns folder when newFileLocation is folder", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "folder",
			"newFileFolderPath": "00 Inbox"
		}`), 0644)

		result := obsidian.DefaultNoteFolder(tmpDir)
		assert.Equal(t, "00 Inbox", result)
	})

	t.Run("Returns empty when newFileLocation is root", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "root"
		}`), 0644)

		result := obsidian.DefaultNoteFolder(tmpDir)
		assert.Equal(t, "", result)
	})

	t.Run("Returns empty when newFileLocation is current", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "current"
		}`), 0644)

		result := obsidian.DefaultNoteFolder(tmpDir)
		assert.Equal(t, "", result)
	})

	t.Run("Returns empty when config is absent", func(t *testing.T) {
		tmpDir := t.TempDir()
		result := obsidian.DefaultNoteFolder(tmpDir)
		assert.Equal(t, "", result)
	})

	t.Run("Returns empty when config is invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`not json`), 0644)

		result := obsidian.DefaultNoteFolder(tmpDir)
		assert.Equal(t, "", result)
	})

	t.Run("Returns empty when folder location set but path is empty", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "folder",
			"newFileFolderPath": ""
		}`), 0644)

		result := obsidian.DefaultNoteFolder(tmpDir)
		assert.Equal(t, "", result)
	})
}

func TestApplyDefaultFolder(t *testing.T) {
	t.Run("Prepends folder when configured and no slash in name", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "folder",
			"newFileFolderPath": "Inbox"
		}`), 0644)

		result := obsidian.ApplyDefaultFolder("my-note", tmpDir)
		assert.Equal(t, "Inbox/my-note", result)
	})

	t.Run("Does not prepend when name contains slash", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
			"newFileLocation": "folder",
			"newFileFolderPath": "Inbox"
		}`), 0644)

		result := obsidian.ApplyDefaultFolder("sub/my-note", tmpDir)
		assert.Equal(t, "sub/my-note", result)
	})

	t.Run("Returns name unchanged when no config", func(t *testing.T) {
		tmpDir := t.TempDir()
		result := obsidian.ApplyDefaultFolder("my-note", tmpDir)
		assert.Equal(t, "my-note", result)
	})
}

func TestReadDailyNotesConfig(t *testing.T) {
	t.Run("Reads full config", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "daily-notes.json"), []byte(`{
			"folder": "02 - Daily Notes",
			"format": "YYYY-MM-DD",
			"template": "Templates/Daily Note"
		}`), 0644)

		config := obsidian.ReadDailyNotesConfig(tmpDir)
		assert.Equal(t, "02 - Daily Notes", config.Folder)
		assert.Equal(t, "YYYY-MM-DD", config.Format)
		assert.Equal(t, "Templates/Daily Note", config.Template)
	})

	t.Run("Returns zero config when file is absent", func(t *testing.T) {
		tmpDir := t.TempDir()
		config := obsidian.ReadDailyNotesConfig(tmpDir)
		assert.Equal(t, obsidian.DailyNotesConfig{}, config)
	})

	t.Run("Returns zero config for invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "daily-notes.json"), []byte(`bad`), 0644)

		config := obsidian.ReadDailyNotesConfig(tmpDir)
		assert.Equal(t, obsidian.DailyNotesConfig{}, config)
	})

	t.Run("Handles partial config", func(t *testing.T) {
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		os.MkdirAll(obsDir, 0755)
		os.WriteFile(filepath.Join(obsDir, "daily-notes.json"), []byte(`{
			"folder": "Daily"
		}`), 0644)

		config := obsidian.ReadDailyNotesConfig(tmpDir)
		assert.Equal(t, "Daily", config.Folder)
		assert.Equal(t, "", config.Format)
		assert.Equal(t, "", config.Template)
	})
}

func TestMomentToGoFormat(t *testing.T) {
	tests := []struct {
		name     string
		moment   string
		expected string
	}{
		{"Default format", "YYYY-MM-DD", "2006-01-02"},
		{"Year month day with slashes", "YYYY/MM/DD", "2006/01/02"},
		{"Short year", "YY-MM-DD", "06-01-02"},
		{"Full month name", "MMMM DD, YYYY", "January 02, 2006"},
		{"Short month name", "MMM DD, YYYY", "Jan 02, 2006"},
		{"Day of week", "dddd, MMMM DD, YYYY", "Monday, January 02, 2006"},
		{"Short day of week", "ddd MMM DD", "Mon Jan 02"},
		{"With time", "YYYY-MM-DD HH:mm", "2006-01-02 15:04"},
		{"Plain text passthrough", "notes", "notes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := obsidian.MomentToGoFormat(tt.moment)
			assert.Equal(t, tt.expected, result)
		})
	}
}
