package cmd

import (
	"fmt"
	"log"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var removeVaultCmd = &cobra.Command{
	Use:     "remove-vault <name>",
	Aliases: []string{"rv"},
	Short:   "Unregister a vault",
	Long:    "Removes a vault from the Obsidian config. Does not delete any files on disk.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := obsidian.RemoveVault(name); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Vault %q removed\n", name)

		if err := obsidian.ClearDefaultIfMatch(name); err != nil {
			fmt.Println("Warning: could not clear default vault:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeVaultCmd)
}
