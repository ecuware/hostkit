package status

import (
	"os/exec"
	"strings"

	"hostkit/internal/config"
)

// Checker checks package installation status
type Checker struct{}

// NewChecker creates a new status checker
func NewChecker() *Checker {
	return &Checker{}
}

// Status represents package installation status
type Status int

const (
	StatusNotInstalled Status = iota
	StatusInstalled
	StatusUpdateAvailable
	StatusInstalling
	StatusFailed
)

// Check checks if a package is installed
func (c *Checker) Check(pkg *config.Config) Status {
	// Try various detection methods
	checkers := []func(*config.Config) bool{
		c.checkService,
		c.checkBinary,
		c.checkPackage,
		c.checkPort,
	}

	for _, checker := range checkers {
		if checker(pkg) {
			// Check for update availability
			if c.hasUpdate(pkg) {
				return StatusUpdateAvailable
			}
			return StatusInstalled
		}
	}

	return StatusNotInstalled
}

// checkService checks if service is running
func (c *Checker) checkService(pkg *config.Config) bool {
	services := []string{
		pkg.ID,
		strings.ReplaceAll(pkg.ID, "-", ""),
		pkg.ID + "d",
	}

	for _, service := range services {
		cmd := exec.Command("systemctl", "is-active", "--quiet", service)
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

// checkBinary checks if binary exists
func (c *Checker) checkBinary(pkg *config.Config) bool {
	binaries := []string{
		pkg.ID,
		strings.ReplaceAll(pkg.ID, "-", ""),
	}

	for _, binary := range binaries {
		cmd := exec.Command("which", binary)
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

// checkPackage checks via package manager
func (c *Checker) checkPackage(pkg *config.Config) bool {
	// Try dpkg (Debian/Ubuntu)
	cmd := exec.Command("dpkg", "-l", pkg.ID)
	if err := cmd.Run(); err == nil {
		return true
	}

	// Try rpm (CentOS/RHEL)
	cmd = exec.Command("rpm", "-q", pkg.ID)
	if err := cmd.Run(); err == nil {
		return true
	}

	return false
}

// checkPort checks if package's default port is listening
func (c *Checker) checkPort(pkg *config.Config) bool {
	if len(pkg.Requirements.Ports) == 0 {
		return false
	}

	for _, port := range pkg.Requirements.Ports {
		cmd := exec.Command("ss", "-tlnp", "|", "grep", string(rune(port)))
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

// hasUpdate checks if update is available
func (c *Checker) hasUpdate(pkg *config.Config) bool {
	// TODO: Implement version comparison
	// For now, return false
	return false
}

// GetStatusIcon returns the icon for a status
func GetStatusIcon(status Status) string {
	switch status {
	case StatusInstalled:
		return "✅"
	case StatusUpdateAvailable:
		return "🔄"
	case StatusInstalling:
		return "⏳"
	case StatusFailed:
		return "❌"
	default:
		return "  "
	}
}

// GetStatusText returns the text for a status
func GetStatusText(status Status) string {
	switch status {
	case StatusInstalled:
		return "Installed"
	case StatusUpdateAvailable:
		return "Update Available"
	case StatusInstalling:
		return "Installing..."
	case StatusFailed:
		return "Failed"
	default:
		return "Not Installed"
	}
}
