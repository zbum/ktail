# ktail Makefile

# Variables
BINARY_NAME=ktail
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(shell go version | awk '{print $$3}')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-386 .
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm .

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=386 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-386.exe .
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-arm64.exe .



# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install .

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Run with help flag
.PHONY: help-run
help-run: build
	@echo "Running $(BINARY_NAME) with help..."
	./$(BINARY_NAME) --help

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Run vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Check for security vulnerabilities
.PHONY: security
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found, installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Generate dependency graph
.PHONY: deps
deps:
	@echo "Generating dependency graph..."
	@if command -v go-mod-graph >/dev/null 2>&1; then \
		go-mod-graph | dot -Tpng -o deps.png; \
	else \
		echo "go-mod-graph not found, installing..."; \
		go install github.com/PaulXu-cn/go-mod-graph@latest; \
		go-mod-graph | dot -Tpng -o deps.png; \
	fi

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify

# Create release package
.PHONY: release
release: clean build-all
	@echo "Creating release package..."
	@cd dist && \
	for binary in $(BINARY_NAME)-*; do \
		if [[ $$binary == *".exe" ]]; then \
			zip $$binary.zip $$binary; \
		else \
			tar -czf $$binary.tar.gz $$binary; \
		fi; \
	done
	@echo "Release packages created in dist/"
	@echo "Total binaries created: $$(ls dist/$(BINARY_NAME)-* | wc -l)"

# Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

# List all supported platforms and architectures
.PHONY: list-platforms
list-platforms:
	@echo "Supported platforms and architectures:"
	@echo "  Linux:     amd64, arm64, 386, arm"
	@echo "  macOS:     amd64, arm64"
	@echo "  Windows:   amd64, 386, arm64"

# Build only for specific architecture (usage: make build-arch GOOS=linux GOARCH=amd64)
.PHONY: build-arch
build-arch:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "Usage: make build-arch GOOS=<os> GOARCH=<arch>"; \
		echo "Example: make build-arch GOOS=linux GOARCH=amd64"; \
		exit 1; \
	fi
	@mkdir -p dist
	@if [ "$(GOOS)" = "windows" ]; then \
		go build $(LDFLAGS) -o dist/$(BINARY_NAME)-$(GOOS)-$(GOARCH).exe .; \
	else \
		go build $(LDFLAGS) -o dist/$(BINARY_NAME)-$(GOOS)-$(GOARCH) .; \
	fi

# Development setup
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	go mod download
	go mod verify
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi

# Quick development build and run
.PHONY: dev
dev: fmt vet build run

# Create a new git tag for release
.PHONY: tag
tag:
	@echo "Creating new tag..."
	@read -p "Enter tag version (e.g., v1.0.0): " version; \
	if [ -z "$$version" ]; then \
		echo "Error: Version cannot be empty"; \
		exit 1; \
	fi; \
	git tag -a $$version -m "Release $$version"; \
	echo "Tag $$version created. Push with: git push origin $$version"

# Push tag to trigger GitHub release
.PHONY: release-tag
release-tag:
	@echo "Pushing tag to trigger GitHub release..."
	@git push origin --tags

# Create and push release tag
.PHONY: release-github
release-github: tag release-tag
	@echo "Release tag created and pushed to GitHub"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for all platforms and architectures"
	@echo "  build-linux   - Build for Linux (amd64, arm64, 386, arm)"
	@echo "  build-darwin  - Build for macOS (amd64, arm64)"
	@echo "  build-windows - Build for Windows (amd64, 386, arm64)"
	@echo "  build-arch    - Build for specific architecture (GOOS=linux GOARCH=amd64)"
	@echo "  list-platforms - List all supported platforms and architectures"
	@echo "  tag           - Create a new git tag for release"
	@echo "  release-tag   - Push tag to trigger GitHub release"
	@echo "  release-github - Create and push release tag (combines tag + release-tag)"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install the binary to GOPATH/bin"
	@echo "  run           - Build and run the application"
	@echo "  help-run      - Build and run with --help flag"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  vet           - Run go vet"
	@echo "  security      - Check for security vulnerabilities"
	@echo "  deps          - Generate dependency graph"
	@echo "  tidy          - Tidy dependencies"
	@echo "  release       - Create release packages"
	@echo "  version       - Show version information"
	@echo "  dev-setup     - Set up development environment"
	@echo "  dev           - Quick development build and run"
	@echo "  help          - Show this help message"
