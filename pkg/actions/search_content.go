package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
)

const (
	searchContentFormatText = "text"
	searchContentFormatJSON = "json"
)

type SearchContentOptions struct {
	UseEditor           bool
	EditorFlagExplicit  bool
	NoInteractive       bool
	Format              string
	InteractiveTerminal bool
	Output              io.Writer
}

type searchContentJSONMatch struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	Content   string `json:"content"`
	MatchType string `json:"match_type"`
}

// SearchNotesContent preserves backward-compatible interactive behavior.
func SearchNotesContent(vault obsidian.VaultManager, note obsidian.NoteManager, uri obsidian.UriManager, fuzzyFinder obsidian.FuzzyFinderManager, searchTerm string, useEditor bool) error {
	return SearchNotesContentWithOptions(vault, note, uri, fuzzyFinder, searchTerm, SearchContentOptions{
		UseEditor:           useEditor,
		EditorFlagExplicit:  useEditor,
		Format:              searchContentFormatText,
		InteractiveTerminal: true,
		Output:              os.Stdout,
	})
}

func SearchNotesContentWithOptions(vault obsidian.VaultManager, note obsidian.NoteManager, uri obsidian.UriManager, fuzzyFinder obsidian.FuzzyFinderManager, searchTerm string, options SearchContentOptions) error {
	format, err := normalizeSearchContentFormat(options.Format)
	if err != nil {
		return err
	}

	nonInteractiveMode := shouldUseNonInteractiveMode(options, format)
	useEditor := options.UseEditor

	if nonInteractiveMode && options.EditorFlagExplicit && options.UseEditor {
		return errors.New("--editor cannot be used with non-interactive search-content output")
	}

	if nonInteractiveMode {
		// If editor mode came from config default rather than explicit flag,
		// prefer non-interactive output for script-friendly behavior.
		useEditor = false
	}

	output := options.Output
	if output == nil {
		output = os.Stdout
	}

	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	matches, err := note.SearchNotesWithSnippets(vaultPath, searchTerm)
	if err != nil {
		return err
	}

	if nonInteractiveMode {
		return printMatches(matches, searchTerm, format, output)
	}

	if len(matches) == 0 {
		fmt.Fprintf(output, "No notes found containing '%s'\n", searchTerm)
		return nil
	}

	if len(matches) == 1 {
		fmt.Fprintf(output, "Opening note: %s\n", matches[0].FilePath)
		if useEditor {
			filePath := filepath.Join(vaultPath, matches[0].FilePath)
			return obsidian.OpenInEditor(filePath)
		}
		obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
			"file":  matches[0].FilePath,
			"vault": vaultName,
		})
		return uri.Execute(obsidianUri)
	}

	displayItems := formatMatchesForDisplay(matches)

	index, err := fuzzyFinder.Find(displayItems, func(i int) string {
		return displayItems[i]
	})
	if err != nil {
		return err
	}

	selectedMatch := matches[index]
	if useEditor {
		filePath := filepath.Join(vaultPath, selectedMatch.FilePath)
		fmt.Fprintf(output, "Opening note: %s\n", selectedMatch.FilePath)
		return obsidian.OpenInEditor(filePath)
	}
	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"file":  selectedMatch.FilePath,
		"vault": vaultName,
	})
	return uri.Execute(obsidianUri)
}

func shouldUseNonInteractiveMode(options SearchContentOptions, format string) bool {
	if options.NoInteractive {
		return true
	}
	if format == searchContentFormatJSON {
		return true
	}
	return !options.InteractiveTerminal
}

func normalizeSearchContentFormat(format string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(format))
	if trimmed == "" {
		return searchContentFormatText, nil
	}

	switch trimmed {
	case searchContentFormatText, searchContentFormatJSON:
		return trimmed, nil
	default:
		return "", fmt.Errorf("invalid format '%s': expected one of text, json", format)
	}
}

func printMatches(matches []obsidian.NoteMatch, searchTerm string, format string, output io.Writer) error {
	switch format {
	case searchContentFormatText:
		if len(matches) == 0 {
			fmt.Fprintf(os.Stderr, "No notes found containing '%s'\n", searchTerm)
			return nil
		}
		for _, match := range matches {
			fmt.Fprintln(output, formatMatchForList(match))
		}
		return nil
	case searchContentFormatJSON:
		result := make([]searchContentJSONMatch, 0, len(matches))
		for _, match := range matches {
			result = append(result, searchContentJSONMatch{
				File:      match.FilePath,
				Line:      match.LineNumber,
				Content:   match.MatchLine,
				MatchType: getMatchType(match),
			})
		}

		encoder := json.NewEncoder(output)
		encoder.SetEscapeHTML(false)
		return encoder.Encode(result)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func formatMatchForList(match obsidian.NoteMatch) string {
	if match.LineNumber > 0 {
		return fmt.Sprintf("%s:%d: %s", match.FilePath, match.LineNumber, match.MatchLine)
	}
	return fmt.Sprintf("%s: %s", match.FilePath, match.MatchLine)
}

func getMatchType(match obsidian.NoteMatch) string {
	if match.LineNumber == 0 {
		return "filename"
	}
	return "content"
}

func formatMatchesForDisplay(matches []obsidian.NoteMatch) []string {
	maxPathLength := calculateMaxPathLength(matches)

	var displayItems []string
	for _, match := range matches {
		displayStr := formatSingleMatch(match, maxPathLength)
		displayItems = append(displayItems, displayStr)
	}

	return displayItems
}

func calculateMaxPathLength(matches []obsidian.NoteMatch) int {
	maxLength := 0
	for _, match := range matches {
		pathWithLine := formatPathWithLine(match)
		if len(pathWithLine) > maxLength {
			maxLength = len(pathWithLine)
		}
	}
	return maxLength
}

func formatPathWithLine(match obsidian.NoteMatch) string {
	if match.LineNumber > 0 {
		return fmt.Sprintf("%s:%d", match.FilePath, match.LineNumber)
	}
	return match.FilePath
}

func formatSingleMatch(match obsidian.NoteMatch, maxPathLength int) string {
	pathWithLine := formatPathWithLine(match)
	if match.LineNumber == 0 {
		// Filename match - show path and indicate it's a filename match
		return fmt.Sprintf("%-*s | %s", maxPathLength, pathWithLine, match.MatchLine)
	}
	// Content match - show path:line | snippet
	return fmt.Sprintf("%-*s | %s", maxPathLength, pathWithLine, match.MatchLine)
}
