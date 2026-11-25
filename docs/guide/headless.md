# Headless / No Obsidian

notesmd-cli works on servers, containers, and CI environments where Obsidian is not installed. When Obsidian is installed it registers vaults automatically; on headless systems you register them via the CLI instead.

## Register a vault

```bash
# Register a vault directory
notesmd-cli add-vault /home/user/vaults/my-brain

# Register and immediately set as default
notesmd-cli add-vault /home/user/vaults/my-brain --set-default
```

## Multiple vaults

```bash
notesmd-cli add-vault /home/user/vaults/personal
notesmd-cli add-vault /home/user/vaults/work
notesmd-cli set-default-vault "personal"

# Target a non-default vault with any command
notesmd-cli create "note" --vault "work"
```

## Manual setup

If you prefer to configure the registry file directly:

1. Create the Obsidian config directory:
   ```bash
   mkdir -p ~/.config/obsidian
   ```

2. Create `~/.config/obsidian/obsidian.json`:
   ```json
   {
     "vaults": {
       "any-unique-id": {
         "path": "/home/user/vaults/my-brain"
       }
     }
   }
   ```

   - The key (`any-unique-id`) can be any string.
   - The CLI uses the **directory name** as the vault name (e.g. `my-brain`).
   - Use an **absolute path** — `~` is not expanded.

::: warning
Fully quit Obsidian before modifying the vault registry. A running Obsidian process will overwrite external changes on exit.
:::
