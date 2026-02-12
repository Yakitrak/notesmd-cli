package cmd

import (
	"fmt"
	"log"

	"github.com/Yakitrak/notesmd-cli/pkg/actions"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [path]",
	Aliases: []string{"ls"},
	Short:   "List files and folders in vault",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		}

		vault := obsidian.Vault{Name: vaultName}
		entries, err := actions.ListEntries(&vault, actions.ListParams{Path: targetPath})
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			fmt.Printf("• %s\n", entry)
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "vault name")
	rootCmd.AddCommand(listCmd)
}
