# Quick Start

## 1. Set your vault

Before running any other command, point the CLI at your Obsidian vault.

```bash
# Use the interactive picker to select a vault Obsidian already knows about
notesmd-cli set-default-vault --default

# Or specify by name
notesmd-cli set-default-vault "my-vault"
```

## 2. Create a note

```bash
notesmd-cli create "Projects/New Idea"
# Created: Projects/New Idea.md
```

Intermediate folders are created automatically. Notes are placed in your configured default new-file folder when no path is given.

## 3. Open a note

```bash
# Open in Obsidian
notesmd-cli open "Projects/New Idea"

# Open in your terminal editor ($EDITOR)
notesmd-cli open "Projects/New Idea" --editor
```

## 4. Search

```bash
# Fuzzy search by filename
notesmd-cli search

# Search by content
notesmd-cli search-content "meeting notes"
```

## 5. Move and rename

```bash
notesmd-cli move "Projects/New Idea" "Archive/Old Idea"
# All wikilinks and markdown links pointing at the note are updated automatically.
```

## 6. Daily note

```bash
notesmd-cli daily
# Creates today's note using your .obsidian/daily-notes.json config and template.
```
