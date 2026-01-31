package actions

import (
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
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

	targetPath := vaultPath
	if strings.TrimSpace(params.Path) != "" {
		validatedPath, err := obsidian.ValidatePath(vaultPath, params.Path)
		if err != nil {
			return nil, err
		}
		targetPath = validatedPath
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, errors.New(obsidian.VaultAccessError)
	}
	if !info.IsDir() {
		return nil, errors.New(obsidian.VaultAccessError)
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, errors.New(obsidian.VaultReadError)
	}

	dirs := make([]string, 0, len(entries))
	files := make([]string, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, name+"/")
			continue
		}
		files = append(files, name)
	}

	sort.Strings(dirs)
	sort.Strings(files)

	return append(dirs, files...), nil
}
