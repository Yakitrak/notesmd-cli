package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"text/tabwriter"

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

		if listVaultsPathOnly {
			for _, v := range vaults {
				fmt.Println(v.Path)
			}
		} else {
			formatVaultsTable(os.Stdout, vaults)
		}
	},
}

// formatVaultsTable writes vaults as aligned columns using tabwriter,
// so that the path column lines up regardless of vault name length.
//
// Example output:
//
//	Notes          /home/user/Notes
//	LongVaultName  /home/user/LongVaultName
//	Work           /home/user/Work
func formatVaultsTable(w io.Writer, vaults []obsidian.VaultInfo) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, v := range vaults {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", v.Name, v.Path)
	}
	_ = tw.Flush()
}

func init() {
	listVaultsCmd.Flags().BoolVar(&listVaultsJSON, "json", false, "output as JSON array")
	listVaultsCmd.Flags().BoolVar(&listVaultsPathOnly, "path-only", false, "output one path per line")
	listVaultsCmd.MarkFlagsMutuallyExclusive("json", "path-only")
	rootCmd.AddCommand(listVaultsCmd)
}
