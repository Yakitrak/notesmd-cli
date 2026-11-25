# Vault Management

## add-vault

Registers a directory as an Obsidian vault. Creates the Obsidian config file (`~/.config/obsidian/obsidian.json`) if it does not exist. Alias: `av`

::: warning
Fully quit Obsidian before adding or removing vaults (macOS: Cmd+Q). A running process will overwrite external registry changes on exit.
:::

```bash
# Register a vault
notesmd-cli add-vault /path/to/vault

# Register and set as default
notesmd-cli add-vault /path/to/vault --set-default
```

## remove-vault

Removes a vault from the Obsidian config. Does not delete any files on disk. If the removed vault was the default, the default is cleared. Alias: `rv`

```bash
# Remove by vault name
notesmd-cli remove-vault "{vault-name}"

# Remove by vault path
notesmd-cli remove-vault /path/to/vault
```

## list-vaults

Lists all registered Obsidian vaults. The default vault is marked with `(default)`. Alias: `lv`

```bash
# List all vaults
notesmd-cli list-vaults

# Output as JSON
notesmd-cli list-vaults --json

# Output only vault paths (useful for scripting)
notesmd-cli list-vaults --path-only

# Show only the default vault
notesmd-cli list-vaults --default

# Get just the default vault path
notesmd-cli list-vaults --default --path-only
```

**Shell tip:** add this to `~/.zshrc` to `cd` into your default vault instantly:

```bash
obs_cd() {
  local result=$(notesmd-cli list-vaults --default --path-only)
  [ -n "$result" ] && cd -- "$result"
}
```

## set-default-vault

Defines the default vault and/or open type for all future commands.

```bash
# Set default vault by name or path
notesmd-cli set-default-vault "{vault-name}"

# Set default open type: 'obsidian' (default) or 'editor'
notesmd-cli set-default-vault --open-type editor

# Set both at once
notesmd-cli set-default-vault "{vault-name}" --open-type editor
```
