---
name: obsidian-notesmd-cli
description: Interact with Obsidian vaults from the terminal to inspect, search, create, move, print, and manage markdown notes without needing the Obsidian app to be open.
---

# Obsidian NotesMD CLI

Use this skill when the user wants to work with Obsidian vaults through `notesmd-cli` instead of editing note files manually. This includes listing vaults, selecting a default vault, opening notes, creating notes, searching note names or content, printing note contents, editing frontmatter, moving notes, and daily-note workflows.

## When to use

Trigger this skill for requests like:

- "Find a note in my vault"
- "Create or rename an Obsidian note from the terminal"
- "Search note content without opening Obsidian"
- "Set or inspect the default vault"
- "Edit frontmatter on a markdown note"
- "Use NotesMD CLI on a headless machine"

Do not use this skill when the task is purely about editing repository source code unless the user is specifically asking about the CLI's note-management behavior.

## Workflow

1. Confirm whether the user wants to operate on the default vault or a specific vault.
2. Prefer read-only inspection commands first: `list-vaults`, `print-default`, `list`, `search-content --no-interactive`, or `print`.
3. For write operations, use the narrowest command that matches the task: `create`, `frontmatter`, `move`, `delete`, or `daily`.
4. If the environment is terminal-only or Obsidian is not installed, prefer editor or non-interactive modes and rely on the manual config described below.
5. After write operations, verify the result with `print`, `list`, `search-content --no-interactive`, or a targeted follow-up command.

## Command guide

Use the checked-out project binary if needed:

```bash
go run . --help
```

Common commands:

```bash
notesmd-cli list-vaults
notesmd-cli list-vaults --json
notesmd-cli print-default
notesmd-cli set-default "Vault Name"
notesmd-cli set-default --open-type editor
```

Inspect notes and folders:

```bash
notesmd-cli list
notesmd-cli list "Projects"
notesmd-cli print "Inbox/Idea.md"
notesmd-cli print "Inbox/Idea.md" --mentions
notesmd-cli search-content "search term" --no-interactive
notesmd-cli search-content "search term" --no-interactive --format json
```

Open or create notes:

```bash
notesmd-cli open "Inbox/Idea.md"
notesmd-cli open "Inbox/Idea.md" --section "Next Steps"
notesmd-cli open "Inbox/Idea.md" --editor
notesmd-cli daily
notesmd-cli create "Inbox/New Note.md" --content "Initial content"
notesmd-cli create "Inbox/New Note.md" --content "Extra line" --append
notesmd-cli create "Inbox/New Note.md" --open --editor
```

Modify notes:

```bash
notesmd-cli frontmatter "Inbox/Idea.md" --print
notesmd-cli frontmatter "Inbox/Idea.md" --edit --key "status" --value "done"
notesmd-cli frontmatter "Inbox/Idea.md" --delete --key "draft"
notesmd-cli move "Inbox/Idea.md" "Projects/Idea.md"
notesmd-cli delete "Inbox/Old Note.md"
```

## Vault selection

- Many commands accept `--vault "<vault-name>"`.
- If no vault is passed, the CLI uses the configured default vault.
- Use `notesmd-cli print-default` before acting if the target vault is unclear.
- Use `notesmd-cli list-vaults` to discover valid vault names instead of guessing.

## Headless setup

If Obsidian is not installed or there is no GUI, the CLI can still work if the config exists at `~/.config/obsidian/obsidian.json`.

Minimal example:

```json
{
  "vaults": {
    "vault-1": {
      "path": "/absolute/path/to/vault"
    }
  }
}
```

Important constraints:

- Use an absolute path.
- The effective vault name is the directory name of the vault path.
- Prefer `--editor` or non-interactive commands in terminal-only environments.
- `daily` reads `.obsidian/daily-notes.json` when present and falls back to defaults otherwise.

## Safety

- Prefer `search-content --no-interactive` for scripts and automation.
- Verify the destination path before using `move` or `delete`.
- Use `create --append` or `frontmatter --edit` for targeted changes instead of overwriting full note contents.
- Do not assume Obsidian is running; this CLI is designed to work directly on disk.

## Verification

After modifications, verify with one of:

```bash
notesmd-cli print "path/to/note.md"
notesmd-cli frontmatter "path/to/note.md" --print
notesmd-cli list "folder"
notesmd-cli search-content "unique text" --no-interactive
```
