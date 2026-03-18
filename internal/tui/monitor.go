package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hostkit/internal/monitor"
)

// SystemMonitorModel represents the system monitor screen
type SystemMonitorModel struct {
	monitor *monitor.Monitor
	stats   *monitor.SystemStats
	width   int
	height  int
	loading bool
	error   error
}

// MonitorTickMsg is sent periodically to update stats
type MonitorTickMsg struct {
	Stats *monitor.SystemStats
	Error error
}

// NewSystemMonitorModel creates a new system monitor model
func NewSystemMonitorModel() SystemMonitorModel {
	return SystemMonitorModel{
		monitor: monitor.NewMonitor(),
		loading: true,
	}
}

// Init initializes the monitor model
func (m SystemMonitorModel) Init() tea.Cmd {
	return tea.Batch(
		m.updateStats(),
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

func (m SystemMonitorModel) updateStats() tea.Cmd {
	return func() tea.Msg {
		stats, err := m.monitor.GetStats()
		return MonitorTickMsg{Stats: stats, Error: err}
	}
}

// Update handles messages
func (m SystemMonitorModel) Update(msg tea.Msg) (SystemMonitorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "r":
			return m, m.updateStats()
		}

	case tickMsg:
		return m, tea.Batch(m.updateStats(), tickCmd())

	case MonitorTickMsg:
		m.loading = false
		if msg.Error != nil {
			m.error = msg.Error
		} else {
			m.stats = msg.Stats
			m.error = nil
		}
	}

	return m, nil
}

// View renders the system monitor
func (m SystemMonitorModel) View() string {
	if m.loading {
		return "Loading system stats..."
	}

	if m.error != nil {
		// Check if it's the "not Linux" error
		errMsg := m.error.Error()
		if strings.Contains(errMsg, "requires Linux") {
			var b strings.Builder
			b.WriteString(titleStyle.Render("⚠️  System Monitor"))
			b.WriteString("\n\n")
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Bold(true).
				Render("Linux Gerekiyor!"))
			b.WriteString("\n\n")
			b.WriteString("System Monitor özelliği sadece Linux işletim sisteminde çalışır.\n")
			b.WriteString("\n")
			b.WriteString("MacOS/Windows kullanıyorsanız, Linux Sandbox ortamını kullanabilirsiniz:\n\n")
			b.WriteString(itemStyle.Render("  make sandbox\n"))
			b.WriteString("\n")
			b.WriteString("Veya:\n")
			b.WriteString(itemStyle.Render("  ./scripts/sandbox.sh\n"))
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("q: Geri dön"))
			return b.String()
		}
		return fmt.Sprintf("Error: %v\n\nPress any key to go back", m.error)
	}

	if m.stats == nil {
		return "No stats available"
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("📊 System Monitor"))
	b.WriteString("\n\n")

	// Update time
	b.WriteString(itemStyle.Render(fmt.Sprintf("Last updated: %s\n\n", m.stats.Timestamp.Format("15:04:05"))))

	// CPU Section
	b.WriteString(categoryStyle.Render("💻 CPU Usage"))
	b.WriteString("\n")
	b.WriteString(renderBar(m.stats.CPU.Usage, 100, "CPU"))
	b.WriteString(fmt.Sprintf("  %.1f%% (User: %.1f%%, System: %.1f%%)\n",
		m.stats.CPU.Usage, m.stats.CPU.User, m.stats.CPU.System))
	b.WriteString(itemStyle.Render(fmt.Sprintf("Cores: %d\n", m.stats.CPU.Cores)))
	b.WriteString("\n")

	// Memory Section
	b.WriteString(categoryStyle.Render("🧠 Memory Usage"))
	b.WriteString("\n")
	b.WriteString(renderBar(m.stats.Memory.UsagePct, 100, "RAM"))
	b.WriteString(fmt.Sprintf("  %.1f%% (%s / %s)\n",
		m.stats.Memory.UsagePct,
		monitor.FormatBytes(m.stats.Memory.Used),
		monitor.FormatBytes(m.stats.Memory.Total)))
	b.WriteString(itemStyle.Render(fmt.Sprintf("Free: %s | Cached: %s\n",
		monitor.FormatBytes(m.stats.Memory.Free),
		monitor.FormatBytes(m.stats.Memory.Cached))))
	b.WriteString("\n")

	// Disk Section
	b.WriteString(categoryStyle.Render("💾 Disk Usage"))
	b.WriteString("\n")
	b.WriteString(renderBar(m.stats.Disk.UsagePct, 100, "Disk"))
	b.WriteString(fmt.Sprintf("  %.1f%% (%s / %s)\n",
		m.stats.Disk.UsagePct,
		monitor.FormatBytes(m.stats.Disk.Used),
		monitor.FormatBytes(m.stats.Disk.Total)))
	b.WriteString(itemStyle.Render(fmt.Sprintf("Free: %s\n", monitor.FormatBytes(m.stats.Disk.Free))))
	b.WriteString("\n")

	// Network Section
	b.WriteString(categoryStyle.Render("🌐 Network"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("Interface: %s\n", m.stats.Network.Interface)))
	b.WriteString(itemStyle.Render(fmt.Sprintf("Download: %s/s (Total: %s)\n",
		monitor.FormatSpeed(m.stats.Network.RxSpeed),
		monitor.FormatBytes(m.stats.Network.RxBytes))))
	b.WriteString(itemStyle.Render(fmt.Sprintf("Upload: %s/s (Total: %s)\n",
		monitor.FormatSpeed(m.stats.Network.TxSpeed),
		monitor.FormatBytes(m.stats.Network.TxBytes))))
	b.WriteString("\n")

	// Load Average
	b.WriteString(categoryStyle.Render("⚡ Load Average"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("1min: %.2f | 5min: %.2f | 15min: %.2f\n",
		m.stats.Load.Load1, m.stats.Load.Load5, m.stats.Load.Load15)))
	b.WriteString("\n")

	// Uptime
	b.WriteString(categoryStyle.Render("⏱️  Uptime"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s\n", m.stats.Uptime)))
	b.WriteString("\n")

	// Footer
	b.WriteString(helpStyle.Render("r: Refresh | q: Back"))

	return b.String()
}

func renderBar(percent, max float64, label string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > max {
		percent = max
	}

	barWidth := 30
	filled := int((percent / max) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	var color lipgloss.Color
	switch {
	case percent < 50:
		color = lipgloss.Color("#04B575") // Green
	case percent < 80:
		color = lipgloss.Color("#FFA500") // Orange
	default:
		color = lipgloss.Color("#FF0000") // Red
	}

	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Width(barWidth)

	filledStyle := lipgloss.NewStyle().
		Background(color).
		Foreground(lipgloss.Color("#FFFFFF"))

	filledStr := strings.Repeat(" ", filled)
	emptyStr := strings.Repeat(" ", barWidth-filled)

	return fmt.Sprintf("[%s%s]", filledStyle.Render(filledStr), barStyle.Render(emptyStr))
}
