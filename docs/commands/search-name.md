# search

Fuzzy search by filename across the vault. Opens an interactive picker; selecting a note opens it in Obsidian (or your editor).

Excluded notes (configured via `Settings > Files & Links > Excluded Files` in Obsidian) won't appear in results.

```bash
# Search in the default vault
notesmd-cli search

# Search in a specific vault
notesmd-cli search --vault "{vault-name}"

# Open selected note in your editor
notesmd-cli search --editor
```
