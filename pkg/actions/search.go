package actions

import (
	"fmt"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

// Search performs a unified tag search with multiple output formats
func Search(
	vault obsidian.VaultManager,
	note obsidian.NoteManager,
	uri obsidian.UriManager,
	tag string,
	format string,
	location string,
	includeSubtags bool,
	contextLines int,
	useEditor bool,
) error {
	vaultName, vaultPath, err := getVaultInfo(vault)
	if err != nil {
		return err
	}

	searcher := &obsidian.TagSearcher{}
	notes, err := note.GetNotesList(vaultPath)
	if err != nil {
		return err
	}

	switch format {
	case "names":
		return searchNames(vaultPath, notes, note, searcher, tag, location, includeSubtags, vaultName, uri, useEditor)
	case "locations":
		return searchLocations(vaultPath, notes, note, searcher, tag, location, includeSubtags, vaultName, uri, useEditor)
	case "lines":
		return searchLines(vaultPath, notes, note, searcher, tag, location, includeSubtags, contextLines, vaultName, uri, useEditor)
	case "blocks":
		return searchBlocks(vaultPath, notes, note, searcher, tag, location, includeSubtags, contextLines, vaultName, uri, useEditor)
	default:
		return fmt.Errorf("unknown format: %s (use: names, locations, lines, blocks)", format)
	}
}

// searchNames returns just unique filenames
func searchNames(vaultPath string, notes []string, note obsidian.NoteManager, searcher *obsidian.TagSearcher,
	tag, location string, includeSubtags bool, vaultName string, uri obsidian.UriManager, useEditor bool) error {

	foundNotes := make(map[string]bool)

	for _, notePath := range notes {
		content, err := note.GetContents(vaultPath, notePath)
		if err != nil {
			continue
		}

		tagMatches := searcher.ExtractTags(content)
		for _, tm := range tagMatches {
			if !matchesTag(tm.Tag, tag, includeSubtags) {
				continue
			}
			if !filterByLocation(tm, location) {
				continue
			}
			foundNotes[notePath] = true
			break
		}
	}

	if len(foundNotes) == 0 {
		fmt.Printf("No notes found with tag '#%s'%s\n", tag, formatLocationMessage(location, includeSubtags))
		return nil
	}

	fmt.Printf("Found %d notes with tag '#%s'%s\n\n", len(foundNotes), tag, formatLocationMessage(location, includeSubtags))

	if len(foundNotes) == 1 {
		for notePath := range foundNotes {
			if useEditor {
				return openMatch(uri, vaultPath, notePath, vaultName, useEditor)
			}
			fmt.Println(notePath)
			return nil
		}
	}

	for notePath := range foundNotes {
		fmt.Println(notePath)
	}
	return nil
}

// searchLocations returns filename + line number
func searchLocations(vaultPath string, notes []string, note obsidian.NoteManager, searcher *obsidian.TagSearcher,
	tag, location string, includeSubtags bool, vaultName string, uri obsidian.UriManager, useEditor bool) error {

	type match struct {
		filePath   string
		lineNumber int
		tag        string
	}
	var matches []match

	for _, notePath := range notes {
		content, err := note.GetContents(vaultPath, notePath)
		if err != nil {
			continue
		}

		tagMatches := searcher.ExtractTags(content)
		for _, tm := range tagMatches {
			if !matchesTag(tm.Tag, tag, includeSubtags) {
				continue
			}
			if !filterByLocation(tm, location) {
				continue
			}
			matches = append(matches, match{
				filePath:   notePath,
				lineNumber: tm.LineNumber,
				tag:        tm.Tag,
			})
			break // Only first match per note
		}
	}

	if len(matches) == 0 {
		fmt.Printf("No notes found with tag '#%s'%s\n", tag, formatLocationMessage(location, includeSubtags))
		return nil
	}

	fmt.Printf("Found %d notes with tag '#%s'%s\n\n", len(matches), tag, formatLocationMessage(location, includeSubtags))

	if len(matches) == 1 && useEditor {
		return openMatch(uri, vaultPath, matches[0].filePath, vaultName, useEditor)
	}

	for _, m := range matches {
		fmt.Printf("📄 %s:%d #%s\n", m.filePath, m.lineNumber, m.tag)
	}
	return nil
}

// searchLines returns full line content
func searchLines(vaultPath string, notes []string, note obsidian.NoteManager, searcher *obsidian.TagSearcher,
	tag, location string, includeSubtags bool, contextLines int, vaultName string, uri obsidian.UriManager, useEditor bool) error {

	type match struct {
		filePath   string
		lineNumber int
		content    string
	}
	var matches []match

	for _, notePath := range notes {
		content, err := note.GetContents(vaultPath, notePath)
		if err != nil {
			continue
		}

		lines := strings.Split(content, "\n")
		tagMatches := searcher.ExtractTags(content)

		for _, tm := range tagMatches {
			if !matchesTag(tm.Tag, tag, includeSubtags) {
				continue
			}
			if !filterByLocation(tm, location) {
				continue
			}

			// Get context lines
			var contextContent []string
			startLine := tm.LineNumber - contextLines - 1
			if startLine < 0 {
				startLine = 0
			}
			endLine := tm.LineNumber + contextLines
			if endLine > len(lines) {
				endLine = len(lines)
			}

			for i := startLine; i < endLine && i < len(lines); i++ {
				contextContent = append(contextContent, lines[i])
			}

			matches = append(matches, match{
				filePath:   notePath,
				lineNumber: tm.LineNumber,
				content:    strings.Join(contextContent, "\n"),
			})
		}
	}

	if len(matches) == 0 {
		fmt.Printf("No lines found with tag '#%s'%s\n", tag, formatLocationMessage(location, includeSubtags))
		return nil
	}

	fmt.Printf("Found %d occurrences of '#%s'%s\n\n", len(matches), tag, formatLocationMessage(location, includeSubtags))

	if len(matches) == 1 && useEditor {
		return openMatch(uri, vaultPath, matches[0].filePath, vaultName, useEditor)
	}

	for _, m := range matches {
		fmt.Printf("📄 %s:%d\n%s\n\n", m.filePath, m.lineNumber, m.content)
	}
	return nil
}

// searchBlocks returns paragraph blocks (until empty line)
func searchBlocks(vaultPath string, notes []string, note obsidian.NoteManager, searcher *obsidian.TagSearcher,
	tag, location string, includeSubtags bool, contextLines int, vaultName string, uri obsidian.UriManager, useEditor bool) error {

	type match struct {
		filePath   string
		lineNumber int
		block      string
	}
	var matches []match

	for _, notePath := range notes {
		content, err := note.GetContents(vaultPath, notePath)
		if err != nil {
			continue
		}

		lines := strings.Split(content, "\n")
		tagMatches := searcher.ExtractTags(content)

		for _, tm := range tagMatches {
			if !matchesTag(tm.Tag, tag, includeSubtags) {
				continue
			}
			if !filterByLocation(tm, location) {
				continue
			}

			// Find the block containing this line
			lineIdx := tm.LineNumber - 1
			if lineIdx < 0 || lineIdx >= len(lines) {
				continue
			}

			// Find block start (previous empty line or start of file)
			blockStart := lineIdx
			for i := lineIdx - 1; i >= 0; i-- {
				if strings.TrimSpace(lines[i]) == "" {
					blockStart = i + 1
					break
				}
				blockStart = i
			}

			// Find block end (next empty line or end of file)
			blockEnd := lineIdx
			for i := lineIdx + 1; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) == "" {
					blockEnd = i - 1
					break
				}
				blockEnd = i
			}

			// Extract block
			var blockLines []string
			for i := blockStart; i <= blockEnd && i < len(lines); i++ {
				blockLines = append(blockLines, lines[i])
			}

			matches = append(matches, match{
				filePath:   notePath,
				lineNumber: tm.LineNumber,
				block:      strings.Join(blockLines, "\n"),
			})
		}
	}

	if len(matches) == 0 {
		fmt.Printf("No blocks found with tag '#%s'%s\n", tag, formatLocationMessage(location, includeSubtags))
		return nil
	}

	fmt.Printf("Found %d blocks with tag '#%s'%s\n\n", len(matches), tag, formatLocationMessage(location, includeSubtags))

	if len(matches) == 1 && useEditor {
		return openMatch(uri, vaultPath, matches[0].filePath, vaultName, useEditor)
	}

	for _, m := range matches {
		fmt.Printf("📄 %s:%d\n%s\n\n", m.filePath, m.lineNumber, m.block)
	}
	return nil
}