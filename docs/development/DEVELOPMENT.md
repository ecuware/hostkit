# Development Setup

This guide helps you set up a development environment for contributing to HostKit.

## Prerequisites

### Required Software

- **Go 1.21+** - Programming language
- **Git** - Version control
- **Make** - Build automation
- **Docker** (optional) - For testing

### Operating System

- **Linux** - Ubuntu 20.04+, Debian 11+, or AlmaLinux 8+
- **Root access** - Required for testing installations
- **Virtual Machine recommended** - For safe testing

## Step-by-Step Setup

### 1. Install Go

**Ubuntu/Debian:**
```bash
# Remove old versions
sudo rm -rf /usr/local/go

# Download and install
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
rm go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

**AlmaLinux:**
```bash
sudo dnf install golang
```

### 2. Install Git

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y git

# AlmaLinux
sudo dnf install -y git
```

### 3. Fork and Clone

1. **Fork on GitHub:**
   - Go to https://github.com/ecuware/hostkit
   - Click "Fork" button

2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/hostkit.git
   cd hostkit
   ```

3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/ecuware/hostkit.git
   git fetch upstream
   ```

### 4. Install Dependencies

```bash
# Download Go modules
go mod download

# Verify
go mod verify
```

### 5. Build the Project

```bash
# Simple build
go build -o hostkit ./cmd/hostkit/main.go

# Or use Make
make build

# Test the binary
./hostkit version
./hostkit list
```

## Development Workflow

### 1. Create a Branch

```bash
# Update main
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/my-feature
```

### 2. Make Changes

Edit files as needed:
- **Go code** - `internal/`, `pkg/`, `cmd/`
- **Packages** - `configs/`
- **Documentation** - `docs/`

### 3. Test Locally

```bash
# Run all tests
go test ./...

# Run specific test
go test ./internal/installer/...

# Test with verbose output
go test -v ./...

# Check race conditions
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. Build and Run

```bash
# Build
go build -o hostkit ./cmd/hostkit/main.go

# Test TUI
./hostkit tui

# Test CLI commands
./hostkit list
./hostkit version
```

### 5. Code Quality Checks

```bash
# Format code
go fmt ./...

# Run linter (install golint first)
go install golang.org/x/lint/golint@latest
golint ./...

# Vet code
go vet ./...

# Check for vulnerabilities (install govulncheck)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## IDE Setup

### VS Code

1. **Install Go extension** - Search "Go" in extensions
2. **Install tools** - Press Ctrl+Shift+P → "Go: Install/Update Tools"
3. **Enable format on save**:
   ```json
   // settings.json
   {
     "editor.formatOnSave": true,
     "go.formatTool": "gofmt",
     "go.lintTool": "golint"
   }
   ```

### GoLand (JetBrains)

1. Open project in GoLand
2. GoLand auto-detects Go SDK
3. Enable:
   - Auto-import
   - Format on save
   - Go vet on save

## Testing Your Changes

### Unit Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestInstaller ./internal/installer/
```

### Integration Testing

Test actual package installations:

```bash
# Build
make build

# Test installation
sudo ./hostkit install nginx

# Verify
systemctl status nginx
curl http://localhost

# Clean up
sudo apt-get remove -y nginx
```

### Manual Testing Checklist

Before submitting PR:

- [ ] `go build` succeeds
- [ ] `go test ./...` passes
- [ ] `go fmt` makes no changes
- [ ] `go vet` reports no issues
- [ ] Application starts without errors
- [ ] TUI displays correctly
- [ ] At least one package installs successfully
- [ ] No panic or crash

## Debugging

### Enable Debug Logging

```bash
# Set debug mode
export HOSTKIT_DEBUG=1

# Run with debug
./hostkit install nginx
```

### Use Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug tests
dlv test ./internal/installer

# Debug main
dlv debug ./cmd/hostkit/main.go
```

### Common Issues

**Issue: "command not found: go"**
```bash
# Solution: Add Go to PATH
export PATH=$PATH:/usr/local/go/bin
```

**Issue: "cannot find package"**
```bash
# Solution: Download modules
go mod download
go mod tidy
```

**Issue: Build fails with undefined symbols**
```bash
# Solution: Clean and rebuild
go clean -cache
go build -o hostkit ./cmd/hostkit/main.go
```

**Issue: Permission denied during installation**
```bash
# Solution: Run with sudo
sudo ./hostkit install nginx
```

## Docker Development (Recommended)

Use Docker for safe testing:

### Dockerfile

```dockerfile
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    golang-go \
    git \
    make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o hostkit ./cmd/hostkit/main.go

CMD ["./hostkit", "tui"]
```

### Usage

```bash
# Build image
docker build -t hostkit-dev .

# Run tests
docker run --rm hostkit-dev go test ./...

# Run TUI
docker run -it --rm hostkit-dev ./hostkit tui
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  hostkit:
    build: .
    volumes:
      - .:/app
      - /var/run/docker.sock:/var/run/docker.sock
    privileged: true
    command: bash
```

```bash
docker-compose run hostkit bash
```

## Performance Testing

### Benchmark Tests

```bash
# Run benchmarks
go test -bench=. ./...

# Run with memory profiling
go test -bench=. -memprofile=mem.out ./...
go tool pprof mem.out

# Run with CPU profiling
go test -bench=. -cpuprofile=cpu.out ./...
go tool pprof cpu.out
```

### Memory Leak Detection

```bash
# Install and run leak test
go test -memprofile=prof.mem ./...
go tool pprof -http=:8080 prof.mem
```

## Continuous Integration

### Pre-commit Hooks

Install pre-commit:

```bash
pip install pre-commit
pre-commit install
```

Create `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      - id: go-vet
      - id: go-lint
      - id: go-imports
      - id: go-unit-tests
```

## Building for Release

### Cross-compilation

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o hostkit-linux-amd64 ./cmd/hostkit/main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o hostkit-linux-arm64 ./cmd/hostkit/main.go

# macOS (if needed for testing)
GOOS=darwin GOARCH=amd64 go build -o hostkit-darwin-amd64 ./cmd/hostkit/main.go
```

### Build Flags

```bash
# Strip debug symbols (smaller binary)
go build -ldflags="-s -w" -o hostkit ./cmd/hostkit/main.go

# Version injection
VERSION=$(git describe --tags --always)
go build -ldflags="-X main.version=$VERSION" -o hostkit ./cmd/hostkit/main.go
```

## Documentation

### Update Documentation

If you change behavior:

1. Update code comments
2. Update relevant docs in `docs/`
3. Update CHANGELOG.md
4. Add examples if needed

### Generate Documentation

```bash
# Install godoc
go install golang.org/x/tools/cmd/godoc@latest

# Serve documentation
godoc -http=:6060

# Open http://localhost:6060/pkg/github.com/ecuware/hostkit/
```

## Next Steps

Now you're ready to contribute:

1. Read [CONTRIBUTING.md](../CONTRIBUTING.md)
2. Check [ARCHITECTURE.md](../architecture/ARCHITECTURE.md)
3. Pick an issue to work on
4. Create your first PR!

## Getting Help

- **Documentation**: Check `docs/` directory
- **Issues**: Search existing issues first
- **Discussions**: Ask questions in GitHub Discussions
- **Code**: Look at similar implementations in the codebase

Happy coding! 🚀
