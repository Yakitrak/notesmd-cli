package actions

import (
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
)

type ListParams struct {
	Path string
}

func ListEntries(vault obsidian.VaultManager, params ListParams) ([]string, error) {
	_, err := vault.DefaultName()
	if err != nil {
		return nil, err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return nil, err
	}

	return obsidian.ListEntries(vaultPath, params.Path)
}
