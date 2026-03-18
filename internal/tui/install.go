package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hostkit/internal/config"
	"hostkit/internal/installer"
)

// InstallState represents installation states
type InstallState int

const (
	InstallStateIdle InstallState = iota
	InstallStateChecking
	InstallStateInstalling
	InstallStateSuccess
	InstallStateError
)

// InstallProgressMsg is sent during installation
type InstallProgressMsg struct {
	Step    string
	Message string
	Percent float64
}

// InstallCompleteMsg is sent when installation finishes
type InstallCompleteMsg struct {
	Error error
}

// InstallModel handles the installation screen
type InstallModel struct {
	config      *config.Config
	state       InstallState
	progress    progress.Model
	logs        []string
	logLimit    int
	percent     float64
	currentStep string
	error       error
	engine      *installer.Engine
	startTime   time.Time
}

// NewInstallModel creates a new installation model
func NewInstallModel(cfg *config.Config) InstallModel {
	prog := progress.New(progress.WithScaledGradient("#7D56F4", "#04B575"))
	prog.Width = 50

	return InstallModel{
		config:   cfg,
		state:    InstallStateIdle,
		progress: prog,
		logs:     []string{},
		logLimit: 100,
		engine:   installer.NewEngine(false, false),
	}
}

// Init initializes the installation model
func (m InstallModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for install screen
func (m InstallModel) Update(msg tea.Msg) (InstallModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" {
			return m, tea.Quit
		}
		if msg.String() == "enter" && m.state == InstallStateIdle {
			return m, m.startInstall()
		}

	case InstallProgressMsg:
		m.currentStep = msg.Step
		m.percent = msg.Percent
		m.addLog(fmt.Sprintf("[%s] %s", msg.Step, msg.Message))
		return m, m.progress.SetPercent(msg.Percent)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case InstallCompleteMsg:
		if msg.Error != nil {
			m.state = InstallStateError
			m.error = msg.Error
			m.addLog(fmt.Sprintf("Error: %v", msg.Error))
		} else {
			m.state = InstallStateSuccess
			m.addLog("Installation completed successfully!")
		}
	}

	return m, nil
}

func (m *InstallModel) startInstall() tea.Cmd {
	m.state = InstallStateInstalling
	m.startTime = time.Now()
	m.addLog("Starting installation...")

	return func() tea.Msg {
		ctx := context.Background()

		opts := installer.InstallOptions{
			Version: m.config.Version.Current,
		}

		err := m.engine.Install(ctx, m.config, opts)
		return InstallCompleteMsg{Error: err}
	}
}

func (m *InstallModel) addLog(message string) {
	timestamp := time.Now().Format("15:04:05")
	log := fmt.Sprintf("[%s] %s", timestamp, message)

	m.logs = append(m.logs, log)

	if len(m.logs) > m.logLimit {
		m.logs = m.logs[len(m.logs)-m.logLimit:]
	}
}

// View renders the installation screen
func (m InstallModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(fmt.Sprintf("Installing %s %s", m.config.Icon, m.config.Name)))
	b.WriteString("\n\n")

	// Status
	switch m.state {
	case InstallStateIdle:
		b.WriteString("Press Enter to start installation or Q to quit\n")
		b.WriteString(itemStyle.Render("This will install "))
		b.WriteString(categoryStyle.Render(m.config.Name))
		b.WriteString(itemStyle.Render(" and its dependencies.\n"))

	case InstallStateInstalling:
		b.WriteString(fmt.Sprintf("Current step: %s\n", m.currentStep))
		b.WriteString(m.progress.View())
		b.WriteString(fmt.Sprintf(" %.0f%%\n", m.percent*100))

	case InstallStateSuccess:
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true).
			Render("✓ Installation completed successfully!"))
		b.WriteString("\n")

	case InstallStateError:
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render("✗ Installation failed"))
		b.WriteString("\n")
		if m.error != nil {
			b.WriteString(fmt.Sprintf("Error: %v\n", m.error))
		}
	}

	b.WriteString("\n")

	// Logs
	if len(m.logs) > 0 {
		b.WriteString(categoryStyle.Render("Installation Logs:"))
		b.WriteString("\n")

		// Show last 10 logs
		start := 0
		if len(m.logs) > 10 {
			start = len(m.logs) - 10
		}

		for _, log := range m.logs[start:] {
			b.WriteString(logStyle.Render(log))
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString("\n")
	if m.state == InstallStateSuccess || m.state == InstallStateError {
		b.WriteString(itemStyle.Render("Press Q to return to main menu"))
	} else if m.state == InstallStateInstalling {
		b.WriteString(itemStyle.Render("Installation in progress..."))
	}

	return b.String()
}

var logStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#888888")).
	Faint(true)
