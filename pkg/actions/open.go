package actions

import (
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
)

type OpenParams struct {
	NoteName string
	Section  string
}

func OpenNote(vault obsidian.VaultManager, uri obsidian.UriManager, params OpenParams) error {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

	fileParam := params.NoteName
	if params.Section != "" {
		fileParam = params.NoteName + "#" + params.Section
	}

	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"vault": vaultName,
		"file":  fileParam,
	})

	err = uri.Execute(obsidianUri)
	if err != nil {
		return err
	}
	return nil
}
