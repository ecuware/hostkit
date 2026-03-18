# Usage Guide

## 🎛️ Interactive TUI

The easiest way to use HostKit is through the interactive Terminal User Interface (TUI).

```bash
hostkit tui
```

### Navigation

- **↑/↓** - Navigate menus
- **Enter** - Select item
- **ESC** - Go back
- **Q** - Quit
- **I** - Install (in package details)
- **R** - Refresh (in system monitor)

### Menu Structure

```
HostKit
├── 🎛️ Panels
│   ├── aaPanel
│   ├── CyberPanel
│   ├── cPanel/WHM
│   └── ...
├── 🗄️ Databases
│   ├── MariaDB
│   ├── MySQL
│   └── ...
├── 🌐 Web Servers
├── 🔒 Security
├── ⚙️ Services
├── 📊 Monitoring
├── 🔐 VPN
├── 🛠️ System Tools
│   ├── System Monitor
│   ├── Service Manager
│   ├── Log Viewer
│   ├── Backup Manager
│   └── Firewall
└── 📜 History
```

## 💻 CLI Commands

### List Packages

```bash
# List all packages
hostkit list

# List by category
hostkit list --category panels
hostkit list --category databases
hostkit list --category webservers
```

### Install Packages

```bash
# Install a single package
hostkit install nginx
hostkit install mariadb
hostkit install aapanel

# The install command will:
# 1. Check system requirements
# 2. Show confirmation dialog
# 3. Install dependencies
# 4. Show progress
# 5. Verify installation
```

### Check Version

```bash
hostkit version
```

## 📊 System Monitor

Access the system monitor from TUI:
**System Tools** → **System Monitor**

### Features

- **Real-time CPU Usage** - With color-coded progress bars
- **Memory Usage** - Total, used, free, cached
- **Disk Usage** - Space monitoring
- **Network Stats** - Download/Upload speeds
- **Load Average** - 1, 5, 15 minute averages
- **Uptime** - System running time

### Color Codes

- 🟢 **Green** < 50% - Normal
- 🟠 **Orange** 50-80% - Warning
- 🔴 **Red** > 80% - Critical

## ⚙️ Service Manager

Manage installed services:

```
System Tools → Service Manager
```

### Available Actions

- **Start** - Start a service
- **Stop** - Stop a service
- **Restart** - Restart a service
- **Enable** - Enable auto-start
- **Disable** - Disable auto-start
- **View Status** - Check if running

### Common Services

- nginx
- mysql / mariadb
- postgresql
- redis
- docker
- fail2ban

## 📜 Log Viewer

View system and service logs:

```
System Tools → Log Viewer
```

### Available Logs

- System logs (`/var/log/syslog`)
- Auth logs (`/var/log/auth.log`)
- Nginx error logs
- MySQL error logs
- Fail2ban logs
- HostKit logs

### Features

- Real-time log tailing
- Filter by keyword
- Export logs

## 💾 Backup Manager

Create and restore configuration backups:

```
System Tools → Backup Manager
```

### Creating Backups

1. Select package to backup
2. Choose backup type (config/data/full)
3. Backup is created automatically

### Restoring Backups

1. View available backups
2. Select backup to restore
3. Confirm restoration

### Backup Location

Backups are stored in: `/var/backups/hostkit/`

## 🔥 Firewall Manager

Manage firewall rules (requires installed firewall):

```
System Tools → Firewall
```

### Supported Firewalls

- **CSF** (ConfigServer Security & Firewall)
- **Fail2ban**
- **UFW** (Uncomplicated Firewall)

### Features

- View firewall status
- Block IP addresses
- Unblock IP addresses
- View blocked IPs list
- Restart firewall

## 📦 Package Installation Flow

1. **Select Package** - Browse categories or search
2. **View Details** - See requirements, dependencies, size
3. **Confirm** - Review installation details
4. **Install** - Watch progress with real-time logs
5. **Verify** - Automatic post-installation checks

### Installation Status Icons

- ✅ **Installed** - Package is installed
- 🔄 **Update Available** - Newer version available
- ⏳ **Installing** - Installation in progress
- ❌ **Failed** - Installation failed

## 🎯 Tips & Tricks

### Quick Search

In any list view, start typing to filter items.

### Keyboard Shortcuts

- **Ctrl+C** - Force quit
- **Tab** - Switch between fields (in forms)
- **Space** - Toggle selection

### Dry Run

To simulate installation without actually installing:

```bash
# Edit source to enable dryRun mode
# Or check the code for dryRun flag
```

### Logs

HostKit logs are stored in:
- `/var/log/hostkit/` - Installation logs
- `/var/log/hostkit/install.log` - Detailed logs

## ❓ FAQ

**Q: Can I install multiple packages at once?**
A: Currently, install them one by one through the TUI.

**Q: Does it work on macOS/Windows?**
A: No, HostKit requires Linux. Most features need `/proc` filesystem access.

**Q: How do I add a custom package?**
A: Create a YAML file in `configs/<category>/` directory.

**Q: Can I uninstall packages?**
A: Yes, from package details screen in TUI.

**Q: Is it safe to use on production servers?**
A: Test on a staging server first. Always backup configurations.