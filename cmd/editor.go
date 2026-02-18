package cmd

import (
	"log"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

// resolveUseEditor returns whether the command should open in the user's editor.
// If --editor was explicitly passed, its value is used. Otherwise, the configured
// default_open_type is consulted ("editor" → true).
func resolveUseEditor(cmd *cobra.Command, vault *obsidian.Vault) bool {
	useEditor, err := cmd.Flags().GetBool("editor")
	if err != nil {
		log.Fatalf("Failed to parse --editor flag: %v", err)
	}
	if !cmd.Flags().Changed("editor") {
		defaultOpenType, configErr := vault.DefaultOpenType()
		if configErr == nil && defaultOpenType == "editor" {
			useEditor = true
		}
	}
	return useEditor
}
