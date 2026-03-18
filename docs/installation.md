# Installation Guide

## 📥 Download Pre-built Binary

### Linux (AMD64)

```bash
# Download latest release
curl -fsSL https://github.com/ecuware/hostkit/releases/latest/download/hostkit-linux-amd64 -o hostkit

# Make executable
chmod +x hostkit

# Move to PATH
sudo mv hostkit /usr/local/bin/

# Verify installation
hostkit version
```

### Linux (ARM64)

```bash
curl -fsSL https://github.com/ecuware/hostkit/releases/latest/download/hostkit-linux-arm64 -o hostkit
chmod +x hostkit
sudo mv hostkit /usr/local/bin/
```

## 🏗️ Build from Source

### Prerequisites

- Go 1.21 or higher
- Git

### Build Steps

```bash
# Clone repository
git clone https://github.com/ecuware/hostkit.git
cd hostkit

# Download dependencies
go mod download

# Build
go build -o hostkit ./cmd/hostkit/main.go

# Optional: Install to system
sudo mv hostkit /usr/local/bin/
```

### Build for Multiple Platforms

```bash
# Build all platforms
make build-all

# Or manually:
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o hostkit-linux-amd64 ./cmd/hostkit/main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o hostkit-linux-arm64 ./cmd/hostkit/main.go
```

## 🛠️ System Requirements

### Minimum Requirements

- **OS:** Ubuntu 20.04+ or Debian 11+
- **RAM:** 1GB
- **Disk:** 5GB free space
- **Access:** Root privileges

### Recommended

- **OS:** Ubuntu 22.04 LTS
- **RAM:** 2GB+
- **Disk:** 20GB+ free space
- **Network:** Stable internet connection

### Supported Distributions

| Distribution | Versions | Support Level |
|-------------|----------|---------------|
| Ubuntu | 20.04, 22.04, 24.04 | ⭐⭐⭐ Full Support |
| Debian | 11, 12 | ⭐⭐⭐ Most Packages |
| AlmaLinux | 8, 9 | ⭐⭐ Limited |

## 🚀 First Run

```bash
# Launch interactive TUI
hostkit tui

# Or use CLI commands
hostkit list                    # List all packages
hostkit install nginx          # Install Nginx
hostkit install mariadb        # Install MariaDB
```

## 🔄 Updating

```bash
# Download latest version (same as installation)
curl -fsSL https://github.com/ecuware/hostkit/releases/latest/download/hostkit-linux-amd64 -o hostkit
chmod +x hostkit
sudo mv hostkit /usr/local/bin/
```

## ❌ Uninstallation

```bash
# Remove binary
sudo rm /usr/local/bin/hostkit

# Remove configuration (optional)
sudo rm -rf /var/lib/hostkit
sudo rm -rf /var/log/hostkit
```

## 🐛 Troubleshooting

### Permission Denied

```bash
# Make sure binary is executable
chmod +x hostkit

# Or use sudo
sudo hostkit tui
```

### Command Not Found

```bash
# Check if in PATH
which hostkit

# If not found, use full path
./hostkit tui
```

### Build Errors

```bash
# Update Go to 1.21+
go version

# Clean and rebuild
rm -f hostkit
go clean -cache
go build -o hostkit ./cmd/hostkit/main.go
```