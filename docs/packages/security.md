# 🔒 Security Tools

HostKit includes essential security tools to protect your server.

## Overview

| Tool | Purpose | Protection Level |
|------|---------|------------------|
| **Fail2ban** | Intrusion prevention | High |
| **CSF** | Firewall management | High |
| **Certbot** | SSL certificates | Medium |

---

## Fail2ban

Intrusion prevention framework that protects against brute-force attacks.

### Features
- Monitors log files for suspicious activity
- Automatically bans malicious IP addresses
- Supports multiple services (SSH, Nginx, Apache, etc.)
- Configurable ban duration and retry limits
- Email notifications

### Install
```bash
hostkit install fail2ban
```

### How It Works
1. Monitors log files for failed login attempts
2. Tracks IP addresses with multiple failures
3. Automatically blocks IPs exceeding threshold
4. Unblocks after specified time

### Configuration Files
- Main config: `/etc/fail2ban/jail.conf`
- Local config: `/etc/fail2ban/jail.local`
- Filters: `/etc/fail2ban/filter.d/`
- Actions: `/etc/fail2ban/action.d/`

### Default Settings
- Max retry: 5 attempts
- Find time: 10 minutes
- Ban time: 10 minutes
- Monitors: SSH, Nginx, Apache

### Basic Commands
```bash
# Start Fail2ban
sudo systemctl start fail2ban

# Stop Fail2ban
sudo systemctl stop fail2ban

# Check status
sudo systemctl status fail2ban

# View active jails
sudo fail2ban-client status

# View specific jail status
sudo fail2ban-client status nginx

# List banned IPs
sudo fail2ban-client status sshd | grep "Banned IP list"

# Unban an IP
sudo fail2ban-client set sshd unbanip 192.168.1.100
```

### Enable Specific Jails
Edit `/etc/fail2ban/jail.local`:

```ini
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3

[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/error.log

[nginx-botsearch]
enabled = true
filter = nginx-botsearch
port = http,https
logpath = /var/log/nginx/access.log
maxretry = 2
```

### Custom Filter Example
Create `/etc/fail2ban/filter.d/custom-app.conf`:

```ini
[Definition]
failregex = ^.*Failed login from <HOST>.*$
            ^.*Authentication failed for .* from <HOST>.*$
ignoreregex = 
```

---

## ConfigServer Security & Firewall (CSF)

Advanced firewall configuration script and intrusion detection.

### Features
- Stateful packet inspection firewall
- Login failure detection
- Process tracking
- Directory watching
- Messenger service
- Port flood protection
- Email alerts

### Install
```bash
hostkit install csf
```

### Configuration Files
- Main config: `/etc/csf/csf.conf`
- Allow list: `/etc/csf/csf.allow`
- Deny list: `/etc/csf/csf.deny`
- Ignore list: `/etc/csf/csf.ignore`

### Basic Commands
```bash
# Start CSF
sudo csf -s

# Stop CSF (flush rules)
sudo csf -f

# Restart CSF
sudo csf -r

# Allow IP
sudo csf -a 192.168.1.100

# Deny IP
sudo csf -d 192.168.1.100

# Remove IP from allow list
sudo csf -ar 192.168.1.100

# Remove IP from deny list
sudo csf -dr 192.168.1.100

# Check CSF status
sudo csf -l

# View temporary bans
sudo csf -t

# Check if IP is blocked
sudo csf -g 192.168.1.100

# Enable CSF testing mode
sudo csf -e

# Disable CSF
sudo csf -x
```

### Common Configuration
Edit `/etc/csf/csf.conf`:

```ini
# Testing mode (set to 0 when ready)
TESTING = "0"

# Default port restrictions
TCP_IN = "20,21,22,25,53,80,110,143,443,465,587,993,995,7080,8080,8088"
TCP_OUT = "20,21,22,25,53,80,110,113,443,587,993,995,8080"
UDP_IN = "20,21,53,80,443"
UDP_OUT = "20,21,53,113,123,443"

# Email alerts
LF_ALERT_TO = "admin@example.com"
LF_ALERT_FROM = "csf@example.com"

# Brute force detection
LF_SSHD = "5"
LF_FTPD = "10"
LF_SMTPAUTH = "10"
LF_POP3D = "10"
LF_IMAPD = "10"
LF_HTACCESS = "10"
LF_NGINX = "10"

# Block duration (seconds)
LF_BLOCK_PERMANENT = "0"
LF_BLOCK_TIME = "3600"
```

### Port Management
```bash
# Allow incoming port
sudo csf --port-add 8080 tcp

# Remove allowed port
sudo csf --port-del 8080 tcp
```

---

## Certbot

Free SSL/TLS certificates from Let's Encrypt.

### Features
- Free SSL certificates
- Automatic renewal
- Wildcard certificates
- Multiple domain support
- Web server integration

### Install
```bash
hostkit install certbot
```

### Using with Nginx
```bash
# Get certificate and auto-configure Nginx
sudo certbot --nginx -d example.com -d www.example.com

# Test automatic renewal
sudo certbot renew --dry-run
```

### Using with Standalone
```bash
# Get certificate (standalone mode)
sudo certbot certonly --standalone -d example.com

# Get wildcard certificate
sudo certbot certonly --manual --preferred-challenges dns -d *.example.com -d example.com
```

### Certificate Locations
- Certificates: `/etc/letsencrypt/live/example.com/`
- Private key: `/etc/letsencrypt/live/example.com/privkey.pem`
- Certificate: `/etc/letsencrypt/live/example.com/cert.pem`
- Full chain: `/etc/letsencrypt/live/example.com/fullchain.pem`

### Automatic Renewal
Certbot sets up automatic renewal via systemd timer. Verify:

```bash
# Check renewal timer
sudo systemctl status certbot.timer

# View renewal log
sudo cat /var/log/letsencrypt/letsencrypt.log

# Test renewal
sudo certbot renew --dry-run
```

### Nginx SSL Configuration
```nginx
server {
    listen 443 ssl http2;
    server_name example.com;

    ssl_certificate /etc/letsencrypt/live/example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;
    ssl_trusted_certificate /etc/letsencrypt/live/example.com/chain.pem;

    # SSL settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    root /var/www/example.com;
    index index.html;
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name example.com;
    return 301 https://$server_name$request_uri;
}
```

### Revoke Certificate
```bash
sudo certbot revoke --cert-path /etc/letsencrypt/live/example.com/cert.pem
```

---

## Security Best Practices

### 1. Keep Software Updated
```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Use Strong Passwords
- Minimum 12 characters
- Mix of uppercase, lowercase, numbers, symbols
- No dictionary words

### 3. Disable Root Login
Edit `/etc/ssh/sshd_config`:
```
PermitRootLogin no
PasswordAuthentication no  # Use SSH keys instead
```

### 4. Configure Firewall Properly
- Only open necessary ports
- Use deny-all policy by default
- Regularly review rules

### 5. Monitor Logs
```bash
# View auth logs
sudo tail -f /var/log/auth.log

# View Nginx access logs
sudo tail -f /var/log/nginx/access.log

# View Nginx error logs
sudo tail -f /var/log/nginx/error.log
```

### 6. Regular Backups
- Back up configurations daily
- Test restoration process
- Store backups off-site

### 7. SSL/TLS Everywhere
- Use HTTPS for all web services
- Enable HSTS
- Use strong cipher suites

---

## Security Checklist

- [x] Fail2ban installed and configured
- [x] Firewall (CSF or UFW) enabled
- [x] SSL certificates installed
- [x] SSH key authentication enabled
- [x] Root login disabled
- [x] Automatic security updates enabled
- [x] Regular backups configured
- [x] Log monitoring active
- [x] Unnecessary services disabled
- [x] File permissions properly set
