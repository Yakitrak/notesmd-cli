# search-content

Searches note content for a term. By default opens an interactive picker. Use `--no-interactive` or `--format json` for scripting.

Excluded folders (configured via `Settings > Files & Links > Excluded Files` in Obsidian) won't be searched.

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
