# puby

A command-line utility for managing Dart and Flutter package dependencies.

## Overview

`puby` is a tool that helps you keep your Dart and Flutter projects up-to-date by scanning your `pubspec.yaml` file and detecting outdated dependencies. It can update both SDK versions (Dart and Flutter) and package dependencies with a single command.

## Installation

### Prerequisites

- Go 1.24 or higher

### Building from source

```bash
# Clone the repository
git clone https://github.com/sunderee/puby.git
cd puby

# Build the binary
go build -o ./bin/puby ./cmd/puby

# Optionally, move to a location in your PATH
mv puby /usr/local/bin/  # Linux/macOS
```

## Usage

```bash
# Check for updates in the current directory (checks all dependencies)
puby

# Check for updates in a specific directory
puby --path=/path/to/flutter/project

# Write updates to pubspec.yaml
puby --write

# Only check specific packages
puby --include=http,path,provider

# Check all packages except specific ones
puby --exclude=flutter_svg

# Enable Flutter SDK update check
puby --flutter

# Consider beta SDK versions
puby --beta

# Show help
puby --help

# Show version
puby --version
```

## Command-line options

| Option | Default | Description |
|--------|---------|-------------|
| `--path` | `pubspec.yaml` | Path to the pubspec.yaml file |
| `--write` | `false` | Write changes to pubspec.yaml (otherwise run in dry-run mode) |
| `--include` | | Comma-separated list of packages to include in update check (if not specified, all packages are checked) |
| `--exclude` | | Comma-separated list of packages to exclude from update check |
| `--flutter` | `false` | Check Flutter SDK version |
| `--beta` | `false` | Use beta versions for SDK updates |
| `--help` | `false` | Show help message |
| `--version` | `false` | Show version information |

## Examples

### Checking for updates (dry run)

```bash
puby
```

Output:
```
Checking for updates in /path/to/pubspec.yaml...
=== SDK Updates ===
Dart SDK: 3.7.2

=== Dependency Updates ===
http    : 0.13.3 → 1.3.0
path    : 1.8.0 → 1.9.1
provider: 6.0.0 → 6.1.4

Running in dry-run mode. Use --write flag to apply changes.
```

### Enabling Flutter SDK check

```bash
puby --flutter
```

Output:
```
Checking for updates in /path/to/pubspec.yaml...
=== SDK Updates ===
Dart SDK: 3.7.2
Flutter SDK: 3.29.2

=== Dependency Updates ===
http    : 0.13.3 → 1.3.0
path    : 1.8.0 → 1.9.1
provider: 6.0.0 → 6.1.4

Running in dry-run mode. Use --write flag to apply changes.
```

### Updating dependencies

```bash
puby --write
```

Output:
```
Checking for updates in /path/to/pubspec.yaml...
=== SDK Updates ===
Dart SDK: 3.7.2

=== Dependency Updates ===
http    : 0.13.3 → 1.3.0
path    : 1.8.0 → 1.9.1
provider: 6.0.0 → 6.1.4

Updates have been written to pubspec.yaml
```

### Selective updates

```bash
puby --include=http,path --write
```

Output:
```
Checking for updates in /path/to/pubspec.yaml...
=== Dependency Updates ===
http: 0.13.3 → 1.3.0
path: 1.8.0 → 1.9.1

Updates have been written to pubspec.yaml
```

## How it works

`puby` works by:

1. Parsing your `pubspec.yaml` file to extract current SDK and dependency versions
2. Fetching the latest SDK versions from the Flutter repository
3. Fetching the latest package versions from pub.dev
4. Comparing current and latest versions to identify updates
5. Presenting the updates in a colorful, readable format
6. Optionally writing the changes back to your `pubspec.yaml` file

It uses regex patterns to make targeted updates to the file, preserving its original structure and formatting.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the GNU General Public License v3.0 - check the [LICENSE](./LICENSE) file for more information.