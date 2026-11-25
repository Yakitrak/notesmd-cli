# list

Lists files and folders at a vault path. If no path is given, lists the vault root.

```bash
# List vault root
notesmd-cli list

# List a subfolder
notesmd-cli list "001 Notes"

# List in a specific vault
notesmd-cli list "001 Notes" --vault "{vault-name}"
```
