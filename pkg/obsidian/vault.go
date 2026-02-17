package obsidian

type CliConfig struct {
	DefaultVaultName string `json:"default_vault_name"`
	DefaultOpenType  string `json:"default_open_type,omitempty"`
}

type ObsidianVaultConfig struct {
	Vaults map[string]struct {
		Path string `json:"path"`
	} `json:"vaults"`
}

type VaultManager interface {
	DefaultName() (string, error)
	SetDefaultName(name string) error
	Path() (string, error)
	DefaultOpenType() (string, error)
}

type Vault struct {
	Name string
}
