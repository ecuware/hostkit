# 🔐 VPN

Virtual Private Network solutions for secure remote access.

## Overview

| Solution | Protocol | Best For |
|----------|----------|----------|
| **WireGuard** | WireGuard | Modern, fast, simple |

---

## WireGuard

Modern, fast, and secure VPN tunnel.

### Features
- Simple configuration
- High performance
- Modern cryptography
- Cross-platform
- Minimal code base
- Easy to audit

### Install
```bash
hostkit install wireguard
```

### What's Installed
- WireGuard kernel module
- WireGuard tools
- wg-quick utility

### Basic Concepts

#### Interface
Virtual network interface (e.g., `wg0`)

#### Peer
A device that connects to the VPN

#### Keys
- **Private Key**: Kept secret on each peer
- **Public Key**: Shared with other peers

### Server Configuration

#### 1. Generate Keys
```bash
# Generate private key
wg genkey | tee privatekey | wg pubkey > publickey

# View keys
cat privatekey
cat publickey
```

#### 2. Create Server Config
```bash
sudo nano /etc/wireguard/wg0.conf
```

```ini
[Interface]
PrivateKey = <server_private_key>
Address = 10.0.0.1/24
ListenPort = 51820
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
DNS = 1.1.1.1, 8.8.8.8

# Client 1
[Peer]
PublicKey = <client1_public_key>
AllowedIPs = 10.0.0.2/32

# Client 2
[Peer]
PublicKey = <client2_public_key>
AllowedIPs = 10.0.0.3/32
```

#### 3. Enable IP Forwarding
```bash
# Temporary
sudo sysctl -w net.ipv4.ip_forward=1

# Permanent
sudo nano /etc/sysctl.conf
# Uncomment: net.ipv4.ip_forward=1

# Apply
sudo sysctl -p
```

#### 4. Start WireGuard
```bash
# Start interface
sudo wg-quick up wg0

# Enable on boot
sudo systemctl enable wg-quick@wg0

# Check status
sudo wg show
```

### Client Configuration

#### 1. Generate Client Keys
```bash
wg genkey | tee client_privatekey | wg pubkey > client_publickey
```

#### 2. Create Client Config
```bash
nano client1.conf
```

```ini
[Interface]
PrivateKey = <client_private_key>
Address = 10.0.0.2/32
DNS = 1.1.1.1, 8.8.8.8

[Peer]
PublicKey = <server_public_key>
Endpoint = <server_ip>:51820
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
```

### Adding Clients

#### On Server
```bash
# Edit config
sudo nano /etc/wireguard/wg0.conf

# Add new peer
[Peer]
PublicKey = <new_client_public_key>
AllowedIPs = 10.0.0.4/32

# Restart
sudo wg-quick down wg0
sudo wg-quick up wg0
```

### Client Setup

#### Linux
```bash
# Install WireGuard
sudo apt install wireguard

# Copy config
sudo cp client1.conf /etc/wireguard/wg0.conf

# Start
sudo wg-quick up wg0

# Enable auto-start
sudo systemctl enable wg-quick@wg0
```

#### Windows
1. Download WireGuard from https://www.wireguard.com/install/
2. Import configuration file
3. Click "Activate"

#### macOS
```bash
# Install via Homebrew
brew install wireguard-tools

# Or download app from App Store
```

#### Android/iOS
1. Install WireGuard app
2. Scan QR code or import file
3. Toggle to connect

### Management Commands

```bash
# Show all interfaces
sudo wg show

# Show specific interface
sudo wg show wg0

# Show detailed info
sudo wg show wg0 dump

# Add peer dynamically
sudo wg set wg0 peer <public_key> allowed-ips 10.0.0.5/32

# Remove peer
sudo wg set wg0 peer <public_key> remove

# Stop interface
sudo wg-quick down wg0

# Start interface
sudo wg-quick up wg0

# Restart interface
sudo wg-quick down wg0 && sudo wg-quick up wg0

# Check logs
sudo journalctl -u wg-quick@wg0 -f
```

### Configuration Examples

#### Split Tunnel (VPN only for specific IPs)
```ini
[Interface]
PrivateKey = <client_private_key>
Address = 10.0.0.2/32

[Peer]
PublicKey = <server_public_key>
Endpoint = <server_ip>:51820
AllowedIPs = 10.0.0.0/24, 192.168.1.0/24
PersistentKeepalive = 25
```

#### Full Tunnel (All traffic through VPN)
```ini
[Interface]
PrivateKey = <client_private_key>
Address = 10.0.0.2/32
DNS = 1.1.1.1

[Peer]
PublicKey = <server_public_key>
Endpoint = <server_ip>:51820
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25
```

#### Site-to-Site
```ini
# Site A
[Interface]
PrivateKey = <site_a_private_key>
Address = 10.0.0.1/24
ListenPort = 51820

[Peer]
PublicKey = <site_b_public_key>
AllowedIPs = 10.0.0.2/32, 192.168.2.0/24
Endpoint = site_b_ip:51820

# Site B
[Interface]
PrivateKey = <site_b_private_key>
Address = 10.0.0.2/24
ListenPort = 51820

[Peer]
PublicKey = <site_a_public_key>
AllowedIPs = 10.0.0.1/32, 192.168.1.0/24
Endpoint = site_a_ip:51820
```

### Firewall Configuration

#### UFW
```bash
# Allow WireGuard port
sudo ufw allow 51820/udp

# Allow forwarding
sudo ufw route allow in on wg0 out on eth0
sudo ufw route allow in on eth0 out on wg0
```

#### iptables
```bash
# Allow port
sudo iptables -A INPUT -p udp --dport 51820 -j ACCEPT

# Enable forwarding
sudo iptables -A FORWARD -i wg0 -j ACCEPT
sudo iptables -A FORWARD -o wg0 -j ACCEPT
sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
```

### Troubleshooting

#### Connection Issues
```bash
# Check if running
sudo wg show

# Check interface
ip addr show wg0

# Check routing
ip route show

# Test connectivity
ping 10.0.0.1

# Check logs
sudo journalctl -u wg-quick@wg0 -n 100
```

#### No Internet Through VPN
1. Check IP forwarding is enabled
2. Verify iptables rules
3. Check DNS configuration
4. Ensure AllowedIPs is correct

#### Slow Performance
- Try different MTU values (default: 1420)
- Check server bandwidth
- Verify CPU usage (encryption overhead)

### Security Best Practices

1. **Keep Private Keys Secret**
   - Never share private keys
   - Store securely

2. **Use Strong Keys**
   - Always generate new keys, don't reuse
   - Keys are 256-bit by default

3. **Limit Access**
   - Use specific AllowedIPs when possible
   - Don't use 0.0.0.0/0 unnecessarily

4. **PersistentKeepalive**
   - Use for clients behind NAT
   - Set to 25 seconds

5. **Rotate Keys Periodically**
   - Change keys every few months
   - Remove unused peers

6. **Monitor Connections**
   ```bash
   sudo wg show
   ```

### QR Code for Mobile

Generate QR code for easy mobile setup:

```bash
# Install qrencode
sudo apt install qrencode

# Generate QR
cat client.conf | qrencode -t ansiutf8

# Or save as image
cat client.conf | qrencode -t png -o client.png
```

### Automation Scripts

#### Create Client Script
```bash
#!/bin/bash
# create-client.sh

CLIENT_NAME=$1
SERVER_IP="your_server_ip"
SERVER_PUBLIC_KEY="your_server_public_key"
CLIENT_IP="10.0.0.$((RANDOM % 200 + 10))"

# Generate keys
CLIENT_PRIVATE=$(wg genkey)
CLIENT_PUBLIC=$(echo $CLIENT_PRIVATE | wg pubkey)

# Create client config
cat > ${CLIENT_NAME}.conf << EOF
[Interface]
PrivateKey = $CLIENT_PRIVATE
Address = $CLIENT_IP/32
DNS = 1.1.1.1

[Peer]
PublicKey = $SERVER_PUBLIC_KEY
Endpoint = $SERVER_IP:51820
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
EOF

# Add to server
echo "" | sudo tee -a /etc/wireguard/wg0.conf
echo "[Peer]" | sudo tee -a /etc/wireguard/wg0.conf
echo "PublicKey = $CLIENT_PUBLIC" | sudo tee -a /etc/wireguard/wg0.conf
echo "AllowedIPs = $CLIENT_IP/32" | sudo tee -a /etc/wireguard/wg0.conf

# Restart
sudo wg-quick down wg0
sudo wg-quick up wg0

echo "Client $CLIENT_NAME created with IP $CLIENT_IP"
```

---

## Comparison with Other VPNs

| Feature | WireGuard | OpenVPN | IPSec |
|---------|-----------|---------|-------|
| Speed | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| Ease of Use | ⭐⭐⭐ | ⭐⭐ | ⭐ |
| Security | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| Code Size | ⭐⭐⭐ (4K LOC) | ⭐ (100K+ LOC) | ⭐⭐ |
| Mobile Battery | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| Auditability | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ |

WireGuard is included because it's modern, fast, and much simpler to configure than traditional VPN solutions.
