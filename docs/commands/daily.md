# Daily Notes

Creates or opens today's daily note directly on disk. **Obsidian does not need to be running.**

If `.obsidian/daily-notes.json` exists in the vault, the CLI reads `folder`, `format` (Moment.js date format, default `YYYY-MM-DD`), and `template` from it. Template variables are expanded the same way Obsidian's daily notes plugin does it:

<div v-pre>

| Variable | Resolves to |
|---|---|
| `{{date}}` | Note title (formatted date) |
| `{{title}}` | Note title |
| `{{time}}` | Current time as `HH:mm` |
| `{{date:FORMAT}}` | Date formatted with the given Moment.js format |
| `{{time:FORMAT}}` | Time formatted with the given Moment.js format |

</div>

If the config is missing or unreadable, defaults are used: vault root, `YYYY-MM-DD` format, no template.

## Examples

```bash
# Create or open today's daily note
notesmd-cli daily

# Use a specific vault
notesmd-cli daily --vault "{vault-name}"

# Open in your editor
notesmd-cli daily --editor

# Append content (creates the note if it doesn't exist)
notesmd-cli daily --content "abcde"

# Append content and open in editor
notesmd-cli daily --content "abcde" --editor
```
