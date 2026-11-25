# Installation

## macOS & Linux — Homebrew

```bash
brew tap yakitrak/yakitrak
brew install yakitrak/yakitrak/notesmd-cli
```

## Windows — Scoop

```powershell
scoop bucket add scoop-yakitrak https://github.com/yakitrak/scoop-yakitrak.git
scoop install notesmd-cli
```

## Arch Linux — AUR

Install the official binary with your preferred AUR helper:

```sh
yay -S notesmd-cli-bin
```

Or build from source via the AUR:

```sh
yay -S notesmd-cli
```

## Build from Source

Requires [Go](https://go.dev/dl/) 1.19 or later.

```bash
git clone https://github.com/yakitrak/notesmd-cli.git
cd notesmd-cli
go build -o notesmd-cli .
sudo install -m 755 notesmd-cli /usr/local/bin/
```

## Verify

```bash
notesmd-cli --version
```
