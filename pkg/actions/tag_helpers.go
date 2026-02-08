package actions

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

// TagSearchParams contains common parameters for tag searches
type TagSearchParams struct {
	Tag            string
	Location       string // "all", "frontmatter", or "inline"
	IncludeSubtags bool
}

// matchesTag checks if a tag matches the search criteria
func matchesTag(tag, searchTag string, includeSubtags bool) bool {
	tagLower := strings.ToLower(tag)
	searchLower := strings.ToLower(searchTag)

	if tagLower == searchLower {
		return true
	}
	if includeSubtags {
		return strings.HasPrefix(tagLower, searchLower+"/")
	}
	return false
}

// filterByLocation checks if a tag location matches the filter criteria
func filterByLocation(tm obsidian.TagMatch, location string) bool {
	if location == "all" {
		return true
	}
	if location == "frontmatter" && tm.Location == obsidian.TagFrontmatter {
		return true
	}
	if location == "inline" && tm.Location == obsidian.TagInline {
		return true
	}
	return false
}

// getVaultInfo retrieves common vault information
func getVaultInfo(vault obsidian.VaultManager) (vaultName, vaultPath string, err error) {
	vaultName, err = vault.DefaultName()
	if err != nil {
		return "", "", err
	}
	vaultPath, err = vault.Path()
	if err != nil {
		return "", "", err
	}
	return vaultName, vaultPath, nil
}

// openMatch opens a matched note in Obsidian or editor
func openMatch(uri obsidian.UriManager, vaultPath, notePath, vaultName string, useEditor bool) error {
	if useEditor {
		filePath := filepath.Join(vaultPath, notePath)
		return obsidian.OpenInEditor(filePath)
	}
	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"file":  notePath,
		"vault": vaultName,
	})
	return uri.Execute(obsidianUri)
}

// formatLocationMessage creates a consistent message based on search parameters
func formatLocationMessage(location string, includeSubtags bool) string {
	var parts []string
	if location != "" && location != "all" {
		parts = append(parts, fmt.Sprintf("in %s", location))
	}
	if includeSubtags {
		parts = append(parts, "including subtags")
	}
	if len(parts) > 0 {
		return fmt.Sprintf(" (%s)", strings.Join(parts, ", "))
	}
	return ""
}