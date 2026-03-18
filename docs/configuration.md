# Configuration Guide

HostKit uses YAML files to define packages. This makes it easy to customize and extend.

## 📁 Configuration Structure

```
configs/
├── panels/          # Hosting control panels
├── databases/       # Database servers
├── webservers/      # Web servers
├── security/        # Security tools
├── services/        # Container/services
├── monitoring/      # Monitoring tools
└── vpn/            # VPN solutions
```

## 📝 Package Configuration Format

### Basic Structure

```yaml
id: unique-id                    # Unique identifier (lowercase, no spaces)
name: "Package Name"            # Display name
category: category              # Category folder name
description: "Description"      # Short description
icon: "🎨"                      # Emoji icon

version:
  current: "1.0"               # Current version
  source:
    type: "github_release"     # github_release, url_scrape, or static
    owner: "username"          # GitHub owner (for github_release)
    repo: "repository"         # GitHub repo (for github_release)

requirements:
  os:
    ubuntu: ["20.04", "22.04"]  # Supported Ubuntu versions
    debian: ["11", "12"]        # Supported Debian versions
  min_ram: "1GB"                # Minimum RAM
  min_disk: "5GB"               # Minimum disk space
  install_size: "2GB"           # Estimated install size
  estimated_time: "3-5 minutes" # Estimated installation time
  ports: [80, 443]              # Required ports
  architecture: ["amd64"]       # Supported architectures

install:
  method: "shell"               # shell, apt, or yum
  script: |                     # Installation script
    # Commands here
  pre_check:                    # Pre-installation checks
    - command: "which wget"
      error_msg: "wget is required"
  post_check:                   # Post-installation checks
    - service: "nginx"
      port: 80
      timeout: 60

dependencies:
  required:                     # Required dependencies
    - package1
    - package2
  optional:                     # Optional dependencies
    - package3

config:
  files:                        # Configuration files
    - path: "/etc/nginx/nginx.conf"
      description: "Main config"
      template: "nginx.conf"

uninstall:
  command: |                    # Uninstallation commands
    apt-get remove -y nginx

license:
  required: false               # License required?
  type: "free"                  # License type

support:
  docs: "https://docs.example.com"
  forum: "https://forum.example.com"
  issues: "https://github.com/issues"
```

## 🎯 Complete Example

### Nginx Configuration

```yaml
id: nginx
name: "Nginx"
category: webserver
description: "High-performance web server and reverse proxy"
icon: "🌐"

version:
  current: "1.25"
  source:
    type: "github_release"
    owner: "nginx"
    repo: "nginx"

requirements:
  os:
    ubuntu: ["20.04", "22.04", "24.04"]
    debian: ["11", "12"]
  install_size: "50MB"
  estimated_time: "30 seconds"
  ports: [80, 443]

install:
  ubuntu:
    method: "apt"
    packages:
      - nginx
      - nginx-extras
  centos:
    method: "yum"
    packages:
      - nginx

config:
  files:
    - path: "/etc/nginx/nginx.conf"
      description: "Main configuration"
    - path: "/etc/nginx/sites-available/default"
      description: "Default site"

commands:
  start: "systemctl start nginx"
  stop: "systemctl stop nginx"
  reload: "nginx -s reload"
  test: "nginx -t"

uninstall:
  command: |
    systemctl stop nginx
    apt-get remove -y nginx

license:
  required: false
  type: "free"

support:
  docs: "https://nginx.org/en/docs/"
```

## 🔧 Custom Package Example

Create your own package:

```yaml
# configs/custom/myapp.yaml
id: myapp
name: "My Application"
category: custom
description: "My custom application"
icon: "🚀"

version:
  current: "2.0"
  source:
    type: "static"

requirements:
  min_ram: "512MB"
  min_disk: "1GB"
  ports: [8080]

install:
  method: "shell"
  script: |
    #!/bin/bash
    wget https://example.com/myapp.tar.gz
    tar xzf myapp.tar.gz
    cd myapp
    ./install.sh

uninstall:
  command: |
    rm -rf /opt/myapp
```

## 📊 Version Detection Methods

### 1. GitHub Release (Auto-update)

```yaml
version:
  source:
    type: "github_release"
    owner: "nodejs"
    repo: "node"
```

### 2. URL Scraping

```yaml
version:
  source:
    type: "url_scrape"
    url: "https://nginx.org/en/download.html"
    regex: "nginx-(\\d+\\.\\d+\\.\\d+)"
```

### 3. Static Version

```yaml
version:
  current: "1.0"
  source:
    type: "static"
```

## 🐛 Troubleshooting

### Package Not Appearing

1. Check YAML syntax: `yamllint configs/category/package.yaml`
2. Verify file extension: `.yaml`
3. Check file location: Correct category folder
4. Verify unique ID: No duplicates allowed

### Installation Fails

1. Check `pre_check` commands
2. Verify OS compatibility
3. Check port availability
4. Review logs: `/var/log/hostkit/`

## 🎨 Best Practices

1. **Use descriptive IDs** - `myapp-v2` not `ma2`
2. **Provide accurate sizes** - Helps users plan resources
3. **Include all OS variants** - Ubuntu, Debian support
4. **Test uninstall** - Ensure clean removal
5. **Add post_check** - Verify service started
6. **Document ports** - Avoid conflicts
7. **Use specific versions** - Avoid "latest" tags

## 🔐 Security Considerations

- Never include passwords in YAML files
- Use environment variables for sensitive data
- Validate all user inputs
- Run commands with least privilege
- Test in isolated environment first

## 📚 More Examples

See existing configs in the `configs/` directory for more examples covering:
- Complex installations
- Multi-step scripts
- Database configurations
- Service management