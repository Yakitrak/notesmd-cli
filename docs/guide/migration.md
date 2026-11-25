# Migrating from obsidian-cli

If you're upgrading from `obsidian-cli` v0.2.3 or earlier, see the full [Migration Guide](https://github.com/Yakitrak/notesmd-cli/blob/main/MIGRATION.md) on GitHub for step-by-step instructions.

## Deprecated commands

The following commands still work but print a deprecation warning to stderr (so existing pipes and scripts are unaffected). They will be removed in the next major version.

| Old command | Replacement |
|---|---|
| `set-default` | `set-default-vault` |
| `print-default` | `list-vaults --default` |
| `print-default --path-only` | `list-vaults --default --path-only` |
