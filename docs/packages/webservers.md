# 🌐 Web Servers

HostKit supports 3 popular web servers for hosting your applications.

## Overview

| Server | Type | Best For | Performance |
|--------|------|----------|-------------|
| **Nginx** | Reverse Proxy, Static | High traffic, microservices | Excellent |
| **OpenLiteSpeed** | Web Server + Cache | WordPress, dynamic content | Excellent |
| **PHP** | Language Runtime | PHP applications | Good |

---

## Nginx

High-performance web server and reverse proxy.

### Features
- Event-driven architecture
- Low memory footprint
- Reverse proxy & load balancing
- SSL/TLS termination
- HTTP/2 support
- Static file serving

### Install
```bash
hostkit install nginx
```

### Version
Latest stable version from official Nginx repository.

### Configuration Files
- Main config: `/etc/nginx/nginx.conf`
- Sites available: `/etc/nginx/sites-available/`
- Sites enabled: `/etc/nginx/sites-enabled/`

### Default Configuration
- Port: 80 (HTTP)
- Document root: `/var/www/html`
- User: www-data

### Basic Commands
```bash
# Start Nginx
sudo systemctl start nginx

# Stop Nginx
sudo systemctl stop nginx

# Restart Nginx
sudo systemctl restart nginx

# Reload configuration
sudo systemctl reload nginx

# Test configuration
sudo nginx -t

# Check status
sudo systemctl status nginx
```

### Virtual Host Example
```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/example.com;
    index index.html index.php;

    location / {
        try_files $uri $uri/ =404;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
    }
}
```

### With PHP-FPM
```bash
# Install PHP-FPM
hostkit install php

# Configure Nginx to use PHP-FPM
# Add to your server block:
location ~ \.php$ {
    include snippets/fastcgi-php.conf;
    fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
}
```

---

## OpenLiteSpeed

High-performance, lightweight web server with built-in cache.

### Features
- Event-driven architecture
- Built-in LSCache
- HTTP/3 support
- WebAdmin GUI
- .htaccess compatible
- PHP suExec

### Install
```bash
hostkit install litespeed
```

### Version
Latest stable version from OpenLiteSpeed repository.

### Configuration Files
- Main config: `/usr/local/lsws/conf/httpd_config.conf`
- Virtual hosts: `/usr/local/lsws/conf/vhosts/`

### Default Configuration
- HTTP Port: 8088
- Admin Port: 7080
- Document root: `/usr/local/lsws/Example/html`

### WebAdmin Console
```
URL: https://YOUR_IP:7080
Username: admin
Password: Set during installation
```

### Basic Commands
```bash
# Start OpenLiteSpeed
sudo systemctl start lsws

# Stop OpenLiteSpeed
sudo systemctl stop lsws

# Restart OpenLiteSpeed
sudo systemctl restart lsws

# Check status
sudo systemctl status lsws
```

### Create Virtual Host
1. Login to WebAdmin Console
2. Go to Configuration → Virtual Hosts
3. Click "Add"
4. Configure:
   - Virtual Host Name: `example.com`
   - Virtual Host Root: `/var/www/example.com`
   - Config File: `$SERVER_ROOT/conf/vhosts/example.com/vhost.conf`

### Enable LSCache
LSCache is built-in and automatically available for WordPress and other supported applications.

---

## PHP

Multi-version PHP with FPM support for running PHP applications.

### Features
- Multiple PHP versions (7.4, 8.0, 8.1, 8.2, 8.3)
- PHP-FPM for better performance
- Common extensions included
- Composer included
- OPCache enabled

### Install
```bash
hostkit install php
```

### Installed Versions
By default, the latest stable PHP version is installed. Multiple versions can be installed side-by-side.

### PHP-FPM Pools
- Main pool: `/etc/php/8.1/fpm/pool.d/www.conf`
- Socket: `/var/run/php/php8.1-fpm.sock`
- Port: 9000 (optional)

### Common Extensions
- mysqli, pdo_mysql - MySQL/MariaDB support
- pgsql, pdo_pgsql - PostgreSQL support
- mongodb - MongoDB support
- redis - Redis support
- gd - Image processing
- curl - HTTP requests
- mbstring - Multibyte strings
- xml - XML processing
- zip - ZIP archives
- intl - Internationalization
- bcmath - Precision math
- opcache - Bytecode cache

### Basic Commands
```bash
# Start PHP-FPM
sudo systemctl start php8.1-fpm

# Stop PHP-FPM
sudo systemctl stop php8.1-fpm

# Restart PHP-FPM
sudo systemctl restart php8.1-fpm

# Check status
sudo systemctl status php8.1-fpm

# View PHP version
php -v

# View installed extensions
php -m

# Check PHP configuration
php --ini
```

### Switch PHP Version
```bash
# List available versions
update-alternatives --list php

# Set default version
sudo update-alternatives --config php
```

### Install Additional Extensions
```bash
# Example: Install Redis extension
sudo apt install php-redis
sudo systemctl restart php8.1-fpm
```

### PHP Configuration
```bash
# Edit PHP configuration
sudo nano /etc/php/8.1/fpm/php.ini

# Important settings to modify:
memory_limit = 256M
upload_max_filesize = 64M
post_max_size = 64M
max_execution_time = 300
```

### Composer
Composer is included with PHP installation:

```bash
# Check Composer version
composer --version

# Create new project
composer create-project laravel/laravel myapp

# Install dependencies
composer install

# Update dependencies
composer update
```

---

## Comparison

| Feature | Nginx | OpenLiteSpeed | Apache* |
|---------|-------|---------------|---------|
| Static Files | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| Dynamic Content | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| Memory Usage | Low | Low | Medium |
| Configuration | Moderate | Easy | Easy |
| .htaccess | ❌ | ✅ | ✅ |
| Built-in Cache | ❌ | ✅ (LSCache) | ❌ |
| HTTP/3 | Via module | ✅ | Via module |

*Apache is not included in HostKit but can be installed through panels like HestiaCP

---

## Choosing the Right Web Server

### For Static Sites
1. **Nginx** - Best for static content
2. **OpenLiteSpeed** - Good alternative

### For WordPress
1. **OpenLiteSpeed** - Best with LSCache
2. **Nginx** + PHP-FPM - Fast and reliable

### For PHP Applications
1. **Nginx** + PHP-FPM - Industry standard
2. **OpenLiteSpeed** - Excellent performance

### For Learning/Ease
1. **OpenLiteSpeed** - WebAdmin GUI
2. **Nginx** - Simple configuration

---

## Combining Web Servers

### Nginx + PHP-FPM
Best combination for most PHP applications:

```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/example.com;
    index index.php index.html;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

### Nginx as Reverse Proxy
Use Nginx in front of other services:

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

---

## SSL/TLS with Let's Encrypt

Install Certbot to automatically get SSL certificates:

```bash
# Install Certbot
hostkit install certbot

# Get certificate for domain
sudo certbot --nginx -d example.com -d www.example.com

# For OpenLiteSpeed
sudo certbot certonly --standalone -d example.com
```

---

## Performance Optimization

### Nginx
```nginx
# In nginx.conf
worker_processes auto;
worker_connections 4096;

gzip on;
gzip_vary on;
gzip_proxied any;
gzip_comp_level 6;
gzip_types text/plain text/css text/xml application/json application/javascript application/rss+xml application/atom+xml image/svg+xml;
```

### PHP-FPM
```ini
; In php-fpm.conf
pm = dynamic
pm.max_children = 50
pm.start_servers = 5
pm.min_spare_servers = 5
pm.max_spare_servers = 35
pm.max_requests = 500
```

### OpenLiteSpeed
Use LSCache plugin for your application (WordPress, Magento, etc.) for best performance.
