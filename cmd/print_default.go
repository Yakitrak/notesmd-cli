package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var printDefaultDeprecatedCmd = &cobra.Command{
	Use:        "print-default",
	Aliases:    []string{"pd"},
	Short:      "prints default vault name and path (deprecated: use list-vaults --default)",
	Args:       cobra.ExactArgs(0),
	Deprecated: "use list-vaults --default instead",
	Run: func(cmd *cobra.Command, args []string) {
		pathOnly, _ := cmd.Flags().GetBool("path-only")

		vault := obsidian.Vault{}
		name, err := vault.DefaultName()
		if err != nil {
			log.Fatal(err)
		}
		path, err := vault.Path()
		if err != nil {
			log.Fatal(err)
		}

		if pathOnly {
			fmt.Print(path)
			return
		}

		openType, _ := vault.DefaultOpenType()

		fmt.Fprintln(os.Stdout, "Default vault name:", name)
		fmt.Fprintln(os.Stdout, "Default vault path:", path)
		fmt.Fprintln(os.Stdout, "Default open type:", openType)
	},
}

func init() {
	printDefaultDeprecatedCmd.Flags().Bool("path-only", false, "print only the vault path")
	rootCmd.AddCommand(printDefaultDeprecatedCmd)
}
