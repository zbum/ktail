# ktail

A simple and powerful tool for real-time Kubernetes pod log tailing.

## Features

- 🎯 **Interactive Selection**: Use fuzzy finder to interactively select namespaces and pods
- 🔄 **Real-time Log Streaming**: Follow logs in real-time with `tail -f` behavior
- 🎨 **Colored Output**: Namespace and pod names are displayed in green for better readability
- 👀 **Watch Mode**: Automatically track logs from newly created pods in a namespace
- 🎛️ **Flexible Options**: Support for custom tail lines, container selection, color disabling, and more

## Prerequisites

- Access to a Kubernetes cluster
- `kubectl` configured to access your cluster

## Installation

### Build from Source

```bash
# Clone repository
git clone <repository-url>
cd ktail

# Build
go build -o ktail

# Make executable and add to PATH (optional)
chmod +x ktail
sudo mv ktail /usr/local/bin/
```

## Usage

### Basic Commands

```bash
# Interactive mode - select namespace and pods to tail logs
ktail

# Tail logs from all pods in a specific namespace
ktail -n my-namespace

# Tail logs from a specific pod
ktail -n my-namespace -p my-pod

# Watch mode - automatically track logs from newly created pods in a namespace
ktail -n my-namespace -w
```

### Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-n, --namespace` | Kubernetes namespace | Interactive selection |
| `-p, --pod` | Pod name | All pods |
| `-c, --container` | Container name | First container |
| `-t, --tail` | Number of lines to show from the end of logs | 100 |
| `-m, --multi` | Enable multi-selection | true |
| `-w, --watch` | Watch mode (when namespace only selected) | false |
| `--no-color` | Disable colored output | false |

### Usage Examples

#### 1. Interactive Mode
```bash
# Interactively select namespace and pods
ktail
```

#### 2. All Pods in Namespace
```bash
# Tail logs from all pods in production namespace
ktail -n production
```

#### 3. Specific Pod
```bash
# Tail logs from a specific pod
ktail -n production -p web-app-7d4f8b9c6-xyz12
```

#### 4. Watch Mode
```bash
# Automatically track logs from newly created pods in namespace
ktail -n production -w
```

#### 5. Specify Container
```bash
# Tail logs from a specific container
ktail -n production -p web-app -c nginx
```

#### 6. Adjust Log Lines
```bash
# Show last 500 lines and follow
ktail -t 500

# Use tail-style flags (follow recent 1000 lines)
ktail -1000f
```

#### 7. Disable Colors
```bash
# Disable colored output
ktail --no-color -n production
```

## Troubleshooting

### Common Issues

1. **kubectl connection error**: Check if `kubectl` is properly configured
2. **Permission error**: Verify you have appropriate permissions for the cluster
3. **Pods not visible**: Check namespace permissions

## License

This project is licensed under the MIT License.


