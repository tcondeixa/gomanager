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

echo "Building release for version: $VERSION"

# Check if tag already exists on GitHub
if gh release view "$VERSION" >/dev/null 2>&1; then
    echo "Release $VERSION already exists. No new version defined."
    exit 0
fi

make build.all

# Generate release notes with commit history
echo "Generating release notes..."
LAST_TAG=$(git tag -l --sort=-version:refname | head -1)

if [ "$LAST_TAG" = "" ]; then
    echo "No previous releases found. Including all commits."
    COMMIT_RANGE="HEAD"
else
    echo "Last release: $LAST_TAG"
    COMMIT_RANGE="$LAST_TAG..HEAD"
fi

# Generate commit list
COMMITS=$(git log "$COMMIT_RANGE" --pretty=format:"- %s (%h)" --reverse)

# Create release notes
RELEASE_NOTES="Release $VERSION

## Changes
$COMMITS"

echo "Creating GitHub release with assets..."
gh release create "$VERSION" \
    --title "Release $VERSION" \
    --notes "$RELEASE_NOTES" \
    gomanager-"$VERSION"-linux-amd64.tar.gz \
    gomanager-"$VERSION"-linux-arm64.tar.gz \
    gomanager-"$VERSION"-darwin-amd64.tar.gz \
    gomanager-"$VERSION"-darwin-arm64.tar.gz

rm -f gomanager-*.tar.gz

echo "Release $VERSION created successfully!"
