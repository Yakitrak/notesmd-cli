# Migration Guide: Obsidian CLI → NotesMD CLI

This guide will help you migrate from `obsidian-cli` (v0.2.3 and earlier) to `notesmd-cli` (v0.3.0+).

## Why the Rename?

With the release of the Official Obsidian CLI, this project has been renamed from "Obsidian CLI" to "NotesMD CLI" to avoid confusion.

**Key difference**: NotesMD CLI works **without requiring Obsidian to be running**, making it perfect for scripting, automation, and terminal-only environments.

## What's Changed?

- **Binary name**: `obsidian-cli` → `notesmd-cli`
- **Package names**: Updated in Homebrew, Scoop, and AUR
- **Config directory**: `~/.config/obsidian-cli` → `~/.config/notesmd-cli`
- **Go module path**: `github.com/Yakitrak/obsidian-cli` → `github.com/Yakitrak/notesmd-cli`

## Migration Steps

### 1. Uninstall the Old Version

Choose the method that matches how you installed `obsidian-cli`:

#### Windows (Scoop)

```bash
scoop uninstall obsidian-cli
```

#### Mac/Linux (Homebrew)

```bash
brew uninstall obsidian-cli
```

#### Arch Linux (AUR)

```bash
# Using yay
yay -R obsidian-cli-bin

# Using paru
paru -R obsidian-cli-bin
```

#### Manual Installation

If you installed manually, remove the binary:

```bash
# Find where it's installed
which obsidian-cli

# Remove it (common locations)
sudo rm /usr/local/bin/obsidian-cli
# or
rm ~/bin/obsidian-cli
```

### 2. Install NotesMD CLI

Follow the installation instructions in the [main README](README.md#install):

#### Windows (Scoop)

```bash
scoop bucket add scoop-yakitrak https://github.com/yakitrak/scoop-yakitrak.git
scoop install notesmd-cli
```

#### Mac/Linux (Homebrew)

```bash
brew tap yakitrak/yakitrak
brew install yakitrak/yakitrak/notesmd-cli
```

#### Arch Linux (AUR)

```bash
# Using yay
yay -S notesmd-cli-bin

# Using paru
paru -S notesmd-cli-bin
```

### 3. Migrate Your Configuration

Your default vault settings and preferences are stored in a config file. You have two options:

#### Option A: Copy the Old Config (Preserves All Settings)

**Linux/Mac:**

```bash
# Create new config directory
mkdir -p ~/.config/notesmd-cli

# Copy your old config
cp ~/.config/obsidian-cli/preferences.json ~/.config/notesmd-cli/preferences.json
```

**Mac (Alternative Location):**
If your config was in your home directory:

```bash
# Create new config directory
mkdir -p ~/.notesmd-cli

# Copy your old config
cp ~/.obsidian-cli/preferences.json ~/.notesmd-cli/preferences.json
```

**Windows:**

```powershell
# Create new config directory
New-Item -ItemType Directory -Force -Path "$env:APPDATA\notesmd-cli"

# Copy your old config
Copy-Item "$env:APPDATA\obsidian-cli\preferences.json" "$env:APPDATA\notesmd-cli\preferences.json"
```

#### Option B: Set Your Default Vault Again (Fresh Start)

If you only had a default vault configured, simply set it again:

```bash
notesmd-cli set-default "{vault-name}"
```

### 4. Verify the Installation

Check that `notesmd-cli` is working:

```bash
# Check version
notesmd-cli --version

# Verify your default vault
notesmd-cli print-default
```

### 5. Update Your Scripts (If Applicable)

If you have any shell scripts, aliases, or automation that references `obsidian-cli`, update them to use `notesmd-cli`:

```bash
# Example: Update a shell alias in ~/.zshrc or ~/.bashrc
# Old: alias obs='obsidian-cli'
# New: alias obs='notesmd-cli'
```

### 6. Clean Up (Optional)

After confirming everything works, you can remove the old config:

**Linux/Mac:**

```bash
rm -rf ~/.config/obsidian-cli
# or
rm -rf ~/.obsidian-cli
```

**Windows:**

```powershell
Remove-Item -Recurse -Force "$env:APPDATA\obsidian-cli"
```

## Troubleshooting

### "Command not found: notesmd-cli"

Make sure you've installed `notesmd-cli` and that it's in your PATH. Try:

```bash
# Reload your shell
exec $SHELL

# Or check if it's installed but not in PATH
which notesmd-cli
```

### Config Not Working

If your settings aren't being recognized:

1. Check the config file exists: `notesmd-cli print-default`
2. Verify the config directory location for your OS
3. Try setting your default vault again: `notesmd-cli set-default "{vault-name}"`

### Need Help?

Open an issue on [GitHub](https://github.com/Yakitrak/notesmd-cli/issues) if you encounter any problems during migration.

## For Developers

If you've been using `obsidian-cli` in your Go projects:

### Update Go Module Imports

```bash
# Update go.mod
go get github.com/Yakitrak/notesmd-cli@latest

# Update imports in your code
# Old: import "github.com/Yakitrak/obsidian-cli/pkg/..."
# New: import "github.com/Yakitrak/notesmd-cli/pkg/..."
```

---

**Welcome to NotesMD CLI!** 🎉 If you have any questions or feedback about the migration, please [open an issue](https://github.com/Yakitrak/notesmd-cli/issues/new/choose).
