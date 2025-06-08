# ScannerCLI

A CLI tool for scanning and managing repositories with SonarQube integration, built with Go and Cobra.

![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Features

- **Scan repositories**: Clone and analyze codebases
- **Project management**: Create/delete SonarQube projects
- **Simple interface**: Easy-to-use commands

## Installation

### From Source

```bash
git clone https://github.com/CervinoB/scannercli.git
cd scannercli
go build -o scannercli
```

### Using Pre-built Binaries

Download the latest release from the [Releases page](https://github.com/CervinoB/scannercli/releases).

## Usage

### Scan Command

```bash
./scannercli scan
```

Scans a repository by:

1. Creating a new SonarQube project
2. Cloning the target repository
3. Preparing for analysis

Example output:

```
scan called
Project created with key: new-project-1749356835
Scan completed
```

### Delete Command

```bash
./scannercli delete
```

Deletes all projects from SonarQube.

Example output:

```
delete called
Deleted project: new-project-1749356835
Deleted project: old-project-1234567890
```

## Configuration

### Environment Variables

Set these before running commands:

```bash
export SONARQUBE_URL="http://your-sonarqube-instance:9000"
export SONARQUBE_TOKEN="your_api_token"
```

### Custom Repository

Modify the hardcoded repository URL in `cmd/scan.go`:

```go
err = git.CloneRepository("YOUR_REPO_URL", "repo/"+projectName)
```

## Development

### Building

```bash
go build -o scannercli
```

### Testing

Run all tests:

```bash
go test ./...
```

### Adding New Commands

Use Cobra generator:

```bash
go run main.go add command [command-name]
```

## Project Structure

```
.
├── cmd/
│   ├── scan.go       # Scan command implementation
│   ├── delete.go     # Delete command implementation
│   └── root.go       # Root command configuration
├── internal/
│   ├── api/          # SonarQube API client
│   └── git/          # Git operations
├── go.mod
├── go.sum
└── main.go
```

## License

MIT - See [LICENSE](LICENSE) for details.
