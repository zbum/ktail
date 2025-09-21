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
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
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
	mkdir -p dist
	@for binary in $(BINARY_NAME)-*; do \
		if [[ $$binary == *".exe" ]]; then \
			zip dist/$$binary.zip $$binary; \
		else \
			tar -czf dist/$$binary.tar.gz $$binary; \
		fi; \
	done
	@echo "Release packages created in dist/"

# Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

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

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for all platforms (Linux, macOS, Windows)"
	@echo "  build-linux   - Build for Linux (amd64, arm64)"
	@echo "  build-darwin  - Build for macOS (amd64, arm64)"
	@echo "  build-windows - Build for Windows (amd64)"
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
