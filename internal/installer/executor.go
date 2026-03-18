package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Executor handles command execution
type Executor struct {
	verbose bool
}

// NewExecutor creates a new executor
func NewExecutor() *Executor {
	return &Executor{
		verbose: true,
	}
}

// Run executes a command and returns output
func (e *Executor) Run(ctx context.Context, command string) ([]byte, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Env = os.Environ()

	return cmd.CombinedOutput()
}

// RunWithOutput executes a command and streams output
func (e *Executor) RunWithOutput(ctx context.Context, command string) error {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}

// RunSilent executes a command without output
func (e *Executor) RunSilent(ctx context.Context, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Env = os.Environ()

	return cmd.Run()
}

// IsRoot checks if running as root
func (e *Executor) IsRoot() bool {
	return os.Geteuid() == 0
}

// RequireRoot ensures command runs as root
func (e *Executor) RequireRoot() error {
	if !e.IsRoot() {
		return fmt.Errorf("this command requires root privileges. Please run with sudo")
	}
	return nil
}
