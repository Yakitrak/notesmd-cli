package actions

import (
	"fmt"
	"os"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
)

type OpenParams struct {
	NoteName  string
	Section   string
	UseEditor bool
}

func OpenNote(vault obsidian.VaultManager, uri obsidian.UriManager, params OpenParams) error {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

	if params.UseEditor {
		if params.Section != "" {
			fmt.Fprintln(os.Stderr, "Warning: --section is ignored when using --editor")
		}
		vaultPath, err := vault.Path()
		if err != nil {
			return err
		}
		filePath, err := obsidian.ValidatePath(vaultPath, obsidian.AddMdSuffix(params.NoteName))
		if err != nil {
			return err
		}
		return obsidian.OpenInEditor(filePath)
	}

	fileParam := params.NoteName
	if params.Section != "" {
		fileParam = params.NoteName + "#" + params.Section
	}

	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"vault": vaultName,
		"file":  fileParam,
	})

	return uri.Execute(obsidianUri)
}
