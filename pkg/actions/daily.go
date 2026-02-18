package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
)

type DailyParams struct {
	UseEditor bool
}

func DailyNote(vault obsidian.VaultManager, uri obsidian.UriManager, params DailyParams) error {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	config := obsidian.ReadDailyNotesConfig(vaultPath)

	// Format today's date using the configured Moment.js format.
	format := config.Format
	if format == "" {
		format = "YYYY-MM-DD"
	}
	noteName := time.Now().Format(obsidian.MomentToGoFormat(format))

	// Prepend configured daily notes folder.
	if config.Folder != "" {
		noteName = config.Folder + "/" + noteName
	}

	notePath, err := obsidian.ValidatePath(vaultPath, obsidian.AddMdSuffix(noteName))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(notePath), 0755); err != nil {
		return fmt.Errorf("failed to create daily note directory: %w", err)
	}

	// Read template content if configured.
	content := ""
	if config.Template != "" {
		templatePath := filepath.Join(vaultPath, obsidian.AddMdSuffix(config.Template))
		if templateContent, readErr := os.ReadFile(templatePath); readErr == nil {
			content = string(templateContent)
		}
	}

	// WriteNoteFile leaves existing files unchanged (no append/overwrite).
	if err := WriteNoteFile(notePath, content, false, false); err != nil {
		return err
	}

	// Open the note.
	if params.UseEditor {
		return obsidian.OpenInEditor(notePath)
	}

	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"vault": vaultName,
		"file":  noteName,
	})
	return uri.Execute(obsidianUri)
}
