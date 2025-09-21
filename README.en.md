# ktail

A Kubernetes log tail utility with interactive namespace and pod selection using fuzzy finder.

## Overview

`ktail` is a tool that provides `tail`-like functionality for Kubernetes pod logs. It allows you to interactively select namespaces and pods using a fuzzy finder for a better user experience, making it easy to follow logs from any pod in your cluster.

## Features

- 🎯 **Interactive Selection**: Use fuzzy finder to interactively select namespaces and pods
- 📊 **Pod Status Display**: Visual indicators for pod status (Running, Pending, Failed, etc.)
- 🔄 **Real-time Log Streaming**: Follow logs in real-time with `tail -f` behavior
- 🎛️ **Flexible Options**: Support for custom tail lines, container selection, and more
- 🚀 **Easy to Use**: Simple CLI interface with sensible defaults

## Prerequisites

- Go 1.19 or later
- Access to a Kubernetes cluster
- `kubectl` configured to access your cluster

**Note**: No external dependencies required! The tool uses the `go-fuzzyfinder` library for interactive selection, so you don't need to install `fzf` separately.

## Installation

### Option 1: Download from GitHub Releases (Recommended)

1. Go to the [Releases page](https://github.com/zbum/ktail/releases)
2. Download the appropriate archive for your platform:
   - **Linux**: `ktail-linux.tar.gz` (includes amd64, arm64, 386, arm, ppc64, ppc64le, mips, mipsle, mips64, mips64le, riscv64, s390x)
   - **macOS**: `ktail-darwin.tar.gz` (includes amd64, arm64)
   - **Windows**: `ktail-windows.zip` (includes amd64, 386, arm64)
   - **FreeBSD**: `ktail-freebsd.tar.gz` (includes amd64, 386, arm64, arm)
   - **NetBSD**: `ktail-netbsd.tar.gz` (includes amd64, 386, arm64, arm)
   - **OpenBSD**: `ktail-openbsd.tar.gz` (includes amd64, 386, arm64, arm)
   - **All platforms**: `ktail-all.tar.gz` (includes all binaries)

3. Extract the archive and move the binary to your PATH:
```bash
# For Linux/macOS
tar -xzf ktail-linux.tar.gz
chmod +x ktail-linux-*
sudo mv ktail-linux-* /usr/local/bin/ktail

# For Windows
# Extract ktail-windows.zip and add to PATH
```

4. Verify installation:
```bash
ktail --help
```

### Option 2: Build from Source

1. Clone the repository:
```bash
git clone <repository-url>
cd ktail
```

2. Build the project:
```bash
# Using Makefile (recommended)
make build

# Or manually
go build -o ktail
```

3. Make it executable and move to PATH (optional):
```bash
chmod +x ktail
sudo mv ktail /usr/local/bin/
```

### Option 3: Build for All Platforms

```bash
# Build for all supported platforms and architectures
make build-all

# Build for specific platform
make build-linux
make build-darwin
make build-windows

# Build for specific architecture
make build-arch GOOS=linux GOARCH=arm64
```

## Usage

### Basic Usage

```bash
# Interactive mode - select namespace and pod using fuzzy finder
./ktail

# Specify namespace and pod directly
./ktail -n my-namespace -p my-pod

# Show help
./ktail --help
```

### Command Line Options

```bash
Usage:
  ktail [flags]

Flags:
  -c, --container string    Container name (if not provided, will use the first container)
  -f, --follow             Follow log output (default: true)
  -h, --help               help for ktail
  -n, --namespace string   Kubernetes namespace (if not provided, will be selected interactively)
  -p, --pod string         Pod name (if not provided, will be selected interactively)
  -t, --tail int           Number of lines to show from the end of logs (default: 100)
```

### Examples

```bash
# Follow logs from a specific pod
./ktail -n production -p web-app-7d4f8b9c6-xyz12

# Show last 50 lines and follow
./ktail -n staging -p api-server -t 50

# Follow logs from a specific container
./ktail -n default -p my-pod -c sidecar-container

# Interactive selection with custom tail lines
./ktail -t 200
```

## Development

### Prerequisites
- Go 1.24 or later
- Make (for using Makefile)

### Running tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Or manually
go test ./...
```

### Building for different platforms
```bash
# Build for all supported platforms
make build-all

# Build for specific platforms
make build-linux
make build-darwin
make build-windows

# Build for specific architecture
make build-arch GOOS=linux GOARCH=arm64

# List all supported platforms
make list-platforms
```

### Code Quality
```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Security check
make security

# Run all quality checks
make dev-setup
```

### Creating a Release
```bash
# Create a new release tag and push to GitHub
make release-github

# Or manually
make tag
make release-tag
```

This will:
1. Create a git tag (e.g., v1.0.0)
2. Push the tag to GitHub
3. Trigger GitHub Actions to build binaries for all platforms
4. Create a GitHub Release with all the binaries attached

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.


