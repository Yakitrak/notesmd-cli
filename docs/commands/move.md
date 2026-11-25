# move

Moves or renames a note and updates every link in the vault that points to it — wikilinks, markdown links, and percent-encoded paths (e.g. `My%20Note.md`). Intermediate destination directories are created automatically.

::: info Partial failure
Link updates are written to temporary files before replacing originals. If a later update fails, the note may already have moved and some links may already be updated. The error message reports this partial completion.
:::

```bash
# Move or rename a note
notesmd-cli move "{current-note-path}" "{new-note-path}"

# Move in a specific vault
notesmd-cli move "{current-note-path}" "{new-note-path}" --vault "{vault-name}"

# Move and open in Obsidian
notesmd-cli move "{current-note-path}" "{new-note-path}" --open

# Move and open in your editor
notesmd-cli move "{current-note-path}" "{new-note-path}" --open --editor
```
