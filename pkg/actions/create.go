package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
)

type CreateParams struct {
	NoteName        string
	ShouldAppend    bool
	ShouldOverwrite bool
	Content         string
	ShouldOpen      bool
	UseEditor       bool
}

func CreateNote(vault obsidian.VaultManager, uri obsidian.UriManager, params CreateParams) error {
	// DefaultName populates vault name from config if not already set (required before Path()).
	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	// Prepend configured default folder when note name has no explicit path.
	params.NoteName = obsidian.ApplyDefaultFolder(params.NoteName, vaultPath)

	// Validate the note path stays within the vault directory.
	notePath, err := obsidian.ValidatePath(vaultPath, obsidian.AddMdSuffix(params.NoteName))
	if err != nil {
		return err
	}

	// Create any intermediate directories the note path requires.
	if err := os.MkdirAll(filepath.Dir(notePath), 0755); err != nil {
		return fmt.Errorf("failed to create note directory: %w", err)
	}

	// Write the file directly to disk — no Obsidian required.
	normalizedContent := NormalizeContent(params.Content)
	if err := WriteNoteFile(notePath, normalizedContent, params.ShouldAppend, params.ShouldOverwrite); err != nil {
		return err
	}

	if !params.ShouldOpen {
		return nil
	}

	if params.UseEditor {
		return obsidian.OpenInEditor(notePath)
	}

	// Open the note in Obsidian via URI.
	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"vault": vaultName,
		"file":  params.NoteName,
	})
	return uri.Execute(obsidianUri)
}

// WriteNoteFile writes content to notePath, respecting append/overwrite semantics.
// If the file already exists and neither flag is set, it is left unchanged.
func WriteNoteFile(notePath, content string, shouldAppend, shouldOverwrite bool) error {
	_, err := os.Stat(notePath)
	fileExists := err == nil

	if fileExists && shouldAppend {
		f, err := os.OpenFile(notePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open note for appending: %w", err)
		}
		defer f.Close()
		_, err = f.WriteString(content)
		return err
	}

	if fileExists && !shouldOverwrite {
		// File exists but no modification requested — leave it as-is.
		return nil
	}

	return os.WriteFile(notePath, []byte(content), 0644)
}

func NormalizeContent(content string) string {
	replacer := strings.NewReplacer(
		"\\n", "\n",
		"\\r", "\r",
		"\\t", "\t",
		"\\\\", "\\",
		"\\\"", "\"",
		"\\'", "'",
	)
	return replacer.Replace(content)
}
