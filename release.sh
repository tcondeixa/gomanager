#!/bin/bash

set -e

make mod.ci
make lint.ci
make

VERSION=$(sed -n 's/var version = "\(v[0-9]*\.[0-9]*\.[0-9]*\)"/\1/p' main.go)

git fetch --tags
if [[ $(git tag --list | grep -c "$VERSION") -eq 0 ]]; then
	echo "Tag $VERSION not found. Creating..."
	git tag -a "$VERSION" -m "$VERSION"
	git push origin "$VERSION"
	goreleaser release
	echo "Release $VERSION created successfully!"
else
	echo "Tag $VERSION already exists. Skipping..."
fi
