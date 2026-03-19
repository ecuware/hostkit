# 🎛️ Hosting Panels

HostKit supports 7 popular hosting control panels.

## Overview

| Panel | Best For | Difficulty | Community |
|-------|----------|------------|-----------|
| **aaPanel** | Beginners | Easy | Large |
| **HestiaCP** | Open source lovers | Easy | Medium |
| **CloudPanel** | PHP developers | Easy | Medium |
| **CyberPanel** | OpenLiteSpeed | Medium | Medium |
| **Plesk** | Professionals | Medium | Large |
| **cPanel/WHM** | Enterprise | Hard | Large |
| **CWP** | AlmaLinux users | Medium | Small |

---

## aaPanel

Simple and lightweight control panel with modern UI.

### Features
- One-click LNMP/LAMP deployment
- SSL certificate management
- Database management
- File manager
- FTP server

### Requirements
- Ubuntu 18.04+ / Debian 10+
- 1GB RAM minimum
- 2GB disk space

### Install
```bash
hostkit install aapanel
```

### Default Access
- URL: `http://YOUR_IP:7800`
- Username: Set during installation
- Password: Set during installation

---

## HestiaCP

Open source hosting panel, forked from VestaCP.

### Features
- Web domains management
- DNS management
- Mail server
- Database management
- Firewall configuration

### Requirements
- Ubuntu 20.04+ / Debian 11+
- 1GB RAM minimum
- 1GB disk space

### Install
```bash
hostkit install hestia
```

### Default Access
- URL: `https://YOUR_IP:8083`
- Username: `admin`
- Password: Set during installation

---

## CloudPanel

Modern PHP hosting panel optimized for performance.

### Features
- PHP 7.x/8.x management
- NGINX reverse proxy
- Redis/Memcached integration
- Let's Encrypt SSL
- File manager

### Requirements
- Ubuntu 22.04+ / Debian 11+
- 1GB RAM minimum
- 1GB disk space

### Install
```bash
hostkit install cloudpanel
```

### Default Access
- URL: `https://YOUR_IP:8443`
- Username: Set during installation
- Password: Set during installation

---

## CyberPanel

Feature-rich panel powered by OpenLiteSpeed web server.

### Features
- OpenLiteSpeed web server
- LSCache integration
- Email server
- FTP server
- DNS management

### Requirements
- Ubuntu 20.04+ (AlmaLinux recommended)
- 2GB RAM minimum
- 3GB disk space

### Install
```bash
hostkit install cyberpanel
```

### Default Access
- URL: `https://YOUR_IP:8090`
- Username: `admin`
- Password: Set during installation

---

## Plesk

Professional hosting management platform.

### Features
- Multi-domain management
- WordPress toolkit
- Git integration
- Docker support
- Extensions marketplace

### Requirements
- Ubuntu 20.04+ / Debian 11+
- 2GB RAM minimum (4GB recommended)
- 3GB disk space

### Install
```bash
hostkit install plesk
```

### Default Access
- URL: `https://YOUR_IP:8443`
- Username: `admin`
- Password: Set during installation

!!! note "License Required"
    Plesk requires a license after trial period. You can use free developer license for limited domains.

---

## cPanel/WHM

Industry-standard hosting panel for enterprise.

### Features
- Complete hosting automation
- DNS clustering
- Email server
- Security tools
- Softaculous integration

### Requirements
- AlmaLinux 8+ only
- 4GB RAM minimum (8GB recommended)
- 5GB disk space

### Install
```bash
hostkit install cpanel
```

### Default Access
- WHM: `https://YOUR_IP:2087`
- cPanel: `https://YOUR_IP:2083`
- Username: `root` (WHM)
- Password: Your root password

!!! warning "AlmaLinux Only"
    cPanel only supports AlmaLinux. Ubuntu and Debian are NOT supported.

!!! note "License Required"
    cPanel requires a paid license. Trial available for testing.

---

## CentOS Web Panel (CWP)

Free Linux control panel for AlmaLinux.

### Features
- Apache/Nginx web server
- MySQL/MariaDB/PostgreSQL
- Email server
- Firewall management
- Backup system

### Requirements
- AlmaLinux 8+ only
- 1GB RAM minimum
- 2GB disk space

### Install
```bash
hostkit install cwp
```

### Default Access
- URL: `http://YOUR_IP:2030`
- Username: `root`
- Password: Your root password

!!! warning "AlmaLinux Only"
    CWP only supports AlmaLinux. Ubuntu and Debian are NOT supported.

---

## Comparison

| Feature | aaPanel | HestiaCP | CloudPanel | CyberPanel | Plesk | cPanel | CWP |
|---------|---------|----------|------------|------------|-------|--------|-----|
| Free | ✅ | ✅ | ✅ | ✅ | Trial | Trial | ✅ |
| Ubuntu | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Debian | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ |
| AlmaLinux | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ |
| Nginx | ✅ | ✅ | ✅ | ❌* | ✅ | ✅ | ✅ |
| Apache | ✅ | ✅ | ❌ | ❌ | ✅ | ✅ | ✅ |
| OpenLiteSpeed | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Email Server | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ✅ |
| DNS Server | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ✅ |
| Database UI | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

*CyberPanel uses OpenLiteSpeed instead of Nginx

---

## Choosing the Right Panel

### For Beginners
1. **aaPanel** - Easiest to use
2. **HestiaCP** - Clean interface

### For Performance
1. **CloudPanel** - Optimized for PHP
2. **CyberPanel** - OpenLiteSpeed speed

### For Enterprise
1. **cPanel/WHM** - Industry standard
2. **Plesk** - Professional features

### For Open Source
1. **HestiaCP** - Community driven
2. **aaPanel** - Free forever
