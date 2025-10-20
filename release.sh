#!/bin/bash

set -e

make mod.ci
make lint.ci
make

VERSION=$(sed -n 's/var version = "\(v[0-9]*\.[0-9]*\.[0-9]*\)"/\1/p' main.go)

if [ "$VERSION" = "" ]; then
    echo "Error: Could not extract version from main.go"
    exit 1
fi

echo "Building release for version: $VERSION"

# Check if tag already exists on GitHub
if gh release view "$VERSION" >/dev/null 2>&1; then
    echo "Release $VERSION already exists. No new version defined."
    exit 0
fi

make build.all

echo "Creating GitHub release with assets..."
gh release create "$VERSION" \
    --title "Release $VERSION" \
    --notes "Release $VERSION" \
    gomanager-"$VERSION"-linux-amd64.tar.gz \
    gomanager-"$VERSION"-linux-arm64.tar.gz \
    gomanager-"$VERSION"-darwin-amd64.tar.gz \
    gomanager-"$VERSION"-darwin-arm64.tar.gz

rm -f gomanager-*.tar.gz

echo "Release $VERSION created successfully!"
