# Frontmatter

View and modify YAML frontmatter in notes. Alias: `fm`

## Examples

```bash
# Print frontmatter of a note
notesmd-cli frontmatter "{note-name}" --print

# Edit a field (creates it if it doesn't exist)
notesmd-cli frontmatter "{note-name}" --edit --key "status" --value "done"

# Delete a field
notesmd-cli frontmatter "{note-name}" --delete --key "draft"

# Use with a specific vault
notesmd-cli frontmatter "{note-name}" --print --vault "{vault-name}"
```
