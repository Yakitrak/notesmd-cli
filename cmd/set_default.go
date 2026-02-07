package cmd

import (
	"fmt"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
	"log"
)

var setDefaultCmd = &cobra.Command{
	Use:     "set-default",
	Aliases: []string{"sd"},
	Short:   "Sets default vault",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		v := obsidian.Vault{}
		err := v.SetDefaultName(input)
		if err != nil {
			// Check if it's a warning about multiple vaults with same name
			if obsidian.IsMultipleVaultsWarning(err) {
				fmt.Println("Warning:", err.Error())
			} else {
				log.Fatal(err)
			}
		}
		path, err := v.Path()
		if err != nil {
			log.Fatal(err)
		}
		
		// Get vault info to display ID
		info, infoErr := v.GetVaultInfo()
		if infoErr == nil {
			fmt.Println("Default vault set to: ", info.Name)
			fmt.Println("Default vault ID:     ", info.ID)
			fmt.Println("Default vault path:   ", path)
		} else {
			// Fallback for backward compatibility
			fmt.Println("Default vault set to: ", v.Name)
			fmt.Println("Default vault path:   ", path)
		}
	},
}

func init() {
	rootCmd.AddCommand(setDefaultCmd)
}
