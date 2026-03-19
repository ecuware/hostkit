# Creating Packages for HostKit

This guide explains how to add new packages to HostKit. Packages are the heart of HostKit - they define how software is installed, configured, and managed.

## Table of Contents

- [Quick Start](#quick-start)
- [Package Structure](#package-structure)
- [Field Reference](#field-reference)
- [Installation Methods](#installation-methods)
- [Examples](#examples)
- [Testing](#testing)
- [Best Practices](#best-practices)

## Quick Start

The easiest way to create a package:

1. Copy an existing package as template:
   ```bash
   cp configs/webservers/nginx.yaml configs/services/myapp.yaml
   ```

2. Edit the file with your app details

3. Test locally:
   ```bash
   go build -o hostkit ./cmd/hostkit/main.go
   ./hostkit list  # See if it appears
   ./hostkit install myapp  # Test installation
   ```

4. Submit a pull request

## Package Structure

```yaml
# 1. Basic Information
id: unique_identifier
name: "Display Name"
category: category_name
description: "Brief description"
icon: "🎨"

# 2. Version Information
version:
  current: "1.0.0"
  source:
    type: github_release
    owner: username
    repo: repository

# 3. Requirements
requirements:
  os:
    ubuntu: ["20.04", "22.04"]
    debian: ["11", "12"]
  min_ram: "1GB"
  min_disk: "5GB"
  install_size: "2GB"
  estimated_time: "3-5 minutes"
  ports: [8080, 443]

# 4. Dependencies
dependencies:
  - nginx
  - mariadb

conflicts:
  - apache2

# 5. Installation
install:
  method: shell
  script: |
    #!/bin/bash
    apt-get update
    apt-get install -y myapp

# 6. Post-Install Instructions
post_install:
  message: |
    Installation complete!
    Access: http://localhost:8080
  show_logs: true
```

## Field Reference

### Required Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `id` | string | Unique identifier | `nginx` |
| `name` | string | Display name | `Nginx` |
| `category` | string | Category folder | `webserver` |
| `description` | string | Brief description | `High-performance web server` |

### Version Section

```yaml
version:
  current: "1.0.0"  # Current known version
  source:
    type: github_release  # Detection method
    owner: nginx          # GitHub username
    repo: nginx           # Repository name
```

**Version Source Types:**

**1. GitHub Releases** (Recommended)
```yaml
version:
  source:
    type: github_release
    owner: nginx
    repo: nginx
    asset_pattern: "nginx-(.*)\\.tar\\.gz"
```

**2. Static Version**
```yaml
version:
  source:
    type: static
    version: "1.25.3"
```

**3. API Endpoint**
```yaml
version:
  source:
    type: api
    url: "https://api.github.com/repos/nginx/nginx/releases/latest"
    json_path: "tag_name"
```

**4. URL Scraping**
```yaml
version:
  source:
    type: url_scrape
    url: "https://nginx.org/en/download.html"
    regex: "Stable version.*?(\\d+\\.\\d+\\.\\d+)"
```

### Requirements Section

```yaml
requirements:
  os:
    ubuntu: ["20.04", "22.04", "24.04"]
    debian: ["11", "12"]
    almalinux: ["8", "9"]
  
  min_ram: "1GB"           # Minimum RAM
  min_disk: "5GB"          # Free space needed
  install_size: "2GB"      # After installation
  estimated_time: "5 min"  # Human-readable
  
  ports: [80, 443]         # Required ports
```

### Dependencies

```yaml
dependencies:
  - nginx      # Install nginx first
  - mariadb    # Then mariadb
  - php        # Then php
```

Installation order follows the dependency list.

### Conflicts

```yaml
conflicts:
  - apache2    # Can't coexist with apache2
  - lighttpd   # Or lighttpd
```

HostKit prevents installation if conflicts exist.

## Installation Methods

### Method 1: Shell Script (Most Common)

```yaml
install:
  method: shell
  script: |
    #!/bin/bash
    set -e  # Exit on error
    
    echo "Installing MyApp..."
    
    # Update package list
    apt-get update
    
    # Install dependencies
    apt-get install -y dependency1 dependency2
    
    # Download application
    wget https://example.com/myapp-latest.deb
    dpkg -i myapp-latest.deb
    
    # Fix any dependency issues
    apt-get install -f -y
    
    # Enable service
    systemctl enable myapp
    systemctl start myapp
    
    # Configure firewall
    ufw allow 8080/tcp
    
    echo "Installation complete!"
```

**Best Practices:**
- Use `set -e` to exit on errors
- Use `sudo` for commands that need it
- Provide clear output messages
- Handle errors gracefully

### Method 2: Remote Script

```yaml
install:
  method: script
  url: "https://get.myapp.com/install.sh"
  checksum: "sha256:abc123..."
```

HostKit downloads and executes the script.

### Method 3: Package Manager

```yaml
install:
  method: package
  packages:
    - myapp
    - myapp-extra
  repository:
    name: myapp-repo
    url: "https://repo.myapp.com/ubuntu"
    key: "https://repo.myapp.com/key.asc"
```

### Method 4: Docker

```yaml
install:
  method: docker
  image: "myapp:latest"
  ports:
    - "8080:8080"
  volumes:
    - "./data:/data"
  environment:
    - "KEY=value"
```

## Examples

### Example 1: Simple Web Application

```yaml
id: mywebapp
name: "My Web App"
category: services
description: "Simple web application with built-in server"
icon: "🌐"

version:
  current: "2.0.0"
  source:
    type: github_release
    owner: myuser
    repo: mywebapp

requirements:
  os:
    ubuntu: ["20.04", "22.04"]
    debian: ["11", "12"]
  min_ram: "512MB"
  min_disk: "1GB"
  ports: [3000]

install:
  method: shell
  script: |
    #!/bin/bash
    set -e
    
    # Create user
    useradd -r -s /bin/false mywebapp || true
    
    # Download latest release
    cd /opt
    wget https://github.com/myuser/mywebapp/releases/latest/download/mywebapp-linux-amd64.tar.gz
    tar -xzf mywebapp-linux-amd64.tar.gz
    rm mywebapp-linux-amd64.tar.gz
    
    # Set permissions
    chown -R mywebapp:mywebapp /opt/mywebapp
    
    # Create systemd service
    cat > /etc/systemd/system/mywebapp.service <<EOF
    [Unit]
    Description=My Web App
    After=network.target
    
    [Service]
    Type=simple
    User=mywebapp
    WorkingDirectory=/opt/mywebapp
    ExecStart=/opt/mywebapp/mywebapp
    Restart=always
    
    [Install]
    WantedBy=multi-user.target
    EOF
    
    # Start service
    systemctl daemon-reload
    systemctl enable mywebapp
    systemctl start mywebapp
    
    # Open port
    ufw allow 3000/tcp || true

dependencies:
  - nginx  # If using nginx as reverse proxy

post_install:
  message: |
    My Web App installed successfully!
    
    Access: http://YOUR_IP:3000
    
    Service commands:
    - systemctl status mywebapp
    - systemctl restart mywebapp
    - journalctl -u mywebapp -f
```

### Example 2: Database with Configuration

```yaml
id: mydatabase
name: "My Database"
category: databases
description: "High-performance database server"
icon: "🗄️"

version:
  current: "1.5.0"
  source:
    type: static
    version: "1.5.0"

requirements:
  min_ram: "2GB"
  min_disk: "10GB"
  ports: [5432]

install:
  method: shell
  script: |
    #!/bin/bash
    set -e
    
    # Install from official repo
    curl -fsSL https://mydb.io/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/mydb.gpg
    echo "deb [signed-by=/usr/share/keyrings/mydb.gpg] https://mydb.io/apt stable main" > /etc/apt/sources.list.d/mydb.list
    
    apt-get update
    apt-get install -y mydatabase
    
    # Generate random password
    DB_PASSWORD=$(openssl rand -base64 32)
    
    # Configure
    cat > /etc/mydatabase/mydatabase.conf <<EOF
    [server]
    port = 5432
    max_connections = 100
    
    [security]
    authentication = md5
    ssl = true
    EOF
    
    # Create initial user
    mydb-setup --user=admin --password="$DB_PASSWORD"
    
    # Save password to file
    echo "admin:$DB_PASSWORD" > /root/.mydb_credentials
    chmod 600 /root/.mydb_credentials
    
    systemctl enable mydatabase
    systemctl start mydatabase

post_install:
  message: |
    Database installed!
    
    Port: 5432
    Credentials saved: /root/.mydb_credentials
    
    Connection string:
    mydb://admin:PASSWORD@localhost:5432
    
    IMPORTANT: Change the default password!
```

### Example 3: Multi-OS Support

```yaml
id: universalapp
name: "Universal App"
category: services

version:
  current: "3.0.0"
  source:
    type: github_release
    owner: universal
    repo: app

requirements:
  os:
    ubuntu: ["20.04", "22.04"]
    debian: ["11", "12"]
    almalinux: ["8", "9"]

install:
  ubuntu:
    method: shell
    script: |
      apt-get update
      apt-get install -y universal-app
  
  debian:
    method: shell
    script: |
      apt-get update
      apt-get install -y universal-app
  
  almalinux:
    method: shell
    script: |
      dnf install -y universal-app
```

## Testing

### Local Testing

1. **Build HostKit**:
   ```bash
   go build -o hostkit ./cmd/hostkit/main.go
   ```

2. **Check if package appears**:
   ```bash
   ./hostkit list | grep mypackage
   ```

3. **Test installation**:
   ```bash
   sudo ./hostkit install mypackage
   ```

4. **Verify installation**:
   ```bash
   which mypackage
   systemctl status mypackage
   ```

### Test Scenarios

Test your package on:
- [ ] Fresh Ubuntu 20.04
- [ ] Fresh Ubuntu 22.04
- [ ] Fresh Debian 11
- [ ] With dependencies already installed
- [ ] Without dependencies (should auto-install)
- [ ] Low resources (minimum RAM/disk)
- [ ] Re-installation (idempotent)

### Automated Testing

Create a test script:

```bash
#!/bin/bash
# test-mypackage.sh

set -e

echo "Testing MyPackage installation..."

# Clean environment
docker run --rm -v $(pwd):/hostkit ubuntu:22.04 bash -c "
  cd /hostkit
  apt-get update
  apt-get install -y golang-go
  go build -o hostkit ./cmd/hostkit/main.go
  ./hostkit install mypackage
"

echo "Test passed!"
```

## Best Practices

### Do's ✅

1. **Use `set -e`** - Exit on errors
2. **Check if already installed** - Make idempotent
3. **Clear error messages** - Help users debug
4. **Open necessary ports** - Configure firewall
5. **Enable services** - Start on boot
6. **Set proper permissions** - Security first
7. **Test thoroughly** - Multiple scenarios
8. **Document well** - Clear post-install instructions

### Don'ts ❌

1. **Don't assume root** - Check and use sudo explicitly
2. **Don't hardcode versions** - Use version detection
3. **Don't ignore errors** - Handle failures gracefully
4. **Don't leave temp files** - Clean up after installation
5. **Don't use interactive commands** - Use `-y` flags
6. **Don't modify unrelated configs** - Stay in your scope

### Security Checklist

- [ ] No hardcoded passwords
- [ ] Generate random credentials
- [ ] Restrict file permissions (600 for sensitive files)
- [ ] Don't log sensitive data
- [ ] Validate input
- [ ] Use HTTPS for downloads
- [ ] Verify checksums when possible

### Performance Tips

1. **Minimize apt-get update** - Do it once at the start
2. **Parallel downloads** - Use `apt-get install -y pkg1 pkg2 pkg3`
3. **Clean up** - Remove temp files
4. **Use caches** - Don't re-download

## Troubleshooting

### Package Not Appearing in List

1. Check YAML syntax: `yamllint configs/mycategory/myapp.yaml`
2. Verify file extension: Must be `.yaml`
3. Check file location: Must be in `configs/<category>/`
4. Restart HostKit: Changes are loaded at startup

### Installation Fails

1. Check script syntax locally
2. Test on fresh VM
3. Review logs: `/var/log/hostkit/`
4. Check dependencies are available

### Debugging

Add debug output:
```yaml
install:
  method: shell
  script: |
    #!/bin/bash
    set -e
    set -x  # Enable debug mode
    
    echo "DEBUG: Starting installation"
    echo "DEBUG: Current user: $(whoami)"
    echo "DEBUG: Working directory: $(pwd)"
    
    # Your installation code
```

## Submitting Your Package

1. **Fork the repository**
2. **Add your package** to appropriate category
3. **Test thoroughly**
4. **Update documentation** if needed
5. **Submit pull request** with:
   - Clear description
   - Test results
   - Screenshots if applicable

See [CONTRIBUTING.md](../CONTRIBUTING.md) for full details.

## Need Help?

- Check existing packages for examples
- Open an issue for discussion
- Join our community

Happy packaging! 🎉
