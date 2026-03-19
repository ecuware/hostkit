# Package Format Specification

This document describes the YAML format for defining packages in HostKit.

## Overview

Each package in HostKit is defined by a YAML file located in `configs/<category>/`. These files tell HostKit how to install, configure, and manage the package.

## File Structure

```yaml
id: unique_package_id
name: "Display Name"
category: panel|database|webserver|security|service|monitoring|vpn
description: "Short description"
icon: "🎨"
version:
  current: "1.0.0"
  source:
    type: github_release|api|static
    # ... source-specific config
requirements:
  os:
    ubuntu: ["20.04", "22.04"]
    debian: ["11", "12"]
  min_ram: "1GB"
  min_disk: "5GB"
  install_size: "2GB"
  estimated_time: "3-5 minutes"
  ports: [80, 443]
install:
  method: shell|script|package
  # ... installation-specific config
dependencies:
  - other_package_id
conflicts:
  - conflicting_package_id
```

## Fields Reference

### Basic Information

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier (lowercase, no spaces) |
| `name` | string | Yes | Display name for UI |
| `category` | string | Yes | Category folder name |
| `description` | string | Yes | Short description (max 100 chars) |
| `icon` | string | No | Emoji icon for UI |

### Version

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version.current` | string | Yes | Current version number |
| `version.source.type` | string | Yes | Version detection method |

#### Version Source Types

**github_release:**
```yaml
version:
  source:
    type: github_release
    owner: nginx
    repo: nginx
    asset_pattern: "nginx-(.*)\\.tar\\.gz"
```

**api:**
```yaml
version:
  source:
    type: api
    url: "https://api.github.com/repos/owner/repo/releases/latest"
    json_path: "tag_name"
```

**static:**
```yaml
version:
  source:
    type: static
    version: "1.25.3"
```

### Requirements

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `requirements.os` | object | Yes | Supported operating systems |
| `requirements.min_ram` | string | Yes | Minimum RAM required |
| `requirements.min_disk` | string | Yes | Minimum free disk space |
| `requirements.install_size` | string | Yes | Installation size |
| `requirements.estimated_time` | string | Yes | Estimated install time |
| `requirements.ports` | array | No | Required ports |

#### OS Support

```yaml
requirements:
  os:
    ubuntu: ["20.04", "22.04", "24.04"]
    debian: ["11", "12"]
    almalinux: ["8", "9"]
```

### Installation Methods

#### Shell Script

```yaml
install:
  method: shell
  script: |
    #!/bin/bash
    set -e
    apt-get update
    apt-get install -y package-name
    systemctl enable package-name
```

#### External Script

```yaml
install:
  method: script
  url: "https://example.com/install.sh"
  checksum: "sha256:abc123..."
```

#### Package Manager

```yaml
install:
  method: package
  packages:
    - package-name
    - package-name-extra
  repository:
    name: repo-name
    url: "https://repo.example.com"
    key: "https://repo.example.com/key.asc"
```

### Dependencies

```yaml
dependencies:
  - nginx          # Install Nginx first
  - mariadb        # Install MariaDB first
```

Dependencies are installed before the package.

### Conflicts

```yaml
conflicts:
  - apache2        # Cannot be installed with Apache
  - mysql          # Cannot be installed with MySQL
```

HostKit will prevent installation if conflicting packages are present.

## Complete Example

```yaml
id: aapanel
name: "aaPanel"
category: panel
description: "Simple and lightweight hosting control panel"
icon: "🎛️"

version:
  current: "6.0"
  source:
    type: static
    version: "latest"

requirements:
  os:
    ubuntu: ["18.04", "20.04", "22.04", "24.04"]
    debian: ["10", "11", "12"]
  min_ram: "1GB"
  min_disk: "2GB"
  install_size: "2GB"
  estimated_time: "3-5 minutes"
  ports: [7800, 8888]

install:
  method: shell
  script: |
    #!/bin/bash
    set -e
    
    # Download and install aaPanel
    URL="https://www.aapanel.com/script/install_6.0_en.sh"
    wget -O install.sh "$URL"
    bash install.sh aapanel
    
    # Wait for installation
    sleep 5
    
    # Enable firewall rules
    if command -v ufw &> /dev/null; then
      ufw allow 7800/tcp
      ufw allow 8888/tcp
    fi

dependencies: []
conflicts: []

post_install:
  message: |
    aaPanel installation complete!
    
    Access URL: http://YOUR_IP:7800
    
    Please save your login credentials shown above.
  show_logs: true
```

## Validation

Package files are validated against the following rules:

1. **Required Fields:** All required fields must be present
2. **ID Format:** Must be lowercase alphanumeric with underscores
3. **Category:** Must match an existing category folder
4. **OS Versions:** Must be valid version strings
5. **Script Safety:** Scripts are checked for dangerous commands

## Adding New Packages

1. Create YAML file in `configs/<category>/`
2. Follow the format specification above
3. Test locally with `hostkit list` and `hostkit install <package>`
4. Submit a pull request

## Category Guidelines

### Panels
- Web-based control panels
- Multi-domain support
- User management features

### Databases
- SQL or NoSQL databases
- Client tools included
- Service management

### Web Servers
- HTTP/HTTPS servers
- Reverse proxy capabilities
- Static file serving

### Security
- Firewall tools
- Intrusion detection
- SSL/TLS management

### Services
- Container platforms
- Supporting services

### Monitoring
- Metrics collection
- Dashboard/visualization
- Alerting capabilities

### VPN
- Network tunneling
- Encryption support
- Multi-client support
