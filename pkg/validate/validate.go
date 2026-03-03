package validate

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
)

// RuleFunc is a function that checks a note and returns validation results.
type RuleFunc func(content, filePath string, config Config) []ValidationResult

// ValidateVault walks a vault directory, validates all markdown notes, and returns a summary.
func ValidateVault(vault obsidian.VaultManager, config Config, rules []RuleFunc) (*ValidationSummary, error) {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return nil, err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return nil, err
	}

	summary := &ValidationSummary{
		Vault: vaultName,
	}

	excludeSet := make(map[string]bool, len(config.ExcludeDirs))
	for _, dir := range config.ExcludeDirs {
		excludeSet[dir] = true
	}

	err = filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil // Skip inaccessible paths
		}

		// Skip excluded directories
		if d.IsDir() {
			if excludeSet[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .md files
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		summary.TotalNotes++

		relPath, err := filepath.Rel(vaultPath, path)
		if err != nil {
			relPath = path
		}

		content, err := os.ReadFile(path)
		if err != nil {
			// Report unreadable files as errors
			if config.ShouldReport(SeverityError) {
				summary.Results = append(summary.Results, ValidationResult{
					File:     relPath,
					Severity: SeverityError,
					Rule:     "file.unreadable",
					Message:  "Could not read file",
				})
			}
			return nil
		}

		// Run all rules on this note
		for _, rule := range rules {
			results := rule(string(content), relPath, config)
			summary.Results = append(summary.Results, results...)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	summary.ComputeSummary()
	return summary, nil
}

// ValidatePath validates a single note at the given path within the vault.
func ValidatePath(vault obsidian.VaultManager, notePath string, config Config, rules []RuleFunc) (*ValidationSummary, error) {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return nil, err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(vaultPath, notePath)
	if !strings.HasSuffix(fullPath, ".md") {
		fullPath += ".md"
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	summary := &ValidationSummary{
		Vault:      vaultName,
		TotalNotes: 1,
	}

	for _, rule := range rules {
		results := rule(string(content), notePath, config)
		summary.Results = append(summary.Results, results...)
	}

	summary.ComputeSummary()
	return summary, nil
}
