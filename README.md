# goinstall

A CLI tool to manage Go installed binaries with tracking and easy updates.

## Features

- 📦 **Install Go packages** with version tracking
- 📋 **List installed packages** with details
- 🔄 **Update packages** to latest versions
- 🗑️ **Uninstall packages** cleanly
- 💾 **Export/import** package lists
- 🎯 **Custom binary names** for installed tools

## Installation

### Option 1: Install from source (recommended)

```bash
go install github.com/tcondeixa/goinstall@latest
```

### Option 2: Download pre-built binaries

Download the latest release for your platform from the [releases page](https://github.com/tcondeixa/goinstall/releases).

### Option 3: Build from source

```bash
git clone https://github.com/tcondeixa/goinstall.git
cd goinstall
make build.local
# Binary will be in bin/goinstall
```

## Quick Start

```bash
# Install a Go package
goinstall install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# List installed packages
goinstall list

# Update all packages to latest
goinstall update

# Update specific package
goinstall update -n golangci-lint

# Uninstall a package
goinstall uninstall golangci-lint
```

## Usage

### Install packages

```bash
# Install latest version
goinstall install github.com/user/tool@latest

# Install specific version
goinstall install github.com/user/tool@v1.2.3

# Install with custom binary name
goinstall install github.com/user/tool@latest --name my-tool

# Install multiple packages
goinstall install pkg1@latest pkg2@v1.0.0
```

### List installed packages

```bash
# List in human-readable format
goinstall list

# List in JSON format
goinstall list --output json
```

### Update packages

```bash
# Update all packages with "latest" version
goinstall update

# Update specific package
goinstall update --name tool-name

# Force update all packages (including pinned versions)
goinstall update --force
```

### Uninstall packages

```bash
# Uninstall one package
goinstall uninstall tool-name

# Uninstall multiple packages
goinstall uninstall tool1 tool2 tool3
```

### Dump packages

```bash
# Export to file (defaults to ~/goinstall.json)
goinstall dump

# Export to custom location
goinstall dump --file /path/to/backup.json
```

## Configuration

goinstall stores its data in:

- **Config directory**: `$HOME/.config/goinstall/` (or `$GOINSTALL_CONFIG_DIR`)
- **Storage file**: `storage.json` (tracks installed packages)

### Environment Variables

- `GOINSTALL_CONFIG_DIR`: Override default config directory

## Examples

```bash
# Install popular Go tools
goinstall install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
goinstall install github.com/air-verse/air@latest
goinstall install github.com/goreleaser/goreleaser@latest

# List what's installed
goinstall list
# Output:
# Installed Packages:
# -------------------
# Name: golangci-lint
# URI: github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# Updated: 2024-01-15 10:30:00
# -------------------

# Update everything
goinstall update

# Clean up unused tools
goinstall uninstall air goreleaser
```

## Shell Completion

Enable shell completion for better UX:

```bash
# Bash
goinstall completion bash > /etc/bash_completion.d/goinstall

# Zsh
goinstall completion zsh > "${fpath[1]}/_goinstall"

# Fish
goinstall completion fish > ~/.config/fish/completions/goinstall.fish
```

## Requirements

- Go 1.25.1 or later
- Standard Go toolchain (`go install` command available)

## License

MIT License - see [LICENSE](LICENSE) file for details.
