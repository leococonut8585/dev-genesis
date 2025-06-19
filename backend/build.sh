#!/bin/bash
set -e

echo "🔨 Building Dev Genesis..."

# Clean previous builds
rm -rf ../build
mkdir -p ../build

# Build for multiple platforms
echo "📦 Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../build/dev-genesis-windows-amd64.exe cmd/server/main.go

echo "📦 Building for Windows ARM64..."
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o ../build/dev-genesis-windows-arm64.exe cmd/server/main.go

echo "📦 Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../build/dev-genesis-linux-amd64 cmd/server/main.go

echo "📦 Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../build/dev-genesis-darwin-amd64 cmd/server/main.go

echo "📦 Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../build/dev-genesis-darwin-arm64 cmd/server/main.go

# Create checksums
cd ../build
echo "🔐 Generating checksums..."
sha256sum * > checksums.txt

echo "✅ Build complete! Files in build/ directory:"
ls -la

cd ../backend