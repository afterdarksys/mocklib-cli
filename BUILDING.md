# Building MockLib CLI

Super simple build instructions for people who aren't that smart (your words, not mine! 😄)

## TL;DR - Just Build It

```bash
./build.sh
```

That's it. Done.

## What You Get

After running `./build.sh`, you'll have:

```
build/
  ├── mocklib-darwin-amd64       (macOS Intel)
  ├── mocklib-darwin-arm64       (macOS Apple Silicon)
  ├── mocklib-linux-amd64        (Linux Intel/AMD)
  ├── mocklib-linux-arm64        (Linux ARM)
  ├── mocklib-windows-amd64.exe  (Windows)
  └── checksums.txt

dist/
  ├── mocklib-darwin-amd64.tar.gz
  ├── mocklib-darwin-arm64.tar.gz
  ├── mocklib-linux-amd64.tar.gz
  ├── mocklib-linux-arm64.tar.gz
  ├── mocklib-windows-amd64.zip
  └── checksums.txt
```

## Even Simpler (Using Make)

```bash
make
```

Or:

```bash
make build          # Build all platforms
make clean          # Clean build artifacts
make release        # Create GitHub release
```

## Testing Your Build

### macOS (Apple Silicon)
```bash
./build/mocklib-darwin-arm64
```

### macOS (Intel)
```bash
./build/mocklib-darwin-amd64
```

### Linux
```bash
./build/mocklib-linux-amd64
```

### Windows
```bash
./build/mocklib-windows-amd64.exe
```

## Building for Specific Platform

```bash
# macOS only
GOOS=darwin GOARCH=arm64 go build -o mocklib

# Linux only
GOOS=linux GOARCH=amd64 go build -o mocklib

# Windows only
GOOS=windows GOARCH=amd64 go build -o mocklib.exe
```

## Creating a GitHub Release

After building:

```bash
# Set version
export VERSION=0.1.0

# Build
./build.sh

# Create release
gh release create v$VERSION ./dist/* \
  --title "v$VERSION" \
  --notes "Release notes here"
```

Or just:

```bash
make release VERSION=0.1.0
```

## Requirements

- Go 1.21 or later
- That's it!

## Troubleshooting

### "Permission denied"

```bash
chmod +x build.sh
./build.sh
```

### "go: command not found"

Install Go: https://go.dev/dl/

### "I'm too dumb for this"

No you're not! Just run:

```bash
./build.sh
```

If that doesn't work, copy/paste the error message to ChatGPT or Claude. They'll help.

## Distribution

Upload the files in `dist/` to:
- GitHub Releases
- Your website
- S3 bucket
- Whatever

People download, extract, and run. No dependencies needed.

## Installation for End Users

### macOS/Linux

```bash
# Download
curl -LO https://github.com/afterdarksys/mocklib-cli/releases/latest/download/mocklib-darwin-arm64.tar.gz

# Extract
tar -xzf mocklib-darwin-arm64.tar.gz

# Make executable
chmod +x mocklib-darwin-arm64

# Move to PATH
sudo mv mocklib-darwin-arm64 /usr/local/bin/mocklib

# Test
mocklib --help
```

### Windows

1. Download `mocklib-windows-amd64.zip`
2. Extract it
3. Add to PATH or run from current directory

## Quick Install Script (for users)

```bash
curl -sSL https://get.mockfactory.io/install.sh | bash
```

(We should create this install script later)

## That's It!

Seriously, just run `./build.sh` and you're done. 🚀
