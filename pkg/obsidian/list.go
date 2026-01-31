package obsidian

import (
	"errors"
	"os"
	"sort"
	"strings"
)

func ListEntries(vaultPath, relativePath string) ([]string, error) {
	targetPath := vaultPath
	if strings.TrimSpace(relativePath) != "" {
		validatedPath, err := ValidatePath(vaultPath, relativePath)
		if err != nil {
			return nil, err
		}
		targetPath = validatedPath
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, errors.New(VaultAccessError)
	}
	if !info.IsDir() {
		return nil, errors.New(VaultAccessError)
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, errors.New(VaultReadError)
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
