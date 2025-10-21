#!/bin/bash

set -e

# Check if the current gh profile is tcondeixa
CURRENT_USER=$(gh config get user -h github.com)
if [ "$CURRENT_USER" != "tcondeixa" ]; then
    echo "Error: Current gh profile is '$CURRENT_USER', expected 'tcondeixa'"
    echo "Please change your gh profile to tcondeixa and try again"
    exit 1
fi

make mod.ci
make lint.ci
make

VERSION=$(sed -n 's/var version = "\(v[0-9]*\.[0-9]*\.[0-9]*\)"/\1/p' main.go)

if [ "$VERSION" = "" ]; then
    echo "Error: Could not extract version from main.go"
    exit 1
fi

# Check if tag already exists on GitHub
if gh release view "$VERSION" >/dev/null 2>&1; then
    echo "Release $VERSION already exists. No new version defined."
    exit 0
fi

goreleaser release

echo "Release $VERSION created successfully!"
