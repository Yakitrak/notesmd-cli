package actions

import (
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

func DailyNote(vault obsidian.VaultManager, uri obsidian.UriManager) error {
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

	obsidianUri := uri.Construct(OnsDailyUrl, map[string]string{
		"vault": vaultIdentifier,
	})

	err = uri.Execute(obsidianUri)
	if err != nil {
		return err
	}
	return nil
}
