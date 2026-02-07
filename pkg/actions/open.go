package actions

import (
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
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

	// Get vault info to use ID if available
	vaultInfo, err := vault.GetVaultInfo()
	vaultIdentifier := vaultName
	if err == nil && vaultInfo.ID != "" {
		// Use vault ID for more reliable vault resolution
		vaultIdentifier = vaultInfo.ID
	}

	fileParam := params.NoteName
	if params.Section != "" {
		fileParam = params.NoteName + "#" + params.Section
	}

	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"vault": vaultIdentifier,
		"file":  fileParam,
	})

	err = uri.Execute(obsidianUri)
	if err != nil {
		return err
	}
	return nil
}
