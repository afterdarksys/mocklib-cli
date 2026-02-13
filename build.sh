#!/bin/bash
# MockLib CLI Build Script
# Builds binaries for macOS, Linux, and Windows

set -e  # Exit on error

VERSION="${VERSION:-0.1.0}"
BUILD_DIR="./build"
DIST_DIR="./dist"

echo "================================================"
echo "  MockLib CLI Build Script"
echo "  Version: $VERSION"
echo "================================================"
echo ""

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf "$BUILD_DIR" "$DIST_DIR"
mkdir -p "$BUILD_DIR" "$DIST_DIR"
echo "   ✓ Cleaned"
echo ""

# Get dependencies
echo "📦 Getting dependencies..."
go mod download
echo "   ✓ Dependencies ready"
echo ""

# Build for each platform
echo "🔨 Building binaries..."
echo ""

# macOS (Intel)
echo "  Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/mocklib-darwin-amd64"
echo "     ✓ mocklib-darwin-amd64"

# macOS (Apple Silicon)
echo "  Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$BUILD_DIR/mocklib-darwin-arm64"
echo "     ✓ mocklib-darwin-arm64"

# Linux (amd64)
echo "  Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/mocklib-linux-amd64"
echo "     ✓ mocklib-linux-amd64"

# Linux (arm64)
echo "  Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o "$BUILD_DIR/mocklib-linux-arm64"
echo "     ✓ mocklib-linux-arm64"

# Windows (amd64)
echo "  Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/mocklib-windows-amd64.exe"
echo "     ✓ mocklib-windows-amd64.exe"

echo ""
echo "✅ All binaries built successfully!"
echo ""

# Create checksums
echo "🔐 Creating checksums..."
cd "$BUILD_DIR"
shasum -a 256 * > checksums.txt
cd ..
echo "   ✓ Checksums created"
echo ""

# Create archives for release
echo "📦 Creating release archives..."
echo ""

for binary in "$BUILD_DIR"/mocklib-*; do
    if [ -f "$binary" ]; then
        basename=$(basename "$binary")
        name="${basename%.*}"  # Remove extension if present

        if [[ "$binary" == *.exe ]]; then
            # Windows - use zip
            echo "  Creating $name.zip..."
            (cd "$BUILD_DIR" && zip -q "../$DIST_DIR/$name.zip" "$basename")
        else
            # Unix - use tar.gz
            echo "  Creating $name.tar.gz..."
            tar -czf "$DIST_DIR/$name.tar.gz" -C "$BUILD_DIR" "$basename"
        fi

        echo "     ✓ $name archive created"
    fi
done

# Copy checksums to dist
cp "$BUILD_DIR/checksums.txt" "$DIST_DIR/"

echo ""
echo "================================================"
echo "  Build Complete! 🎉"
echo "================================================"
echo ""
echo "Binaries location: $BUILD_DIR/"
echo "Release archives:  $DIST_DIR/"
echo ""
echo "Built binaries:"
ls -lh "$BUILD_DIR"/ | grep mocklib | awk '{print "  " $9, "(" $5 ")"}'
echo ""
echo "Release archives:"
ls -lh "$DIST_DIR"/ | grep -E '\.(tar\.gz|zip)$' | awk '{print "  " $9, "(" $5 ")"}'
echo ""
echo "To test locally:"
echo "  ./build/mocklib-darwin-arm64 --version"
echo "  ./build/mocklib-darwin-amd64 --version"
echo ""
echo "To create GitHub release:"
echo "  gh release create v$VERSION ./dist/* --title \"v$VERSION\" --notes \"Release notes here\""
echo ""
