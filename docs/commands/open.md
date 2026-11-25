# Open & Print

## open

Opens a note in Obsidian or your editor. The note name can be a bare name or a path relative to the vault root.

```bash
# Open in Obsidian
notesmd-cli open "{note-name}"

# Open in a specific vault
notesmd-cli open "{note-name}" --vault "{vault-name}"

# Open at a specific heading (case-sensitive)
notesmd-cli open "{note-name}" --section "{heading-text}"

# Open in your editor
notesmd-cli open "{note-name}" --editor
```

## print

Prints the contents of a note to stdout.

```bash
# Print a note in the default vault
notesmd-cli print "{note-name}"

# Print by path
notesmd-cli print "{note-path}"

# Print from a specific vault
notesmd-cli print "{note-name}" --vault "{vault-name}"
```
