package obsidian

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ObsidianAppConfig represents relevant fields from .obsidian/app.json.
type ObsidianAppConfig struct {
	NewFileLocation   string `json:"newFileLocation"`
	NewFileFolderPath string `json:"newFileFolderPath"`
}

// DailyNotesConfig represents relevant fields from .obsidian/daily-notes.json.
type DailyNotesConfig struct {
	Folder   string `json:"folder"`
	Format   string `json:"format"`
	Template string `json:"template"`
}

// DefaultNoteFolder reads the configured default folder for new notes from
// .obsidian/app.json. Returns "" if not configured or unreadable (caller
// should use vault root).
func DefaultNoteFolder(vaultPath string) string {
	data, err := os.ReadFile(filepath.Join(vaultPath, ".obsidian", "app.json"))
	if err != nil {
		return ""
	}

	var config ObsidianAppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return ""
	}

	if config.NewFileLocation == "folder" && config.NewFileFolderPath != "" {
		return config.NewFileFolderPath
	}

	return ""
}

// ReadDailyNotesConfig reads the daily notes plugin config from the vault.
// Returns zero-value config if unreadable.
func ReadDailyNotesConfig(vaultPath string) DailyNotesConfig {
	data, err := os.ReadFile(filepath.Join(vaultPath, ".obsidian", "daily-notes.json"))
	if err != nil {
		return DailyNotesConfig{}
	}

	var config DailyNotesConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DailyNotesConfig{}
	}

	return config
}

// ApplyDefaultFolder prepends the configured default note folder to noteName
// when noteName has no explicit path (no "/"). If the note name already
// contains a "/", it is treated as an explicit path and returned unchanged.
// Falls back to the original name if no default folder is configured.
func ApplyDefaultFolder(noteName, vaultPath string) string {
	if strings.Contains(noteName, "/") {
		return noteName
	}
	if folder := DefaultNoteFolder(vaultPath); folder != "" {
		return folder + "/" + noteName
	}
	return noteName
}

// MomentToGoFormat converts a Moment.js date format string to a Go time layout.
// It uses a two-pass approach with placeholders to avoid cascading replacements
// (e.g., replacing "a" inside "January").
func MomentToGoFormat(momentFmt string) string {
	// Order matters: longer tokens must be replaced before shorter ones.
	replacements := []struct {
		moment string
		goFmt  string
	}{
		{"YYYY", "2006"},
		{"YY", "06"},
		{"MMMM", "January"},
		{"MMM", "Jan"},
		{"MM", "01"},
		{"M", "1"},
		{"DD", "02"},
		{"D", "2"},
		{"dddd", "Monday"},
		{"ddd", "Mon"},
		{"dd", "Mo"},
		{"HH", "15"},
		{"hh", "03"},
		{"h", "3"},
		{"mm", "04"},
		{"ss", "05"},
		{"A", "PM"},
		{"a", "pm"},
	}

	// Pass 1: replace Moment tokens with unique placeholders.
	result := momentFmt
	for i, r := range replacements {
		placeholder := fmt.Sprintf("\x00%d\x00", i)
		result = strings.Replace(result, r.moment, placeholder, -1)
	}

	// Pass 2: replace placeholders with Go format strings.
	for i, r := range replacements {
		placeholder := fmt.Sprintf("\x00%d\x00", i)
		result = strings.Replace(result, placeholder, r.goFmt, -1)
	}

	return result
}
