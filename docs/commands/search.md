# Search

## search

Fuzzy search by filename across the vault. Opens an interactive picker; selecting a note opens it in Obsidian (or your editor).

```bash
# Search in the default vault
notesmd-cli search

# Search in a specific vault
notesmd-cli search --vault "{vault-name}"

# Open selected note in your editor
notesmd-cli search --editor
```

## search-content

Searches note content for a term. By default opens an interactive picker. Use `--no-interactive` or `--format json` for scripting.

```bash
# Interactive content search
notesmd-cli search-content "search term"

# Search in a specific vault
notesmd-cli search-content "search term" --vault "{vault-name}"

# Open selected note in your editor
notesmd-cli search-content "search term" --editor

# Non-interactive grep-style output
notesmd-cli search-content "search term" --no-interactive

# JSON output (implies non-interactive)
notesmd-cli search-content "search term" --format json

# Paginated JSON (default page size: 25, max: 100)
notesmd-cli search-content "search term" --format json --page 1 --page-size 50
```

## Excluded files

The CLI respects Obsidian's **Excluded Files** setting (`Settings > Files & Links > Excluded Files`).

- `search` — excluded notes won't appear in the picker.
- `search-content` — excluded folders won't be searched.

All other commands (`open`, `move`, `print`, etc.) still access excluded files by name.
