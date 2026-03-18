# Contributing Guide

Thank you for your interest in contributing to HostKit! 🎉

## 🤝 How to Contribute

### Reporting Bugs

1. **Check existing issues** - Search before creating new ones
2. **Use issue templates** - Select bug report template
3. **Provide details:**
   - HostKit version
   - OS and version
   - Steps to reproduce
   - Expected vs actual behavior
   - Error messages/logs

### Suggesting Features

1. **Check roadmap** - See if already planned
2. **Open discussion** - Use GitHub Discussions for ideas
3. **Create issue** - Use feature request template

### Contributing Code

#### 1. Fork and Clone

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR_USERNAME/hostkit.git
cd hostkit
```

#### 2. Create Branch

```bash
git checkout -b feature/amazing-feature
# or
git checkout -b fix/bug-description
```

#### 3. Make Changes

Follow our coding standards:
- Use `gofmt` for formatting
- Add comments for complex logic
- Keep functions small and focused
- Write clear commit messages

#### 4. Test

```bash
# Build
go build -o hostkit ./cmd/hostkit/main.go

# Run
./hostkit tui

# Test your changes thoroughly
```

#### 5. Commit

```bash
git add .
git commit -m "feat: add amazing feature"
```

**Commit Message Format:**
```
type: description

[optional body]

[optional footer]
```

**Types:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `style:` - Formatting
- `refactor:` - Code restructuring
- `test:` - Tests
- `chore:` - Maintenance

#### 6. Push and Create PR

```bash
git push origin feature/amazing-feature
```

Then create a Pull Request on GitHub.

## 📝 Code Standards

### Go Code Style

```go
// Good: Clear function name and comment
// InstallPackage installs a package with given configuration
func InstallPackage(cfg *config.Config) error {
    // Implementation
}

// Bad: Unclear name, no comment
func install(c *config.Config) error {
    // Implementation
}
```

### Package YAML Standards

```yaml
# Use clear, lowercase IDs
id: mypackage

# Descriptive names
name: "My Package"

# Accurate requirements
requirements:
  min_ram: "1GB"  # Not "1000MB"
  
# Tested installation scripts
install:
  script: |
    # Each command on new line
    echo "Step 1"
    echo "Step 2"
```

## 🎁 Adding New Packages

Want to add support for a new software package? Great!

### Step 1: Check Existing Issues

Search for existing requests or create one to discuss.

### Step 2: Create Package YAML

```bash
# Create file
touch configs/category/mypackage.yaml

# Edit with your favorite editor
vim configs/category/mypackage.yaml
```

### Step 3: Follow Template

```yaml
id: mypackage
name: "My Package"
category: category
description: "Clear description of what it does"
icon: "🎨"

version:
  current: "1.0"
  source:
    type: "static"

requirements:
  os:
    ubuntu: ["20.04", "22.04"]
    debian: ["11", "12"]
  min_ram: "512MB"
  min_disk: "1GB"
  install_size: "500MB"
  estimated_time: "1-2 minutes"
  ports: [8080]

install:
  method: "shell"
  script: |
    # Your installation commands
    wget https://example.com/install.sh
    bash install.sh
```

### Step 4: Test

```bash
# Build
make build

# List (should show your package)
./hostkit list

# Try to install (use --dry-run if possible)
./hostkit install mypackage
```

### Step 5: Document

Update relevant docs:
- README.md - Add to package list
- docs/packages.md - Add details
- CHANGELOG.md - Mention in "Unreleased"

## 🏗️ Project Structure

```
hostkit/
├── cmd/hostkit/        # Main entry point
├── internal/
│   ├── config/         # YAML parsing
│   ├── tui/           # Terminal UI
│   ├── installer/     # Installation logic
│   ├── monitor/       # System monitoring
│   ├── service/       # Service management
│   ├── backup/        # Backup system
│   ├── firewall/      # Firewall management
│   ├── status/        # Status checking
│   ├── history/       # Installation history
│   └── checker/       # System checks
├── pkg/
│   └── detector/      # Version detection
├── configs/           # Package definitions
└── docs/             # Documentation
```

## 🎯 Development Workflow

### Running Locally

```bash
# 1. Clone
git clone https://github.com/ecuware/hostkit.git
cd hostkit

# 2. Install dependencies
go mod download

# 3. Build
go build -o hostkit ./cmd/hostkit/main.go

# 4. Run
./hostkit tui
```

### Adding a Feature

Example: Adding a new command

```go
// cmd/hostkit/main.go
rootCmd.AddCommand(newMyCommand())

func newMyCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "mycommand",
        Short: "Short description",
        Run: func(cmd *cobra.Command, args []string) {
            // Implementation
        },
    }
}
```

## ✅ Checklist Before PR

- [ ] Code builds without errors
- [ ] Follows Go formatting (`gofmt`)
- [ ] No unnecessary imports
- [ ] Comments added for complex logic
- [ ] Tested on target OS
- [ ] Documentation updated
- [ ] Commit messages are clear
- [ ] Branch is up to date with main

## 🐛 Debugging Tips

### Build Issues

```bash
# Clean and rebuild
go clean -cache
rm -f hostkit
go build -v -o hostkit ./cmd/hostkit/main.go
```

### Runtime Issues

```bash
# Enable verbose logging
./hostkit --debug tui

# Check logs
tail -f /var/log/hostkit/install.log
```

## 📞 Getting Help

- **Discord:** [Join our community](https://discord.gg/hostkit)
- **GitHub Discussions:** Q&A and ideas
- **Issues:** Bugs and feature requests

## 🏆 Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Added to Hall of Fame

Thank you for contributing! 🎉