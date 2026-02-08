package cmd

import (
	"log"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var searchFormat string
var searchLocation string
var includeSubtags bool
var contextLines int

var searchCmd = &cobra.Command{
	Use:     "search [tag]",
	Short:   "Search for tags in your vault",
	Long: `Search for notes containing specific tags with flexible output formats.

Formats:
  names      - Just filenames (one per line, deduplicated)
  locations  - Filename with line number (default)
  lines      - Full line content containing the tag
  blocks     - Paragraph/block until empty line

Examples:
  obsidian-cli search programming
  obsidian-cli search programming --format=names
  obsidian-cli search daily --format=blocks -c 2
  obsidian-cli search work -s -l frontmatter`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		vault := obsidian.Vault{Name: vaultName}
		note := obsidian.Note{}
		uri := obsidian.Uri{}

		tag := args[0]
		if len(tag) > 0 && tag[0] == '#' {
			tag = tag[1:]
		}

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatalf("failed to retrieve 'format' flag: %v", err)
		}

		location, err := cmd.Flags().GetString("location")
		if err != nil {
			log.Fatalf("failed to retrieve 'location' flag: %v", err)
		}

		includeSubtags, err := cmd.Flags().GetBool("subtags")
		if err != nil {
			log.Fatalf("failed to retrieve 'subtags' flag: %v", err)
		}

		contextLines, err := cmd.Flags().GetInt("context")
		if err != nil {
			log.Fatalf("failed to retrieve 'context' flag: %v", err)
		}

		useEditor, err := cmd.Flags().GetBool("editor")
		if err != nil {
			log.Fatalf("failed to retrieve 'editor' flag: %v", err)
		}

		err = actions.Search(&vault, &note, &uri, tag, format, location, includeSubtags, contextLines, useEditor)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	searchCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	searchCmd.Flags().StringVarP(&searchFormat, "format", "f", "locations", "output format: names, locations, lines, blocks")
	searchCmd.Flags().StringVarP(&searchLocation, "location", "l", "all", "tag location: frontmatter, inline, or all")
	searchCmd.Flags().BoolP("subtags", "s", false, "include subtags (e.g., searching 'programming' also finds 'programming/philosophy')")
	searchCmd.Flags().IntVarP(&contextLines, "context", "c", 0, "context lines around match (for lines/blocks format)")
	searchCmd.Flags().BoolP("editor", "e", false, "open in editor instead of Obsidian")
	rootCmd.AddCommand(searchCmd)
}