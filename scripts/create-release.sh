#!/bin/bash
set -e

VERSION="v1.0.0"
REPO="leococonut8585/dev-genesis"

echo "📝 Creating GitHub Release $VERSION..."

# Create release using GitHub CLI
gh release create $VERSION \
    --repo $REPO \
    --title "Dev Genesis $VERSION - Initial Release" \
    --notes-file ../RELEASE_NOTES.md \
    --draft \
    ../build/DevGenesisInstaller.exe \
    ../build/dev-genesis-windows-amd64.exe \
    ../build/dev-genesis-windows-arm64.exe \
    ../build/dev-genesis-linux-amd64 \
    ../build/dev-genesis-darwin-amd64 \
    ../build/dev-genesis-darwin-arm64 \
    ../build/checksums.txt

echo "✅ Draft release created! Review and publish at:"
echo "https://github.com/$REPO/releases"