# Contributing to HostKit

First off, thank you for considering contributing to HostKit! It's people like you that make HostKit such a great tool.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Workflow](#development-workflow)
- [Style Guidelines](#style-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project and everyone participating in it is governed by our commitment to:
- Being respectful and inclusive
- Welcoming newcomers
- Focusing on constructive feedback
- Prioritizing user safety and privacy

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Linux environment (Ubuntu 20.04+ or Debian 11+ recommended)
- Root access for testing installations

### Setting Up Your Environment

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/hostkit.git
   cd hostkit
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Build the project**:
   ```bash
   go build -o hostkit ./cmd/hostkit/main.go
   ```

5. **Run tests**:
   ```bash
   go test ./...
   ```

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report, please check the existing issues. When creating a bug report, include:

- **Clear title** - Summarize the issue
- **Steps to reproduce** - Detailed steps
- **Expected behavior** - What should happen
- **Actual behavior** - What actually happens
- **System info** - OS, version, architecture
- **Logs** - Relevant error messages

**Example:**
```markdown
**Bug:** Installation fails on Ubuntu 22.04 with "permission denied"

**Steps:**
1. Run `hostkit install nginx`
2. Enter sudo password
3. See error

**Expected:** Nginx should install successfully

**Actual:** Error: "permission denied: /etc/nginx/nginx.conf"

**System:** Ubuntu 22.04 LTS, x86_64, HostKit v1.0.0
```

### Suggesting Enhancements

Enhancement suggestions are welcome! Include:
- Use case description
- Expected behavior
- Possible implementation approach
- Willingness to contribute

### Adding New Packages

This is one of the easiest ways to contribute! See [Creating Packages](development/CREATING_PACKAGES.md) for detailed instructions.

**Quick package addition checklist:**
- [ ] YAML config file in `configs/<category>/`
- [ ] Tested on Ubuntu 20.04+ or Debian 11+
- [ ] Installation script works
- [ ] Post-install instructions clear
- [ ] Dependencies and conflicts defined

### Improving Documentation

Documentation improvements are always welcome:
- Fix typos
- Clarify instructions
- Add examples
- Translate to other languages

### Code Contributions

Areas where code contributions are especially welcome:
- **New installers** - Support for new operating systems
- **TUI improvements** - Better user interface
- **Performance optimizations** - Faster operations
- **Bug fixes** - Squash those bugs!
- **Tests** - Increase test coverage

## Development Workflow

### Branch Naming

- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `docs/description` - Documentation
- `refactor/description` - Code refactoring

**Examples:**
- `feature/add-redis-cluster`
- `bugfix/fix-mysql-install-ubuntu-24`
- `docs/improve-installation-guide`

### Making Changes

1. **Create a branch**:
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes**

3. **Test locally**:
   ```bash
   go test ./...
   go build -o hostkit ./cmd/hostkit/main.go
   ./hostkit list
   ```

4. **Commit your changes** (see [Commit Messages](#commit-messages))

5. **Push to your fork**:
   ```bash
   git push origin feature/my-feature
   ```

6. **Create a Pull Request**

## Style Guidelines

### Go Code Style

Follow standard Go conventions:
- Use `gofmt` for formatting
- Use `golint` for linting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Keep functions small and focused
- Write clear, descriptive variable names
- Add comments for exported functions

**Example:**
```go
// InstallPackage installs a package by ID
// Returns an error if the package is not found or installation fails
func InstallPackage(id string, config *Config) error {
    // Implementation
}
```

### YAML Style (Package Configs)

- Use 2 spaces for indentation
- Use descriptive keys
- Add comments for complex sections
- Keep lines under 100 characters

**Example:**
```yaml
id: mypackage
name: "My Package"
description: "Short, clear description"

requirements:
  min_ram: "1GB"  # Minimum for basic operation
  ports: [8080, 443]  # Required ports
```

## Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Formatting, semicolons, etc.
- `refactor` - Code restructuring
- `test` - Adding tests
- `chore` - Maintenance tasks

**Scopes:**
- `installer` - Installation engine
- `tui` - Terminal UI
- `package` - Package configs
- `docs` - Documentation
- `cluster` - Cluster management

**Examples:**
```
feat(installer): Add parallel installation support

Allow installing multiple packages simultaneously
with dependency resolution. This improves speed
by 40% for batch operations.

Closes #123
```

```
fix(package): Correct nginx port configuration

Nginx was defaulting to port 8080 instead of 80
on fresh installs. Fixed the default config template.

Fixes #456
```

## Pull Request Process

1. **Update documentation** if needed
2. **Add tests** for new functionality
3. **Ensure all tests pass**
4. **Update CHANGELOG.md** with your changes
5. **Link related issues** in the PR description
6. **Request review** from maintainers

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation
- [ ] Breaking change

## Testing
- [ ] Tested on Ubuntu 20.04
- [ ] Tested on Ubuntu 22.04
- [ ] Tested on Debian 11
- [ ] All tests pass

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated

## Related Issues
Fixes #123
Related to #456
```

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Credited in documentation

## Questions?

- Open an issue for discussion
- Join GitHub Discussions
- Check existing documentation

Thank you for contributing to HostKit! 🎉
