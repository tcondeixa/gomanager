# gomanager

A CLI tool to manage Go installed binaries with tracking and easy updates.

## Features

- 📦 **Install Go packages** with version tracking
- 📋 **List installed packages** with details
- 🔄 **Update packages** to latest versions
- 🗑️ **Uninstall packages** cleanly
- 💾 **Export/Import** package lists
- 🎯 **Custom binary names** for installed tools

## Installation

### Option 1: Install via Homebrew

```bash
brew tap tcondeixa/taps
brew install --cask gomanager
```

This will install the binary with shell completions for Bash, Zsh, and Fish.
See the [homebrew-taps repository](https://github.com/tcondeixa/homebrew-taps) for more details.

### Option 2: Install from source

```bash
go install github.com/tcondeixa/gomanager@latest
```

### Option 3: Download pre-built binaries

Download the latest release for your OS and architecture from the
[github releases](https://github.com/tcondeixa/gomanager/releases).

### Option 4: Build from source

```bash
git clone https://github.com/tcondeixa/gomanager.git
cd gomanager
make
# Binary will be in bin/gomanager
```

## Quick Start

```bash
# Install a Go package
gomanager install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# List installed packages
gomanager list

# Update all packages to latest
gomanager update

# Update specific package
gomanager update -n golangci-lint

# Uninstall a package
gomanager uninstall golangci-lint
```

## Usage

### Install packages

```bash
# Install latest version
gomanager install github.com/user/tool@latest

# Install specific version
gomanager install github.com/user/tool@v1.2.3

# Install with custom binary name
gomanager install github.com/user/tool@latest --name my-tool

# Install multiple packages
gomanager install pkg1@latest pkg2@v1.0.0
```

### List installed packages

```bash
# List in human-readable format
gomanager list

# List in JSON format
gomanager list --output json
```

### Update packages

```bash
# Update all packages with "latest" version
gomanager update

# Update specific package
gomanager update --name tool-name

# Force update all packages (including pinned versions)
gomanager update --force
```

### Uninstall packages

```bash
# Uninstall one package
gomanager uninstall tool-name

# Uninstall multiple packages
gomanager uninstall tool1 tool2 tool3
```

### Export packages

```bash
# Export to file (defaults to ~/gomanager.json)
gomanager export

# Export to custom location
gomanager export --file /path/to/backup.json
```

### Import packages

```bash
# Import from file (defaults to ~/gomanager.json)
gomanager import

# Export to custom location
gomanager import --file /path/to/backup.json
```

## Configuration

gomanager stores its data in the default config directory depending on the OS.
The **config directory** can set by defining the `$gomanager_CONFIG_DIR` environment variable.

## Examples

```bash
# Install popular Go tools
gomanager install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
gomanager install github.com/air-verse/air@latest
gomanager install github.com/goreleaser/goreleaser@latest

# List what's installed
gomanager list
# Output:
# Installed Packages:
# -------------------
# Name: golangci-lint
# URI: github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# Updated: 2024-01-15 10:30:00
# -------------------

# Update everything
gomanager update

# Clean up unused tools
gomanager uninstall air goreleaser
```

## Shell Completion

Enable shell completion for better UX.
The Homebrew installation includes completions automatically.
The following commands will provide instruction for the different OS and Shells.

```bash
gomanager completion --help
```

## Requirements

- Go 1.25.1 or later
- Standard Go toolchain (`go install` command available)

## License

MIT License - see [LICENSE](LICENSE) file for details.
