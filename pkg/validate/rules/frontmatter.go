package rules

import (
	"fmt"
	"strings"

	"github.com/Yakitrak/notesmd-cli/pkg/frontmatter"
	"github.com/Yakitrak/notesmd-cli/pkg/validate"
)

// deprecatedSingularKeys maps deprecated singular keys to their correct plural forms.
var deprecatedSingularKeys = map[string]string{
	"tag":      "tags",
	"alias":    "aliases",
	"cssclass": "cssclasses",
}

// listTypeKeys are frontmatter keys that must be YAML lists, not strings.
var listTypeKeys = []string{"tags", "aliases", "cssclasses"}

// CheckFrontmatter runs all frontmatter validation rules on a note's content.
// The filePath is used only for result reporting.
func CheckFrontmatter(content, filePath string, config validate.Config) []validate.ValidationResult {
	var results []validate.ValidationResult

	// Rule: structure.empty-note — check before frontmatter rules
	trimmed := strings.TrimSpace(content)
	if len(trimmed) == 0 {
		if config.ShouldReport(validate.SeverityWarning) {
			results = append(results, validate.ValidationResult{
				File:     filePath,
				Severity: validate.SeverityWarning,
				Rule:     "structure.empty-note",
				Message:  "Note is empty",
			})
		}
		return results
	}

	// Rule: frontmatter.missing
	if !frontmatter.HasFrontmatter(content) {
		if config.ShouldReport(validate.SeverityWarning) {
			results = append(results, validate.ValidationResult{
				File:     filePath,
				Severity: validate.SeverityWarning,
				Rule:     "frontmatter.missing",
				Message:  "Note has no frontmatter",
			})
		}
		return results
	}

	// Rule: frontmatter.parse-error
	fm, _, err := frontmatter.Parse(content)
	if err != nil {
		if config.ShouldReport(validate.SeverityError) {
			results = append(results, validate.ValidationResult{
				File:     filePath,
				Severity: validate.SeverityError,
				Rule:     "frontmatter.parse-error",
				Message:  "YAML frontmatter parse error",
			})
		}
		return results
	}

	if fm == nil {
		return results
	}

	// Rule: frontmatter.deprecated-singular
	for singular, plural := range deprecatedSingularKeys {
		if _, ok := fm[singular]; ok {
			if config.ShouldReport(validate.SeverityWarning) {
				results = append(results, validate.ValidationResult{
					File:       filePath,
					Severity:   validate.SeverityWarning,
					Rule:       "frontmatter.deprecated-singular",
					Message:    fmt.Sprintf("Deprecated property '%s' found", singular),
					Suggestion: fmt.Sprintf("Use '%s' instead of '%s'", plural, singular),
				})
			}
		}
	}

	// Rule: frontmatter.list-type
	for _, key := range listTypeKeys {
		val, ok := fm[key]
		if !ok {
			continue
		}
		if !isList(val) {
			if config.ShouldReport(validate.SeverityError) {
				results = append(results, validate.ValidationResult{
					File:       filePath,
					Severity:   validate.SeverityError,
					Rule:       "frontmatter.list-type",
					Message:    fmt.Sprintf("Property '%s' should be a list, not a scalar", key),
					Suggestion: fmt.Sprintf("Change '%s: value' to '%s:\\n  - value'", key, key),
				})
			}
		}
	}

	// Rule: frontmatter.tags-hash
	if tags, ok := fm["tags"]; ok {
		checkTagsForHash(tags, filePath, config, &results)
	}

	return results
}

// isList returns true if the value is a YAML list (Go slice).
func isList(val interface{}) bool {
	switch val.(type) {
	case []interface{}:
		return true
	case []string:
		return true
	default:
		return false
	}
}

// checkTagsForHash checks if any tag in the tags list starts with '#'.
func checkTagsForHash(tags interface{}, filePath string, config validate.Config, results *[]validate.ValidationResult) {
	var tagStrings []string

	switch v := tags.(type) {
	case []interface{}:
		for _, t := range v {
			if s, ok := t.(string); ok {
				tagStrings = append(tagStrings, s)
			}
		}
	case []string:
		tagStrings = v
	case string:
		// Single tag as string — also check for hash
		tagStrings = []string{v}
	default:
		return
	}

	for _, tag := range tagStrings {
		if strings.HasPrefix(tag, "#") {
			if config.ShouldReport(validate.SeverityWarning) {
				*results = append(*results, validate.ValidationResult{
					File:       filePath,
					Severity:   validate.SeverityWarning,
					Rule:       "frontmatter.tags-hash",
					Message:    fmt.Sprintf("Tag '%s' has '#' prefix in YAML frontmatter", tag),
					Suggestion: "Remove '#' prefix from tags in frontmatter (Obsidian adds it automatically)",
				})
			}
			return // Report once per file, not per tag
		}
	}
}
