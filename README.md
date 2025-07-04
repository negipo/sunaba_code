# sunaba_code

A Go CLI tool that runs Claude Code in a sandboxed environment using macOS `sandbox-exec`.

## Overview

sunaba_code (砂場コード - "sandbox code" in Japanese) provides a secure way to run Claude Code with restricted file system access. By default, it only allows writing to the current directory, with no network access.

## Features

- 🔒 Sandboxed execution using macOS `sandbox-exec`
- 📁 Configurable write permissions (default: current directory only)
- 🌐 Optional network access control
- 🎯 Automatic Claude Code executable detection
- ⚡ Simple and lightweight

## Installation

```bash
go install github.com/negipo/sunaba_code/cmd/sunaba_code@latest
```

Or build from source:

```bash
git clone https://github.com/negipo/sunaba_code.git
cd sunaba_code
go build -o sunaba_code ./cmd/sunaba_code
```

## Usage

### Basic Usage

Run Claude Code with default settings (current directory writable):

```bash
sunaba_code
```

### Allow Writing to Specific Directories

```bash
sunaba_code -w ~/projects -w /tmp
```

### Enable Network Access

```bash
sunaba_code --network
```

### Run a Different Command

```bash
sunaba_code -- /usr/bin/python script.py
```

### Command Line Options

- `-w, --writable`: Paths where write access is allowed (can be specified multiple times)
- `--network`: Allow network access
- `--ports`: Specific ports to allow (requires --network flag)
- `-h, --help`: Show help message

## How It Works

sunaba_code creates a sandbox profile for macOS `sandbox-exec` that:

1. Denies all operations by default
2. Allows reading from all files
3. Allows writing only to specified directories
4. Optionally allows network access
5. Allows process forking and execution

This ensures that Claude Code (or any other command) can only modify files in explicitly allowed directories.

## Examples

### Development Workflow

```bash
# Work on a project with access only to the project directory
cd ~/projects/my-app
sunaba_code
```

### Multiple Projects

```bash
# Allow access to multiple project directories
sunaba_code -w ~/projects/frontend -w ~/projects/backend
```

### With Network Access for API Development

```bash
# Enable network access for making API calls
sunaba_code --network
```

## Requirements

- macOS (uses `sandbox-exec` which is macOS-specific)
- Go 1.21 or later (for building from source)
- Claude Code installed and accessible in PATH

## Security Considerations

- By default, write access is restricted to the current directory only
- Network access is disabled by default
- The sandbox profile is created temporarily and cleaned up after execution
- All file reads are allowed (following the principle of least surprise for development tools)

## Credits

Inspired by [maccha](https://github.com/kazuho/maccha) by Kazuho Oku.

## License

MIT License - see LICENSE file for details.
