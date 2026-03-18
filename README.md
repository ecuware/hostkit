# 🎛️ HostKit

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Platform](https://img.shields.io/badge/platform-Linux-blue.svg)]()

> **All-in-one hosting server management toolkit** written in Go. Install, configure, and manage hosting panels, databases, web servers, and security tools with an interactive TUI.

![HostKit Demo](https://via.placeholder.com/800x400/2d3748/ffffff?text=HostKit+TUI+Demo)

## ✨ Features

🎯 **Interactive TUI** - Beautiful terminal interface powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)  
📦 **22+ Packages** - Panels, databases, web servers, security tools, and more  
🔧 **One-Click Install** - Automated installation with dependency resolution  
🔄 **Auto-Version Detection** - Always installs the latest versions  
📊 **System Monitor** - Real-time CPU, RAM, Disk, and Network monitoring  
🛡️ **Security Tools** - Firewall management and intrusion prevention  
💾 **Backup System** - Automated configuration backups  
📜 **YAML-Driven** - Easy to extend and customize  

## 🚀 Quick Start

### Installation

**Linux (AMD64)**
```bash
curl -fsSL https://github.com/ecuware/hostkit/releases/latest/download/hostkit-linux-amd64 -o hostkit
chmod +x hostkit
sudo mv hostkit /usr/local/bin/
```

**Build from Source**
```bash
git clone https://github.com/ecuware/hostkit.git
cd hostkit
go build -o hostkit ./cmd/hostkit/main.go
sudo mv hostkit /usr/local/bin/
```

### Usage

```bash
# Launch interactive TUI
hostkit tui

# List all packages
hostkit list

# Install a specific package
hostkit install nginx
hostkit install mariadb
hostkit install aapanel
```

## 📦 Supported Packages

### 🎛️ Hosting Panels
| Package | Description | Install Size | Time |
|---------|-------------|--------------|------|
| **aaPanel** | Simple and lightweight control panel | 2GB | 3-5 min |
| **CyberPanel** | OpenLiteSpeed powered panel | 3GB | 10-15 min |
| **cPanel/WHM** | Industry-standard hosting panel | 5GB | 15-20 min |
| **Plesk** | Professional hosting management | 3GB | 10-15 min |
| **HestiaCP** | Open source hosting panel | 1GB | 3-5 min |
| **CloudPanel** | Modern PHP hosting panel | 1GB | 3-5 min |
| **CentOS Web Panel** | Free Linux control panel | 2GB | 5-10 min |

### 🗄️ Databases
| Package | Description | Install Size | Time |
|---------|-------------|--------------|------|
| **MariaDB** | MySQL-compatible database | 200MB | 1-2 min |
| **MySQL** | World's most popular database | 400MB | 1-2 min |
| **PostgreSQL** | Advanced relational database | 300MB | 1-2 min |
| **MongoDB** | Document-oriented NoSQL | 800MB | 2-3 min |
| **Redis** | In-memory data store | 50MB | 30 sec |

### 🌐 Web Servers
| Package | Description | Install Size | Time |
|---------|-------------|--------------|------|
| **Nginx** | High-performance web server | 50MB | 30 sec |
| **OpenLiteSpeed** | Lightweight HTTP server | 150MB | 1-2 min |
| **PHP** | Multi-version PHP with FPM | 200MB | 1-2 min |

### 🔒 Security Tools
- **Fail2ban** - Intrusion prevention framework
- **CSF** - ConfigServer Security & Firewall  
- **Certbot** - Free SSL certificates (Let's Encrypt)

### ⚙️ Services
- **Docker** - Container platform
- **Portainer** - Docker management UI

### 📊 Monitoring
- **Netdata** - Real-time performance monitoring

### 🔐 VPN
- **WireGuard** - Modern, fast, secure VPN tunnel

## 🖥️ System Monitor

HostKit includes a built-in system monitor accessible from the TUI:

```
📊 System Monitor

💻 CPU Usage
[████████████████████] 45.2%
Cores: 4

🧠 Memory Usage  
[████████████░░░░░░░░] 62.5% (8.2 GB / 16 GB)

💾 Disk Usage
[██████░░░░░░░░░░░░░░] 28.3%

🌐 Network
Download: 2.5 MB/s | Upload: 512 KB/s

⚡ Load Average
1min: 0.52 | 5min: 0.48 | 15min: 0.45

⏱️ Uptime: 15d 8h 32m
```

## 🛠️ System Requirements

HostKit supports modern Linux distributions. Based on package availability:

**Minimum:**
- Ubuntu 20.04+ or Debian 11+
- 1GB RAM
- 5GB disk space
- Root access

**Recommended:**
- Ubuntu 22.04 LTS or Debian 12
- 2GB+ RAM
- 20GB+ disk space

**Supported Distributions:**

| OS | Versions | Package Support | Recommendation |
|----|----------|-----------------|----------------|
| **Ubuntu** | 20.04, 22.04, 24.04 | ⭐⭐⭐ All packages | ✅ **Highly Recommended** |
| **Debian** | 11, 12 | ⭐⭐⭐ Most packages | ✅ Recommended |
| **AlmaLinux** | 8, 9 | ⭐⭐ Limited (Panels only) | ⚠️ Partial support |

*Note: CentOS 7/8 reached End-of-Life and is no longer maintained. Some legacy packages may still support it, but we strongly recommend migrating to **Ubuntu 22.04 LTS** for the best compatibility and security.*

## 📋 Configuration

Packages are defined in YAML format. Example:

```yaml
id: mypackage
name: "My Package"
category: panel
description: "Custom package description"
icon: "🎨"

version:
  current: "1.0"
  source:
    type: "github_release"
    owner: "username"
    repo: "repository"

requirements:
  os:
    ubuntu: ["20.04", "22.04"]
    centos: ["7", "8"]
  min_ram: "1GB"
  min_disk: "5GB"
  install_size: "2GB"
  estimated_time: "3-5 minutes"
  ports: [8080, 443]

install:
  method: "shell"
  script: |
    wget -O install.sh https://example.com/install.sh
    bash install.sh
```

### Adding Custom Packages

1. Create a YAML file in `configs/<category>/`
2. Define your package configuration
3. Run `hostkit list` - your package appears automatically!

## 🏗️ Development

```bash
# Clone repository
git clone https://github.com/ecuware/hostkit.git
cd hostkit

# Install dependencies
go mod download

# Build
go build -o hostkit ./cmd/hostkit/main.go

# Run
./hostkit tui

# Build for all platforms
make build-all
```

### Project Structure

```
hostkit/
├── cmd/hostkit/           # Main application entry
├── internal/
│   ├── config/            # YAML configuration parser
│   ├── tui/               # Terminal UI (Bubble Tea)
│   ├── installer/         # Installation engine
│   │   ├── engine.go      # Core installation logic
│   │   ├── resolver.go    # Dependency resolution
│   │   └── executor.go    # Command execution
│   ├── monitor/           # System monitoring
│   ├── service/           # Service management
│   ├── backup/            # Backup management
│   ├── firewall/          # Firewall management
│   ├── status/            # Installation status
│   └── history/           # Installation history
├── pkg/
│   └── detector/          # Version detection
└── configs/               # Package definitions
    ├── panels/            # 7 hosting panels
    ├── databases/         # 5 database servers
    ├── webservers/        # 3 web servers
    ├── security/          # 3 security tools
    ├── services/          # 2 services
    ├── monitoring/        # 1 monitoring tool
    └── vpn/               # 1 VPN solution
```

## 🗺️ Roadmap

- [ ] **Web UI Dashboard** - Browser-based management interface
- [ ] **SSH Remote Management** - Manage multiple servers
- [ ] **WordPress Toolkit** - One-click WP installation
- [ ] **Email Server Stack** - Postfix, Dovecot, Roundcube
- [ ] **SSL Automation** - Auto-renewal and deployment
- [ ] **Backup Automation** - Scheduled backups to cloud
- [ ] **Ansible Integration** - Infrastructure as code
- [ ] **Plugin System** - Third-party package support

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Adding New Packages

Want to add support for a new package? It's easy!

1. Create a YAML file in `configs/<category>/`
2. Test it locally
3. Submit a PR

See [ADDING_PACKAGES.md](ADDING_PACKAGES.md) for detailed guide.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Cobra](https://github.com/spf13/cobra) - CLI framework  
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- All the amazing open source projects that make this possible!

## 💬 Support

- 📖 **Documentation**: [GitHub Pages](https://ecuware.github.io/hostkit)
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/ecuware/hostkit/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/ecuware/hostkit/discussions)

---

<p align="center">Made with ❤️ by the HostKit team</p>