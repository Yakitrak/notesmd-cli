package rules_test

import (
	"testing"

	"github.com/Yakitrak/notesmd-cli/pkg/validate"
	"github.com/Yakitrak/notesmd-cli/pkg/validate/rules"
	"github.com/stretchr/testify/assert"
)

func defaultConfig() validate.Config {
	return validate.DefaultConfig()
}

func TestCheckFrontmatter_EmptyNote(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
		rule    string
	}{
		{"empty string", "", 1, "structure.empty-note"},
		{"whitespace only", "   \n\n  ", 1, "structure.empty-note"},
		{"non-empty without frontmatter", "# Hello\nWorld", 1, "frontmatter.missing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rules.CheckFrontmatter(tt.content, "test.md", defaultConfig())
			assert.Len(t, results, tt.want)
			if tt.want > 0 {
				assert.Equal(t, tt.rule, results[0].Rule)
				assert.Equal(t, validate.SeverityWarning, results[0].Severity)
			}
		})
	}
}

func TestCheckFrontmatter_Missing(t *testing.T) {
	t.Run("note without frontmatter", func(t *testing.T) {
		content := "# Title\n\nSome content"
		results := rules.CheckFrontmatter(content, "test.md", defaultConfig())
		assert.Len(t, results, 1)
		assert.Equal(t, "frontmatter.missing", results[0].Rule)
		assert.Equal(t, validate.SeverityWarning, results[0].Severity)
	})

	t.Run("note with frontmatter", func(t *testing.T) {
		content := "---\ntitle: Test\n---\n# Title"
		results := rules.CheckFrontmatter(content, "test.md", defaultConfig())
		assert.Empty(t, results)
	})
}

func TestCheckFrontmatter_ParseError(t *testing.T) {
	t.Run("invalid YAML", func(t *testing.T) {
		content := "---\ntitle: [invalid yaml\n---\n"
		results := rules.CheckFrontmatter(content, "test.md", defaultConfig())
		assert.Len(t, results, 1)
		assert.Equal(t, "frontmatter.parse-error", results[0].Rule)
		assert.Equal(t, validate.SeverityError, results[0].Severity)
	})
}

func TestCheckFrontmatter_DeprecatedSingular(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"tags is valid", "---\ntags:\n  - go\n---\n", 0},
		{"tag triggers warning", "---\ntag: go\n---\n", 1},
		{"alias triggers warning", "---\nalias: test\n---\n", 1},
		{"cssclass triggers warning", "---\ncssclass: wide\n---\n", 1},
		{"aliases is valid", "---\naliases:\n  - test\n---\n", 0},
		{"cssclasses is valid", "---\ncssclasses:\n  - wide\n---\n", 0},
		{"both tag and alias trigger warnings", "---\ntag: go\nalias: test\n---\n", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rules.CheckFrontmatter(tt.content, "test.md", defaultConfig())
			deprecatedResults := filterByRule(results, "frontmatter.deprecated-singular")
			assert.Len(t, deprecatedResults, tt.want)
			for _, r := range deprecatedResults {
				assert.Equal(t, validate.SeverityWarning, r.Severity)
				assert.NotEmpty(t, r.Suggestion)
			}
		})
	}
}

func TestCheckFrontmatter_ListType(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"tags as list is valid", "---\ntags:\n  - go\n  - cli\n---\n", 0},
		{"tags as string is error", "---\ntags: go\n---\n", 1},
		{"aliases as string is error", "---\naliases: test\n---\n", 1},
		{"cssclasses as string is error", "---\ncssclasses: wide\n---\n", 1},
		{"missing key is ok", "---\ntitle: Test\n---\n", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rules.CheckFrontmatter(tt.content, "test.md", defaultConfig())
			listResults := filterByRule(results, "frontmatter.list-type")
			assert.Len(t, listResults, tt.want)
			for _, r := range listResults {
				assert.Equal(t, validate.SeverityError, r.Severity)
			}
		})
	}
}

func TestCheckFrontmatter_TagsHash(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"tags without hash", "---\ntags:\n  - go\n  - cli\n---\n", 0},
		{"tags with hash prefix", "---\ntags:\n  - \"#go\"\n  - cli\n---\n", 1},
		{"single tag string with hash", "---\ntags: \"#go\"\n---\n", 1}, // tags as string also triggers hash check
		{"no tags key", "---\ntitle: Test\n---\n", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rules.CheckFrontmatter(tt.content, "test.md", defaultConfig())
			hashResults := filterByRule(results, "frontmatter.tags-hash")
			assert.Len(t, hashResults, tt.want)
		})
	}
}

func TestCheckFrontmatter_SeverityFilter(t *testing.T) {
	t.Run("error-only severity skips warnings", func(t *testing.T) {
		content := "---\ntag: go\n---\n" // deprecated-singular is a warning
		config := validate.Config{
			ExcludeDirs: validate.DefaultExcludeDirs(),
			MinSeverity: validate.SeverityError,
		}
		results := rules.CheckFrontmatter(content, "test.md", config)
		assert.Empty(t, results)
	})

	t.Run("error severity reports parse errors", func(t *testing.T) {
		content := "---\ntitle: [broken\n---\n"
		config := validate.Config{
			ExcludeDirs: validate.DefaultExcludeDirs(),
			MinSeverity: validate.SeverityError,
		}
		results := rules.CheckFrontmatter(content, "test.md", config)
		assert.Len(t, results, 1)
		assert.Equal(t, "frontmatter.parse-error", results[0].Rule)
	})
}

func TestCheckFrontmatter_FilePathInResults(t *testing.T) {
	t.Run("file path is preserved in results", func(t *testing.T) {
		content := "# No frontmatter"
		results := rules.CheckFrontmatter(content, "daily/2024-01-15.md", defaultConfig())
		assert.Len(t, results, 1)
		assert.Equal(t, "daily/2024-01-15.md", results[0].File)
	})
}

func filterByRule(results []validate.ValidationResult, rule string) []validate.ValidationResult {
	var filtered []validate.ValidationResult
	for _, r := range results {
		if r.Rule == rule {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
