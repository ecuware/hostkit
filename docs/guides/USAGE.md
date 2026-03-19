# HostKit Usage Guide

Complete guide for using HostKit to manage your hosting server.

## Table of Contents

- [Installation](#installation)
- [First Steps](#first-steps)
- [Using the TUI](#using-the-tui)
- [CLI Commands](#cli-commands)
- [Package Management](#package-management)
- [System Monitoring](#system-monitoring)
- [Service Management](#service-management)
- [Cluster Management](#cluster-management)
- [Backup and Restore](#backup-and-restore)
- [Tips and Tricks](#tips-and-tricks)

## Installation

### Quick Install

```bash
# Download latest release
curl -fsSL https://github.com/ecuware/hostkit/releases/latest/download/hostkit-linux-amd64 -o hostkit
chmod +x hostkit
sudo mv hostkit /usr/local/bin/

# Verify installation
hostkit version
```

### Build from Source

```bash
git clone https://github.com/ecuware/hostkit.git
cd hostkit
go build -o hostkit ./cmd/hostkit/main.go
sudo mv hostkit /usr/local/bin/
```

## First Steps

### View Available Packages

```bash
# List all packages
hostkit list

# List by category
hostkit list --category panels
hostkit list --category databases
hostkit list --category webservers
```

### Check System Requirements

```bash
# View system info
hostkit system info

# Check compatibility
hostkit system check
```

## Using the TUI

The TUI (Terminal User Interface) is the easiest way to use HostKit.

### Launch TUI

```bash
hostkit tui
```

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate items |
| `←` / `→` | Navigate categories |
| `Enter` | Select/Confirm |
| `Esc` | Go back |
| `Q` | Quit |
| `Tab` | Switch tabs |
| `/` | Search |

### TUI Sections

```
HostKit
├── 🎛️ Panels          # Hosting control panels
├── 🗄️ Databases       # Database servers
├── 🌐 Web Servers     # Web servers
├── 🔒 Security        # Security tools
├── ⚙️ Services        # Additional services
├── 📊 Monitoring      # Monitoring tools
├── 🔐 VPN            # VPN solutions
├── 🖥️ Cluster         # Multi-server management
└── 🛠️ System Tools   # System utilities
```

### System Tools

Navigate to `System Tools` for:
- **System Monitor** - Real-time stats
- **Service Manager** - Start/stop services
- **Log Viewer** - View system logs
- **Backup Manager** - Backup configurations
- **Firewall** - Manage firewall rules
- **History** - View installation history

## CLI Commands

### Package Commands

```bash
# Install a package
hostkit install nginx
hostkit install mariadb
hostkit install aapanel

# Uninstall a package
hostkit uninstall nginx

# Check if installed
hostkit status nginx

# View package details
hostkit show nginx
```

### List Commands

```bash
# List all packages
hostkit list

# List with details
hostkit list --verbose

# List installed only
hostkit list --installed

# List by category
hostkit list --category databases

# Search packages
hostkit list --search mysql
```

### System Commands

```bash
# View system info
hostkit system info

# Check system requirements
hostkit system check

# View logs
hostkit logs

# Clear cache
hostkit cache clear

# Update HostKit
hostkit self-update
```

### Service Commands

```bash
# List services
hostkit service list

# Start a service
hostkit service start nginx

# Stop a service
hostkit service stop nginx

# Restart a service
hostkit service restart nginx

# View service status
hostkit service status nginx

# Enable auto-start
hostkit service enable nginx

# Disable auto-start
hostkit service disable nginx
```

## Package Management

### Installing Packages

```bash
# Simple install
sudo hostkit install nginx

# Multiple packages (if supported)
sudo hostkit install nginx mariadb php

# With specific version (if available)
sudo hostkit install nginx --version 1.25.3
```

### Installation Process

1. **Check Requirements** - RAM, disk space, OS compatibility
2. **Resolve Dependencies** - Install required packages first
3. **Download** - Get installation files
4. **Install** - Run installation scripts
5. **Configure** - Apply default configuration
6. **Verify** - Check if installation successful

### Post-Installation

After installation, HostKit shows:
- Access URLs
- Default credentials
- Important configuration notes
- Service commands

**Example:**
```
✓ aaPanel installed successfully!

Admin Interface: http://YOUR_IP:7800
Username: Set during installation
Password: Set during installation

IMPORTANT: Change default credentials immediately!
```

### Uninstalling Packages

```bash
# Remove package
sudo hostkit uninstall nginx

# Remove with dependencies
sudo hostkit uninstall nginx --remove-deps

# Force remove (if normal uninstall fails)
sudo hostkit uninstall nginx --force
```

### Viewing Package Info

```bash
# Show package details
hostkit show nginx

# Output includes:
# - Description
# - Version
# - Requirements
# - Dependencies
# - Installation size
# - Estimated time
```

## System Monitoring

### Real-time Monitor

```bash
# Launch monitor
hostkit monitor

# Or from TUI: System Tools → System Monitor
```

### Displayed Metrics

```
📊 System Monitor

💻 CPU Usage: 45.2%
[████████████░░░░░░░░] 4 Cores

🧠 Memory: 8.2 GB / 16 GB (62.5%)
[████████████░░░░░░░░]

💾 Disk: 28.3% used
/ partition: 150 GB / 500 GB

🌐 Network:
Download: 2.5 MB/s
Upload: 512 KB/s

⚡ Load Average: 0.52 0.48 0.45
⏱️ Uptime: 15d 8h 32m
```

### Color Indicators

- 🟢 **Green** (< 50%) - Normal
- 🟠 **Orange** (50-80%) - Warning
- 🔴 **Red** (> 80%) - Critical

## Service Management

### Viewing Services

```bash
# List all services
hostkit service list

# List running services
hostkit service list --running

# List enabled services
hostkit service list --enabled
```

### Managing Services

```bash
# Start service
sudo hostkit service start nginx

# Stop service
sudo hostkit service stop nginx

# Restart service
sudo hostkit service restart nginx

# Reload configuration
sudo hostkit service reload nginx

# View logs
hostkit service logs nginx

# Follow logs in real-time
hostkit service logs nginx --follow
```

### Common Services

| Service | Description |
|---------|-------------|
| `nginx` | Web server |
| `mysql` / `mariadb` | Database |
| `postgresql` | PostgreSQL database |
| `redis` | Redis cache |
| `docker` | Docker daemon |
| `fail2ban` | Intrusion prevention |

## Cluster Management

### Setup Cluster

```bash
# Install cluster manager
sudo hostkit install cluster-manager

# Configure servers
sudo nano /etc/hostkit/cluster/servers.yaml
```

### Configuration Example

```yaml
clusters:
  - name: production
    servers:
      - name: web-01
        host: 192.168.1.101
        user: root
        tags: [web]
      - name: web-02
        host: 192.168.1.102
        user: root
        tags: [web]
      - name: db-01
        host: 192.168.1.201
        user: root
        tags: [database]
```

### Cluster Commands

```bash
# List all servers
hostkit cluster list

# Check server status
hostkit cluster status

# Execute command on all servers
hostkit cluster exec "apt update" --all

# Execute on specific group
hostkit cluster exec "systemctl restart nginx" --tags web

# Install on multiple servers
hostkit cluster install nginx --tags web

# Parallel execution (faster)
hostkit cluster exec "reboot" --all --parallel

# Upload file to servers
hostkit cluster upload /local/config.conf /etc/nginx/ --all

# Download from servers
hostkit cluster download /var/log/nginx/error.log ./logs/ --servers web-01,web-02
```

### SSH Key Setup

```bash
# Copy SSH key to remote servers
ssh-copy-id -i /root/.ssh/cluster_keys/master_key.pub root@SERVER_IP

# Or use HostKit
hostkit cluster setup-keys
```

## Backup and Restore

### Creating Backups

```bash
# Backup all configurations
sudo hostkit backup create

# Backup specific package
sudo hostkit backup create nginx

# Backup to specific location
sudo hostkit backup create --destination /backup/hostkit/
```

### Backup Contents

Backups typically include:
- Configuration files
- Database dumps (if applicable)
- SSL certificates
- Custom modifications

### Restoring Backups

```bash
# List available backups
hostkit backup list

# Restore from backup
sudo hostkit backup restore backup-name

# Restore specific package from backup
sudo hostkit backup restore backup-name --package nginx
```

### Automated Backups

Add to crontab for daily backups:

```bash
# Edit crontab
sudo crontab -e

# Add line for daily backup at 2 AM
0 2 * * * /usr/local/bin/hostkit backup create --quiet
```

## Tips and Tricks

### Keyboard Shortcuts

**TUI Navigation:**
- `Ctrl+C` - Force quit
- `F5` - Refresh
- `F1` - Help
- `/` - Search mode
- `Esc` - Cancel/Back

### Quick Commands

```bash
# Install multiple packages quickly
sudo hostkit install nginx && sudo hostkit install mariadb && sudo hostkit install php

# Check all services status
hostkit service list --running

# Monitor in background
hostkit monitor &
```

### Environment Variables

```bash
# Debug mode
export HOSTKIT_DEBUG=1
hostkit install nginx

# Custom config directory
export HOSTKIT_CONFIG_DIR=/custom/path
hostkit list

# Skip confirmation
export HOSTKIT_NO_CONFIRM=1
hostkit install nginx
```

### Aliases

Add to `~/.bashrc`:

```bash
# HostKit aliases
alias hk='hostkit'
alias hki='sudo hostkit install'
alias hku='sudo hostkit uninstall'
alias hkl='hostkit list'
alias hks='hostkit service'
```

Then use:
```bash
hki nginx      # Install nginx
hkl            # List packages
hks status     # Service status
```

### Troubleshooting

**Installation fails:**
```bash
# Check logs
hostkit logs

# Try with debug
debug=1 hostkit install nginx

# Check requirements
hostkit system check
```

**Service won't start:**
```bash
# Check status
hostkit service status nginx

# View logs
hostkit service logs nginx

# Check configuration
nginx -t
```

**Connection issues:**
```bash
# Check firewall
hostkit firewall status

# Open port
sudo ufw allow 80/tcp
```

### Best Practices

1. **Always backup** before major changes
2. **Test on staging** before production
3. **Monitor resources** during installations
4. **Keep system updated** - Run `apt update` regularly
5. **Check logs** when things go wrong
6. **Use tags** for cluster organization
7. **Document custom configurations**

### Getting Help

```bash
# Show help
hostkit --help
hostkit install --help

# Show version
hostkit version

# View logs
hostkit logs
```

## Advanced Usage

### Custom Package Installation

Create your own package YAML:

```bash
# Create custom package
sudo mkdir -p /etc/hostkit/custom-packages
sudo nano /etc/hostkit/custom-packages/myapp.yaml

# Add to HostKit
hostkit package add /etc/hostkit/custom-packages/myapp.yaml
```

### Script Integration

```bash
#!/bin/bash
# deploy.sh - Automated deployment script

set -e

echo "Deploying web stack..."

# Install packages
hostkit install nginx
hostkit install mariadb
hostkit install php

# Configure
# ... your configuration ...

# Start services
hostkit service start nginx
hostkit service start mariadb
hostkit service start php-fpm

echo "Deployment complete!"
```

---

For more information, visit the [documentation](../README.md) or [GitHub repository](https://github.com/ecuware/hostkit).
