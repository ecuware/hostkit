package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"hostkit/internal/config"
	"hostkit/internal/installer"
)

// DialogType represents dialog types
type DialogType int

const (
	DialogConfirmInstall DialogType = iota
	DialogConfirmUninstall
	DialogShowDependencies
	DialogError
	DialogSuccess
)

// DialogModel represents a dialog
type DialogModel struct {
	dialogType   DialogType
	config       *config.Config
	dependencies []string
	message      string
	selected     int // 0 = yes/ok, 1 = no/cancel
}

// NewConfirmDialog creates a confirmation dialog
func NewConfirmDialog(cfg *config.Config, deps []string) DialogModel {
	return DialogModel{
		dialogType:   DialogConfirmInstall,
		config:       cfg,
		dependencies: deps,
		selected:     0,
	}
}

// NewErrorDialog creates an error dialog
func NewErrorDialog(message string) DialogModel {
	return DialogModel{
		dialogType: DialogError,
		message:    message,
		selected:   0,
	}
}

// View renders the dialog
func (d DialogModel) View() string {
	var b strings.Builder

	// Box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(2, 4).
		Width(60)

	switch d.dialogType {
	case DialogConfirmInstall:
		b.WriteString(titleStyle.Render("⚠️  Confirm Installation"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("You are about to install %s %s\n\n", d.config.Icon, d.config.Name))
		b.WriteString(descriptionStyle.Render(d.config.Description))
		b.WriteString("\n\n")

		if len(d.dependencies) > 0 {
			b.WriteString(categoryStyle.Render("Required Dependencies:"))
			b.WriteString("\n")
			for _, dep := range d.dependencies {
				b.WriteString(fmt.Sprintf("  • %s\n", dep))
			}
			b.WriteString("\n")
		}

		if d.config.Requirements.MinRAM != "" {
			b.WriteString(fmt.Sprintf("Required RAM: %s\n", d.config.Requirements.MinRAM))
		}
		if d.config.Requirements.MinDisk != "" {
			b.WriteString(fmt.Sprintf("Required Disk: %s\n", d.config.Requirements.MinDisk))
		}

		b.WriteString("\n")
		b.WriteString(d.renderButtons("Install", "Cancel"))

	case DialogError:
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render("❌ Error"))
		b.WriteString("\n\n")
		b.WriteString(d.message)
		b.WriteString("\n\n")
		b.WriteString(d.renderButtons("OK", ""))
	}

	return boxStyle.Render(b.String())
}

func (d DialogModel) renderButtons(yesText, noText string) string {
	var b strings.Builder

	yesStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2)

	noStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2)

	if d.selected == 0 {
		b.WriteString(yesStyle.Render("[" + yesText + "]"))
	} else {
		b.WriteString(yesStyle.Render(" " + yesText + " "))
	}

	if noText != "" {
		b.WriteString("  ")
		if d.selected == 1 {
			b.WriteString(noStyle.Render("[" + noText + "]"))
		} else {
			b.WriteString(noStyle.Render(" " + noText + " "))
		}
	}

	return b.String()
}

// IsConfirmed returns true if yes is selected
func (d DialogModel) IsConfirmed() bool {
	return d.selected == 0
}

// Next moves selection to next button
func (d *DialogModel) Next() {
	if d.selected < 1 {
		d.selected++
	}
}

// Prev moves selection to previous button
func (d *DialogModel) Prev() {
	if d.selected > 0 {
		d.selected--
	}
}

// ResolveDependencies resolves and returns all dependencies for a package
func ResolveDependencies(cfg *config.Config, configs map[string]*config.Config) []string {
	resolver := installer.NewDependencyResolver(configs)
	resolved, err := resolver.Resolve(cfg.ID)
	if err != nil {
		return []string{}
	}

	// Remove the package itself from dependencies
	var deps []string
	for _, dep := range resolved {
		if dep != cfg.ID {
			if depCfg, ok := configs[dep]; ok {
				deps = append(deps, depCfg.Name)
			}
		}
	}

	return deps
}
