package firewall

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Manager handles firewall operations
type Manager struct {
	firewallType string // csf, fail2ban, ufw
}

// NewManager creates a new firewall manager
func NewManager() *Manager {
	return &Manager{}
}

// DetectFirewall detects installed firewall
func (m *Manager) DetectFirewall() string {
	// Check for CSF
	if _, err := os.Stat("/usr/sbin/csf"); err == nil {
		m.firewallType = "csf"
		return "csf"
	}

	// Check for Fail2ban
	cmd := exec.Command("which", "fail2ban-client")
	if err := cmd.Run(); err == nil {
		m.firewallType = "fail2ban"
		return "fail2ban"
	}

	// Check for UFW
	cmd = exec.Command("which", "ufw")
	if err := cmd.Run(); err == nil {
		m.firewallType = "ufw"
		return "ufw"
	}

	return ""
}

// GetStatus gets firewall status
func (m *Manager) GetStatus() (map[string]string, error) {
	switch m.firewallType {
	case "csf":
		return m.getCSFStatus()
	case "fail2ban":
		return m.getFail2banStatus()
	case "ufw":
		return m.getUFWStatus()
	default:
		return nil, fmt.Errorf("no firewall detected")
	}
}

func (m *Manager) getCSFStatus() (map[string]string, error) {
	status := make(map[string]string)

	// Check if CSF is running
	cmd := exec.Command("csf", "-l")
	output, err := cmd.Output()
	if err != nil {
		status["status"] = "stopped"
	} else {
		status["status"] = "running"
		// Parse output for additional info
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "DROP") {
				status["mode"] = "enabled"
				break
			}
		}
	}

	// Get version
	cmd = exec.Command("csf", "-v")
	output, _ = cmd.Output()
	if len(output) > 0 {
		status["version"] = strings.TrimSpace(string(output))
	}

	return status, nil
}

func (m *Manager) getFail2banStatus() (map[string]string, error) {
	status := make(map[string]string)

	// Check status
	cmd := exec.Command("fail2ban-client", "status")
	output, err := cmd.Output()
	if err != nil {
		status["status"] = "stopped"
	} else {
		status["status"] = "running"
		// Count jails
		if strings.Contains(string(output), "Jail list:") {
			jails := regexp.MustCompile(`Jail list:\s*(.+)`).FindStringSubmatch(string(output))
			if len(jails) > 1 {
				status["jails"] = strings.TrimSpace(jails[1])
			}
		}
	}

	return status, nil
}

func (m *Manager) getUFWStatus() (map[string]string, error) {
	status := make(map[string]string)

	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()
	if err != nil {
		status["status"] = "stopped"
	} else {
		outputStr := string(output)
		if strings.Contains(outputStr, "Status: active") {
			status["status"] = "active"
		} else {
			status["status"] = "inactive"
		}

		// Count rules
		rules := strings.Count(outputStr, "\n") - 2 // Header lines
		if rules > 0 {
			status["rules"] = fmt.Sprintf("%d", rules)
		}
	}

	return status, nil
}

// BlockIP blocks an IP address
func (m *Manager) BlockIP(ip string) error {
	switch m.firewallType {
	case "csf":
		cmd := exec.Command("csf", "-d", ip)
		return cmd.Run()
	case "fail2ban":
		// Fail2ban doesn't have direct IP blocking, use iptables
		cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
		return cmd.Run()
	case "ufw":
		cmd := exec.Command("ufw", "deny", "from", ip)
		return cmd.Run()
	default:
		return fmt.Errorf("no firewall detected")
	}
}

// UnblockIP unblocks an IP address
func (m *Manager) UnblockIP(ip string) error {
	switch m.firewallType {
	case "csf":
		cmd := exec.Command("csf", "-dr", ip)
		return cmd.Run()
	case "fail2ban":
		cmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
		return cmd.Run()
	case "ufw":
		cmd := exec.Command("ufw", "delete", "deny", "from", ip)
		return cmd.Run()
	default:
		return fmt.Errorf("no firewall detected")
	}
}

// AllowIP allows an IP address
func (m *Manager) AllowIP(ip string) error {
	switch m.firewallType {
	case "csf":
		cmd := exec.Command("csf", "-a", ip)
		return cmd.Run()
	case "ufw":
		cmd := exec.Command("ufw", "allow", "from", ip)
		return cmd.Run()
	default:
		return fmt.Errorf("operation not supported for %s", m.firewallType)
	}
}

// GetBlockedIPs returns list of blocked IPs
func (m *Manager) GetBlockedIPs() ([]string, error) {
	switch m.firewallType {
	case "csf":
		return m.getCSFBlockedIPs()
	case "fail2ban":
		return m.getFail2banBlockedIPs()
	case "ufw":
		return m.getUFWBlockedIPs()
	default:
		return nil, fmt.Errorf("no firewall detected")
	}
}

func (m *Manager) getCSFBlockedIPs() ([]string, error) {
	file, err := os.Open("/etc/csf/csf.deny")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				ips = append(ips, fields[0])
			}
		}
	}

	return ips, scanner.Err()
}

func (m *Manager) getFail2banBlockedIPs() ([]string, error) {
	cmd := exec.Command("fail2ban-client", "status", "sshd")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ips []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Banned IP list:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				ipList := strings.TrimSpace(parts[1])
				ips = strings.Fields(ipList)
			}
		}
	}

	return ips, nil
}

func (m *Manager) getUFWBlockedIPs() ([]string, error) {
	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ips []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "DENY") || strings.Contains(line, "REJECT") {
			fields := strings.Fields(line)
			for _, field := range fields {
				if net.ParseIP(field) != nil {
					ips = append(ips, field)
				}
			}
		}
	}

	return ips, nil
}

// RestartFirewall restarts the firewall
func (m *Manager) RestartFirewall() error {
	switch m.firewallType {
	case "csf":
		cmd := exec.Command("csf", "-r")
		return cmd.Run()
	case "fail2ban":
		cmd := exec.Command("systemctl", "restart", "fail2ban")
		return cmd.Run()
	case "ufw":
		cmd := exec.Command("ufw", "reload")
		return cmd.Run()
	default:
		return fmt.Errorf("no firewall detected")
	}
}
