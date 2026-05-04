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
	Page                int
	PageSize            int
}

type searchContentJSONMatch struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	Content   string `json:"content"`
	MatchType string `json:"match_type"`
}

type searchContentPaginatedJSON struct {
	Page            int                      `json:"page"`
	PageSize        int                      `json:"page_size"`
	TotalResults    int                      `json:"total_results"`
	ReturnedResults int                      `json:"returned_results"`
	HasMore         bool                     `json:"has_more"`
	Results         []searchContentJSONMatch `json:"results"`
}

const (
	defaultPageSize = 25
	maxPageSize     = 100
)

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
		return printMatches(matches, searchTerm, format, output, options)
	}

	if len(matches) == 0 {
		_, _ = fmt.Fprintf(output, "No notes found containing '%s'\n", searchTerm)
		return nil
	}

	if len(matches) == 1 {
		_, _ = fmt.Fprintf(output, "Opening note: %s\n", matches[0].FilePath)
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
		_, _ = fmt.Fprintf(output, "Opening note: %s\n", selectedMatch.FilePath)
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

func paginateMatches(matches []obsidian.NoteMatch, options SearchContentOptions) ([]obsidian.NoteMatch, int, int, bool) {
	page := options.Page
	pageSize := options.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	total := len(matches)
	start := (page - 1) * pageSize
	if start >= total {
		return nil, page, pageSize, false
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return matches[start:end], page, pageSize, end < total
}

func isPaginationRequested(options SearchContentOptions) bool {
	return options.Page > 0 || options.PageSize > 0
}

func printMatches(matches []obsidian.NoteMatch, searchTerm string, format string, output io.Writer, options SearchContentOptions) error {
	paginate := isPaginationRequested(options)

	switch format {
	case searchContentFormatText:
		if len(matches) == 0 {
			fmt.Fprintf(os.Stderr, "No notes found containing '%s'\n", searchTerm)
			return nil
		}

		displayMatches := matches
		if paginate {
			var page, pageSize int
			var hasMore bool
			displayMatches, page, pageSize, hasMore = paginateMatches(matches, options)
			_ = hasMore
			for _, match := range displayMatches {
				_, _ = fmt.Fprintln(output, formatMatchForList(match))
			}
			total := len(matches)
			totalPages := (total + pageSize - 1) / pageSize
			_, _ = fmt.Fprintf(output, "-- Page %d/%d (%d of %d results) --\n", page, totalPages, len(displayMatches), total)
			return nil
		}

		for _, match := range displayMatches {
			_, _ = fmt.Fprintln(output, formatMatchForList(match))
		}
		return nil
	case searchContentFormatJSON:
		if paginate {
			pageMatches, page, pageSize, hasMore := paginateMatches(matches, options)
			result := make([]searchContentJSONMatch, 0, len(pageMatches))
			for _, match := range pageMatches {
				result = append(result, searchContentJSONMatch{
					File:      match.FilePath,
					Line:      match.LineNumber,
					Content:   match.MatchLine,
					MatchType: getMatchType(match),
				})
			}
			paginated := searchContentPaginatedJSON{
				Page:            page,
				PageSize:        pageSize,
				TotalResults:    len(matches),
				ReturnedResults: len(result),
				HasMore:         hasMore,
				Results:         result,
			}
			encoder := json.NewEncoder(output)
			encoder.SetEscapeHTML(false)
			return encoder.Encode(paginated)
		}

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
