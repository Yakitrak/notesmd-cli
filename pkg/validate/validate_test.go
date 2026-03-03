package validate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/validate"
	"github.com/Yakitrak/notesmd-cli/pkg/validate/rules"
	"github.com/stretchr/testify/assert"
)

func allRules() []validate.RuleFunc {
	return []validate.RuleFunc{
		rules.CheckFrontmatter,
		rules.CheckStructure,
	}
}

func createNote(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.MkdirAll(filepath.Dir(path), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
}

func TestValidateVault_Integration(t *testing.T) {
	t.Run("mixed vault with valid and invalid notes", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := &mocks.MockVaultOperator{Name: "TestVault", PathValue: tmpDir}

		// Good note - valid frontmatter
		createNote(t, tmpDir, "good.md", "---\ntitle: Good Note\ntags:\n  - go\n---\n# Good Note\n\nContent here.\n")

		// Bad YAML - parse error
		createNote(t, tmpDir, "bad-yaml.md", "---\ntitle: [broken yaml\n---\n")

		// Deprecated property
		createNote(t, tmpDir, "deprecated.md", "---\ntag: oldstyle\n---\n# Note\n")

		// Empty note
		createNote(t, tmpDir, "empty.md", "")

		// Obsidian config dir - should be skipped
		createNote(t, tmpDir, ".obsidian/app.json", `{"key": "value"}`)

		// Trash dir - should be skipped
		createNote(t, tmpDir, ".trash/deleted.md", "---\ntitle: Deleted\n---\n")

		config := validate.DefaultConfig()
		summary, err := validate.ValidateVault(vault, config, allRules())
		assert.NoError(t, err)

		// Should count 4 notes (not .obsidian/app.json or .trash/deleted.md)
		assert.Equal(t, 4, summary.TotalNotes)
		assert.Equal(t, "TestVault", summary.Vault)

		// Should have findings
		assert.True(t, summary.HasErrors(), "should have errors from bad-yaml.md")
		assert.Greater(t, summary.Summary.Warnings, 0, "should have warnings")

		// Verify specific rules were triggered
		ruleMap := make(map[string]int)
		for _, r := range summary.Results {
			ruleMap[r.Rule]++
		}

		assert.Equal(t, 1, ruleMap["frontmatter.parse-error"], "bad YAML should trigger parse error")
		assert.Equal(t, 1, ruleMap["structure.empty-note"], "empty note should be detected")
		assert.Equal(t, 1, ruleMap["frontmatter.deprecated-singular"], "deprecated tag should be detected")
	})

	t.Run("vault with only valid notes", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := &mocks.MockVaultOperator{Name: "CleanVault", PathValue: tmpDir}

		createNote(t, tmpDir, "note1.md", "---\ntitle: Note 1\ntags:\n  - test\n---\n# Note 1\n\nContent.\n")
		createNote(t, tmpDir, "note2.md", "---\ntitle: Note 2\n---\n## Section\n\nMore content.\n")

		config := validate.DefaultConfig()
		summary, err := validate.ValidateVault(vault, config, allRules())
		assert.NoError(t, err)
		assert.Equal(t, 2, summary.TotalNotes)
		assert.Empty(t, summary.Results)
		assert.False(t, summary.HasErrors())
	})

	t.Run("empty vault", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := &mocks.MockVaultOperator{Name: "EmptyVault", PathValue: tmpDir}

		config := validate.DefaultConfig()
		summary, err := validate.ValidateVault(vault, config, allRules())
		assert.NoError(t, err)
		assert.Equal(t, 0, summary.TotalNotes)
		assert.Empty(t, summary.Results)
	})

	t.Run("vault path error", func(t *testing.T) {
		vault := &mocks.MockVaultOperator{
			Name:      "BadVault",
			PathError: assert.AnError,
		}

		config := validate.DefaultConfig()
		_, err := validate.ValidateVault(vault, config, allRules())
		assert.Error(t, err)
	})
}

func TestValidatePath_Integration(t *testing.T) {
	t.Run("validate single note", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := &mocks.MockVaultOperator{Name: "TestVault", PathValue: tmpDir}

		createNote(t, tmpDir, "daily/2024-01-15.md", "---\ntag: daily\n---\n# January 15\n")

		config := validate.DefaultConfig()
		summary, err := validate.ValidatePath(vault, "daily/2024-01-15.md", config, allRules())
		assert.NoError(t, err)
		assert.Equal(t, 1, summary.TotalNotes)

		deprecatedResults := filterByRule(summary.Results, "frontmatter.deprecated-singular")
		assert.Len(t, deprecatedResults, 1)
	})

	t.Run("validate note without .md extension", func(t *testing.T) {
		tmpDir := t.TempDir()
		vault := &mocks.MockVaultOperator{Name: "TestVault", PathValue: tmpDir}

		createNote(t, tmpDir, "note.md", "---\ntitle: Test\n---\n# Note\n")

		config := validate.DefaultConfig()
		summary, err := validate.ValidatePath(vault, "note", config, allRules())
		assert.NoError(t, err)
		assert.Equal(t, 1, summary.TotalNotes)
		assert.Empty(t, summary.Results)
	})
}

func TestValidationSummary_JSON(t *testing.T) {
	t.Run("JSON output is valid", func(t *testing.T) {
		summary := &validate.ValidationSummary{
			Vault:      "TestVault",
			TotalNotes: 10,
			Results: []validate.ValidationResult{
				{
					File:     "test.md",
					Severity: validate.SeverityError,
					Rule:     "frontmatter.parse-error",
					Message:  "YAML parse error",
				},
			},
		}
		summary.ComputeSummary()

		jsonStr, err := summary.ToJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, `"vault": "TestVault"`)
		assert.Contains(t, jsonStr, `"frontmatter.parse-error"`)
	})
}

func TestConfig_ShouldReport(t *testing.T) {
	tests := []struct {
		name        string
		minSeverity validate.Severity
		check       validate.Severity
		expected    bool
	}{
		{"error reports error", validate.SeverityError, validate.SeverityError, true},
		{"error skips warning", validate.SeverityError, validate.SeverityWarning, false},
		{"error skips info", validate.SeverityError, validate.SeverityInfo, false},
		{"warning reports error", validate.SeverityWarning, validate.SeverityError, true},
		{"warning reports warning", validate.SeverityWarning, validate.SeverityWarning, true},
		{"warning skips info", validate.SeverityWarning, validate.SeverityInfo, false},
		{"info reports everything", validate.SeverityInfo, validate.SeverityError, true},
		{"info reports warning", validate.SeverityInfo, validate.SeverityWarning, true},
		{"info reports info", validate.SeverityInfo, validate.SeverityInfo, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := validate.Config{MinSeverity: tt.minSeverity}
			assert.Equal(t, tt.expected, config.ShouldReport(tt.check))
		})
	}
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
