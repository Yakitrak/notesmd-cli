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
var listVaultsDefault bool

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

		defaultName := resolveDefaultVaultName()

		if listVaultsDefault {
			runListVaultsDefault(vaults, defaultName)
			return
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults registered. Use add-vault to register one.")
			return
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
			formatVaultsTable(os.Stdout, vaults, defaultName)
		}
	},
}

func runListVaultsDefault(vaults []obsidian.VaultInfo, defaultName string) {
	if defaultName == "" {
		fmt.Println("No default vault set. Use set-default-vault to set one.")
		return
	}

	var defaultVault *obsidian.VaultInfo
	for _, v := range vaults {
		if v.Name == defaultName {
			defaultVault = &v
			break
		}
	}

	if defaultVault == nil {
		fmt.Printf("Default vault %q is set but not found in registered vaults.\n", defaultName)
		return
	}

	if listVaultsJSON {
		output, err := json.MarshalIndent(defaultVault, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(output))
		return
	}

	if listVaultsPathOnly {
		fmt.Print(defaultVault.Path)
		return
	}

	vault := obsidian.Vault{Name: defaultName}
	openType, _ := vault.DefaultOpenType()

	fmt.Println("Default vault name:", defaultVault.Name)
	fmt.Println("Default vault path:", defaultVault.Path)
	fmt.Println("Default open type:", openType)
}

// formatVaultsTable writes vaults as aligned columns using tabwriter,
// so that the path column lines up regardless of vault name length.
// The default vault is marked with (default).
//
// Example output:
//
//	Notes          /home/user/Notes  (default)
//	LongVaultName  /home/user/LongVaultName
//	Work           /home/user/Work
func formatVaultsTable(w io.Writer, vaults []obsidian.VaultInfo, defaultName string) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, v := range vaults {
		if v.Name == defaultName {
			_, _ = fmt.Fprintf(tw, "%s\t%s\t(default)\n", v.Name, v.Path)
		} else {
			_, _ = fmt.Fprintf(tw, "%s\t%s\n", v.Name, v.Path)
		}
	}
	_ = tw.Flush()
}

func resolveDefaultVaultName() string {
	vault := obsidian.Vault{}
	name, err := vault.DefaultName()
	if err != nil {
		return ""
	}
	return name
}

func init() {
	listVaultsCmd.Flags().BoolVar(&listVaultsJSON, "json", false, "output as JSON array")
	listVaultsCmd.Flags().BoolVar(&listVaultsPathOnly, "path-only", false, "output one path per line")
	listVaultsCmd.Flags().BoolVar(&listVaultsDefault, "default", false, "show only the default vault")
	listVaultsCmd.MarkFlagsMutuallyExclusive("json", "path-only")
	rootCmd.AddCommand(listVaultsCmd)
}
