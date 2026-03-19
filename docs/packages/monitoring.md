# 📊 Monitoring

Real-time system and service monitoring tools.

## Overview

| Tool | Type | Best For |
|------|------|----------|
| **Netdata** | System & Apps | Real-time monitoring, alerts |

---

## Netdata

High-performance, real-time infrastructure monitoring.

### Features
- Real-time metrics (per second)
- 1000+ metrics out of the box
- Interactive visualizations
- Powerful alerting
- Distributed architecture
- Low resource usage
- Zero configuration

### Install
```bash
hostkit install netdata
```

### Access
```
URL: http://YOUR_IP:19999
```

### What Gets Monitored

#### System Metrics
- CPU usage (per core, total)
- Memory (RAM, swap, available)
- Disk I/O (per disk)
- Network (per interface)
- Load average
- Processes
- System uptime

#### Application Metrics
- Web servers (Nginx, Apache)
- Databases (MySQL, PostgreSQL)
- Containers (Docker)
- Message queues
- And many more...

### Dashboard Overview

#### Main Dashboard
1. **System Overview** - CPU, RAM, Disk, Network
2. **Applications** - Running apps and their resources
3. **Network** - Bandwidth usage, connections
4. **Disks** - I/O operations, space usage

#### Navigation
- Left sidebar: All metrics categories
- Top bar: Time range selector
- Charts: Click and drag to zoom

### Key Metrics

#### CPU
- Total usage
- Per-core usage
- User vs system time
- Interrupts and context switches

#### Memory
- Used vs free
- Cached and buffers
- Swap usage
- Available memory

#### Disks
- Read/write operations
- I/O wait time
- Disk space usage
- IOPS

#### Network
- Incoming/outgoing bandwidth
- Packets per second
- Errors and drops
- TCP connections

### Configuration

#### Main Config
Location: `/etc/netdata/netdata.conf`

```ini
[global]
    # Update frequency
    update every = 1
    
    # History duration
    history = 86400
    
    # Bind to
    bind to = 0.0.0.0
    
    # Port
    port = 19999

[web]
    # Web mode
    mode = static-threaded
    
    # ACL
    allow connections from = localhost 10.* 192.168.* 172.16.* 172.17.* 172.18.* 172.19.* 172.20.* 172.21.* 172.22.* 172.23.* 172.24.* 172.25.* 172.26.* 172.27.* 172.28.* 172.29.* 172.30.* 172.31.*
```

#### Enable Authentication
1. Create password file:
```bash
sudo htpasswd -c /etc/netdata/users admin
```

2. Edit config:
```ini
[web]
    allow netdata.conf from = *
    enable gzip compression = yes
    respect do not track policy = no
```

3. Restart:
```bash
sudo systemctl restart netdata
```

### Alerts

#### Default Alerts
Netdata includes 100+ pre-configured alerts:
- High CPU usage
- Low memory
- Disk full
- Network errors
- Service down

#### View Alerts
- Dashboard: Red bell icon at top
- Command line:
```bash
sudo cat /var/cache/netdata/health-alarm-log.db
```

#### Configure Alerts
Location: `/etc/netdata/health.d/`

Example: Custom CPU alert
```bash
# Create custom alert
sudo nano /etc/netdata/health.d/cpu_usage.conf
```

```yaml
alarm: cpu_usage_high
    on: system.cpu
lookup: average -3m unaligned of user,system,softirq,irq,guest
  units: %
  every: 10s
   warn: $this > 70
   crit: $this > 85
   info: CPU utilization is high
     to: sysadmin
```

Restart to apply:
```bash
sudo systemctl restart netdata
```

### Integrations

#### Slack Notifications
```bash
# Edit health config
sudo nano /etc/netdata/health_alarm_notify.conf
```

```bash
# Enable Slack
SEND_SLACK="YES"
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
DEFAULT_RECIPIENT_SLACK="#alerts"
```

#### Discord Notifications
```bash
SEND_DISCORD="YES"
DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK"
DEFAULT_RECIPIENT_DISCORD="#alerts"
```

#### Email Notifications
```bash
SEND_EMAIL="YES"
EMAIL_SENDER="netdata@example.com"
DEFAULT_RECIPIENT_EMAIL="admin@example.com"
```

### Exporting Metrics

#### Prometheus
```bash
# Enable Prometheus exporter
sudo nano /etc/netdata/netdata.conf
```

```ini
[exporting:global]
    enabled = yes
    send configured labels = yes
    send automatic labels = yes
    update every = 10

[exporting:prometheus:exporter]
    enabled = yes
    send names instead of ids = yes
    send configured labels = yes
    send automatic labels = yes
```

Access at: `http://YOUR_IP:19999/api/v1/allmetrics?format=prometheus`

### Command Line Tools

#### View Metrics
```bash
# Get all metrics
curl http://localhost:19999/api/v1/info

# Get specific chart
curl http://localhost:19999/api/v1/data?chart=system.cpu

# Get current alarms
curl http://localhost:19999/api/v1/alarms
```

### Troubleshooting

#### Check Status
```bash
sudo systemctl status netdata
```

#### View Logs
```bash
sudo journalctl -u netdata -f
```

#### Restart
```bash
sudo systemctl restart netdata
```

#### Update
```bash
# Netdata updates automatically via cron
# Or manually:
sudo /usr/sbin/netdata-updater
```

### Performance Impact

Netdata is designed for minimal impact:
- **CPU**: ~1% of a single core
- **RAM**: ~50MB for 24h of metrics
- **Disk**: Minimal (mostly in memory)
- **Network**: ~1KB/s

### Custom Dashboards

Create custom HTML dashboards:

```html
<!DOCTYPE html>
<html>
<head>
    <title>My Dashboard</title>
    <script src="http://YOUR_IP:19999/dashboard.js"></script>
</head>
<body>
    <div data-netdata="system.cpu"
         data-chart-library="dygraph"
         data-width="100%"
         data-height="200px">
    </div>
    
    <div data-netdata="system.ram"
         data-chart-library="dygraph"
         data-width="100%"
         data-height="200px">
    </div>
</body>
</html>
```

---

## Monitoring Best Practices

### 1. Set Baselines
- Monitor for a week to establish normal patterns
- Document expected resource usage

### 2. Configure Alerts
- Set realistic thresholds
- Use multiple notification channels
- Don't over-alert (alert fatigue)

### 3. Regular Reviews
- Check dashboards weekly
- Review alert effectiveness
- Adjust thresholds as needed

### 4. Keep It Simple
- Start with default dashboards
- Add custom metrics gradually
- Focus on key metrics

### 5. Documentation
- Document what each metric means
- Record troubleshooting steps
- Maintain runbooks

---

## Alternative Monitoring Tools

While not included in HostKit, you may consider:

| Tool | Type | Complexity |
|------|------|------------|
| **Grafana + Prometheus** | Time-series | High |
| **Zabbix** | Enterprise | High |
| **Nagios** | Infrastructure | Medium |
| **Cacti** | Network | Medium |
| **Munin** | System | Low |

Netdata is included because it provides excellent monitoring with zero configuration and minimal resource usage.
