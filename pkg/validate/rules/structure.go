package rules

import (
	"fmt"
	"strings"

	"github.com/Yakitrak/notesmd-cli/pkg/validate"
)

// CheckStructure runs structural validation rules on note content.
// The filePath is used only for result reporting.
func CheckStructure(content, filePath string, config validate.Config) []validate.ValidationResult {
	var results []validate.ValidationResult

	// Rule: structure.heading-hierarchy
	results = append(results, checkHeadingHierarchy(content, filePath, config)...)

	return results
}

// checkHeadingHierarchy detects skipped heading levels (e.g. H1 -> H3 without H2).
func checkHeadingHierarchy(content, filePath string, config validate.Config) []validate.ValidationResult {
	if !config.ShouldReport(validate.SeverityWarning) {
		return nil
	}

	var results []validate.ValidationResult
	lines := strings.Split(content, "\n")
	lastLevel := 0
	inFrontmatter := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip frontmatter block
		if trimmed == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			continue
		}

		// Skip code blocks
		if strings.HasPrefix(trimmed, "```") {
			continue
		}

		// Detect ATX headings (# H1, ## H2, etc.)
		if !strings.HasPrefix(trimmed, "#") {
			continue
		}

		level := 0
		for _, ch := range trimmed {
			if ch == '#' {
				level++
			} else {
				break
			}
		}

		// Must have a space after the hashes to be a valid heading
		if level >= len(trimmed) || trimmed[level] != ' ' {
			continue
		}

		if level < 1 || level > 6 {
			continue
		}

		if lastLevel > 0 && level > lastLevel+1 {
			results = append(results, validate.ValidationResult{
				File:       filePath,
				Severity:   validate.SeverityWarning,
				Rule:       "structure.heading-hierarchy",
				Message:    fmt.Sprintf("Heading level skipped: H%d to H%d", lastLevel, level),
				Suggestion: fmt.Sprintf("Add an H%d heading before this H%d", lastLevel+1, level),
				Line:       i + 1,
			})
		}

		lastLevel = level
	}

	return results
}
