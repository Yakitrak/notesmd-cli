package cmd

import (
	"fmt"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/spf13/cobra"
	"log"
)

var setDefaultCmd = &cobra.Command{
	Use:     "set-default",
	Aliases: []string{"sd"},
	Short:   "Sets default vault",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		v := obsidian.Vault{Name: name}
		err := v.SetDefaultName(name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Default vault set to: ", name)
		path, err := v.Path()
		if err != nil {
			// Path resolution is best-effort: the name is saved; Obsidian's
			// config file may not be present or may not contain this vault yet.
			log.Printf("Note: could not resolve vault path (%v)", err)
			return
		}
		fmt.Println("Default vault path set to: ", path)
	},
}

func init() {
	rootCmd.AddCommand(setDefaultCmd)
}
