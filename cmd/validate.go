package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/Yakitrak/notesmd-cli/pkg/validate"
	"github.com/Yakitrak/notesmd-cli/pkg/validate/rules"
	"github.com/spf13/cobra"
)

var outputFormat string
var minSeverity string

var validateCmd = &cobra.Command{
	Use:     "validate [path]",
	Aliases: []string{"val"},
	Short:   "Validate notes in a vault",
	Long: `Validate notes in an Obsidian vault for common issues.

Checks frontmatter syntax, deprecated properties, tag formatting,
empty notes, and heading hierarchy.

If a path is provided, only that note is validated.
Otherwise, all notes in the vault are validated.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}

		config := validate.DefaultConfig()

		// Parse severity flag
		switch strings.ToLower(minSeverity) {
		case "error":
			config.MinSeverity = validate.SeverityError
		case "warning":
			config.MinSeverity = validate.SeverityWarning
		case "info":
			config.MinSeverity = validate.SeverityInfo
		default:
			log.Fatalf("Invalid severity: %s (use error, warning, or info)", minSeverity)
		}

		// Build the list of rules
		allRules := []validate.RuleFunc{
			rules.CheckFrontmatter,
			rules.CheckStructure,
		}

		var summary *validate.ValidationSummary
		var err error

		if len(args) == 1 {
			summary, err = validate.ValidatePath(&vault, args[0], config, allRules)
		} else {
			summary, err = validate.ValidateVault(&vault, config, allRules)
		}

		if err != nil {
			log.Fatal(err)
		}

		// Output results
		switch strings.ToLower(outputFormat) {
		case "json":
			out, jsonErr := summary.ToJSON()
			if jsonErr != nil {
				log.Fatal(jsonErr)
			}
			fmt.Println(out)
		case "summary":
			printSummaryOnly(summary)
		default:
			printTextOutput(summary)
		}

		// Exit code 1 if errors found
		if summary.HasErrors() {
			os.Exit(1)
		}
	},
}

func printTextOutput(s *validate.ValidationSummary) {
	fmt.Printf("Validating vault: %s (%d notes)\n\n", s.Vault, s.TotalNotes)

	// Group results by file
	grouped := make(map[string][]validate.ValidationResult)
	var fileOrder []string
	for _, r := range s.Results {
		if _, seen := grouped[r.File]; !seen {
			fileOrder = append(fileOrder, r.File)
		}
		grouped[r.File] = append(grouped[r.File], r)
	}

	for _, file := range fileOrder {
		for _, r := range grouped[file] {
			label := severityLabel(r.Severity)
			fmt.Printf("%-5s %s\n", label, r.File)
			msg := r.Message
			if r.Line > 0 {
				msg = fmt.Sprintf("%s (line %d)", msg, r.Line)
			}
			fmt.Printf("      %s\n", msg)
			if r.Suggestion != "" {
				fmt.Printf("      -> %s\n", r.Suggestion)
			}
			fmt.Println()
		}
	}

	printSummaryLine(s)
}

func printSummaryOnly(s *validate.ValidationSummary) {
	fmt.Printf("Vault: %s\n", s.Vault)
	printSummaryLine(s)
}

func printSummaryLine(s *validate.ValidationSummary) {
	fmt.Printf("Summary: %d notes, %d errors, %d warnings\n",
		s.TotalNotes, s.Summary.Errors, s.Summary.Warnings)
}

func severityLabel(s validate.Severity) string {
	switch s {
	case validate.SeverityError:
		return "ERR"
	case validate.SeverityWarning:
		return "WARN"
	case validate.SeverityInfo:
		return "INFO"
	default:
		return "???"
	}
}

func init() {
	validateCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	validateCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "output format: text, json, summary")
	validateCmd.Flags().StringVarP(&minSeverity, "severity", "s", "warning", "minimum severity: error, warning, info")
	rootCmd.AddCommand(validateCmd)
}
