# HostKit

HostKit is an **all-in-one hosting server management toolkit** that simplifies the installation and management of hosting panels, databases, web servers, and security tools through an interactive terminal interface.

<div class="grid cards" markdown>

-   :material-rocket-launch:{ .lg .middle } __Quick Start__

    ---

    Get HostKit up and running in under 5 minutes with our simple installation guide.

    [:octicons-arrow-right-24: Installation](installation.md)

-   :material-package-variant:{ .lg .middle } __22+ Packages__

    ---

    One-click installation for panels, databases, web servers, and security tools.

    [:octicons-arrow-right-24: Browse Packages](packages.md)

-   :material-monitor-dashboard:{ .lg .middle } __System Monitor__

    ---

    Real-time monitoring of CPU, RAM, disk, and network usage.

    [:octicons-arrow-right-24: Learn More](usage.md#system-monitor)

-   :material-shield-check:{ .lg .middle } __Security First__

    ---

    Built-in firewall, intrusion detection, and SSL certificate management.

    [:octicons-arrow-right-24: Security Tools](packages/security.md)

</div>

---

## ✨ Features

- **Interactive TUI** - Beautiful terminal interface powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **22+ Packages** - Panels, databases, web servers, security tools, and more
- **One-Click Install** - Automated installation with dependency resolution
- **Auto-Version Detection** - Always installs the latest versions
- **System Monitor** - Real-time CPU, RAM, Disk, and Network monitoring
- **Service Manager** - Start, stop, and restart services
- **Backup System** - Automated configuration backups
- **Firewall Manager** - IP blocking and security rules
- **YAML-Driven** - Easy to extend and customize

---

## 🚀 Quick Start

### Installation

```bash
# Download latest release
curl -fsSL https://github.com/ecuware/hostkit/releases/latest/download/hostkit-linux-amd64 -o hostkit
chmod +x hostkit
sudo mv hostkit /usr/local/bin/

# Launch interactive TUI
hostkit tui
```

### Basic Usage

```bash
# List all packages
hostkit list

# Install a package
hostkit install nginx
hostkit install mariadb
hostkit install aapanel
```

---

## 📦 Supported Packages

### Hosting Panels (7)
- aaPanel, CyberPanel, cPanel/WHM, Plesk, HestiaCP, CloudPanel, CWP

### Databases (5)
- MariaDB, MySQL, PostgreSQL, MongoDB, Redis

### Web Servers (3)
- Nginx, OpenLiteSpeed, PHP

### Security (3)
- Fail2ban, CSF, Certbot

### Services (2)
- Docker, Portainer

### Monitoring (1)
- Netdata

### VPN (1)
- WireGuard

---

## 💻 System Requirements

**Minimum:**
- Ubuntu 20.04+ or Debian 11+
- 1GB RAM
- 5GB disk space
- Root access

**Recommended:**
- Ubuntu 22.04 LTS
- 2GB+ RAM
- 20GB+ disk space

---

## 🎯 Why HostKit?

### For System Administrators
- Save hours of manual configuration
- Consistent, repeatable installations
- Easy to script and automate

### For Developers
- Quick development environment setup
- Same environment as production
- Easy to reset and rebuild

### For Beginners
- No complex commands to remember
- Interactive interface guides you
- Safe, tested installation scripts

---

## 📚 Documentation Structure

<div class="grid" markdown>

:material-book-open-page-variant: **[Getting Started](installation.md)**
: Installation, basic usage, and configuration

:material-package: **[Packages](packages.md)**
: Detailed documentation for all supported packages

:material-code-braces: **[Development](contributing.md)**
: Contributing guidelines and development setup

:material-api: **[API Reference](api/package-format.md)**
: Package format specification

</div>

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](contributing.md) for details.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/ecuware/hostkit/blob/main/LICENSE) file for details.

---

## 💬 Support

- 📖 **Documentation**: [GitHub Pages](https://ecuware.github.io/hostkit)
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/ecuware/hostkit/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/ecuware/hostkit/discussions)
