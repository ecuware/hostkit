package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"hostkit/internal/config"
)

// Engine handles package installation
type Engine struct {
	executor    *Executor
	logger      *Logger
	dryRun      bool
	interactive bool
}

// NewEngine creates a new installer engine
func NewEngine(dryRun, interactive bool) *Engine {
	return &Engine{
		executor:    NewExecutor(),
		logger:      NewLogger(),
		dryRun:      dryRun,
		interactive: interactive,
	}
}

// InstallOptions contains installation options
type InstallOptions struct {
	Version   string
	SkipDeps  bool
	Force     bool
	OSName    string
	OSVersion string
}

// Install installs a package
func (e *Engine) Install(ctx context.Context, cfg *config.Config, opts InstallOptions) error {
	// Check if already installed
	if !opts.Force && e.IsInstalled(cfg) {
		return fmt.Errorf("%s is already installed. Use --force to reinstall", cfg.Name)
	}

	e.logger.Info("Starting installation of %s", cfg.Name)
	startTime := time.Now()

	// Pre-installation checks
	if err := e.runPreChecks(ctx, cfg); err != nil {
		return fmt.Errorf("pre-installation check failed: %w", err)
	}

	// Install dependencies
	if !opts.SkipDeps {
		if err := e.installDependencies(ctx, cfg); err != nil {
			return fmt.Errorf("dependency installation failed: %w", err)
		}
	}

	// Perform installation based on method
	var err error
	switch cfg.Install.Method {
	case "shell":
		err = e.installShell(ctx, cfg, opts)
	case "apt":
		err = e.installAPT(ctx, cfg, opts)
	case "yum":
		err = e.installYUM(ctx, cfg, opts)
	default:
		return fmt.Errorf("unknown installation method: %s", cfg.Install.Method)
	}

	if err != nil {
		e.logger.Error("Installation failed: %v", err)
		return err
	}

	// Post-installation checks
	if err := e.runPostChecks(ctx, cfg); err != nil {
		e.logger.Warn("Post-installation check failed: %v", err)
	}

	duration := time.Since(startTime)
	e.logger.Success("Successfully installed %s in %v", cfg.Name, duration)

	return nil
}

// IsInstalled checks if a package is already installed
func (e *Engine) IsInstalled(cfg *config.Config) bool {
	// Check common installation indicators
	checks := []string{
		fmt.Sprintf("which %s", cfg.ID),
		fmt.Sprintf("systemctl is-active %s", cfg.ID),
		fmt.Sprintf("dpkg -l | grep -i %s", cfg.ID),
	}

	for _, check := range checks {
		cmd := exec.Command("bash", "-c", check)
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

func (e *Engine) runPreChecks(ctx context.Context, cfg *config.Config) error {
	for _, check := range cfg.Install.PreCheck {
		e.logger.Info("Running pre-check: %s", check.Command)

		if e.dryRun {
			e.logger.Info("[DRY RUN] Would run: %s", check.Command)
			continue
		}

		output, err := e.executor.Run(ctx, check.Command)
		if err != nil {
			return fmt.Errorf("%s: %w", check.ErrorMsg, err)
		}

		e.logger.Debug("Pre-check output: %s", string(output))
	}
	return nil
}

func (e *Engine) installDependencies(ctx context.Context, cfg *config.Config) error {
	if len(cfg.Dependencies.Required) == 0 && len(cfg.Dependencies.Optional) == 0 {
		return nil
	}

	e.logger.Info("Installing dependencies for %s", cfg.Name)

	// Install required dependencies
	for _, dep := range cfg.Dependencies.Required {
		e.logger.Info("Installing required dependency: %s", dep)
		// TODO: Load dep config and install
		_ = dep
	}

	// Optional dependencies are not auto-installed unless specified
	if e.interactive {
		for _, dep := range cfg.Dependencies.Optional {
			e.logger.Info("Optional dependency available: %s", dep)
			// TODO: Prompt user in TUI
		}
	}

	return nil
}

func (e *Engine) installShell(ctx context.Context, cfg *config.Config, opts InstallOptions) error {
	script := cfg.Install.Script

	// Replace template variables
	script = strings.ReplaceAll(script, "{{ .Version }}", opts.Version)
	script = strings.ReplaceAll(script, "{{ .OSName }}", opts.OSName)
	script = strings.ReplaceAll(script, "{{ .OSVersion }}", opts.OSVersion)

	if e.dryRun {
		e.logger.Info("[DRY RUN] Would execute shell script:\n%s", script)
		return nil
	}

	e.logger.Info("Executing installation script...")

	// Write script to temp file
	tmpFile, err := os.CreateTemp("", "hostkit-install-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(script); err != nil {
		return fmt.Errorf("failed to write script: %w", err)
	}
	tmpFile.Close()

	// Execute script
	cmd := exec.CommandContext(ctx, "bash", tmpFile.Name())
	cmd.Stdout = e.logger
	cmd.Stderr = e.logger
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installation script failed: %w", err)
	}

	return nil
}

func (e *Engine) installAPT(ctx context.Context, cfg *config.Config, opts InstallOptions) error {
	osInstall := cfg.GetOSInstall("ubuntu")
	if osInstall == nil {
		return fmt.Errorf("no APT installation configuration for this OS")
	}

	if e.dryRun {
		e.logger.Info("[DRY RUN] Would install packages via APT: %v", osInstall.Packages)
		return nil
	}

	// Update package list
	e.logger.Info("Updating package list...")
	if err := e.executor.RunWithOutput(ctx, "apt-get update"); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	// Install packages
	for _, pkg := range osInstall.Packages {
		e.logger.Info("Installing package: %s", pkg)
		cmd := fmt.Sprintf("apt-get install -y %s", pkg)
		if err := e.executor.RunWithOutput(ctx, cmd); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}

	// Run post-install commands
	if osInstall.PostInstall != "" {
		e.logger.Info("Running post-installation steps...")
		if err := e.executor.RunWithOutput(ctx, osInstall.PostInstall); err != nil {
			e.logger.Warn("Post-installation script failed: %v", err)
		}
	}

	return nil
}

func (e *Engine) installYUM(ctx context.Context, cfg *config.Config, opts InstallOptions) error {
	osInstall := cfg.GetOSInstall("centos")
	if osInstall == nil {
		return fmt.Errorf("no YUM installation configuration for this OS")
	}

	if e.dryRun {
		e.logger.Info("[DRY RUN] Would install packages via YUM: %v", osInstall.Packages)
		return nil
	}

	// Install packages
	for _, pkg := range osInstall.Packages {
		e.logger.Info("Installing package: %s", pkg)
		cmd := fmt.Sprintf("yum install -y %s", pkg)
		if err := e.executor.RunWithOutput(ctx, cmd); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}

	return nil
}

func (e *Engine) runPostChecks(ctx context.Context, cfg *config.Config) error {
	for _, check := range cfg.Install.PostCheck {
		e.logger.Info("Running post-check: %s", check.Service)

		if e.dryRun {
			e.logger.Info("[DRY RUN] Would check service: %s", check.Service)
			continue
		}

		// Check service status
		cmd := fmt.Sprintf("systemctl is-active %s", check.Service)
		if err := e.executor.RunWithOutput(ctx, cmd); err != nil {
			return fmt.Errorf("service %s is not running: %w", check.Service, err)
		}

		// Check port if specified
		if check.Port > 0 {
			cmd := fmt.Sprintf("ss -tlnp | grep ':%d'", check.Port)
			if err := e.executor.RunWithOutput(ctx, cmd); err != nil {
				return fmt.Errorf("port %d is not listening: %w", check.Port, err)
			}
		}
	}

	return nil
}

// Uninstall removes a package
func (e *Engine) Uninstall(ctx context.Context, cfg *config.Config, force bool) error {
	if !force {
		// TODO: Prompt for confirmation in TUI
	}

	e.logger.Info("Uninstalling %s", cfg.Name)

	if e.dryRun {
		e.logger.Info("[DRY RUN] Would run uninstall command:\n%s", cfg.Uninstall.Command)
		return nil
	}

	if err := e.executor.RunWithOutput(ctx, cfg.Uninstall.Command); err != nil {
		return fmt.Errorf("uninstall failed: %w", err)
	}

	e.logger.Success("Successfully uninstalled %s", cfg.Name)
	return nil
}
