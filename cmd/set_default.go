package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var setDefaultCmd = &cobra.Command{
	Use:     "set-default",
	Aliases: []string{"sd"},
	Short:   "Sets default vault and/or open type",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		openType, err := cmd.Flags().GetString("open-type")
		if err != nil {
			log.Fatalf("Failed to parse --open-type flag: %v", err)
		}

		if len(args) == 0 && openType == "" {
			log.Fatal("Please provide a vault name or use --open-type to set the default open type")
		}

		if len(args) > 0 {
			name := args[0]
			v := obsidian.Vault{Name: name}
			if err := v.SetDefaultName(name); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Default vault set to:", name)
			path, err := v.Path()
			if err != nil {
				// Path resolution is best-effort: the name is saved; Obsidian's
				// config file may not be present or may not contain this vault yet.
				fmt.Fprintln(os.Stderr, "Note: could not resolve vault path:", err)
			} else {
				fmt.Println("Default vault path set to:", path)
			}
		}

		if openType != "" {
			if openType != "obsidian" && openType != "editor" {
				log.Fatalf("Invalid open type %q: must be 'obsidian' or 'editor'", openType)
			}
			v := obsidian.Vault{}
			if err := v.SetDefaultOpenType(openType); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Default open type set to:", openType)
		}
	},
}

func init() {
	setDefaultCmd.Flags().String("open-type", "", "default open type: 'obsidian' (default) or 'editor'")
	rootCmd.AddCommand(setDefaultCmd)
}
