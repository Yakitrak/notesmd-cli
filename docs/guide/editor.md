# Editor Flag

The `open`, `daily`, `search`, `search-content`, `create`, and `move` commands all support `--editor` (short: `-e`), which opens notes in your terminal editor instead of Obsidian.

## Setup

The editor is determined by the `EDITOR` environment variable. If it is not set, `vim` is used as a fallback.

```bash
# Add to ~/.zshrc or ~/.bashrc to make it permanent
export EDITOR="code"   # VSCode
export EDITOR="vim"
export EDITOR="nano"
export EDITOR="subl"   # Sublime Text
```

For GUI editors (VSCode, Sublime Text, Atom, TextMate), the CLI automatically appends `--wait` so the command blocks until you close the file.

## Usage

```bash
notesmd-cli open "note" --editor
notesmd-cli daily --editor
notesmd-cli search --editor
notesmd-cli search-content "term" --editor
notesmd-cli create "note" --open --editor
notesmd-cli move "old" "new" --open --editor
```

## Set editor as default

To avoid passing `--editor` on every command, configure it as the default open type once:

```bash
notesmd-cli set-default-vault --open-type editor
```

After this, any command that supports `--open` will use your editor automatically.
