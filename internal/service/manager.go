package service

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Manager handles service operations
type Manager struct{}

// NewManager creates a new service manager
func NewManager() *Manager {
	return &Manager{}
}

// Status represents service status
type Status int

const (
	StatusUnknown Status = iota
	StatusRunning
	StatusStopped
	StatusFailed
	StatusNotInstalled
)

// ServiceInfo holds service information
type ServiceInfo struct {
	Name        string
	Status      Status
	Enabled     bool
	PID         string
	Uptime      string
	Memory      string
	Description string
}

// GetStatus returns service status
func (m *Manager) GetStatus(serviceName string) (Status, error) {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	if err != nil {
		status := strings.TrimSpace(string(output))
		if status == "inactive" || status == "failed" {
			return StatusStopped, nil
		}
		return StatusUnknown, fmt.Errorf("failed to get status: %w", err)
	}

	status := strings.TrimSpace(string(output))
	if status == "active" {
		return StatusRunning, nil
	}
	return StatusStopped, nil
}

// Start starts a service
func (m *Manager) Start(ctx context.Context, serviceName string) error {
	cmd := exec.CommandContext(ctx, "systemctl", "start", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	return nil
}

// Stop stops a service
func (m *Manager) Stop(ctx context.Context, serviceName string) error {
	cmd := exec.CommandContext(ctx, "systemctl", "stop", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}
	return nil
}

// Restart restarts a service
func (m *Manager) Restart(ctx context.Context, serviceName string) error {
	cmd := exec.CommandContext(ctx, "systemctl", "restart", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}
	return nil
}

// Enable enables auto-start
func (m *Manager) Enable(ctx context.Context, serviceName string) error {
	cmd := exec.CommandContext(ctx, "systemctl", "enable", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}
	return nil
}

// Disable disables auto-start
func (m *Manager) Disable(ctx context.Context, serviceName string) error {
	cmd := exec.CommandContext(ctx, "systemctl", "disable", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable service: %w", err)
	}
	return nil
}

// GetInfo gets detailed service information
func (m *Manager) GetInfo(serviceName string) (*ServiceInfo, error) {
	info := &ServiceInfo{Name: serviceName}

	// Get status
	status, err := m.GetStatus(serviceName)
	if err != nil {
		return nil, err
	}
	info.Status = status

	// Check if enabled
	cmd := exec.Command("systemctl", "is-enabled", serviceName)
	if err := cmd.Run(); err == nil {
		info.Enabled = true
	}

	// Get additional info via systemctl show
	cmd = exec.Command("systemctl", "show", serviceName,
		"--property=MainPID,MemoryCurrent,Description")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MainPID=") {
				info.PID = strings.TrimPrefix(line, "MainPID=")
			} else if strings.HasPrefix(line, "MemoryCurrent=") {
				mem := strings.TrimPrefix(line, "MemoryCurrent=")
				if mem != "[not set]" {
					info.Memory = mem
				}
			} else if strings.HasPrefix(line, "Description=") {
				info.Description = strings.TrimPrefix(line, "Description=")
			}
		}
	}

	return info, nil
}

// GetStatusText returns human-readable status
func GetStatusText(status Status) string {
	switch status {
	case StatusRunning:
		return "Running"
	case StatusStopped:
		return "Stopped"
	case StatusFailed:
		return "Failed"
	case StatusNotInstalled:
		return "Not Installed"
	default:
		return "Unknown"
	}
}

// GetStatusIcon returns icon for status
func GetStatusIcon(status Status) string {
	switch status {
	case StatusRunning:
		return "🟢"
	case StatusStopped:
		return "🔴"
	case StatusFailed:
		return "❌"
	default:
		return "⚪"
	}
}

// PostInstallConfig holds post-install configuration
type PostInstallConfig struct {
	Type    string            `yaml:"type"`
	Command string            `yaml:"command"`
	Inputs  map[string]string `yaml:"inputs,omitempty"`
	Message string            `yaml:"message"`
}

// RunPostInstall runs post-install configuration
func (m *Manager) RunPostInstall(ctx context.Context, config PostInstallConfig) error {
	if config.Type == "mysql_secure" {
		return m.runMySQLSecureInstallation(ctx, config)
	}

	if config.Command != "" {
		cmd := exec.CommandContext(ctx, "bash", "-c", config.Command)
		return cmd.Run()
	}

	return nil
}

func (m *Manager) runMySQLSecureInstallation(ctx context.Context, config PostInstallConfig) error {
	// Generate secure installation script
	rootPassword := config.Inputs["root_password"]
	if rootPassword == "" {
		rootPassword = generateSecurePassword()
	}

	script := fmt.Sprintf(`
#!/bin/bash
mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '%s';"
mysql -e "DELETE FROM mysql.user WHERE User='';"
mysql -e "DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');"
mysql -e "DROP DATABASE IF EXISTS test;"
mysql -e "FLUSH PRIVILEGES;"
`, rootPassword)

	cmd := exec.CommandContext(ctx, "bash", "-c", script)
	return cmd.Run()
}

func generateSecurePassword() string {
	// Simple password generation
	return "HostKit" + "@" + fmt.Sprintf("%d", 1000+9999)
}
