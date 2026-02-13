.PHONY: build clean test release help

VERSION ?= 0.1.0

help:
	@echo "MockLib CLI - Build Commands"
	@echo "============================="
	@echo ""
	@echo "make build     - Build for all platforms"
	@echo "make clean     - Clean build artifacts"
	@echo "make test      - Run tests"
	@echo "make release   - Create GitHub release (VERSION=0.1.0)"
	@echo "make help      - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make release VERSION=1.0.0"
	@echo ""

build:
	@echo "Building MockLib CLI v$(VERSION)..."
	@./build.sh

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/ dist/
	@echo "✓ Clean"

test:
	@echo "Running tests..."
	@go test -v ./...

release: build
	@echo "Creating GitHub release v$(VERSION)..."
	@gh release create v$(VERSION) ./dist/* \
		--title "v$(VERSION)" \
		--notes "MockLib CLI v$(VERSION)" \
		--draft
	@echo "✓ Release draft created at https://github.com/afterdarksys/mocklib-cli/releases"
	@echo ""
	@echo "Edit the release notes and publish when ready!"

install-local: build
	@echo "Installing mocklib locally..."
	@cp build/mocklib-darwin-arm64 /usr/local/bin/mocklib
	@chmod +x /usr/local/bin/mocklib
	@echo "✓ Installed to /usr/local/bin/mocklib"
	@echo ""
	@echo "Test it:"
	@echo "  mocklib --help"
