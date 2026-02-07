package obsidian

type CliConfig struct {
	DefaultVaultName string `json:"default_vault_name"`
	DefaultVaultID   string `json:"default_vault_id,omitempty"`
}

type ObsidianVaultConfig struct {
	Vaults map[string]struct {
		Path string `json:"path"`
	} `json:"vaults"`
}

type VaultInfo struct {
	ID   string
	Name string
	Path string
}

type VaultManager interface {
	DefaultName() (string, error)
	SetDefaultName(name string) error
	Path() (string, error)
	GetVaultInfo() (*VaultInfo, error)
}

type Vault struct {
	Name string
	ID   string
}
