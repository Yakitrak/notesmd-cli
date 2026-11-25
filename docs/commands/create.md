# create

Creates a note directly on disk. **Obsidian does not need to be running.** Intermediate directories are created automatically.

When no path separator (`/`) is in the note name, the CLI reads `.obsidian/app.json` to check for a configured default new-file folder. If configured, the note is placed there. If the config is missing, the note is created at the vault root.

If the note already exists and neither `--overwrite` nor `--append` is passed, the file is left unchanged.

```bash
# Create an empty note
notesmd-cli create "{note-name}"

# Create in a specific vault
notesmd-cli create "{note-name}" --vault "{vault-name}"

# Create with content
notesmd-cli create "{note-name}" --content "abcde"

# Overwrite an existing note
notesmd-cli create "{note-name}" --content "abcde" --overwrite

# Append to an existing note
notesmd-cli create "{note-name}" --content "abcde" --append

# Create and open in Obsidian
notesmd-cli create "{note-name}" --content "abcde" --open

# Create and open in your editor
notesmd-cli create "{note-name}" --content "abcde" --open --editor
```
