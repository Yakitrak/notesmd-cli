package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

var listVaultsJSON bool
var listVaultsPathOnly bool

var listVaultsCmd = &cobra.Command{
	Use:     "list-vaults",
	Aliases: []string{"lv"},
	Short:   "lists all registered Obsidian vaults",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		vaults, err := obsidian.ListVaults()
		if err != nil {
			log.Fatal(err)
		}

		sort.Slice(vaults, func(i, j int) bool {
			return vaults[i].Name < vaults[j].Name
		})

		if listVaultsJSON {
			output, err := json.MarshalIndent(vaults, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(output))
			return
		}

		for _, v := range vaults {
			if listVaultsPathOnly {
				fmt.Println(v.Path)
			} else {
				fmt.Printf("%s\t%s\n", v.Name, v.Path)
			}
		}
	},
}

func init() {
	listVaultsCmd.Flags().BoolVar(&listVaultsJSON, "json", false, "output as JSON array")
	listVaultsCmd.Flags().BoolVar(&listVaultsPathOnly, "path-only", false, "output one path per line")
	rootCmd.AddCommand(listVaultsCmd)
}
