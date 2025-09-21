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

1. Clone the repository:
```bash
git clone <repository-url>
cd ktail
```

2. Build the project:
```bash
go build -o ktail
```

3. Make it executable and move to PATH (optional):
```bash
chmod +x ktail
sudo mv ktail /usr/local/bin/
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

### Running tests
```bash
go test ./...
```

### Building for different platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o ktail-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o ktail.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o ktail-macos
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.


