package obsidian

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func AddMdSuffix(str string) string {
	if !strings.HasSuffix(str, ".md") {
		return str + ".md"
	}
	return str
}

func RemoveMdSuffix(str string) string {
	if strings.HasSuffix(str, ".md") {
		return strings.TrimSuffix(str, ".md")
	}
	return str
}

// normalizePathSeparators converts backslashes to forward slashes for cross-platform consistency.
// Obsidian uses forward slashes in links regardless of OS.
func normalizePathSeparators(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// wikiLinkPatterns returns the three wikilink patterns for a note name:
// [[name]], [[name|, [[name#
func wikiLinkPatterns(name string) [3]string {
	return [3]string{
		"[[" + name + "]]",
		"[[" + name + "|",
		"[[" + name + "#",
	}
}

func GenerateNoteLinkTexts(noteName string) [3]string {
	noteName = filepath.Base(noteName)
	noteName = RemoveMdSuffix(noteName)
	return wikiLinkPatterns(noteName)
}

// GenerateBacklinkSearchPatterns creates patterns to find links pointing to a note.
func GenerateBacklinkSearchPatterns(notePath string) []string {
	normalized := normalizePathSeparators(notePath)
	pathNoExt := RemoveMdSuffix(normalized)
	baseName := RemoveMdSuffix(path.Base(normalized))

	// Wikilinks (basename patterns)
	basePatterns := wikiLinkPatterns(baseName)
	patterns := append([]string{}, basePatterns[:]...)

	// Path-based wikilinks (only if path differs from basename)
	if pathNoExt != baseName {
		pathPatterns := wikiLinkPatterns(pathNoExt)
		patterns = append(patterns, pathPatterns[:]...)
	}

	// Markdown links
	mdPath := AddMdSuffix(normalized)
	patterns = append(patterns,
		"]("+mdPath+")",
		"]("+pathNoExt+")",
		"](./"+mdPath+")",
		"](./"+pathNoExt+")",
	)

	return patterns
}

// GenerateLinkReplacements creates all replacement patterns for updating links when moving a note.
// It normalizes path separators to forward slashes for cross-platform consistency,
// as Obsidian uses forward slashes in links regardless of operating system.
// This handles:
// - Simple wikilinks: [[note]], [[note|alias]], [[note#heading]]
// - Path-based wikilinks: [[folder/note]], [[folder/note|alias]], [[folder/note#heading]]
// - Markdown links: [text](folder/note.md), [text](./folder/note.md)
func GenerateLinkReplacements(oldNotePath, newNotePath string) map[string]string {
	replacements := make(map[string]string)

	// Normalize paths to forward slashes for consistent matching
	oldNormalized := normalizePathSeparators(oldNotePath)
	newNormalized := normalizePathSeparators(newNotePath)

	// Get basename without .md extension (use path.Base on normalized paths for cross-platform consistency)
	oldBase := RemoveMdSuffix(path.Base(oldNormalized))
	newBase := RemoveMdSuffix(path.Base(newNormalized))

	// Get full path without .md extension
	oldPathNoExt := RemoveMdSuffix(oldNormalized)
	newPathNoExt := RemoveMdSuffix(newNormalized)

	// 1. Simple wikilinks (basename only) - for backward compatibility
	replacements["[["+oldBase+"]]"] = "[[" + newBase + "]]"
	replacements["[["+oldBase+"|"] = "[[" + newBase + "|"
	replacements["[["+oldBase+"#"] = "[[" + newBase + "#"

	// 2. Path-based wikilinks (only if path differs from basename)
	if oldPathNoExt != oldBase {
		replacements["[["+oldPathNoExt+"]]"] = "[[" + newPathNoExt + "]]"
		replacements["[["+oldPathNoExt+"|"] = "[[" + newPathNoExt + "|"
		replacements["[["+oldPathNoExt+"#"] = "[[" + newPathNoExt + "#"
	}

	// 3. Markdown links (various formats)
	oldMd := AddMdSuffix(oldNormalized)
	newMd := AddMdSuffix(newNormalized)

	// Standard markdown link: [text](folder/note.md)
	replacements["]("+oldMd+")"] = "](" + newMd + ")"
	replacements["]("+oldPathNoExt+")"] = "](" + newPathNoExt + ")"

	// Relative markdown link: [text](./folder/note.md)
	replacements["](./"+oldMd+")"] = "](./" + newMd + ")"
	replacements["](./"+oldPathNoExt+")"] = "](./" + newPathNoExt + ")"

	return replacements
}

func ReplaceContent(content []byte, replacements map[string]string) []byte {
	for o, n := range replacements {
		content = bytes.ReplaceAll(content, []byte(o), []byte(n))
	}
	return content
}

func ShouldSkipDirectoryOrFile(info os.FileInfo) bool {
	isDirectory := info.IsDir()
	isHidden := info.Name()[0] == '.'
	isNonMarkdownFile := filepath.Ext(info.Name()) != ".md"
	if isDirectory || isHidden || isNonMarkdownFile {
		return true
	}
	return false
}

// OpenInEditor opens the specified file path in the user's preferred editor.
// It supports common GUI editors with appropriate wait flags and handles
// EDITOR values that contain arguments (e.g., "code -w").
func OpenInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default fallback
	}

	// Split EDITOR into command and any user-provided arguments.
	parts := strings.Fields(editor)
	editorBin := parts[0]
	userArgs := parts[1:]

	// Build arguments: user-provided args + auto-detected wait flag + file path.
	var args []string
	args = append(args, userArgs...)

	editorLower := strings.ToLower(filepath.Base(editorBin))
	needsWait := strings.Contains(editorLower, "code") ||
		strings.Contains(editorLower, "vscode") ||
		strings.Contains(editorLower, "subl") ||
		strings.Contains(editorLower, "atom") ||
		strings.Contains(editorLower, "mate")

	// Only add --wait if the user hasn't already included it.
	if needsWait && !containsWaitFlag(userArgs) {
		args = append(args, "--wait")
	}

	args = append(args, filePath)

	cmd := exec.Command(editorBin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open file in editor '%s': %w", editor, err)
	}

	return nil
}

// containsWaitFlag checks if any of the args already include a wait-style flag.
func containsWaitFlag(args []string) bool {
	for _, a := range args {
		if a == "--wait" || a == "-w" {
			return true
		}
	}
	return false
}
